import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Edit3, Plus, RotateCcw, Trash2 } from 'lucide-react'
import { useState } from 'react'
import { Link, useLocation } from 'react-router'
import { articleKeys, deleteArticle, listArticles, updateArticle } from './api'
import { PageHeader } from './PageHeader'
import type { Article, ArticleStatus } from './types'

const tabs: { label: string; status: ArticleStatus }[] = [
  { label: 'Published', status: 'publish' },
  { label: 'Drafts', status: 'draft' },
  { label: 'Trashed', status: 'thrash' },
]

export function PostsPage() {
  const [status, setStatus] = useState<ArticleStatus>('publish')
  const location = useLocation()
  const queryClient = useQueryClient()
  const query = useQuery({ queryKey: articleKeys.list({ status, limit: 50, offset: 0 }), queryFn: () => listArticles(status, 50, 0) })
  const mutation = useMutation({
    mutationFn: async ({ article, permanent }: { article: Article; permanent: boolean }) => {
      if (permanent) await deleteArticle(article.id)
      else await updateArticle(article.id, { status: 'thrash' })
    },
    onSuccess: async () => queryClient.invalidateQueries({ queryKey: articleKeys.all }),
  })
  const handleTrash = (article: Article) => {
    const permanent = article.status === 'thrash'
    const message = permanent ? `Permanently delete “${article.title}”? This cannot be undone.` : `Move “${article.title}” to Trashed?`
    if (window.confirm(message)) mutation.mutate({ article, permanent })
  }
  return (
    <>
      <PageHeader title="All posts" description="Review articles across each publishing state." action={<Link to="/posts/new" className="inline-flex h-10 items-center gap-2 bg-brand px-4 text-sm font-semibold text-white hover:bg-brand-dark"><Plus size={17} /> Add New</Link>} />
      {location.state?.message && <div role="status" className="mb-5 border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-900">{location.state.message}</div>}
      {mutation.isError && <div role="alert" className="mb-5 border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">The article action failed. Please try again.</div>}
      <div role="tablist" aria-label="Article status" className="flex border-b border-line">
        {tabs.map((tab) => <button key={tab.status} role="tab" aria-selected={status === tab.status} onClick={() => setStatus(tab.status)} className={`min-h-11 border-b-2 px-4 text-sm font-semibold ${status === tab.status ? 'border-brand text-brand' : 'border-transparent text-zinc-500 hover:text-ink'}`}>{tab.label}</button>)}
      </div>
      <div className="mt-5 overflow-hidden border border-line bg-panel">
        {query.isPending ? <TableState>Loading articles...</TableState> : query.isError ? <TableState tone="error"><span>Failed to load articles.</span><button onClick={() => query.refetch()} className="inline-flex items-center gap-2 font-semibold"><RotateCcw size={15} /> Retry</button></TableState> : query.data.length === 0 ? <TableState>No {status === 'thrash' ? 'trashed' : status === 'publish' ? 'published' : 'draft'} articles found.</TableState> : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[620px] border-collapse text-left">
              <thead className="bg-zinc-50 text-xs font-semibold uppercase text-zinc-500"><tr><th className="px-5 py-3">Title</th><th className="px-5 py-3">Category</th><th className="w-36 px-5 py-3">Action</th></tr></thead>
              <tbody>{query.data.map((article) => <tr key={article.id} className="border-t border-line"><td className="px-5 py-4"><p className="max-w-xl font-medium">{article.title}</p><p className="mt-1 text-xs text-zinc-500">Updated {new Date(article.updated_date).toLocaleDateString()}</p></td><td className="px-5 py-4 text-sm text-zinc-600">{article.category}</td><td className="px-5 py-4"><div className="flex items-center gap-1"><Link to={`/posts/${article.id}/edit`} aria-label={`Edit ${article.title}`} title="Edit" className="grid size-9 place-items-center text-zinc-500 hover:bg-emerald-50 hover:text-brand"><Edit3 size={17} /></Link><button aria-label={`${article.status === 'thrash' ? 'Delete' : 'Trash'} ${article.title}`} title={article.status === 'thrash' ? 'Delete permanently' : 'Move to trash'} disabled={mutation.isPending} onClick={() => handleTrash(article)} className="grid size-9 place-items-center text-zinc-500 hover:bg-red-50 hover:text-red-700 disabled:opacity-50"><Trash2 size={17} /></button></div></td></tr>)}</tbody>
            </table>
          </div>
        )}
      </div>
    </>
  )
}

function TableState({ children, tone = 'neutral' }: { children: React.ReactNode; tone?: 'neutral' | 'error' }) {
  return <div role={tone === 'error' ? 'alert' : 'status'} className={`flex min-h-40 flex-col items-center justify-center gap-3 p-8 text-sm ${tone === 'error' ? 'text-red-700' : 'text-zinc-500'}`}>{children}</div>
}
