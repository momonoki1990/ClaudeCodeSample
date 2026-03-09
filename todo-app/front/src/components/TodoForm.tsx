import { useState } from 'react'

type Props = {
  onAdd: (text: string) => void
}

export function TodoForm({ onAdd }: Props) {
  const [text, setText] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const trimmed = text.trim()
    if (!trimmed) return
    onAdd(trimmed)
    setText('')
  }

  return (
    <form onSubmit={handleSubmit} style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
      <input
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="新しい Todo を入力"
        style={{ flex: 1, padding: '6px 8px', fontSize: 16 }}
      />
      <button type="submit" style={{ padding: '6px 16px', fontSize: 16 }}>
        追加
      </button>
    </form>
  )
}
