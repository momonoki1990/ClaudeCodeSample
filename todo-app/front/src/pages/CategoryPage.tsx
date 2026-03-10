import { useState } from 'react'
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
import { Category } from '../types'
import { CategoryEditModal } from '../components/CategoryEditModal'

type Props = {
  tabs: Category[]
  onAdd: (name: string) => void
  onRename: (id: number, name: string) => void
  onDelete: (id: number) => void
  onReorder: (ids: number[]) => void
}

type SortableItemProps = {
  category: Category
  onEdit: (category: Category) => void
}

function SortableItem({ category, onEdit }: SortableItemProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } =
    useSortable({ id: category.id })

  return (
    <li
      ref={setNodeRef}
      style={{
        display: 'flex',
        alignItems: 'center',
        gap: 8,
        padding: '10px 0',
        borderBottom: '1px solid #eee',
        fontSize: 16,
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.5 : 1,
        background: isDragging ? '#f9f9f9' : 'transparent',
      }}
    >
      <span
        onClick={() => onEdit(category)}
        style={{ flex: 1, cursor: 'pointer' }}
      >
        {category.name}
      </span>
      <span
        {...attributes}
        {...listeners}
        style={{
          cursor: 'grab',
          padding: '4px 6px',
          color: '#bbb',
          touchAction: 'none',
          userSelect: 'none',
          fontSize: 18,
          lineHeight: 1,
        }}
        aria-label="ドラッグして並び替え"
      >
        ⠿
      </span>
    </li>
  )
}

export function CategoryPage({ tabs, onAdd, onRename, onDelete, onReorder }: Props) {
  const [editing, setEditing] = useState<Category | null>(null)

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
    useSensor(TouchSensor, { activationConstraint: { delay: 150, tolerance: 5 } }),
  )

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const input = e.currentTarget.elements.namedItem('name') as HTMLInputElement
    const trimmed = input.value.trim()
    if (!trimmed) return
    onAdd(trimmed)
    input.value = ''
  }

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    if (!over || active.id === over.id) return
    const oldIndex = tabs.findIndex((t) => t.id === active.id)
    const newIndex = tabs.findIndex((t) => t.id === over.id)
    const reordered = arrayMove(tabs, oldIndex, newIndex)
    onReorder(reordered.map((t) => t.id))
  }

  return (
    <div>
      <h2 style={{ marginTop: 0 }}>カテゴリ管理</h2>
      <form onSubmit={handleSubmit} style={{ display: 'flex', gap: 8, marginBottom: 24 }}>
        <input
          name="name"
          placeholder="カテゴリ名を入力"
          style={{ flex: 1, padding: '6px 8px', fontSize: 16 }}
        />
        <button type="submit" style={{ padding: '6px 16px', fontSize: 16 }}>
          追加
        </button>
      </form>
      {tabs.length === 0 ? (
        <p style={{ color: '#888' }}>カテゴリがありません</p>
      ) : (
        <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
          <SortableContext items={tabs.map((t) => t.id)} strategy={verticalListSortingStrategy}>
            <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
              {tabs.map((tab) => (
                <SortableItem key={tab.id} category={tab} onEdit={setEditing} />
              ))}
            </ul>
          </SortableContext>
        </DndContext>
      )}
      {editing && (
        <CategoryEditModal
          category={editing}
          onSave={onRename}
          onDelete={onDelete}
          onClose={() => setEditing(null)}
        />
      )}
    </div>
  )
}
