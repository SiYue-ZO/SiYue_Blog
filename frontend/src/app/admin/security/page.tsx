"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { KeyRound, Loader2 } from "lucide-react";

import { api } from "@/lib/api";

export default function SecurityPage() {
  const router = useRouter();
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmation, setConfirmation] = useState("");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  async function submit(event: React.FormEvent) {
    event.preventDefault();
    if (newPassword !== confirmation) {
      setError("两次输入的新密码不一致。");
      return;
    }
    setSaving(true);
    setError("");
    try {
      await api.changePassword(currentPassword, newPassword);
      router.replace("/admin/login?password=changed");
      router.refresh();
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "修改密码失败");
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="max-w-2xl space-y-6">
      <div><p className="eyebrow">Account security</p><h1 className="mt-2 text-3xl font-bold">账号安全</h1><p className="mt-2 text-muted-foreground">修改密码后会撤销所有刷新会话，并要求重新登录。</p></div>
      <form onSubmit={submit} className="surface-card space-y-5 p-6 md:p-8">
        <div className="flex items-center gap-3"><KeyRound className="text-primary" /><h2 className="text-xl font-semibold">修改密码</h2></div>
        {error && <p className="rounded-lg bg-destructive/10 p-3 text-sm text-destructive" role="alert">{error}</p>}
        <div className="space-y-2"><label htmlFor="current-password" className="label">当前密码</label><input id="current-password" type="password" autoComplete="current-password" required maxLength={128} className="input" value={currentPassword} onChange={(event) => setCurrentPassword(event.target.value)} /></div>
        <div className="space-y-2"><label htmlFor="new-password" className="label">新密码</label><input id="new-password" type="password" autoComplete="new-password" required minLength={12} maxLength={128} className="input" value={newPassword} onChange={(event) => setNewPassword(event.target.value)} /><p className="text-xs text-muted-foreground">至少 12 个字符，建议使用密码管理器生成。</p></div>
        <div className="space-y-2"><label htmlFor="confirm-password" className="label">确认新密码</label><input id="confirm-password" type="password" autoComplete="new-password" required minLength={12} maxLength={128} className="input" value={confirmation} onChange={(event) => setConfirmation(event.target.value)} /></div>
        <button disabled={saving} className="button-primary">{saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <KeyRound className="h-4 w-4" />}{saving ? "修改中…" : "修改并退出"}</button>
      </form>
    </div>
  );
}
