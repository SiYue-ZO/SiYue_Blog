"use client";

/* eslint-disable @next/next/no-img-element */

import { motion } from "framer-motion";
import Link from "next/link";
import { Calendar, Clock } from "lucide-react";

import type { PublicPost } from "@/types";
import { formatDate } from "@/lib/utils";

interface LatestPostsProps {
  posts: PublicPost[];
}

export function LatestPosts({ posts }: LatestPostsProps) {
  return (
    <section className="mx-auto w-full max-w-7xl px-4 py-16 sm:px-6 lg:px-8 md:py-24">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5 }}
        className="space-y-8"
      >
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-3xl font-bold">最新文章</h2>
            <p className="mt-2 text-muted-foreground">
              分享我的技术思考和经验
            </p>
          </div>
          <Link
            href="/blog"
            className="text-sm text-primary hover:underline"
          >
            查看全部 →
          </Link>
        </div>

        {posts.length ? <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {posts.map((post, index) => (
            <motion.article
              key={post.id}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.3, delay: index * 0.1 }}
              whileHover={{ y: -4 }}
              className="group relative overflow-hidden rounded-lg border border-border/50 bg-card p-6 transition-all hover:border-primary/50 hover:shadow-lg"
            >
              <Link href={`/blog/${post.slug}`} className="block space-y-3">
                {post.cover_image && <img src={post.cover_image} alt="" loading="lazy" referrerPolicy="no-referrer" className="mb-4 aspect-video w-full rounded-xl border border-border/60 object-cover" />}
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Calendar size={14} />
                  <time dateTime={post.published_at}>
                    {formatDate(post.published_at)}
                  </time>
                  <span className="mx-1">·</span>
                  <Clock size={14} />
                  <span>{post.view_count} 次阅读</span>
                </div>

                <h3 className="text-xl font-semibold line-clamp-2 group-hover:text-primary transition-colors">
                  {post.title}
                </h3>

                <p className="text-sm text-muted-foreground line-clamp-3">
                  {post.excerpt}
                </p>

                <div className="flex flex-wrap gap-2">
                  {post.tags?.map((tag) => (
                    <span
                      key={tag.id}
                      className="inline-flex items-center rounded-md bg-primary/10 px-2 py-1 text-xs font-medium text-primary"
                    >
                      {tag.name}
                    </span>
                  ))}
                </div>
              </Link>
            </motion.article>
          ))}
        </div> : <div className="empty-state"><h3>文章即将发布</h3><p>第一篇内容准备好后会出现在这里。</p></div>}
      </motion.div>
    </section>
  );
}
