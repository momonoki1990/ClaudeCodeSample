import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { Category, Todo } from '../types'

type Props = {
  todo: Todo
  categories: Category[]
  onSave: (id: number, text: string, categoryId: number | null) => void
  onDelete: (id: number) => void
  onClose: () => void
}

export function TodoEditModal({ todo, categories, onSave, onDelete, onClose }: Props) {
  const [text, setText] = useState(todo.text)
  const [categoryId, setCategoryId] = useState<number | null>(todo.category_id)

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [onClose])

  const handleSave = () => {
    const trimmed = text.trim()
    if (!trimmed) return
    onSave(todo.id, trimmed, categoryId)
    onClose()
  }

  const handleDelete = () => {
    onDelete(todo.id)
    onClose()
  }

  return (
    <div
      onClick={onClose}
      style={{
        position: 'fixed',
        inset: 0,
        background: 'rgba(0,0,0,0.4)',
        zIndex: 200,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        padding: 16,
      }}
    >
      <div
        onClick={(e) => e.stopPropagation()}
        style={{
          background: '#fff',
          borderRadius: 8,
          padding: 24,
          width: '100%',
          maxWidth: 420,
          boxShadow: '0 4px 24px rgba(0,0,0,0.2)',
        }}
      >
        <h3 style={{ margin: '0 0 16px' }}>タスクを編集</h3>

        <label style={{ display: 'block', marginBottom: 12 }}>
          <div style={{ fontSize: 13, color: '#666', marginBottom: 4 }}>名前</div>
          <input
            value={text}
            onChange={(e) => setText(e.target.value)}
            autoFocus
            style={{ width: '100%', padding: '6px 8px', fontSize: 16, boxSizing: 'border-box' }}
          />
        </label>

        <label style={{ display: 'block', marginBottom: 24 }}>
          <div style={{ fontSize: 13, color: '#666', marginBottom: 4 }}>カテゴリ</div>
          <select
            value={categoryId ?? ''}
            onChange={(e) => setCategoryId(e.target.value === '' ? null : Number(e.target.value))}
            style={{ width: '100%', padding: '6px 8px', fontSize: 16, boxSizing: 'border-box' }}
          >
            <option value="">カテゴリなし</option>
            {categories.filter((cat) => !cat.is_system).map((cat) => (
              <option key={cat.id} value={cat.id}>{cat.name}</option>
            ))}
          </select>
          <Link
            to="/categories"
            onClick={onClose}
            style={{ display: 'inline-block', marginTop: 6, fontSize: 13, color: '#555' }}
          >
            + カテゴリを追加
          </Link>
        </label>

        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <button
            onClick={handleDelete}
            style={{ padding: '6px 16px', cursor: 'pointer', color: '#c00', background: 'none', border: '1px solid #c00', borderRadius: 4 }}
          >
            削除
          </button>
          <div style={{ display: 'flex', gap: 8 }}>
            <button
              onClick={onClose}
              style={{ padding: '6px 16px', cursor: 'pointer', background: 'none', border: '1px solid #ccc', borderRadius: 4 }}
            >
              キャンセル
            </button>
            <button
              onClick={handleSave}
              style={{ padding: '6px 16px', cursor: 'pointer', background: '#333', color: '#fff', border: 'none', borderRadius: 4 }}
            >
              保存
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
