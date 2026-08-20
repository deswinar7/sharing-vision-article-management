export type ArticleStatus = 'publish' | 'draft' | 'thrash'

export interface Article {
  id: number
  title: string
  content: string
  category: string
  created_date: string
  updated_date: string
  status: ArticleStatus
}

export interface ArticleInput {
  title: string
  content: string
  category: string
  status: ArticleStatus
}

export type ArticlePatch = Partial<ArticleInput>

