package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/mail"
	"os"
	"time"
	"todo-api/mailer"
	"todo-api/model"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ValidationError はバリデーションエラーを表す型。ハンドラーで 400 を返すために使用する。
type ValidationError string

func (e ValidationError) Error() string { return string(e) }

func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return ValidationError("メールアドレスの形式が正しくありません")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ValidationError("パスワードは8文字以上で入力してください")
	}
	if len(password) > 72 {
		return ValidationError("パスワードは72文字以下で入力してください")
	}
	return nil
}

type UserService struct {
	repo          model.UserRepository
	refreshRepo   model.RefreshTokenRepository
	verifyRepo    model.EmailVerificationTokenRepository
	passwordReset model.PasswordResetTokenRepository
	mailer        mailer.Mailer
	appURL        string
}

func NewUserService(
	repo model.UserRepository,
	refreshRepo model.RefreshTokenRepository,
	verifyRepo model.EmailVerificationTokenRepository,
	passwordReset model.PasswordResetTokenRepository,
	m mailer.Mailer,
	appURL string,
) *UserService {
	return &UserService{repo: repo, refreshRepo: refreshRepo, verifyRepo: verifyRepo, passwordReset: passwordReset, mailer: m, appURL: appURL}
}

func (s *UserService) Register(email, password string) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user, err := s.repo.Create(email, string(hash))
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ValidationError("このメールアドレスはすでに登録されています")
		}
		return err
	}

	rawToken, err := generateRawToken()
	if err != nil {
		return err
	}
	tokenHash := hashToken(rawToken)
	if err := s.verifyRepo.Create(user.ID, tokenHash, time.Now().Add(24*time.Hour)); err != nil {
		return err
	}

	link := s.appURL + "/verify-email?token=" + rawToken
	_ = s.mailer.Send(email, "メールアドレスの確認", "以下のリンクをクリックして登録を完了してください:\n\n"+link)
	return nil
}

func (s *UserService) VerifyEmail(rawToken string) error {
	tokenHash := hashToken(rawToken)
	evt, err := s.verifyRepo.FindByHash(tokenHash)
	if err != nil {
		return errors.New("invalid or expired token")
	}
	if err := s.repo.SetEmailVerified(evt.UserID); err != nil {
		return err
	}
	return s.verifyRepo.DeleteByHash(tokenHash)
}

func (s *UserService) Login(email, password string) (accessToken, refreshToken string, err error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}
	if !user.EmailVerified {
		return "", "", errors.New("email not verified")
	}

	accessToken, err = s.generateAccessToken(user.ID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.createRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *UserService) Refresh(rawRefreshToken string) (newAccess, newRefresh string, err error) {
	tokenHash := hashToken(rawRefreshToken)
	rt, err := s.refreshRepo.FindByHash(tokenHash)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	if err := s.refreshRepo.DeleteByHash(tokenHash); err != nil {
		return "", "", err
	}

	newAccess, err = s.generateAccessToken(rt.UserID)
	if err != nil {
		return "", "", err
	}

	newRefresh, err = s.createRefreshToken(rt.UserID)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}

func (s *UserService) Logout(rawRefreshToken string) error {
	tokenHash := hashToken(rawRefreshToken)
	return s.refreshRepo.DeleteByHash(tokenHash)
}

func (s *UserService) generateAccessToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

func (s *UserService) createRefreshToken(userID int) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	rawHex := hex.EncodeToString(raw)
	tokenHash := hashToken(rawHex)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.refreshRepo.Create(userID, tokenHash, expiresAt); err != nil {
		return "", err
	}
	return rawHex, nil
}

func (s *UserService) RequestPasswordReset(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		// ユーザーが存在しない場合もエラーを返さない（メールアドレス列挙防止）
		return nil
	}
	// 既存トークンを削除（1ユーザー1トークン）
	_ = s.passwordReset.DeleteByUserID(user.ID)

	rawToken, err := generateRawToken()
	if err != nil {
		return err
	}
	tokenHash := hashToken(rawToken)
	if err := s.passwordReset.Create(user.ID, tokenHash, time.Now().Add(1*time.Hour)); err != nil {
		return err
	}

	link := s.appURL + "/reset-password?token=" + rawToken
	_ = s.mailer.Send(email, "パスワードの再設定", "以下のリンクからパスワードを再設定してください（有効期限: 1時間）:\n\n"+link)
	return nil
}

func (s *UserService) ConfirmPasswordReset(rawToken, newPassword string) error {
	if err := validatePassword(newPassword); err != nil {
		return err
	}
	tokenHash := hashToken(rawToken)
	prt, err := s.passwordReset.FindByHash(tokenHash)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.repo.UpdatePassword(prt.UserID, string(hash)); err != nil {
		return err
	}
	return s.passwordReset.DeleteByHash(tokenHash)
}

func generateRawToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
