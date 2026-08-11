"use client";

import { useCallback, useEffect, useState } from "react";
import { Archive, Check, Loader2, Mail, Trash2 } from "lucide-react";

import { api } from "@/lib/api";
import { formatDate } from "@/lib/utils";
import type { ContactMessage } from "@/types";

export default function AdminMessagesPage() {
  const [messages, setMessages] = useState<ContactMessage[]>([]);
  const [status, setStatus] = useState("");
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const load = useCallback(() => { setLoading(true); api.getMessages({ page, page_size: 20, status }).then((response) => { setMessages(response.data.list); setTotalPages(response.data.total_pages); }).catch((cause) => setError(cause instanceof Error ? cause.message : "加载失败")).finally(() => setLoading(false)); }, [page, status]);
  useEffect(() => { const timer = window.setTimeout(load, 0); return () => window.clearTimeout(timer); }, [load]);
  async function update(message: ContactMessage, next: ContactMessage["status"]) { try { await api.updateMessage(message.id, next); setMessages((items) => items.map((item) => item.id === message.id ? { ...item, status: next } : item)); } catch (cause) { setError(cause instanceof Error ? cause.message : "更新失败"); } }
  async function remove(message: ContactMessage) { if (!window.confirm(`永久删除来自 ${message.name} 的来信？`)) return; try { await api.deleteMessage(message.id); setMessages((items) => items.filter((item) => item.id !== message.id)); } catch (cause) { setError(cause instanceof Error ? cause.message : "删除失败"); } }
  return (
    <div className="space-y-6"><div className="flex flex-wrap items-end justify-between gap-4"><div><p className="eyebrow">Inbox</p><h1 className="mt-2 text-3xl font-bold">读者来信</h1><p className="mt-2 text-muted-foreground">联系表单提交的消息会安全地保存在这里。</p></div><select className="input w-auto" aria-label="筛选消息状态" value={status} onChange={(e) => { setStatus(e.target.value); setPage(1); }}><option value="">全部</option><option value="unread">未读</option><option value="read">已读</option><option value="archived">已归档</option></select></div>
      {error && <p className="rounded-xl bg-destructive/10 p-4 text-sm text-destructive" role="alert">{error}</p>}
      {loading ? <div className="flex min-h-56 items-center justify-center"><Loader2 className="animate-spin text-primary" /></div> : messages.length ? <div className="space-y-4">{messages.map((message) => <article key={message.id} className={`surface-card p-5 ${message.status === "unread" ? "border-primary/40" : ""}`}><div className="flex flex-wrap items-start justify-between gap-3"><div><div className="flex items-center gap-2"><h2 className="font-semibold">{message.name}</h2><span className="badge">{{ unread: "未读", read: "已读", archived: "已归档" }[message.status]}</span></div><a className="mt-1 block text-sm text-primary hover:underline" href={`mailto:${message.email}`}>{message.email}</a></div><time className="text-xs text-muted-foreground" dateTime={message.created_at}>{formatDate(message.created_at)}</time></div><p className="mt-4 whitespace-pre-wrap text-sm leading-7 text-foreground/85">{message.message}</p><div className="mt-5 flex flex-wrap gap-2">{message.status === "unread" && <button onClick={() => update(message, "read")} className="button-secondary"><Check size={15} />标为已读</button>}{message.status !== "archived" && <button onClick={() => update(message, "archived")} className="button-secondary"><Archive size={15} />归档</button>}<button onClick={() => remove(message)} className="button-secondary text-destructive"><Trash2 size={15} />删除</button></div></article>)}</div> : <div className="empty-state"><Mail className="h-8 w-8" /><h2>收件箱是空的</h2><p>新的联系消息会显示在这里。</p></div>}
      {totalPages > 1 && <nav className="flex items-center justify-center gap-3" aria-label="来信分页"><button className="button-secondary" disabled={page <= 1} onClick={() => setPage((value) => value - 1)}>上一页</button><span className="text-sm text-muted-foreground">{page} / {totalPages}</span><button className="button-secondary" disabled={page >= totalPages} onClick={() => setPage((value) => value + 1)}>下一页</button></nav>}
    </div>
  );
}
