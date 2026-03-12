import { useEffect, useState } from 'react'
import { ConfirmModal } from './ConfirmModal'
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
  horizontalListSortingStrategy,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { Category } from '../types'

type Props = {
  tabs: Category[]
  filterCategoryId: number | null
  onFilterChange: (id: number | null) => void
  onAdd: (name: string) => void
  onRename: (id: number, name: string) => void
  onDelete: (id: number) => void
  onReorder: (ids: number[]) => void
}

type TabItemProps = {
  tab: Category
  isActive: boolean
  isEditing: boolean
  onFilter: () => void
  onStartEdit: () => void
  onFinishEdit: (name: string) => void
  onDelete: () => void
}

function TabItem({ tab, isActive, isEditing, onFilter, onStartEdit, onFinishEdit, onDelete }: TabItemProps) {
  const [name, setName] = useState(tab.name)
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: tab.id,
  })

  useEffect(() => { setName(tab.name) }, [tab.name])

  return (
    <div
      ref={setNodeRef}
      style={{
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.4 : 1,
        display: 'inline-flex',
        alignItems: 'center',
        marginBottom: -2,
      }}
    >
      <span
        {...attributes}
        {...listeners}
        onDoubleClick={tab.is_system ? undefined : onStartEdit}
        onClick={onFilter}
        style={{
          padding: '8px 6px 8px 16px',
          cursor: 'grab',
          borderBottom: isActive ? '2px solid #333' : '2px solid transparent',
          fontWeight: isActive ? 'bold' : 'normal',
          whiteSpace: 'nowrap',
          fontSize: 14,
          userSelect: 'none',
          touchAction: 'none',
        }}
      >
        {isEditing ? (
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            onBlur={() => onFinishEdit(name.trim() || tab.name)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') (e.target as HTMLInputElement).blur()
              if (e.key === 'Escape') { setName(tab.name); (e.target as HTMLInputElement).blur() }
            }}
            onClick={(e) => e.stopPropagation()}
            autoFocus
            style={{ width: 80, padding: '2px 4px', fontSize: 14, border: '1px solid #aaa', borderRadius: 3 }}
          />
        ) : (
          tab.name
        )}
      </span>
      {!tab.is_system && (
        <button
          onClick={(e) => { e.stopPropagation(); onDelete() }}
          style={{
            background: 'none',
            border: 'none',
            cursor: 'pointer',
            padding: '0 10px 0 0',
            color: '#ccc',
            fontSize: 14,
            lineHeight: 1,
            marginBottom: -2,
            borderBottom: isActive ? '2px solid #333' : '2px solid transparent',
          }}
          aria-label={`${tab.name}を削除`}
        >
          ×
        </button>
      )}
    </div>
  )
}

export function CategoryTabs({ tabs, filterCategoryId, onFilterChange, onAdd, onRename, onDelete, onReorder }: Props) {
  const [editingId, setEditingId] = useState<number | null>(null)
  const [addingNew, setAddingNew] = useState(false)
  const [newName, setNewName] = useState('')
  const [confirmTarget, setConfirmTarget] = useState<{ id: number; name: string } | null>(null)

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
    useSensor(TouchSensor, { activationConstraint: { delay: 250, tolerance: 5 } }),
  )

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    if (!over || active.id === over.id) return
    const oldIndex = tabs.findIndex((t) => t.id === active.id)
    const newIndex = tabs.findIndex((t) => t.id === over.id)
    const reordered = arrayMove(tabs, oldIndex, newIndex)
    onReorder(reordered.map((t) => t.id))
  }

  const handleAddSubmit = () => {
    const trimmed = newName.trim()
    if (trimmed) onAdd(trimmed)
    setNewName('')
    setAddingNew(false)
  }

  const handleDelete = (id: number, name: string) => {
    setConfirmTarget({ id, name })
  }

  return (
    <div style={{ display: 'flex', alignItems: 'center', borderBottom: '2px solid #eee', marginBottom: 16 }}>
      <div style={{ flex: 1, overflowX: 'auto', display: 'flex', alignItems: 'center' }}>
        <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
          <SortableContext items={tabs.map((t) => t.id)} strategy={horizontalListSortingStrategy}>
            <div style={{ display: 'inline-flex' }}>
              {tabs.map((tab) => (
                <TabItem
                  key={tab.id}
                  tab={tab}
                  isActive={tab.is_system ? filterCategoryId === null : filterCategoryId === tab.id}
                  isEditing={editingId === tab.id}
                  onFilter={() => onFilterChange(tab.is_system ? null : tab.id)}
                  onStartEdit={() => setEditingId(tab.id)}
                  onFinishEdit={(name) => { onRename(tab.id, name); setEditingId(null) }}
                  onDelete={() => handleDelete(tab.id, tab.name)}
                />
              ))}
            </div>
          </SortableContext>
        </DndContext>

        {addingNew ? (
          <input
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            onBlur={handleAddSubmit}
            onKeyDown={(e) => {
              if (e.key === 'Enter') handleAddSubmit()
              if (e.key === 'Escape') { setNewName(''); setAddingNew(false) }
            }}
            autoFocus
            placeholder="カテゴリ名"
            style={{ width: 100, padding: '4px 6px', fontSize: 14, border: '1px solid #aaa', borderRadius: 4, marginLeft: 4 }}
          />
        ) : (
          <button
            onClick={() => setAddingNew(true)}
            style={{
              background: 'none',
              border: 'none',
              cursor: 'pointer',
              padding: '8px 10px',
              fontSize: 16,
              color: '#aaa',
              lineHeight: 1,
            }}
            aria-label="カテゴリを追加"
          >
            ＋
          </button>
        )}
      </div>
      {confirmTarget && (
        <ConfirmModal
          title="カテゴリを削除"
          message={`「${confirmTarget.name}」を削除しますか？\nこのカテゴリのタスクもすべて削除されます。`}
          onConfirm={() => { onDelete(confirmTarget.id); setConfirmTarget(null) }}
          onCancel={() => setConfirmTarget(null)}
        />
      )}
    </div>
  )
}
