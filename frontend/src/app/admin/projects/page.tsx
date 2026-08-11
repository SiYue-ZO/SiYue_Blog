"use client";

import { useEffect, useState } from "react";
import { Edit3, Loader2, Plus, Save, Trash2, X } from "lucide-react";

import { api } from "@/lib/api";
import type { Project } from "@/types";

type ProjectForm = Omit<Project, "id" | "created_at" | "updated_at">;
const empty: ProjectForm = { title: "", description: "", repo_url: "", demo_url: "", image_url: "", tech_stack: [], is_featured: false, is_visible: true, sort_order: 0 };

export default function AdminProjectsPage() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [form, setForm] = useState<ProjectForm>(empty);
  const [editing, setEditing] = useState<number | null>(null);
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  useEffect(() => { api.getAdminProjects().then((response) => setProjects(response.data)).catch((cause) => setError(cause instanceof Error ? cause.message : "加载失败")).finally(() => setLoading(false)); }, []);
  function edit(project?: Project) { if (project) { const { id, created_at, updated_at, ...values } = project; void created_at; void updated_at; setEditing(id); setForm(values); } else { setEditing(null); setForm(empty); } setOpen(true); }
  function set<K extends keyof ProjectForm>(key: K, value: ProjectForm[K]) { setForm((current) => ({ ...current, [key]: value })); }
  async function save(event: React.FormEvent) { event.preventDefault(); setSaving(true); setError(""); try { const response = editing ? await api.updateProject(editing, form) : await api.createProject(form); setProjects((items) => editing ? items.map((item) => item.id === editing ? response.data : item) : [...items, response.data]); setOpen(false); } catch (cause) { setError(cause instanceof Error ? cause.message : "保存失败"); } finally { setSaving(false); } }
  async function remove(project: Project) { if (!window.confirm(`删除项目“${project.title}”？`)) return; try { await api.deleteProject(project.id); setProjects((items) => items.filter((item) => item.id !== project.id)); } catch (cause) { setError(cause instanceof Error ? cause.message : "删除失败"); } }
  return (
    <div className="space-y-6"><div className="flex flex-wrap items-end justify-between gap-4"><div><p className="eyebrow">Portfolio</p><h1 className="mt-2 text-3xl font-bold">项目管理</h1></div><button onClick={() => edit()} className="button-primary"><Plus size={17} />添加项目</button></div>
      {error && <p className="rounded-xl bg-destructive/10 p-4 text-sm text-destructive" role="alert">{error}</p>}
      {open && <form onSubmit={save} className="surface-card space-y-5 p-6"><div className="flex items-center justify-between"><h2 className="text-xl font-semibold">{editing ? "编辑项目" : "添加项目"}</h2><button type="button" onClick={() => setOpen(false)} className="rounded-lg p-2 hover:bg-accent" aria-label="关闭表单"><X size={18} /></button></div><div className="grid gap-4 sm:grid-cols-2"><div className="space-y-2"><label className="label">名称</label><input className="input" required maxLength={255} value={form.title} onChange={(e) => set("title", e.target.value)} /></div><div className="space-y-2"><label className="label">技术栈（逗号分隔）</label><input className="input" value={(form.tech_stack ?? []).join(", ")} onChange={(e) => set("tech_stack", e.target.value.split(/[,，]/).map((item) => item.trim()).filter(Boolean))} /></div></div><div className="space-y-2"><label className="label">简介</label><textarea className="input" rows={3} value={form.description} onChange={(e) => set("description", e.target.value)} /></div><div className="grid gap-4 sm:grid-cols-3"><div className="space-y-2"><label className="label">仓库地址</label><input type="url" className="input" value={form.repo_url} onChange={(e) => set("repo_url", e.target.value)} /></div><div className="space-y-2"><label className="label">演示地址</label><input type="url" className="input" value={form.demo_url} onChange={(e) => set("demo_url", e.target.value)} /></div><div className="space-y-2"><label className="label">图片地址</label><input type="url" className="input" value={form.image_url} onChange={(e) => set("image_url", e.target.value)} /></div></div><div className="flex flex-wrap items-center gap-5"><label className="flex items-center gap-2 text-sm"><input type="checkbox" className="accent-primary" checked={form.is_visible} onChange={(e) => set("is_visible", e.target.checked)} />公开展示</label><label className="flex items-center gap-2 text-sm"><input type="checkbox" className="accent-primary" checked={form.is_featured} onChange={(e) => set("is_featured", e.target.checked)} />首页精选</label><label className="flex items-center gap-2 text-sm">排序<input type="number" className="input w-24" value={form.sort_order} onChange={(e) => set("sort_order", Number(e.target.value))} /></label><button disabled={saving} className="button-primary ml-auto">{saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save size={16} />}保存</button></div><p className="text-xs text-muted-foreground">关闭“公开展示”后，项目仍保留在后台，但不会出现在项目页或首页精选区域。</p></form>}
      {loading ? <div className="flex min-h-56 items-center justify-center"><Loader2 className="animate-spin text-primary" /></div> : projects.length ? <div className="grid gap-4 md:grid-cols-2">{projects.map((project) => <article key={project.id} className="surface-card p-5"><div className="flex items-start justify-between gap-3"><div><div className="flex flex-wrap items-center gap-2"><h2 className="font-semibold">{project.title}</h2>{!project.is_visible && <span className="badge">已隐藏</span>}{project.is_featured && <span className="badge">精选</span>}</div><p className="mt-2 line-clamp-2 text-sm text-muted-foreground">{project.description}</p></div></div><div className="mt-4 flex flex-wrap gap-2">{(project.tech_stack ?? []).map((tech) => <span key={tech} className="tag">{tech}</span>)}</div><div className="mt-5 flex gap-2"><button onClick={() => edit(project)} className="button-secondary"><Edit3 size={15} />编辑</button><button onClick={() => remove(project)} className="button-secondary text-destructive"><Trash2 size={15} />删除</button></div></article>)}</div> : <div className="empty-state"><h2>暂无项目</h2><p>添加真实作品后将展示在前台。</p></div>}
    </div>
  );
}
