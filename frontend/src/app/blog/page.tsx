import type { Metadata } from "next";
import Link from "next/link";
import { Calendar, Clock, Search } from "lucide-react";

import { getAllPosts, getAllTags } from "@/lib/posts";
import { formatDate } from "@/lib/utils";

export const metadata: Metadata = {
  title: "博客",
  description: "阅读我的技术文章和思考",
};

export default function BlogPage() {
  const posts = getAllPosts();
  const tags = getAllTags();

  return (
    <div className="mx-auto w-full max-w-7xl px-4 py-12 sm:px-6 lg:px-8 md:py-16">
      <div className="space-y-8">
        <div className="space-y-4">
          <h1 className="text-4xl font-bold">博客</h1>
          <p className="text-lg text-muted-foreground">
            分享我的技术思考和开发经验
          </p>
        </div>

        {/* Search */}
        <div className="relative max-w-md">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <input
            type="text"
            placeholder="搜索文章..."
            className="w-full rounded-md border border-border bg-background py-2 pl-10 pr-4 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50"
          />
        </div>

        {/* Tags */}
        <div className="flex flex-wrap gap-2">
          <Link
            href="/blog"
            className="inline-flex items-center rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground"
          >
            全部
          </Link>
          {tags.map((tag) => (
            <Link
              key={tag.name}
              href={`/tags/${tag.name.toLowerCase()}`}
              className="inline-flex items-center gap-1.5 rounded-md border border-border bg-background px-3 py-1.5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"
            >
              {tag.name}
              <span className="text-xs text-muted-foreground">
                {tag.count}
              </span>
            </Link>
          ))}
        </div>

        {/* Posts */}
        {posts.length > 0 ? (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {posts.map((post) => (
              <article
                key={post.slug}
                className="group relative overflow-hidden rounded-lg border border-border/50 bg-card p-6 transition-all hover:border-primary/50 hover:shadow-lg"
              >
                <Link
                  href={`/blog/${post.slug}`}
                  className="block space-y-3"
                >
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Calendar size={14} />
                    <time dateTime={post.date}>
                      {formatDate(post.date)}
                    </time>
                    <span className="mx-1">·</span>
                    <Clock size={14} />
                    <span>{post.readingTime} 分钟阅读</span>
                  </div>

                  <h2 className="text-xl font-semibold line-clamp-2 group-hover:text-primary transition-colors">
                    {post.title}
                  </h2>

                  {post.excerpt && (
                    <p className="text-sm text-muted-foreground line-clamp-3">
                      {post.excerpt}
                    </p>
                  )}

                  {post.tags && post.tags.length > 0 && (
                    <div className="flex flex-wrap gap-2">
                      {post.tags.map((tag) => (
                        <span
                          key={tag}
                          className="inline-flex items-center rounded-md bg-primary/10 px-2 py-1 text-xs font-medium text-primary"
                        >
                          {tag}
                        </span>
                      ))}
                    </div>
                  )}
                </Link>
              </article>
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <p className="text-muted-foreground">暂无文章</p>
            <p className="mt-2 text-sm text-muted-foreground">
              在 <code>src/content/posts/</code> 目录下添加 MDX 文件开始写作
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
