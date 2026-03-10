import { Link } from 'react-router-dom'

type Props = {
  isOpen: boolean
  onClose: () => void
}

export function Drawer({ isOpen, onClose }: Props) {
  return (
    <>
      <div
        onClick={onClose}
        style={{
          position: 'fixed',
          inset: 0,
          background: 'rgba(0,0,0,0.3)',
          zIndex: 100,
          opacity: isOpen ? 1 : 0,
          pointerEvents: isOpen ? 'auto' : 'none',
          transition: 'opacity 0.25s ease',
        }}
      />
      <div
        style={{
          position: 'fixed',
          top: 0,
          right: 0,
          bottom: 0,
          width: 240,
          background: '#fff',
          boxShadow: '-2px 0 8px rgba(0,0,0,0.15)',
          transform: isOpen ? 'translateX(0)' : 'translateX(100%)',
          transition: 'transform 0.25s ease',
          zIndex: 101,
          display: 'flex',
          flexDirection: 'column',
        }}
      >
        <div style={{ display: 'flex', justifyContent: 'flex-end', padding: '12px 16px' }}>
          <button
            onClick={onClose}
            style={{ background: 'none', border: 'none', fontSize: 24, cursor: 'pointer', lineHeight: 1 }}
          >
            ×
          </button>
        </div>
        <nav style={{ display: 'flex', flexDirection: 'column' }}>
          {[
            { to: '/', label: 'HOME' },
            { to: '/categories', label: 'カテゴリ管理' },
          ].map(({ to, label }) => (
            <Link
              key={to}
              to={to}
              onClick={onClose}
              style={{
                padding: '14px 24px',
                textDecoration: 'none',
                color: '#333',
                fontSize: 16,
                borderBottom: '1px solid #eee',
              }}
            >
              {label}
            </Link>
          ))}
        </nav>
      </div>
    </>
  )
}
