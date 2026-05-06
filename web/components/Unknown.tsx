export default function Unknown({ name }: { name: string }) {
  return (
    <div style={{
      padding: 40, textAlign: 'center', color: '#999',
      border: '2px dashed #ddd', borderRadius: 8, margin: 20,
    }}>
      <p>⚠️ Unknown component: <strong>{name}</strong></p>
    </div>
  )
}
