"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { Edit3, Loader2, Plus, Search, Trash2 } from "lucide-react";

import { api } from "@/lib/api";
import { formatDate } from "@/lib/utils";
import type { Post } from "@/types";

export default function AdminPostsPage() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [status, setStatus] = useState("");
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const load = useCallback(() => { setLoading(true); api.getAdminPosts({ page, page_size: 20, status, search }).then((response) => { setPosts(response.data.list); setTotalPages(response.data.total_pages); }).catch((cause) => setError(cause instanceof Error ? cause.message : "加载失败")).finally(() => setLoading(false)); }, [page, search, status]);
  useEffect(() => { const timer = window.setTimeout(load, 250); return () => window.clearTimeout(timer); }, [load]);
  async function remove(post: Post) { if (!window.confirm(`确定删除《${post.title}》吗？此操作不可撤销。`)) return; try { await api.deletePost(post.id); setPosts((items) => items.filter((item) => item.id !== post.id)); } catch (cause) { setError(cause instanceof Error ? cause.message : "删除失败"); } }
  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-4"><div><p className="eyebrow">Posts</p><h1 className="mt-2 text-3xl font-bold">文章管理</h1></div><Link href="/admin/posts/new" className="button-primary"><Plus size={17} />新建文章</Link></div>
      <div className="surface-card flex flex-wrap gap-3 p-4"><div className="relative min-w-56 flex-1"><Search className="absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" /><input aria-label="搜索文章" className="input pl-10" value={search} onChange={(e) => { setSearch(e.target.value); setPage(1); }} placeholder="搜索标题、别名或摘要" /></div><select aria-label="筛选状态" className="input w-auto min-w-32" value={status} onChange={(e) => { setStatus(e.target.value); setPage(1); }}><option value="">全部状态</option><option value="published">已发布</option><option value="draft">草稿</option><option value="archived">已归档</option></select></div>
      {error && <p className="rounded-xl bg-destructive/10 p-4 text-sm text-destructive" role="alert">{error}</p>}
      <div className="surface-card overflow-hidden">{loading ? <div className="flex min-h-56 items-center justify-center"><Loader2 className="animate-spin text-primary" /></div> : posts.length ? <div className="divide-y divide-border">{posts.map((post) => <article key={post.id} className="flex flex-col gap-4 p-5 sm:flex-row sm:items-center"><div className="min-w-0 flex-1"><div className="flex flex-wrap items-center gap-2"><h2 className="truncate font-semibold">{post.title}</h2><span className="badge">{{ published: "已发布", draft: "草稿", archived: "已归档" }[post.status]}</span>{post.is_pinned && <span className="badge">精选</span>}</div><p className="mt-1 truncate text-sm text-muted-foreground">/blog/{post.slug} · 更新于 {formatDate(post.updated_at)}</p></div><div className="flex gap-2"><Link href={`/admin/posts/${post.id}`} className="button-secondary"><Edit3 size={15} />编辑</Link><button onClick={() => remove(post)} className="button-secondary text-destructive"><Trash2 size={15} /><span className="sr-only sm:not-sr-only">删除</span></button></div></article>)}</div> : <div className="empty-state m-5"><h2>还没有文章</h2><p>新建一篇草稿开始写作。</p></div>}</div>
      {totalPages > 1 && <nav className="flex items-center justify-center gap-3" aria-label="文章管理分页"><button className="button-secondary" disabled={page <= 1} onClick={() => setPage((value) => value - 1)}>上一页</button><span className="text-sm text-muted-foreground">{page} / {totalPages}</span><button className="button-secondary" disabled={page >= totalPages} onClick={() => setPage((value) => value + 1)}>下一页</button></nav>}
    </div>
  );
}
