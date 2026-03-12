import { useRef, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";

export function ResetPasswordPage() {
  const [searchParams] = useSearchParams();
  const [newPassword, setNewPassword] = useState("");
  const [status, setStatus] = useState<"idle" | "success" | "error">("idle");
  const [error, setError] = useState("");
  const submitting = useRef(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (submitting.current) return;
    submitting.current = true;

    const token = searchParams.get("token") ?? "";
    const res = await fetch("/api/auth/password-reset/confirm", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ token, new_password: newPassword }),
    });
    if (res.ok) {
      setStatus("success");
    } else {
      const data = await res.json().catch(() => ({}));
      setError(data.message || "エラーが発生しました");
      setStatus("error");
    }
    submitting.current = false;
  };

  if (status === "success") {
    return (
      <div style={{ maxWidth: 400, margin: "80px auto", padding: "0 16px" }}>
        <h1 style={{ fontSize: 24, marginBottom: 16 }}>パスワードを変更しました</h1>
        <p><Link to="/login">ログイン</Link></p>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: 400, margin: "80px auto", padding: "0 16px" }}>
      <h1 style={{ fontSize: 24, marginBottom: 24 }}>パスワードの再設定</h1>
      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: "block", marginBottom: 4 }}>新しいパスワード</label>
          <input
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            required
            style={{ width: "100%", padding: "8px", border: "1px solid #ccc", borderRadius: 4, boxSizing: "border-box" }}
          />
        </div>
        {status === "error" && (
          <p style={{ color: "red", marginBottom: 16 }}>{error}</p>
        )}
        <button
          type="submit"
          style={{ width: "100%", padding: "10px", background: "#333", color: "#fff", border: "none", borderRadius: 4, cursor: "pointer" }}
        >
          変更する
        </button>
      </form>
    </div>
  );
}
