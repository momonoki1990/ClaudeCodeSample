import { useState } from 'react'
import { Category } from '../types'

type Props = {
  categories: Category[]
  onAdd: (name: string) => void
  onDelete: (id: number) => void
}

export function CategoryManager({ categories, onAdd, onDelete }: Props) {
  const [name, setName] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const trimmed = name.trim()
    if (!trimmed) return
    onAdd(trimmed)
    setName('')
  }

  return (
    <div style={{ marginBottom: 24, padding: 12, background: '#f9f9f9', borderRadius: 8 }}>
      <h2 style={{ fontSize: 16, marginBottom: 8, marginTop: 0 }}>カテゴリ管理</h2>
      <form onSubmit={handleSubmit} style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="カテゴリ名を入力"
          style={{ flex: 1, padding: '4px 8px', fontSize: 14 }}
        />
        <button type="submit" style={{ padding: '4px 12px', fontSize: 14 }}>
          追加
        </button>
      </form>
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6 }}>
        {categories.map((cat) => (
          <span
            key={cat.id}
            style={{
              display: 'inline-flex',
              alignItems: 'center',
              gap: 4,
              padding: '2px 8px',
              background: '#e0e7ff',
              borderRadius: 12,
              fontSize: 13,
            }}
          >
            {cat.name}
            <button
              onClick={() => onDelete(cat.id)}
              style={{ background: 'none', border: 'none', cursor: 'pointer', padding: 0, fontSize: 13, color: '#666' }}
            >
              ×
            </button>
          </span>
        ))}
        {categories.length === 0 && <span style={{ color: '#aaa', fontSize: 13 }}>カテゴリがありません</span>}
      </div>
    </div>
  )
}
