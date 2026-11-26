/* eslint-disable @next/next/no-img-element */

import type { Metadata } from "next";
import { notFound } from "next/navigation";
import Link from "next/link";
import { Calendar, ArrowLeft, Tag } from "lucide-react";

import { getPostBySlug, getPosts } from "@/lib/content";
import { formatDate } from "@/lib/utils";
import { siteConfig } from "@/lib/constants";
import { MarkdownContent } from "@/components/blog/MarkdownContent";
import { ReadingProgress } from "@/components/blog/ReadingProgress";
import { TableOfContents } from "@/components/blog/TableOfContents";
import { PostViewCounter } from "@/components/blog/PostViewCounter";

export const dynamic = "force-dynamic";
type Props = { params: Promise<{ slug: string }> };

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const post = await getPostBySlug(slug);
  if (!post) return { title: "文章未找到" };
  const url = `${siteConfig.url}/blog/${post.slug}`;
  return {
    title: post.title,
    description: post.excerpt,
    alternates: { canonical: url },
    openGraph: {
      title: post.title, description: post.excerpt, type: "article", url,
      publishedTime: post.published_at || undefined, modifiedTime: post.updated_at,
      tags: post.tags?.map((tag) => tag.name), images: post.cover_image ? [post.cover_image] : undefined,
    },
    twitter: { card: "summary_large_image", title: post.title, description: post.excerpt },
  };
}

export default async function BlogPostPage({ params }: Props) {
  const { slug } = await params;
  const post = await getPostBySlug(slug);
  if (!post?.content) notFound();
  const all = await getPosts({ page: 1, page_size: 50 });
  const tagIds = new Set(post.tags?.map((tag) => tag.id));
  const related = all.list.filter((item) => item.slug !== slug && item.tags?.some((tag) => tagIds.has(tag.id))).slice(0, 3);
  const published = post.published_at || post.updated_at;
  const jsonLd = {
    "@context": "https://schema.org", "@type": "BlogPosting", headline: post.title,
    description: post.excerpt, datePublished: published, dateModified: post.updated_at,
    mainEntityOfPage: `${siteConfig.url}/blog/${post.slug}`,
    author: { "@type": "Person", name: post.author?.username || siteConfig.author.name },
  };

  return (
    <>
      <ReadingProgress />
      <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: JSON.stringify(jsonLd).replace(/</g, "\\u003c") }} />
      <div className="page-shell py-10 md:py-16">
        <div className="flex gap-12">
          <article className="min-w-0 max-w-3xl flex-1">
            <Link href="/blog" className="text-link mb-8"><ArrowLeft size={16} />返回文章列表</Link>
            <header className="mb-10 space-y-5">
              <h1 className="text-4xl font-bold tracking-tight lg:text-5xl">{post.title}</h1>
              <p className="text-lg leading-8 text-muted-foreground">{post.excerpt}</p>
              <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
                <span className="inline-flex items-center gap-1.5"><Calendar size={14} /><time dateTime={published}>{formatDate(published)}</time></span>
                <PostViewCounter key={post.id} postId={post.id} slug={post.slug} initialCount={post.view_count} />
              </div>
              <div className="flex flex-wrap gap-2">{post.tags?.map((tag) => <Link key={tag.id} href={`/blog?tag=${encodeURIComponent(tag.slug)}`} className="tag"><Tag size={12} />{tag.name}</Link>)}</div>
            </header>
            {post.cover_image && <img src={post.cover_image} alt={`${post.title} 封面`} referrerPolicy="no-referrer" className="mb-10 aspect-video w-full rounded-2xl border border-border/60 object-cover shadow-sm" />}
            <MarkdownContent content={post.content} />
            {related.length > 0 && (
              <section className="mt-16 border-t border-border pt-8">
                <h2 className="text-xl font-semibold">继续阅读</h2>
                <div className="mt-5 grid gap-4 sm:grid-cols-2">{related.map((item) => <Link key={item.id} href={`/blog/${item.slug}`} className="surface-card p-4"><h3 className="font-medium">{item.title}</h3><p className="mt-2 line-clamp-2 text-sm text-muted-foreground">{item.excerpt}</p></Link>)}</div>
              </section>
            )}
          </article>
          <aside className="hidden w-64 shrink-0 xl:block"><TableOfContents content={post.content} /></aside>
        </div>
      </div>
    </>
  );
}
