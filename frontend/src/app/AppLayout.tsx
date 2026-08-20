import { FileText, LayoutList, Plus, Radio } from 'lucide-react'
import { NavLink, Outlet } from 'react-router'

const navigation = [
  { to: '/posts', label: 'All Posts', icon: LayoutList },
  { to: '/posts/new', label: 'Add New', icon: Plus },
  { to: '/preview', label: 'Preview', icon: Radio },
]

export function AppLayout() {
  return (
    <div className="min-h-screen lg:grid lg:grid-cols-[240px_1fr]">
      <aside className="border-b border-line bg-panel lg:min-h-screen lg:border-r lg:border-b-0">
        <div className="flex h-16 items-center gap-3 border-b border-line px-5">
          <span className="grid size-9 place-items-center bg-brand text-white"><FileText size={19} /></span>
          <div>
            <p className="text-sm font-semibold">Article Desk</p>
            <p className="text-xs text-zinc-500">Editorial workspace</p>
          </div>
        </div>
        <nav aria-label="Main navigation" className="flex gap-1 overflow-x-auto p-3 lg:flex-col">
          {navigation.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/posts'}
              className={({ isActive }) => `flex min-h-10 shrink-0 items-center gap-3 px-3 text-sm font-medium transition-colors ${isActive ? 'bg-emerald-50 text-brand' : 'text-zinc-600 hover:bg-zinc-100 hover:text-ink'}`}
            >
              <Icon size={18} /> {label}
            </NavLink>
          ))}
        </nav>
      </aside>
      <main className="min-w-0 p-4 sm:p-6 lg:p-8">
        <div className="mx-auto max-w-6xl"><Outlet /></div>
      </main>
    </div>
  )
}

