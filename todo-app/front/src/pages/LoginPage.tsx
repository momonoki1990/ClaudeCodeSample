import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";

export function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    const res = await fetch("/api/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    if (!res.ok) {
      setError("ログインに失敗しました");
      return;
    }
    navigate("/");
  };

  return (
    <div style={{ maxWidth: 400, margin: "80px auto", padding: "0 16px" }}>
      <h1 style={{ fontSize: 24, marginBottom: 24 }}>ログイン</h1>
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
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: "block", marginBottom: 4 }}>パスワード</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            style={{ width: "100%", padding: "8px", border: "1px solid #ccc", borderRadius: 4, boxSizing: "border-box" }}
          />
        </div>
        {error && <p style={{ color: "red", marginBottom: 16 }}>{error}</p>}
        <button
          type="submit"
          style={{ width: "100%", padding: "10px", background: "#333", color: "#fff", border: "none", borderRadius: 4, cursor: "pointer" }}
        >
          ログイン
        </button>
      </form>
      <p style={{ marginTop: 16, textAlign: "center" }}>
        アカウントがない方は <Link to="/register">新規登録</Link>
      </p>
      <p style={{ marginTop: 8, textAlign: "center" }}>
        <Link to="/forgot-password">パスワードをお忘れの方</Link>
      </p>
    </div>
  );
}
