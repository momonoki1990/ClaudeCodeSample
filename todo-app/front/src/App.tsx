import { useEffect, useState } from "react";
import { Todo } from "./types";
import { TodoForm } from "./components/TodoForm";
import { TodoList } from "./components/TodoList";

async function fetchTodos(): Promise<Todo[]> {
  const res = await fetch("/api/todos");
  if (!res.ok) throw new Error("fetch failed");
  return res.json();
}

export default function App() {
  const [todos, setTodos] = useState<Todo[]>([]);

  const reload = () => fetchTodos().then(setTodos).catch(console.error);

  useEffect(() => {
    reload();
  }, []);

  const addTodo = async (text: string) => {
    await fetch("/api/todos", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ text }),
    });
    reload();
  };

  const toggleTodo = async (id: number) => {
    const todo = todos.find((t) => t.id === id);
    if (!todo) return;
    await fetch(`/api/todos/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ done: !todo.done }),
    });
    reload();
  };

  const deleteTodo = async (id: number) => {
    await fetch(`/api/todos/${id}`, { method: "DELETE" });
    reload();
  };

  return (
    <div
      style={{ maxWidth: 480, margin: "40px auto", fontFamily: "sans-serif" }}
    >
      <h1>Todo</h1>
      <TodoForm onAdd={addTodo} />
      <TodoList todos={todos} onToggle={toggleTodo} onDelete={deleteTodo} />
    </div>
  );
}
