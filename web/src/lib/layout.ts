export type ModuleKey = 'overview' | 'resources' | 'posts' | 'ai-lab'
export interface Layout { order: ModuleKey[]; hidden: ModuleKey[]; density: 'comfortable' | 'compact'; theme: 'light' | 'dark' }
export const moduleLabels: Record<ModuleKey, string> = { overview: '本周概览', resources: '资源库', posts: '开发日志', 'ai-lab': 'AI 实验室' }
export const defaultLayout: Layout = { order: ['overview', 'resources', 'posts', 'ai-lab'], hidden: [], density: 'comfortable', theme: 'light' }
export function moveModule(layout: Layout, key: ModuleKey, offset: number): Layout { const order = [...layout.order]; const from = order.indexOf(key); const to = from + offset; if (from < 0 || to < 0 || to >= order.length) return layout; [order[from], order[to]] = [order[to], order[from]]; return { ...layout, order } }
