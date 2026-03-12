import { useEffect, useMemo, useState } from "react";
import { Link, Route, Routes, useNavigate } from "react-router-dom";
import { Category, Todo } from "./types";
import { apiFetch } from "./api";
import { Drawer } from "./components/Drawer";
import { TodoForm } from "./components/TodoForm";
import { TodoList } from "./components/TodoList";
import { CategoryTabs } from "./components/CategoryTabs";

async function fetchTodos(categoryId: number | null): Promise<Todo[]> {
  const url = categoryId ? `/api/todos?category_id=${categoryId}` : "/api/todos";
  const res = await apiFetch(url);
  if (!res.ok) throw new Error("fetch failed");
  return res.json();
}

async function fetchCategories(): Promise<Category[]> {
  const res = await apiFetch("/api/categories");
  if (!res.ok) throw new Error("fetch failed");
  return res.json();
}

export default function App() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [filterCategoryId, setFilterCategoryId] = useState<number | null>(null);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const navigate = useNavigate();

  const tabs = useMemo(() => categories, [categories]);

  const reloadTodos = (categoryId: number | null = filterCategoryId) =>
    fetchTodos(categoryId).then(setTodos).catch(console.error);

  const reloadCategories = () =>
    fetchCategories().then(setCategories).catch(console.error);

  useEffect(() => {
    reloadTodos();
    reloadCategories();
  }, []);

  const handleLogout = async () => {
    await fetch("/api/auth/logout", { method: "POST" });
    navigate("/login");
  };

  const addTodo = async (text: string) => {
    await apiFetch("/api/todos", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ text, category_id: filterCategoryId }),
    });
    reloadTodos();
  };

  const toggleTodo = async (id: number) => {
    const todo = todos.find((t) => t.id === id);
    if (!todo) return;
    await apiFetch(`/api/todos/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ text: todo.text, done: !todo.done, category_id: todo.category_id }),
    });
    reloadTodos();
  };

  const updateTodo = async (id: number, text: string, categoryId: number | null) => {
    const todo = todos.find((t) => t.id === id);
    if (!todo) return;
    await apiFetch(`/api/todos/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ text, done: todo.done, category_id: categoryId }),
    });
    reloadTodos();
  };

  const deleteTodo = async (id: number) => {
    await apiFetch(`/api/todos/${id}`, { method: "DELETE" });
    reloadTodos();
  };

  const reorderTodos = async (ids: number[]) => {
    await apiFetch("/api/todos/reorder", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ids }),
    });
    reloadTodos();
  };

  const deleteDoneTodos = async () => {
    const url = filterCategoryId
      ? `/api/todos/done?category_id=${filterCategoryId}`
      : "/api/todos/done";
    await apiFetch(url, { method: "DELETE" });
    reloadTodos();
  };

  const renameCategory = async (id: number, name: string) => {
    await apiFetch(`/api/categories/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name }),
    });
    reloadCategories();
  };

  const addCategory = async (name: string) => {
    await apiFetch("/api/categories", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name }),
    });
    reloadCategories();
  };

  const deleteCategory = async (id: number) => {
    await apiFetch(`/api/categories/${id}`, { method: "DELETE" });
    reloadCategories();
    if (filterCategoryId === id) {
      setFilterCategoryId(null);
      reloadTodos(null);
    } else {
      reloadTodos();
    }
  };

  const reorderTabs = async (orderedIds: number[]) => {
    await apiFetch("/api/categories/reorder", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ids: orderedIds }),
    });
    reloadCategories();
  };

  const handleFilterChange = (categoryId: number | null) => {
    setFilterCategoryId(categoryId);
    reloadTodos(categoryId);
  };

  return (
    <div style={{ fontFamily: "sans-serif" }}>
      <header
        style={{
          position: "sticky",
          top: 0,
          background: "#fff",
          borderBottom: "1px solid #eee",
          padding: "12px 24px",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          zIndex: 10,
        }}
      >
        <Link to="/" style={{ fontSize: 20, fontWeight: "bold", textDecoration: "none", color: "inherit" }}>Todo</Link>
        <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
          <button
            onClick={handleLogout}
            style={{ background: "none", border: "1px solid #ccc", borderRadius: 4, cursor: "pointer", padding: "4px 12px", fontSize: 14 }}
          >
            ログアウト
          </button>
          <button
            onClick={() => setDrawerOpen(true)}
            style={{ background: "none", border: "none", cursor: "pointer", padding: 4 }}
            aria-label="メニューを開く"
          >
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="3" y1="6" x2="21" y2="6" />
              <line x1="3" y1="12" x2="21" y2="12" />
              <line x1="3" y1="18" x2="21" y2="18" />
            </svg>
          </button>
        </div>
      </header>

      <Drawer isOpen={drawerOpen} onClose={() => setDrawerOpen(false)} />

      <main style={{ maxWidth: 480, margin: "32px auto", padding: "0 16px" }}>
        <Routes>
          <Route
            path="/"
            element={
              <>
                <TodoForm onAdd={addTodo} />
                <CategoryTabs
                  tabs={tabs}
                  filterCategoryId={filterCategoryId}
                  onFilterChange={handleFilterChange}
                  onAdd={addCategory}
                  onRename={renameCategory}
                  onDelete={deleteCategory}
                  onReorder={reorderTabs}
                />
                <TodoList
                  todos={todos}
                  categories={categories}
                  currentCategoryName={tabs.find((t) => t.id === filterCategoryId)?.name ?? 'すべて'}
                  onToggle={toggleTodo}
                  onUpdate={updateTodo}
                  onDelete={deleteTodo}
                  onDeleteDone={deleteDoneTodos}
                  onReorder={reorderTodos}
                />
              </>
            }
          />
        </Routes>
      </main>
    </div>
  );
}
