import { useState } from "react";
import { Link } from "react-router-dom";

export function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [sent, setSent] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await fetch("/api/auth/password-reset/request", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email }),
    });
    setSent(true);
  };

  if (sent) {
    return (
      <div style={{ maxWidth: 400, margin: "80px auto", padding: "0 16px" }}>
        <h1 style={{ fontSize: 24, marginBottom: 16 }}>メールを送信しました</h1>
        <p>パスワード再設定用のリンクをお送りしました。メールをご確認ください。</p>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: 400, margin: "80px auto", padding: "0 16px" }}>
      <h1 style={{ fontSize: 24, marginBottom: 24 }}>パスワードを忘れた方</h1>
      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: "block", marginBottom: 4 }}>メールアドレス</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            style={{ width: "100%", padding: "8px", border: "1px solid #ccc", borderRadius: 4, boxSizing: "border-box" }}
          />
        </div>
        <button
          type="submit"
          style={{ width: "100%", padding: "10px", background: "#333", color: "#fff", border: "none", borderRadius: 4, cursor: "pointer" }}
        >
          送信
        </button>
      </form>
      <p style={{ marginTop: 16, textAlign: "center" }}>
        <Link to="/login">ログインに戻る</Link>
      </p>
    </div>
  );
}
