import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { ArrowLeft, ArrowRight, CalendarDays } from 'lucide-react'
import { useState } from 'react'
import { articleKeys, listArticles } from './api'
import { PageHeader } from './PageHeader'

const pageSize = 5

export function PreviewPage() {
  const [page, setPage] = useState(1)
  const offset = (page - 1) * pageSize
  const query = useQuery({ queryKey: articleKeys.list({ status: 'publish', limit: pageSize, offset }), queryFn: () => listArticles('publish', pageSize, offset), placeholderData: keepPreviousData })
  return (
    <>
      <PageHeader title="Published preview" description="A reader-facing view of published articles." />
      {query.isPending ? <PreviewState>Loading published articles...</PreviewState> : query.isError ? <PreviewState tone="error"><span>Failed to load published articles.</span><button className="font-semibold" onClick={() => query.refetch()}>Retry</button></PreviewState> : query.data.length === 0 && page === 1 ? <PreviewState>No published articles found.</PreviewState> : (
        <>
          <div className="divide-y divide-line border-y border-line bg-panel">
            {query.data.map((article) => (
              <article key={article.id} className="px-5 py-8 sm:px-8">
                <div className="mb-3 flex flex-wrap items-center gap-3 text-xs font-medium text-zinc-500"><span className="bg-emerald-50 px-2 py-1 text-brand">{article.category}</span><span className="inline-flex items-center gap-1.5"><CalendarDays size={14} /> {new Date(article.created_date).toLocaleDateString(undefined, { dateStyle: 'long' })}</span></div>
                <h2 className="text-xl font-semibold sm:text-2xl">{article.title}</h2>
                <p className="mt-4 whitespace-pre-wrap text-[15px] leading-7 text-zinc-700">{article.content}</p>
              </article>
            ))}
          </div>
          <nav aria-label="Preview pagination" className="mt-5 flex items-center justify-between">
            <button disabled={page === 1 || query.isFetching} onClick={() => setPage((value) => value - 1)} className="inline-flex h-10 items-center gap-2 border border-line bg-white px-3 text-sm font-semibold disabled:opacity-40"><ArrowLeft size={16} /> Previous</button>
            <span className="text-sm text-zinc-500">Page {page}</span>
            <button disabled={query.data.length < pageSize || query.isFetching} onClick={() => setPage((value) => value + 1)} className="inline-flex h-10 items-center gap-2 border border-line bg-white px-3 text-sm font-semibold disabled:opacity-40">Next <ArrowRight size={16} /></button>
          </nav>
        </>
      )}
    </>
  )
}

function PreviewState({ children, tone = 'neutral' }: { children: React.ReactNode; tone?: 'neutral' | 'error' }) {
  return <div role={tone === 'error' ? 'alert' : 'status'} className={`flex min-h-56 flex-col items-center justify-center gap-3 border border-line bg-white p-8 text-sm ${tone === 'error' ? 'text-red-700' : 'text-zinc-500'}`}>{children}</div>
}

