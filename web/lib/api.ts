export interface SectionResp {
  component: string
  props: Record<string, unknown>
}

export interface RenderResp {
  slug: string
  status: string
  lang: string
  sections: SectionResp[]
}

export async function fetchRender(slug: string): Promise<RenderResp> {
  const base = process.env.API_BASE_URL || 'http://localhost:8080'
  const res = await fetch(`${base}/api/v1/render/${slug}`, {
    next: { revalidate: 60 },
  })
  if (!res.ok) {
    if (res.status === 404) throw new Error('NOT_FOUND')
    throw new Error(`fetch render: ${res.status}`)
  }
  return res.json()
}
