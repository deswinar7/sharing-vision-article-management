const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080').replace(/\/$/, '')

export interface ApiErrorBody {
  error?: { code?: string; message?: string; fields?: Record<string, string> }
}

export class ApiError extends Error {
  constructor(public status: number, message: string, public fields?: Record<string, string>) {
    super(message)
  }
}

export async function apiRequest<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: { 'Content-Type': 'application/json', ...options?.headers },
  })
  if (!response.ok) {
    let body: ApiErrorBody = {}
    try { body = await response.json() as ApiErrorBody } catch { /* Empty or non-JSON error response. */ }
    throw new ApiError(response.status, body.error?.message || 'The request could not be completed', body.error?.fields)
  }
  if (response.status === 204) return undefined as T
  return response.json() as Promise<T>
}

