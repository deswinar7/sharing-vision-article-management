import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { useNavigate } from 'react-router'
import { ApiError } from '../../lib/api'
import { ArticleForm } from './ArticleForm'
import { articleKeys, createArticle } from './api'
import { PageHeader } from './PageHeader'

export function NewArticlePage() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [serverError, setServerError] = useState<string>()
  const mutation = useMutation({
    mutationFn: createArticle,
    onSuccess: async (article) => {
      await queryClient.invalidateQueries({ queryKey: articleKeys.all })
      navigate('/posts', { state: { message: article.status === 'publish' ? 'Article published.' : 'Draft saved.' } })
    },
    onError: (error) => setServerError(error instanceof ApiError ? error.message : 'Unable to save the article.'),
  })
  return (
    <>
      <PageHeader title="Add new article" description="Create a publish-ready article or keep it as a draft." />
      <div className="max-w-3xl"><ArticleForm onSubmit={async (input) => { setServerError(undefined); await mutation.mutateAsync(input) }} isSubmitting={mutation.isPending} serverError={serverError} /></div>
    </>
  )
}

