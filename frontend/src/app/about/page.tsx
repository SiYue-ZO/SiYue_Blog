import Image from "next/image";
import Link from "next/link";
import { Code2, Layers3, ServerCog, Sparkles } from "lucide-react";

import { siteConfig } from "@/lib/constants";
import { pageMetadata } from "@/lib/metadata";

export const metadata = pageMetadata("关于", `了解 ${siteConfig.author.name} 与这个博客`, "/about");

const focus = [
  { icon: Code2, title: "清晰的实现", text: "偏爱易读、可测试、能长期维护的代码，而不只是短期可用。" },
  { icon: ServerCog, title: "可靠的系统", text: "从 API、数据模型到部署与监控，关注完整的软件生命周期。" },
  { icon: Layers3, title: "端到端体验", text: "把工程质量、性能、可访问性与视觉体验视为同一个产品问题。" },
  { icon: Sparkles, title: "持续学习", text: "通过写作整理知识，也公开记录真实的实践、取舍与复盘。" },
];

export default function AboutPage() {
  return (
    <div className="page-shell py-12 md:py-20">
      <div className="mx-auto max-w-4xl space-y-16">
        <section className="grid items-center gap-10 md:grid-cols-[13rem_1fr]">
          <Image src={siteConfig.author.avatar} alt={`${siteConfig.author.name} 的头像`} width={208} height={208} className="aspect-square rounded-3xl border border-border object-cover shadow-xl" priority />
          <div><p className="eyebrow">About me</p><h1 className="mt-3 text-4xl font-bold">你好，我是 {siteConfig.author.name}</h1><p className="mt-5 text-lg leading-8 text-muted-foreground">{siteConfig.author.bio}</p><p className="mt-4 leading-8 text-muted-foreground">这个博客是我的长期数字花园：分享做过的项目、解决过的问题，以及那些值得被写下来的思考。内容以真实实践为主，并会持续修订。</p><div className="mt-7 flex flex-wrap gap-3"><Link href="/blog" className="button-primary">阅读文章</Link><Link href="/contact" className="button-secondary">联系我</Link></div></div>
        </section>
        <section><div className="max-w-2xl"><p className="eyebrow">Principles</p><h2 className="mt-3 text-3xl font-bold">我关心的事情</h2></div><div className="mt-8 grid gap-5 sm:grid-cols-2">{focus.map(({ icon: Icon, title, text }) => <article key={title} className="surface-card p-6"><Icon className="h-6 w-6 text-primary" /><h3 className="mt-4 text-lg font-semibold">{title}</h3><p className="mt-2 text-sm leading-6 text-muted-foreground">{text}</p></article>)}</div></section>
        <section className="surface-card bg-primary/5 p-7 md:p-9"><p className="eyebrow">Now</p><h2 className="mt-3 text-2xl font-bold">这个站点也在持续进化</h2><p className="mt-3 leading-7 text-muted-foreground">博客本身是一个开放的全栈项目，包含内容管理、搜索、SEO、安全与容器化部署。你可以在 GitHub 查看源码，也欢迎交流改进建议。</p><a href={siteConfig.author.github} target="_blank" rel="noopener noreferrer" className="button-secondary mt-6">访问 GitHub</a></section>
      </div>
    </div>
  );
}
