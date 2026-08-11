import Link from "next/link";
import { Tags } from "lucide-react";

import { getTags } from "@/lib/content";
import { pageMetadata } from "@/lib/metadata";

export const metadata = pageMetadata("标签", "按标签浏览文章", "/tags");
export const dynamic = "force-dynamic";

export default async function TagsPage() {
  const tags = await getTags();
  return (
    <div className="page-shell py-12 md:py-16">
      <div className="space-y-8">
        <div className="space-y-3"><p className="eyebrow">Topics</p><h1 className="text-4xl font-bold">标签</h1><p className="text-lg text-muted-foreground">按主题探索所有文章</p></div>
        {tags.length ? (
          <div className="flex flex-wrap gap-3">{tags.map((tag) => <Link key={tag.id} href={`/blog?tag=${encodeURIComponent(tag.slug)}`} className="surface-card inline-flex items-center gap-2 px-4 py-3 text-sm font-medium hover:border-primary/50"><span>{tag.name}</span><span className="badge">{tag.post_count ?? 0}</span></Link>)}</div>
        ) : <div className="empty-state"><Tags className="h-8 w-8" /><h2>暂无标签</h2><p>发布文章后，标签会显示在这里。</p></div>}
      </div>
    </div>
  );
}
