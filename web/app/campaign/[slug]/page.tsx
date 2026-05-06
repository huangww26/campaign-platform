import { fetchRender } from '@/lib/api'
import { get } from '@/lib/registry'
import { notFound } from 'next/navigation'
import Unknown from '@/components/Unknown'

export default async function Page({
  params,
}: {
  params: Promise<{ slug: string }>
}) {
  const { slug } = await params

  let data
  try {
    data = await fetchRender(slug)
  } catch (e: unknown) {
    if (e instanceof Error && e.message === 'NOT_FOUND') notFound()
    throw e
  }

  return (
    <main lang={data.lang}>
      {data.sections.map((s, i) => {
        const Comp = get(s.component)
        if (!Comp) return <Unknown key={i} name={s.component} />
        return <Comp key={i} props={s.props} ctx={{ lang: data.lang }} />
      })}
    </main>
  )
}
