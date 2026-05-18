import type { Metadata } from "next";
import Link from "next/link";

import { getAllTags } from "@/lib/posts";

export const metadata: Metadata = {
  title: "标签",
  description: "按标签浏览文章",
};

export default function TagsPage() {
  const tags = getAllTags();

  return (
    <div className="mx-auto w-full max-w-7xl px-4 py-12 sm:px-6 lg:px-8 md:py-16">
      <div className="space-y-8">
        <div className="space-y-4">
          <h1 className="text-4xl font-bold">标签</h1>
          <p className="text-lg text-muted-foreground">按标签浏览文章</p>
        </div>

        {tags.length > 0 ? (
          <div className="flex flex-wrap gap-3">
            {tags.map((tag) => (
              <Link
                key={tag.name}
                href={`/tags/${tag.name.toLowerCase()}`}
                className="inline-flex items-center gap-2 rounded-lg border border-border/50 bg-card px-4 py-3 text-sm font-medium transition-all hover:border-primary/50 hover:bg-accent"
              >
                <span>{tag.name}</span>
                <span className="rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary">
                  {tag.count}
                </span>
              </Link>
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <p className="text-muted-foreground">暂无标签</p>
          </div>
        )}
      </div>
    </div>
  );
}
