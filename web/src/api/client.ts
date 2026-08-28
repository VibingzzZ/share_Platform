import type { Layout } from '../lib/layout'
export interface Resource { title: string; type: string; description: string }
export interface Post { title: string; summary: string; publishedAt?: string }
export interface AIProduct { name: string; summary: string; url: string; tags?: string[] }
export interface Dashboard { resources: Resource[]; posts: Post[]; aiProducts: AIProduct[] }
const headers = { 'X-User-ID': '00000000-0000-0000-0000-000000000001' }
const apiBase = (import.meta.env.VITE_API_BASE_URL || '').replace(/\/$/, '')
async function request<T>(path: string, init?: RequestInit): Promise<T> { const response = await fetch(`${apiBase}/api/v1${path}`, { ...init, headers: { ...headers, ...init?.headers } }); if (!response.ok) throw new Error(`Request failed: ${response.status}`); return response.json() as Promise<T> }
export const api = { dashboard: () => request<Dashboard>('/dashboard'), layout: () => request<Layout>('/layout'), saveLayout: (layout: Layout) => request<Layout>('/layout', { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(layout) }) }
