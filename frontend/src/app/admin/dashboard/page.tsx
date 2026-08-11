"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Eye, FileText, FolderOpen, Loader2, Mail, Plus, Tags } from "lucide-react";

import { api } from "@/lib/api";
import type { DashboardStats } from "@/types";

export default function AdminDashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [error, setError] = useState("");
  useEffect(() => { api.getStats().then((response) => setStats(response.data)).catch((cause) => setError(cause instanceof Error ? cause.message : "加载失败")); }, []);
  const cards = stats ? [
    { title: "文章", value: stats.total_posts, detail: `${stats.published_posts} 已发布 · ${stats.draft_posts} 草稿`, icon: FileText },
    { title: "标签", value: stats.total_tags, detail: "内容分类", icon: Tags },
    { title: "项目", value: stats.total_projects, detail: "作品展示", icon: FolderOpen },
    { title: "阅读", value: stats.total_views, detail: "累计文章浏览", icon: Eye },
    { title: "未读来信", value: stats.unread_messages, detail: "等待处理", icon: Mail },
  ] : [];
  return (
    <div className="space-y-8">
      <div className="flex flex-wrap items-end justify-between gap-4"><div><p className="eyebrow">Dashboard</p><h1 className="mt-2 text-3xl font-bold">内容概览</h1><p className="mt-2 text-muted-foreground">管理文章、项目和读者来信。</p></div><Link href="/admin/posts/new" className="button-primary"><Plus size={17} />新建文章</Link></div>
      {error && <p className="rounded-xl bg-destructive/10 p-4 text-destructive" role="alert">{error}</p>}
      {!stats && !error ? <div className="flex min-h-40 items-center justify-center"><Loader2 className="animate-spin text-primary" /></div> : <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">{cards.map(({ title, value, detail, icon: Icon }) => <div key={title} className="surface-card p-5"><div className="flex items-center justify-between"><p className="text-sm text-muted-foreground">{title}</p><Icon className="h-5 w-5 text-primary" /></div><p className="mt-3 text-3xl font-bold">{value}</p><p className="mt-2 text-xs text-muted-foreground">{detail}</p></div>)}</div>}
      <section className="surface-card p-6"><h2 className="text-lg font-semibold">快捷入口</h2><div className="mt-4 grid gap-3 sm:grid-cols-2 lg:grid-cols-4">{[{ href: "/admin/posts", label: "管理文章" }, { href: "/admin/tags", label: "管理标签" }, { href: "/admin/projects", label: "管理项目" }, { href: "/admin/messages", label: "查看来信" }].map((item) => <Link key={item.href} href={item.href} className="button-secondary justify-start">{item.label}</Link>)}</div></section>
    </div>
  );
}
