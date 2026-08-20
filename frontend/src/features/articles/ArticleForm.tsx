import { zodResolver } from '@hookform/resolvers/zod'
import { ArrowLeft, FileCheck, Save } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { Link } from 'react-router'
import { z } from 'zod'
import type { ArticleInput, ArticleStatus } from './types'

const articleFormSchema = z.object({
  title: z.string().trim().min(20, 'Title must contain at least 20 characters').max(200, 'Title must contain at most 200 characters'),
  content: z.string().trim().min(200, 'Content must contain at least 200 characters'),
  category: z.string().trim().min(3, 'Category must contain at least 3 characters').max(100, 'Category must contain at most 100 characters'),
})

type FormValues = z.infer<typeof articleFormSchema>

interface Props {
  initialValues?: FormValues
  onSubmit: (input: ArticleInput) => Promise<void>
  isSubmitting: boolean
  serverError?: string
}

export function ArticleForm({ initialValues, onSubmit, isSubmitting, serverError }: Props) {
  const { register, handleSubmit, formState: { errors } } = useForm<FormValues>({
    resolver: zodResolver(articleFormSchema),
    defaultValues: initialValues ?? { title: '', content: '', category: '' },
  })

  const submitWithStatus = (status: ArticleStatus) => handleSubmit((values) => onSubmit({ ...values, status }))

  return (
    <form className="space-y-6" noValidate>
      {serverError && <div role="alert" className="border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">{serverError}</div>}
      <div>
        <label htmlFor="title" className="mb-2 block text-sm font-semibold">Title</label>
        <input id="title" className="h-11 w-full border border-line bg-white px-3 outline-none focus:border-brand" aria-invalid={!!errors.title} {...register('title')} />
        {errors.title && <p role="alert" className="mt-1.5 text-sm text-red-700">{errors.title.message}</p>}
      </div>
      <div>
        <label htmlFor="category" className="mb-2 block text-sm font-semibold">Category</label>
        <input id="category" className="h-11 w-full border border-line bg-white px-3 outline-none focus:border-brand" aria-invalid={!!errors.category} {...register('category')} />
        {errors.category && <p role="alert" className="mt-1.5 text-sm text-red-700">{errors.category.message}</p>}
      </div>
      <div>
        <div className="mb-2 flex items-center justify-between gap-3">
          <label htmlFor="content" className="text-sm font-semibold">Content</label>
          <span className="text-xs text-zinc-500">Minimum 200 characters</span>
        </div>
        <textarea id="content" rows={12} className="w-full resize-y border border-line bg-white p-3 leading-7 outline-none focus:border-brand" aria-invalid={!!errors.content} {...register('content')} />
        {errors.content && <p role="alert" className="mt-1.5 text-sm text-red-700">{errors.content.message}</p>}
      </div>
      <div className="flex flex-wrap items-center gap-3 border-t border-line pt-5">
        <button type="button" disabled={isSubmitting} onClick={submitWithStatus('publish')} className="inline-flex h-10 items-center gap-2 bg-brand px-4 text-sm font-semibold text-white hover:bg-brand-dark disabled:opacity-60">
          <FileCheck size={17} /> Publish
        </button>
        <button type="button" disabled={isSubmitting} onClick={submitWithStatus('draft')} className="inline-flex h-10 items-center gap-2 border border-line bg-white px-4 text-sm font-semibold hover:bg-zinc-50 disabled:opacity-60">
          <Save size={17} /> Save as Draft
        </button>
        <Link to="/posts" className="ml-auto inline-flex h-10 items-center gap-2 px-2 text-sm font-medium text-zinc-600 hover:text-ink"><ArrowLeft size={17} /> Cancel</Link>
      </div>
    </form>
  )
}
