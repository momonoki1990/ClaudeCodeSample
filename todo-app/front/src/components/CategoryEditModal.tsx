import { useEffect, useState } from 'react'
import { Category } from '../types'

type Props = {
  category: Category
  onSave: (id: number, name: string) => void
  onDelete: (id: number) => void
  onClose: () => void
}

export function CategoryEditModal({ category, onSave, onDelete, onClose }: Props) {
  const [name, setName] = useState(category.name)

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [onClose])

  const handleSave = () => {
    const trimmed = name.trim()
    if (!trimmed) return
    onSave(category.id, trimmed)
    onClose()
  }

  const handleDelete = () => {
    onDelete(category.id)
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
        <h3 style={{ margin: '0 0 16px' }}>カテゴリを編集</h3>

        <label style={{ display: 'block', marginBottom: 24 }}>
          <div style={{ fontSize: 13, color: '#666', marginBottom: 4 }}>名前</div>
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            autoFocus
            onKeyDown={(e) => { if (e.key === 'Enter') handleSave() }}
            style={{ width: '100%', padding: '6px 8px', fontSize: 16, boxSizing: 'border-box' }}
          />
        </label>

        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          {!category.is_system ? (
            <button
              onClick={handleDelete}
              style={{ padding: '6px 16px', cursor: 'pointer', color: '#c00', background: 'none', border: '1px solid #c00', borderRadius: 4 }}
            >
              削除
            </button>
          ) : (
            <span style={{ fontSize: 12, color: '#c00' }}>標準設定のため削除できません</span>
          )}
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
