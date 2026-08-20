import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { Navigate, useNavigate, useParams } from 'react-router'
import { ApiError } from '../../lib/api'
import { ArticleForm } from './ArticleForm'
import { articleKeys, getArticle, updateArticle } from './api'
import { PageHeader } from './PageHeader'

export function EditArticlePage() {
  const id = Number(useParams().id)
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [serverError, setServerError] = useState<string>()
  const query = useQuery({ queryKey: articleKeys.detail(id), queryFn: () => getArticle(id), enabled: Number.isInteger(id) && id > 0 })
  const mutation = useMutation({
    mutationFn: (input: Parameters<typeof updateArticle>[1]) => updateArticle(id, input),
    onSuccess: async (article) => {
      await queryClient.invalidateQueries({ queryKey: articleKeys.all })
      navigate('/posts', { state: { message: article.status === 'publish' ? 'Article updated and published.' : 'Draft updated.' } })
    },
    onError: (error) => setServerError(error instanceof ApiError ? error.message : 'Unable to update the article.'),
  })
  if (!Number.isInteger(id) || id < 1) return <Navigate to="/posts" replace />
  if (query.isPending) return <StateMessage>Loading article...</StateMessage>
  if (query.isError || !query.data) return <StateMessage tone="error">Failed to load this article.</StateMessage>
  return (
    <>
      <PageHeader title="Edit article" description="Update the article and choose its next publishing state." />
      <div className="max-w-3xl">
        <ArticleForm
          initialValues={{ title: query.data.title, content: query.data.content, category: query.data.category }}
          onSubmit={async (input) => { setServerError(undefined); await mutation.mutateAsync(input) }}
          isSubmitting={mutation.isPending}
          serverError={serverError}
        />
      </div>
    </>
  )
}

function StateMessage({ children, tone = 'neutral' }: { children: React.ReactNode; tone?: 'neutral' | 'error' }) {
  return <div role={tone === 'error' ? 'alert' : 'status'} className={`border p-5 text-sm ${tone === 'error' ? 'border-red-200 bg-red-50 text-red-800' : 'border-line bg-white text-zinc-600'}`}>{children}</div>
}

