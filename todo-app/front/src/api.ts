let isRefreshing = false;

export async function apiFetch(path: string, options: RequestInit = {}): Promise<Response> {
  const res = await fetch(path, options);
  if (res.status === 401 && !isRefreshing) {
    isRefreshing = true;
    try {
      const refreshRes = await fetch("/api/auth/refresh", { method: "POST" });
      if (!refreshRes.ok) {
        window.location.href = "/login";
        throw new Error("Unauthorized");
      }
      return fetch(path, options);
    } finally {
      isRefreshing = false;
    }
  }
  if (res.status === 401) {
    window.location.href = "/login";
    throw new Error("Unauthorized");
  }
  return res;
}
