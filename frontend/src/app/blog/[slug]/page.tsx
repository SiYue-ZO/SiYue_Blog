import type { Metadata } from "next";
import { notFound } from "next/navigation";
import Link from "next/link";
import { Calendar, Clock, ArrowLeft, Tag } from "lucide-react";
import { MDXRemote } from "next-mdx-remote/rsc";

import { getPostBySlug, getPostSlugs, getAllPosts } from "@/lib/posts";
import { formatDate } from "@/lib/utils";
import { mdxComponents } from "@/components/mdx";
import { ReadingProgress } from "@/components/blog/ReadingProgress";
import { TableOfContents } from "@/components/blog/TableOfContents";

type Props = {
  params: Promise<{ slug: string }>;
};

export async function generateStaticParams() {
  const slugs = getPostSlugs();
  return slugs.map((slug) => ({ slug }));
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const post = getPostBySlug(slug);

  if (!post) {
    return { title: "文章未找到" };
  }

  return {
    title: post.meta.title,
    description: post.meta.excerpt || post.content.substring(0, 160),
    openGraph: {
      title: post.meta.title,
      description: post.meta.excerpt,
      type: "article",
      publishedTime: post.meta.date,
      tags: post.meta.tags,
    },
  };
}

export default async function BlogPostPage({ params }: Props) {
  const { slug } = await params;
  const post = getPostBySlug(slug);

  if (!post) {
    notFound();
  }

  // Get related posts
  const allPosts = getAllPosts();
  const relatedPosts = allPosts
    .filter(
      (p) =>
        p.slug !== slug &&
        p.tags?.some((tag) => post.meta.tags?.includes(tag))
    )
    .slice(0, 3);

  return (
    <>
      <ReadingProgress />
      <div className="mx-auto w-full max-w-7xl px-4 py-12 sm:px-6 lg:px-8 md:py-16">
        <div className="flex gap-12">
          {/* Main content */}
          <article className="flex-1 max-w-3xl">
            <Link
              href="/blog"
              className="mb-8 inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              <ArrowLeft size={16} />
              返回博客列表
            </Link>

            <header className="mb-10 space-y-4">
              <h1 className="text-4xl font-bold tracking-tight lg:text-5xl">
                {post.meta.title}
              </h1>

              <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
                <div className="flex items-center gap-1.5">
                  <Calendar size={14} />
                  <time dateTime={post.meta.date}>
                    {formatDate(post.meta.date)}
                  </time>
                </div>
                <div className="flex items-center gap-1.5">
                  <Clock size={14} />
                  <span>{post.meta.readingTime} 分钟阅读</span>
                </div>
              </div>

              {post.meta.tags && post.meta.tags.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {post.meta.tags.map((tag) => (
                    <Link
                      key={tag}
                      href={`/tags/${tag.toLowerCase()}`}
                      className="inline-flex items-center gap-1 rounded-md bg-primary/10 px-2.5 py-1 text-xs font-medium text-primary hover:bg-primary/20 transition-colors"
                    >
                      <Tag size={12} />
                      {tag}
                    </Link>
                  ))}
                </div>
              )}
            </header>

            {/* MDX Content */}
            <div className="prose prose-invert max-w-none">
              <MDXRemote source={post.content} components={mdxComponents} />
            </div>

            {/* Related posts */}
            {relatedPosts.length > 0 && (
              <div className="mt-16 border-t border-border/50 pt-8">
                <h3 className="text-xl font-semibold mb-6">相关文章</h3>
                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                  {relatedPosts.map((related) => (
                    <Link
                      key={related.slug}
                      href={`/blog/${related.slug}`}
                      className="group rounded-lg border border-border/50 bg-card p-4 transition-all hover:border-primary/50 hover:shadow-md"
                    >
                      <h4 className="font-medium line-clamp-2 group-hover:text-primary transition-colors">
                        {related.title}
                      </h4>
                      <p className="mt-2 text-sm text-muted-foreground line-clamp-2">
                        {related.excerpt}
                      </p>
                    </Link>
                  ))}
                </div>
              </div>
            )}
          </article>

          {/* Sidebar - Table of Contents */}
          <aside className="hidden xl:block w-64 shrink-0">
            <TableOfContents content={post.content} />
          </aside>
        </div>
      </div>
    </>
  );
}
