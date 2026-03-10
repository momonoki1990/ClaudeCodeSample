import { useState } from 'react'
import { Category, Todo } from '../types'
import { TodoEditModal } from './TodoEditModal'

type Props = {
  todo: Todo
  categories: Category[]
  onToggle: (id: number) => void
  onUpdate: (id: number, text: string, categoryId: number | null) => void
  onDelete: (id: number) => void
}

export function TodoItem({ todo, categories, onToggle, onUpdate, onDelete }: Props) {
  const [modalOpen, setModalOpen] = useState(false)

  return (
    <>
      <input
        type="checkbox"
        checked={todo.done}
        onChange={() => onToggle(todo.id)}
        style={{ cursor: 'pointer', flexShrink: 0 }}
      />
      <span
        onClick={() => setModalOpen(true)}
        style={{
          flex: 1,
          textDecoration: todo.done ? 'line-through' : 'none',
          color: todo.done ? '#aaa' : 'inherit',
          cursor: 'pointer',
        }}
      >
        {todo.text}
      </span>
      {modalOpen && (
        <TodoEditModal
          todo={todo}
          categories={categories}
          onSave={onUpdate}
          onDelete={onDelete}
          onClose={() => setModalOpen(false)}
        />
      )}
    </>
  )
}
