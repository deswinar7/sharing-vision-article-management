import { Navigate, Route, Routes } from 'react-router'
import { AppLayout } from './AppLayout'
import { EditArticlePage } from '../features/articles/EditArticlePage'
import { NewArticlePage } from '../features/articles/NewArticlePage'
import { PostsPage } from '../features/articles/PostsPage'
import { PreviewPage } from '../features/articles/PreviewPage'

export function App() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route index element={<Navigate to="/posts" replace />} />
        <Route path="posts" element={<PostsPage />} />
        <Route path="posts/new" element={<NewArticlePage />} />
        <Route path="posts/:id/edit" element={<EditArticlePage />} />
        <Route path="preview" element={<PreviewPage />} />
        <Route path="*" element={<Navigate to="/posts" replace />} />
      </Route>
    </Routes>
  )
}

