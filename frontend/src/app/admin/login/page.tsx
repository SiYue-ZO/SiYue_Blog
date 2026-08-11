"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Loader2, Lock, Mail } from "lucide-react";

import { api } from "@/lib/api";

export default function AdminLoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(event: React.FormEvent) {
    event.preventDefault();
    setError("");
    setLoading(true);
    try {
      const response = await api.login(email, password);
      if (response.data.user.role !== "admin") {
        await api.logout();
        setError("此账号没有管理员权限。");
        return;
      }
      router.replace("/admin/dashboard");
      router.refresh();
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "邮箱或密码不正确。");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-[calc(100vh-4rem)] items-center justify-center px-4 py-12">
      <div className="surface-card w-full max-w-md p-7 md:p-9">
        <div className="text-center"><p className="eyebrow">Private area</p><h1 className="mt-3 text-3xl font-bold">管理后台</h1><p className="mt-2 text-sm text-muted-foreground">使用管理员账号继续</p></div>
        <form onSubmit={handleSubmit} className="mt-8 space-y-5">
          {error && <p className="rounded-lg bg-destructive/10 p-3 text-sm text-destructive" role="alert">{error}</p>}
          <div className="space-y-2"><label htmlFor="email" className="label">邮箱</label><div className="relative"><Mail className="absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" /><input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} autoComplete="username" required className="input pl-10" /></div></div>
          <div className="space-y-2"><label htmlFor="password" className="label">密码</label><div className="relative"><Lock className="absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" /><input id="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} autoComplete="current-password" required className="input pl-10" /></div></div>
          <button type="submit" disabled={loading} className="button-primary w-full">{loading ? <><Loader2 className="h-4 w-4 animate-spin" />登录中…</> : "登录"}</button>
        </form>
      </div>
    </div>
  );
}
