import type { Metadata } from "next";
import { Mail, MapPin, Clock } from "lucide-react";

export const metadata: Metadata = {
  title: "联系",
  description: "与我取得联系",
};

export default function ContactPage() {
  return (
    <div className="mx-auto w-full max-w-7xl px-4 py-12 sm:px-6 lg:px-8 md:py-16">
      <div className="mx-auto max-w-2xl space-y-8">
        <div className="space-y-4 text-center">
          <h1 className="text-4xl font-bold">联系我</h1>
          <p className="text-lg text-muted-foreground">
            有想法或想合作？欢迎随时联系我
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-3">
          <div className="flex flex-col items-center gap-3 rounded-lg border border-border/50 bg-card p-6 text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
              <Mail className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h3 className="font-semibold">邮箱</h3>
              <p className="text-sm text-muted-foreground">
                contact@siyue.dev
              </p>
            </div>
          </div>

          <div className="flex flex-col items-center gap-3 rounded-lg border border-border/50 bg-card p-6 text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
              <MapPin className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h3 className="font-semibold">位置</h3>
              <p className="text-sm text-muted-foreground">中国</p>
            </div>
          </div>

          <div className="flex flex-col items-center gap-3 rounded-lg border border-border/50 bg-card p-6 text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
              <Clock className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h3 className="font-semibold">响应时间</h3>
              <p className="text-sm text-muted-foreground">24 小时内</p>
            </div>
          </div>
        </div>

        <form className="space-y-6 rounded-lg border border-border/50 bg-card p-6 md:p-8">
          <div className="space-y-2">
            <label htmlFor="name" className="text-sm font-medium">
              姓名
            </label>
            <input
              id="name"
              type="text"
              placeholder="你的姓名"
              className="w-full rounded-md border border-border bg-background px-4 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50"
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="email" className="text-sm font-medium">
              邮箱
            </label>
            <input
              id="email"
              type="email"
              placeholder="your@email.com"
              className="w-full rounded-md border border-border bg-background px-4 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50"
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="message" className="text-sm font-medium">
              消息
            </label>
            <textarea
              id="message"
              rows={5}
              placeholder="写下你想说的..."
              className="w-full rounded-md border border-border bg-background px-4 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50"
            />
          </div>

          <button
            type="submit"
            className="w-full rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground shadow transition-colors hover:bg-primary/90"
          >
            发送消息
          </button>
        </form>
      </div>
    </div>
  );
}
