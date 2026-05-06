import { ActivityComponent, register } from '@/lib/registry'

interface HBProps {
  bg_image?: string
  title?: string | Record<string, string>
  subtitle?: string | Record<string, string>
  cta_text?: string | Record<string, string>
  cta_link?: string
}

function txt(v: string | Record<string, string> | undefined, lang: string): string {
  if (!v) return ''
  if (typeof v === 'string') return v
  return v[lang] || v['en'] || ''
}

const HeroBanner: ActivityComponent = ({ props, ctx }) => {
  const p = props as HBProps
  return (
    <section style={{
      position: 'relative', width: '100%', minHeight: 400,
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      color: '#fff', overflow: 'hidden',
    }}>
      {p.bg_image && <img src={p.bg_image} alt="" style={{
        position: 'absolute', width: '100%', height: '100%', objectFit: 'cover',
      }} />}
      <div style={{
        position: 'relative', zIndex: 1, textAlign: 'center', padding: 40,
        background: 'rgba(0,0,0,0.3)', borderRadius: 12,
      }}>
        <h1 style={{ fontSize: 36, margin: '0 0 16px' }}>{txt(p.title, ctx.lang)}</h1>
        {p.subtitle && <p style={{ fontSize: 18, margin: '0 0 24px', opacity: 0.9 }}>{txt(p.subtitle, ctx.lang)}</p>}
        {p.cta_text && <a href={p.cta_link || '#'} style={{
          display: 'inline-block', padding: '12px 32px',
          background: '#007AFF', color: '#fff', borderRadius: 8,
          textDecoration: 'none', fontSize: 16, fontWeight: 600,
        }}>{txt(p.cta_text, ctx.lang)}</a>}
      </div>
    </section>
  )
}

register('hero_banner', HeroBanner)
export default HeroBanner
