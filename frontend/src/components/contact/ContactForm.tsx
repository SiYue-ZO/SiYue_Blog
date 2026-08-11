"use client";

import { useState } from "react";
import { CheckCircle2, Loader2, Send } from "lucide-react";

import { api } from "@/lib/api";

export function ContactForm() {
  const [state, setState] = useState<"idle" | "submitting" | "success" | "error">("idle");
  const [error, setError] = useState("");

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formElement = event.currentTarget;
    setState("submitting");
    setError("");
    const form = new FormData(formElement);
    try {
      await api.submitContact({
        name: String(form.get("name") || ""),
        email: String(form.get("email") || ""),
        message: String(form.get("message") || ""),
        website: String(form.get("website") || ""),
      });
      formElement.reset();
      setState("success");
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "发送失败，请稍后再试。");
      setState("error");
    }
  }

  if (state === "success") {
    return (
      <div className="surface-card flex min-h-80 flex-col items-center justify-center p-8 text-center" role="status">
        <CheckCircle2 className="h-12 w-12 text-emerald-500" />
        <h2 className="mt-5 text-2xl font-semibold">消息已收到</h2>
        <p className="mt-2 text-muted-foreground">感谢来信，我会尽快回复。</p>
        <button type="button" onClick={() => setState("idle")} className="button-secondary mt-6">再发一条</button>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="surface-card space-y-5 p-6 md:p-8">
      <div className="grid gap-5 sm:grid-cols-2">
        <div className="space-y-2"><label htmlFor="name" className="label">姓名</label><input id="name" name="name" required maxLength={100} autoComplete="name" className="input" placeholder="如何称呼你" /></div>
        <div className="space-y-2"><label htmlFor="email" className="label">邮箱</label><input id="email" name="email" required type="email" maxLength={255} autoComplete="email" className="input" placeholder="you@example.com" /></div>
      </div>
      <div className="hidden" aria-hidden="true"><label htmlFor="website">Website</label><input id="website" name="website" tabIndex={-1} autoComplete="off" /></div>
      <div className="space-y-2"><label htmlFor="message" className="label">消息</label><textarea id="message" name="message" required minLength={10} maxLength={5000} rows={7} className="input resize-y" placeholder="写下你想聊的内容…" /></div>
      {error && <p className="rounded-lg bg-destructive/10 p-3 text-sm text-destructive" role="alert">{error}</p>}
      <button type="submit" disabled={state === "submitting"} className="button-primary w-full">
        {state === "submitting" ? <><Loader2 className="h-4 w-4 animate-spin" />发送中…</> : <><Send className="h-4 w-4" />发送消息</>}
      </button>
    </form>
  );
}
