import { apiRequest } from '../../lib/api'
import type { Article, ArticleInput, ArticlePatch, ArticleStatus } from './types'

export const articleKeys = {
  all: ['articles'] as const,
  list: (filters: { status: ArticleStatus; limit: number; offset: number }) => [...articleKeys.all, 'list', filters] as const,
  detail: (id: number) => [...articleKeys.all, 'detail', id] as const,
}

export function listArticles(status: ArticleStatus, limit: number, offset: number) {
  return apiRequest<Article[]>(`/article/${limit}/${offset}?status=${status}`)
}

export function getArticle(id: number) {
  return apiRequest<Article>(`/article/${id}`)
}

export function createArticle(input: ArticleInput) {
  return apiRequest<Article>('/article/', { method: 'POST', body: JSON.stringify(input) })
}

export function updateArticle(id: number, input: ArticlePatch) {
  return apiRequest<Article>(`/article/${id}`, { method: 'PUT', body: JSON.stringify(input) })
}

export function deleteArticle(id: number) {
  return apiRequest<void>(`/article/${id}`, { method: 'DELETE' })
}

