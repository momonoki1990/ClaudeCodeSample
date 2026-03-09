import { Todo } from '../types'

type Props = {
  todo: Todo
  onToggle: (id: number) => void
  onDelete: (id: number) => void
}

export function TodoItem({ todo, onToggle, onDelete }: Props) {
  return (
    <li
      style={{
        display: 'flex',
        alignItems: 'center',
        gap: 8,
        padding: '8px 0',
        borderBottom: '1px solid #eee',
      }}
    >
      <input
        type="checkbox"
        checked={todo.done}
        onChange={() => onToggle(todo.id)}
        style={{ cursor: 'pointer' }}
      />
      <span
        style={{
          flex: 1,
          textDecoration: todo.done ? 'line-through' : 'none',
          color: todo.done ? '#aaa' : 'inherit',
        }}
      >
        {todo.text}
      </span>
      <button onClick={() => onDelete(todo.id)} style={{ padding: '2px 10px', cursor: 'pointer' }}>
        削除
      </button>
    </li>
  )
}
