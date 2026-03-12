import { useEffect, useRef, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";

export function VerifyEmailPage() {
  const [searchParams] = useSearchParams();
  const [status, setStatus] = useState<"loading" | "success" | "error">("loading");
  const called = useRef(false);

  useEffect(() => {
    if (called.current) return;
    called.current = true;

    const token = searchParams.get("token");
    if (!token) {
      setStatus("error");
      return;
    }
    fetch(`/api/auth/verify-email?token=${encodeURIComponent(token)}`)
      .then((res) => {
        if (res.ok) {
          setStatus("success");
        } else {
          setStatus("error");
        }
      })
      .catch(() => setStatus("error"));
  }, []);

  if (status === "loading") {
    return <div style={{ maxWidth: 400, margin: "80px auto", padding: "0 16px" }}>確認中...</div>;
  }

  if (status === "success") {
    return (
      <div style={{ maxWidth: 400, margin: "80px auto", padding: "0 16px" }}>
        <h1 style={{ fontSize: 24, marginBottom: 16 }}>認証完了</h1>
        <p>メールアドレスの確認が完了しました。</p>
        <p style={{ marginTop: 16 }}>
          <Link to="/login">ログイン</Link>
        </p>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: 400, margin: "80px auto", padding: "0 16px" }}>
      <h1 style={{ fontSize: 24, marginBottom: 16 }}>認証エラー</h1>
      <p>リンクが無効または期限切れです。</p>
    </div>
  );
}
