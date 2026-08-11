/* eslint-disable @next/next/no-img-element */

import Link from "next/link";
import { Calendar, Eye, Search, FileText } from "lucide-react";

import { getPosts, getTags } from "@/lib/content";
import { pageMetadata } from "@/lib/metadata";
import { formatDate } from "@/lib/utils";

export const metadata = pageMetadata("博客", "阅读我的技术文章和思考", "/blog");
export const dynamic = "force-dynamic";

type Props = { searchParams: Promise<{ q?: string; tag?: string; page?: string }> };

export default async function BlogPage({ searchParams }: Props) {
  const query = await searchParams;
  const page = Math.max(1, Number(query.page) || 1);
  const [posts, tags] = await Promise.all([
    getPosts({ page, page_size: 9, search: query.q?.trim(), tag: query.tag }),
    getTags(),
  ]);

  const pageHref = (target: number) => {
    const params = new URLSearchParams();
    if (query.q) params.set("q", query.q);
    if (query.tag) params.set("tag", query.tag);
    params.set("page", String(target));
    return `/blog?${params}`;
  };

  return (
    <div className="page-shell py-12 md:py-16">
      <div className="space-y-8">
        <div className="space-y-3">
          <p className="eyebrow">Writing</p>
          <h1 className="text-4xl font-bold">博客</h1>
          <p className="text-lg text-muted-foreground">记录技术实践，也记录持续思考。</p>
        </div>

        <form action="/blog" className="relative max-w-xl" role="search">
          <label htmlFor="blog-search" className="sr-only">搜索文章</label>
          <Search className="absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <input id="blog-search" name="q" defaultValue={query.q} placeholder="搜索标题与摘要…" className="input pl-11 pr-24" />
          {query.tag && <input type="hidden" name="tag" value={query.tag} />}
          <button className="absolute right-1.5 top-1.5 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground">搜索</button>
        </form>

        <div className="flex flex-wrap gap-2" aria-label="文章标签">
          <Link href={query.q ? `/blog?q=${encodeURIComponent(query.q)}` : "/blog"} className={!query.tag ? "tag-active" : "tag"}>全部</Link>
          {tags.map((tag) => (
            <Link key={tag.id} href={`/blog?tag=${encodeURIComponent(tag.slug)}${query.q ? `&q=${encodeURIComponent(query.q)}` : ""}`} className={query.tag === tag.slug ? "tag-active" : "tag"}>
              {tag.name}<span className="opacity-60">{tag.post_count ?? 0}</span>
            </Link>
          ))}
        </div>

        {posts.list.length ? (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {posts.list.map((post) => (
              <article key={post.id} className="surface-card group p-6">
                <Link href={`/blog/${post.slug}`} className="flex h-full flex-col">
                  {post.cover_image && <img src={post.cover_image} alt="" loading="lazy" referrerPolicy="no-referrer" className="mb-5 aspect-video w-full rounded-xl border border-border/60 object-cover" />}
                  <div className="flex items-center gap-3 text-xs text-muted-foreground">
                    <span className="inline-flex items-center gap-1"><Calendar size={13} />{formatDate(post.published_at)}</span>
                    <span className="inline-flex items-center gap-1"><Eye size={13} />{post.view_count}</span>
                  </div>
                  <h2 className="mt-4 text-xl font-semibold transition-colors group-hover:text-primary">{post.title}</h2>
                  <p className="mt-3 line-clamp-3 flex-1 text-sm leading-6 text-muted-foreground">{post.excerpt}</p>
                  <div className="mt-5 flex flex-wrap gap-2">{post.tags?.map((tag) => <span key={tag.id} className="badge">{tag.name}</span>)}</div>
                </Link>
              </article>
            ))}
          </div>
        ) : (
          <div className="empty-state"><FileText className="h-8 w-8" /><h2>没有找到文章</h2><p>换个关键词或标签试试。</p></div>
        )}

        {posts.total_pages > 1 && (
          <nav className="flex items-center justify-center gap-3" aria-label="文章分页">
            {page > 1 && <Link className="button-secondary" href={pageHref(page - 1)}>上一页</Link>}
            <span className="text-sm text-muted-foreground">{page} / {posts.total_pages}</span>
            {page < posts.total_pages && <Link className="button-secondary" href={pageHref(page + 1)}>下一页</Link>}
          </nav>
        )}
      </div>
    </div>
  );
}
