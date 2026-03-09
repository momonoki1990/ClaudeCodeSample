import { Todo } from '../types'
import { TodoItem } from './TodoItem'

type Props = {
  todos: Todo[]
  onToggle: (id: number) => void
  onDelete: (id: number) => void
}

export function TodoList({ todos, onToggle, onDelete }: Props) {
  if (todos.length === 0) {
    return <p style={{ color: '#888' }}>Todo がありません</p>
  }

  return (
    <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
      {todos.map((todo) => (
        <TodoItem key={todo.id} todo={todo} onToggle={onToggle} onDelete={onDelete} />
      ))}
    </ul>
  )
}
