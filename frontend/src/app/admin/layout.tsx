"use client";

import { useEffect, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import Link from "next/link";
import { Contact, ExternalLink, FileText, FolderOpen, Home, LayoutDashboard, Loader2, LogOut, Mail, Menu, PanelBottom, ShieldCheck, Tags, X } from "lucide-react";

import { api } from "@/lib/api";
import { cn } from "@/lib/utils";

const navItems = [
  { title: "概览", href: "/admin/dashboard", icon: LayoutDashboard },
  { title: "首页", href: "/admin/home", icon: Home },
  { title: "文章", href: "/admin/posts", icon: FileText },
  { title: "标签", href: "/admin/tags", icon: Tags },
  { title: "项目", href: "/admin/projects", icon: FolderOpen },
  { title: "联系页", href: "/admin/contact", icon: Contact },
  { title: "页脚", href: "/admin/footer", icon: PanelBottom },
  { title: "来信", href: "/admin/messages", icon: Mail },
  { title: "安全", href: "/admin/security", icon: ShieldCheck },
];

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const [ready, setReady] = useState(pathname === "/admin/login");
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [loggingOut, setLoggingOut] = useState(false);

  useEffect(() => {
    if (pathname === "/admin/login") return;
    let active = true;
    api.me().then((response) => {
      if (!active) return;
      if (response.data.role !== "admin") throw new Error("forbidden");
      setReady(true);
    }).catch(() => router.replace("/admin/login"));
    return () => { active = false; };
  }, [pathname, router]);

  useEffect(() => {
    if (!sidebarOpen) return;
    function closeOnEscape(event: KeyboardEvent) {
      if (event.key === "Escape") setSidebarOpen(false);
    }
    document.addEventListener("keydown", closeOnEscape);
    return () => document.removeEventListener("keydown", closeOnEscape);
  }, [sidebarOpen]);

  async function logout() {
    if (loggingOut) return;
    setLoggingOut(true);
    try { await api.logout(); } finally { router.replace("/admin/login"); router.refresh(); }
  }

  if (pathname === "/admin/login") return <>{children}</>;
  if (!ready) return <div className="flex min-h-[60vh] items-center justify-center"><Loader2 className="h-6 w-6 animate-spin text-primary" /><span className="sr-only">正在验证登录状态</span></div>;

  return (
    <div className="min-h-[calc(100vh-4rem)] bg-muted/20 md:flex">
      {sidebarOpen && <button className="fixed inset-0 z-30 bg-black/40 md:hidden" onClick={() => setSidebarOpen(false)} aria-label="关闭侧栏" />}
      <div className="sticky top-16 z-20 flex h-14 items-center justify-between border-b border-border bg-background/95 px-4 backdrop-blur md:hidden">
        <button className="inline-flex min-h-10 items-center gap-2 rounded-xl border border-border px-3 text-sm font-medium" onClick={() => setSidebarOpen(true)} aria-label="打开后台菜单" aria-expanded={sidebarOpen} aria-controls="admin-sidebar"><Menu size={18} />管理菜单</button>
        <a href="/" target="_blank" rel="noopener noreferrer" className="inline-flex items-center gap-1.5 text-sm text-muted-foreground"><ExternalLink size={15} />查看站点</a>
      </div>
      <aside id="admin-sidebar" className={cn("fixed inset-y-16 left-0 z-40 flex w-64 flex-col border-r border-border bg-background/95 p-4 backdrop-blur transition-transform md:translate-x-0", sidebarOpen ? "translate-x-0" : "-translate-x-full")}>
        <p className="mb-4 px-3 text-xs font-semibold uppercase tracking-widest text-muted-foreground">Content Studio</p>
        <nav className="flex flex-1 flex-col gap-1 overflow-y-auto" aria-label="后台导航">
          {navItems.map(({ title, href, icon: Icon }) => <Link key={href} href={href} onClick={() => setSidebarOpen(false)} aria-current={pathname.startsWith(href) ? "page" : undefined} className={cn("flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition", pathname.startsWith(href) ? "bg-primary/10 text-primary" : "text-muted-foreground hover:bg-accent hover:text-foreground")}><Icon size={18} />{title}</Link>)}
        </nav>
        <a href="/" target="_blank" rel="noopener noreferrer" className="mb-1 hidden items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-muted-foreground transition hover:bg-accent hover:text-foreground md:flex"><ExternalLink size={18} />查看站点</a>
        <button onClick={logout} disabled={loggingOut} className="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-muted-foreground transition hover:bg-accent hover:text-foreground disabled:opacity-60">{loggingOut ? <Loader2 size={18} className="animate-spin" /> : <LogOut size={18} />}{loggingOut ? "正在退出" : "退出登录"}</button>
      </aside>
      {sidebarOpen && <button className="fixed right-4 top-20 z-50 rounded-full bg-background p-2 text-foreground shadow-lg md:hidden" onClick={() => setSidebarOpen(false)} aria-label="关闭后台菜单"><X size={20} /></button>}
      <div className="min-w-0 flex-1 md:ml-64"><div className="mx-auto max-w-6xl p-4 sm:p-5 md:p-8">{children}</div></div>
    </div>
  );
}
