import type { Metadata } from "next";
import Link from "next/link";
import { Calendar, Clock, ArrowLeft } from "lucide-react";

import { getPostsByTag, getAllTags } from "@/lib/posts";
import { formatDate } from "@/lib/utils";

type Props = {
  params: Promise<{ tag: string }>;
};

export async function generateStaticParams() {
  const tags = getAllTags();
  return tags.map((tag) => ({ tag: tag.name.toLowerCase() }));
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { tag } = await params;
  return {
    title: `标签: ${tag}`,
    description: `浏览所有关于 ${tag} 的文章`,
  };
}

export default async function TagPage({ params }: Props) {
  const { tag } = await params;
  const posts = getPostsByTag(tag);

  return (
    <div className="mx-auto w-full max-w-7xl px-4 py-12 sm:px-6 lg:px-8 md:py-16">
      <div className="space-y-8">
        <div className="space-y-4">
          <Link
            href="/tags"
            className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            <ArrowLeft size={16} />
            返回标签列表
          </Link>
          <h1 className="text-4xl font-bold">标签: {tag}</h1>
          <p className="text-lg text-muted-foreground">
            共 {posts.length} 篇文章
          </p>
        </div>

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
                </Link>
              </article>
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <p className="text-muted-foreground">暂无相关文章</p>
          </div>
        )}
      </div>
    </div>
  );
}
