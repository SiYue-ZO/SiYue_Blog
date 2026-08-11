"use client";

import { useEffect, useState } from "react";
import { Loader2, Plus, Trash2 } from "lucide-react";

import { api } from "@/lib/api";
import { slugify } from "@/lib/utils";
import type { Tag } from "@/types";

export default function AdminTagsPage() {
  const [tags, setTags] = useState<Tag[]>([]);
  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");
  const [edited, setEdited] = useState(false);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  useEffect(() => { api.getAdminTags().then((response) => setTags(response.data)).catch((cause) => setError(cause instanceof Error ? cause.message : "加载失败")).finally(() => setLoading(false)); }, []);
  async function create(event: React.FormEvent) { event.preventDefault(); setSaving(true); setError(""); try { const response = await api.createTag({ name, slug }); setTags((items) => [...items, response.data]); setName(""); setSlug(""); setEdited(false); } catch (cause) { setError(cause instanceof Error ? cause.message : "创建失败"); } finally { setSaving(false); } }
  async function remove(tag: Tag) { if (!window.confirm(`删除标签“${tag.name}”？文章不会被删除。`)) return; try { await api.deleteTag(tag.id); setTags((items) => items.filter((item) => item.id !== tag.id)); } catch (cause) { setError(cause instanceof Error ? cause.message : "删除失败"); } }
  return (
    <div className="space-y-6"><div><p className="eyebrow">Taxonomy</p><h1 className="mt-2 text-3xl font-bold">标签管理</h1><p className="mt-2 text-muted-foreground">使用稳定、清晰的标签组织内容。</p></div>
      {error && <p className="rounded-xl bg-destructive/10 p-4 text-sm text-destructive" role="alert">{error}</p>}
      <form onSubmit={create} className="surface-card grid gap-4 p-5 sm:grid-cols-[1fr_1fr_auto] sm:items-end"><div className="space-y-2"><label htmlFor="tag-name" className="label">名称</label><input id="tag-name" className="input" required maxLength={50} value={name} onChange={(e) => { setName(e.target.value); if (!edited) setSlug(slugify(e.target.value)); }} placeholder="例如：Next.js" /></div><div className="space-y-2"><label htmlFor="tag-slug" className="label">别名</label><input id="tag-slug" className="input" required maxLength={50} value={slug} onChange={(e) => { setEdited(true); setSlug(slugify(e.target.value)); }} placeholder="nextjs" /></div><button disabled={saving} className="button-primary">{saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus size={17} />}添加</button></form>
      <div className="surface-card overflow-hidden">{loading ? <div className="flex min-h-48 items-center justify-center"><Loader2 className="animate-spin text-primary" /></div> : tags.length ? <div className="divide-y divide-border">{tags.map((tag) => <div key={tag.id} className="flex items-center gap-4 p-5"><div className="min-w-0 flex-1"><p className="font-medium">{tag.name}</p><p className="text-sm text-muted-foreground">{tag.slug} · {tag.post_count ?? 0} 篇已发布文章</p></div><button onClick={() => remove(tag)} className="button-secondary text-destructive"><Trash2 size={15} />删除</button></div>)}</div> : <div className="empty-state m-5"><h2>暂无标签</h2><p>在上方创建第一个标签。</p></div>}</div>
    </div>
  );
}
