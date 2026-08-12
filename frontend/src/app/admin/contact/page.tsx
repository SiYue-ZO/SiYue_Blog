"use client";

import { useEffect, useState } from "react";
import { ArrowDown, ArrowUp, Loader2, Plus, Save, Trash2 } from "lucide-react";

import { api } from "@/lib/api";
import { defaultSocialLinks } from "@/lib/constants";
import type { SocialLink, SocialPlatform } from "@/types";

const platformOptions: { value: SocialPlatform; label: string }[] = [
  { value: "github", label: "GitHub" },
  { value: "email", label: "邮箱" },
  { value: "bilibili", label: "哔哩哔哩" },
  { value: "twitter", label: "Twitter / X" },
  { value: "linkedin", label: "LinkedIn" },
  { value: "youtube", label: "YouTube" },
  { value: "weibo", label: "微博" },
  { value: "website", label: "个人网站" },
  { value: "other", label: "其他" },
];

export default function AdminContactPage() {
  const [links, setLinks] = useState<SocialLink[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");

  useEffect(() => {
    api.getAdminSettings()
      .then((response) => setLinks(response.data.social_links ?? defaultSocialLinks))
      .catch((cause) => setError(cause instanceof Error ? cause.message : "加载联系页配置失败"))
      .finally(() => setLoading(false));
  }, []);

  function update(index: number, values: Partial<SocialLink>) {
    setLinks((current) => current.map((item, itemIndex) => itemIndex === index ? { ...item, ...values } : item));
  }

  function move(index: number, direction: -1 | 1) {
    const target = index + direction;
    if (target < 0 || target >= links.length) return;
    setLinks((current) => {
      const next = [...current];
      [next[index], next[target]] = [next[target], next[index]];
      return next;
    });
  }

  async function save(event: React.FormEvent) {
    event.preventDefault();
    setSaving(true); setError(""); setNotice("");
    try {
      const response = await api.updateSettings({ social_links: links });
      setLinks(response.data.social_links ?? []);
      setNotice("联系页社交账号已更新。");
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "保存失败");
    } finally {
      setSaving(false);
    }
  }

  if (loading) return <div className="flex min-h-72 items-center justify-center"><Loader2 className="animate-spin text-primary" /><span className="sr-only">正在加载</span></div>;

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div><p className="eyebrow">Contact page</p><h1 className="mt-2 text-3xl font-bold">联系页</h1><p className="mt-2 text-sm text-muted-foreground">管理公开联系页展示的社交账号及排列顺序。</p></div>
        <a href="/contact" target="_blank" rel="noopener noreferrer" className="button-secondary">查看联系页</a>
      </div>
      {error && <p className="rounded-xl bg-destructive/10 p-4 text-sm text-destructive" role="alert">{error}</p>}
      {notice && <p className="rounded-xl bg-primary/10 p-4 text-sm text-primary" role="status">{notice}</p>}
      <form onSubmit={save} className="surface-card space-y-5 p-5 md:p-6">
        <div className="flex items-center justify-between gap-4"><div><h2 className="text-xl font-semibold">社交账号</h2><p className="mt-1 text-xs text-muted-foreground">最多添加 20 项；链接支持 HTTPS 和 mailto。</p></div><button type="button" className="button-secondary" disabled={links.length >= 20} onClick={() => setLinks((items) => [...items, { platform: "other", label: "", href: "https://" }])}><Plus size={16} />添加</button></div>
        {links.length ? <div className="space-y-3">{links.map((link, index) => (
          <fieldset key={index} className="rounded-xl border border-border p-4">
            <legend className="px-2 text-xs font-semibold text-muted-foreground">账号 {index + 1}</legend>
            <div className="grid gap-4 md:grid-cols-[160px_1fr_1.5fr_auto] md:items-end">
              <label className="space-y-2"><span className="label">平台</span><select className="input" value={link.platform} onChange={(event) => update(index, { platform: event.target.value as SocialPlatform })}>{platformOptions.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}</select></label>
              <label className="space-y-2"><span className="label">显示名称</span><input className="input" required maxLength={40} placeholder="例如：我的 GitHub" value={link.label} onChange={(event) => update(index, { label: event.target.value })} /></label>
              <label className="space-y-2"><span className="label">链接</span><input className="input" required maxLength={512} placeholder="https://… 或 mailto:…" value={link.href} onChange={(event) => update(index, { href: event.target.value })} /></label>
              <div className="flex gap-1"><button type="button" className="rounded-lg border p-2 disabled:opacity-30" disabled={index === 0} onClick={() => move(index, -1)} aria-label={`上移账号 ${index + 1}`}><ArrowUp size={16} /></button><button type="button" className="rounded-lg border p-2 disabled:opacity-30" disabled={index === links.length - 1} onClick={() => move(index, 1)} aria-label={`下移账号 ${index + 1}`}><ArrowDown size={16} /></button><button type="button" className="rounded-lg border p-2 text-destructive" onClick={() => setLinks((items) => items.filter((_, itemIndex) => itemIndex !== index))} aria-label={`删除账号 ${index + 1}`}><Trash2 size={16} /></button></div>
            </div>
          </fieldset>
        ))}</div> : <div className="rounded-xl border border-dashed p-8 text-center text-sm text-muted-foreground">暂不展示社交账号。点击“添加”创建一个入口。</div>}
        <div className="flex justify-end"><button className="button-primary" disabled={saving}>{saving ? <Loader2 size={16} className="animate-spin" /> : <Save size={16} />}{saving ? "保存中" : "保存设置"}</button></div>
      </form>
    </div>
  );
}
