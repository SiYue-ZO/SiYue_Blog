"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Loader2, Save } from "lucide-react";

import { api } from "@/lib/api";
import { slugify } from "@/lib/utils";
import type { Post, Tag } from "@/types";

type FormState = {
  title: string; slug: string; excerpt: string; content: string; cover_image: string;
  status: Post["status"]; is_pinned: boolean; tag_ids: number[];
};

const emptyForm: FormState = { title: "", slug: "", excerpt: "", content: "", cover_image: "", status: "draft", is_pinned: false, tag_ids: [] };

export function PostEditor({ postId }: { postId?: number }) {
  const router = useRouter();
  const [form, setForm] = useState<FormState>(emptyForm);
  const [tags, setTags] = useState<Tag[]>([]);
  const [loading, setLoading] = useState(Boolean(postId));
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [slugEdited, setSlugEdited] = useState(Boolean(postId));

  useEffect(() => {
    Promise.all([api.getAdminTags(), postId ? api.getAdminPost(postId) : Promise.resolve(null)])
      .then(([tagResponse, postResponse]) => {
        setTags(tagResponse.data);
        if (postResponse) {
          const post = postResponse.data;
          setForm({
            title: post.title, slug: post.slug, excerpt: post.excerpt || "", content: post.content || "",
            cover_image: post.cover_image || "", status: post.status, is_pinned: post.is_pinned,
            tag_ids: post.tags?.map((tag) => tag.id) || [],
          });
        }
      })
      .catch((cause) => setError(cause instanceof Error ? cause.message : "加载失败"))
      .finally(() => setLoading(false));
  }, [postId]);

  function set<K extends keyof FormState>(key: K, value: FormState[K]) { setForm((current) => ({ ...current, [key]: value })); }
  function updateTitle(title: string) { setForm((current) => ({ ...current, title, slug: slugEdited ? current.slug : slugify(title) })); }
  function toggleTag(id: number) { set("tag_ids", form.tag_ids.includes(id) ? form.tag_ids.filter((item) => item !== id) : [...form.tag_ids, id]); }

  async function submit(event: React.FormEvent) {
    event.preventDefault();
    setSaving(true); setError("");
    try {
      if (postId) await api.updatePost(postId, form);
      else await api.createPost(form);
      router.push("/admin/posts");
      router.refresh();
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "保存失败");
    } finally { setSaving(false); }
  }

  if (loading) return <div className="flex min-h-80 items-center justify-center"><Loader2 className="animate-spin text-primary" /></div>;

  return (
    <form onSubmit={submit} className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4"><div><Link href="/admin/posts" className="text-link"><ArrowLeft size={16} />返回文章</Link><h1 className="mt-3 text-3xl font-bold">{postId ? "编辑文章" : "新建文章"}</h1></div><button className="button-primary" disabled={saving}>{saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}{saving ? "保存中…" : "保存文章"}</button></div>
      {error && <p className="rounded-xl bg-destructive/10 p-4 text-sm text-destructive" role="alert">{error}</p>}
      <div className="grid gap-6 xl:grid-cols-[minmax(0,1fr)_19rem]">
        <div className="space-y-5">
          <section className="surface-card space-y-5 p-6">
            <div className="space-y-2"><label className="label" htmlFor="title">标题</label><input id="title" className="input text-lg font-semibold" required maxLength={255} value={form.title} onChange={(e) => updateTitle(e.target.value)} placeholder="一篇好文章从清晰的标题开始" /></div>
            <div className="space-y-2"><label className="label" htmlFor="slug">链接别名</label><div className="flex items-center rounded-xl border border-input bg-background/80 pl-4 focus-within:border-primary focus-within:ring-3 focus-within:ring-primary/15"><span className="text-sm text-muted-foreground">/blog/</span><input id="slug" className="min-w-0 flex-1 bg-transparent px-1 py-3 text-sm outline-none" required maxLength={255} value={form.slug} onChange={(e) => { setSlugEdited(true); set("slug", slugify(e.target.value)); }} /></div></div>
            <div className="space-y-2"><label className="label" htmlFor="excerpt">摘要</label><textarea id="excerpt" className="input" rows={3} maxLength={500} required value={form.excerpt} onChange={(e) => set("excerpt", e.target.value)} placeholder="用于列表和搜索引擎的简短介绍" /><p className="text-right text-xs text-muted-foreground">{form.excerpt.length}/500</p></div>
          </section>
          <section className="surface-card p-6"><div className="flex items-center justify-between gap-3"><label className="label" htmlFor="content">正文（Markdown）</label><span className="text-xs text-muted-foreground">{form.content.length.toLocaleString()} 字符</span></div><textarea id="content" className="input mt-3 min-h-[34rem] resize-y font-mono text-sm leading-6" required={form.status === "published"} value={form.content} onChange={(e) => set("content", e.target.value)} placeholder={"## 小标题\n\n开始写作…"} /><p className="mt-2 text-xs text-muted-foreground">草稿可暂时留空，发布前必须填写正文。</p></section>
        </div>
        <aside className="space-y-5">
          <section className="surface-card space-y-4 p-5"><h2 className="font-semibold">发布设置</h2><div className="space-y-2"><label htmlFor="status" className="label">状态</label><select id="status" className="input" value={form.status} onChange={(e) => set("status", e.target.value as Post["status"])}><option value="draft">草稿</option><option value="published">发布</option><option value="archived">归档</option></select></div><label className="flex cursor-pointer items-center gap-3 text-sm"><input type="checkbox" checked={form.is_pinned} onChange={(e) => set("is_pinned", e.target.checked)} className="h-4 w-4 accent-primary" />设为精选文章</label></section>
          <section className="surface-card space-y-3 p-5"><h2 className="font-semibold">封面</h2><input type="url" className="input" value={form.cover_image} onChange={(e) => set("cover_image", e.target.value)} placeholder="https://…" /><p className="text-xs leading-5 text-muted-foreground">可填写 HTTPS 图片地址；留空则使用站点默认封面。</p></section>
          <section className="surface-card p-5"><h2 className="font-semibold">标签</h2><div className="mt-3 space-y-2">{tags.length ? tags.map((tag) => <label key={tag.id} className="flex cursor-pointer items-center gap-3 rounded-lg px-2 py-1.5 text-sm hover:bg-accent"><input type="checkbox" checked={form.tag_ids.includes(tag.id)} onChange={() => toggleTag(tag.id)} className="h-4 w-4 accent-primary" />{tag.name}</label>) : <p className="text-sm text-muted-foreground">请先创建标签。</p>}</div></section>
        </aside>
      </div>
    </form>
  );
}
