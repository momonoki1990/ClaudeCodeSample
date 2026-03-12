type Props = {
  title?: string
  message: string
  confirmLabel?: string
  onConfirm: () => void
  onCancel: () => void
}

export function ConfirmModal({ title = '確認', message, confirmLabel = '削除', onConfirm, onCancel }: Props) {
  return (
    <div
      onClick={onCancel}
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
          maxWidth: 360,
          boxShadow: '0 4px 24px rgba(0,0,0,0.2)',
        }}
      >
        <h3 style={{ margin: '0 0 16px' }}>{title}</h3>
        <p style={{ margin: '0 0 24px', fontSize: 15, lineHeight: 1.6, whiteSpace: 'pre-wrap' }}>{message}</p>
        <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
          <button
            onClick={onCancel}
            style={{ padding: '6px 16px', cursor: 'pointer', background: 'none', border: '1px solid #ccc', borderRadius: 4 }}
          >
            キャンセル
          </button>
          <button
            onClick={onConfirm}
            style={{ padding: '6px 16px', cursor: 'pointer', background: '#c00', color: '#fff', border: 'none', borderRadius: 4 }}
          >
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  )
}
