import {
  DndContext,
  PointerSensor,
  TouchSensor,
  closestCenter,
  useSensor,
  useSensors,
  type DragEndEvent,
} from '@dnd-kit/core'
import {
  SortableContext,
  arrayMove,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { Category, Todo } from '../types'
import { TodoItem } from './TodoItem'

type SortableItemProps = {
  todo: Todo
  categories: Category[]
  onToggle: (id: number) => void
  onUpdate: (id: number, text: string, categoryId: number | null) => void
  onDelete: (id: number) => void
}

function SortableTodoItem({ todo, categories, onToggle, onUpdate, onDelete }: SortableItemProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } =
    useSortable({ id: todo.id })

  return (
    <li
      ref={setNodeRef}
      style={{
        display: 'flex',
        alignItems: 'center',
        gap: 8,
        padding: '8px 0',
        borderBottom: '1px solid #eee',
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.5 : 1,
        background: isDragging ? '#f9f9f9' : 'transparent',
      }}
    >
      <TodoItem
        todo={todo}
        categories={categories}
        onToggle={onToggle}
        onUpdate={onUpdate}
        onDelete={onDelete}
      />
      <span
        {...attributes}
        {...listeners}
        style={{
          cursor: 'grab',
          padding: '4px 6px',
          color: '#ccc',
          touchAction: 'none',
          userSelect: 'none',
          fontSize: 18,
          lineHeight: 1,
          flexShrink: 0,
        }}
        aria-label="ドラッグして並び替え"
      >
        ⠿
      </span>
    </li>
  )
}

type Props = {
  todos: Todo[]
  categories: Category[]
  onToggle: (id: number) => void
  onUpdate: (id: number, text: string, categoryId: number | null) => void
  onDelete: (id: number) => void
  onDeleteDone: () => void
  onReorder: (ids: number[]) => void
}

export function TodoList({ todos, categories, onToggle, onUpdate, onDelete, onDeleteDone, onReorder }: Props) {
  const hasDone = todos.some((t) => t.done)

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
    useSensor(TouchSensor, { activationConstraint: { delay: 150, tolerance: 5 } }),
  )

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    if (!over || active.id === over.id) return
    const oldIndex = todos.findIndex((t) => t.id === active.id)
    const newIndex = todos.findIndex((t) => t.id === over.id)
    const reordered = arrayMove(todos, oldIndex, newIndex)
    onReorder(reordered.map((t) => t.id))
  }

  return (
    <>
      {todos.length === 0 ? (
        <p style={{ color: '#888' }}>Todo がありません</p>
      ) : (
        <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
          <SortableContext items={todos.map((t) => t.id)} strategy={verticalListSortingStrategy}>
            <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
              {todos.map((todo) => (
                <SortableTodoItem
                  key={todo.id}
                  todo={todo}
                  categories={categories}
                  onToggle={onToggle}
                  onUpdate={onUpdate}
                  onDelete={onDelete}
                />
              ))}
            </ul>
          </SortableContext>
        </DndContext>
      )}
      {hasDone && (
        <div style={{ marginTop: 16, textAlign: 'right' }}>
          <button
            onClick={onDeleteDone}
            style={{
              padding: '6px 14px',
              background: 'none',
              border: '1px solid #ccc',
              borderRadius: 6,
              cursor: 'pointer',
              color: '#888',
              fontSize: 13,
            }}
          >
            完了済みを削除
          </button>
        </div>
      )}
    </>
  )
}
