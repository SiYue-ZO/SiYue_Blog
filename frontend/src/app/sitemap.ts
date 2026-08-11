import type { MetadataRoute } from "next";

import { getPosts } from "@/lib/content";
import { siteConfig } from "@/lib/constants";

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const firstPage = await getPosts({ page: 1, page_size: 100 });
  const remainingPages = firstPage.total_pages > 1
    ? await Promise.all(Array.from({ length: firstPage.total_pages - 1 }, (_, index) => getPosts({ page: index + 2, page_size: 100 })))
    : [];
  const posts = [firstPage, ...remainingPages].flatMap((result) => result.list);
  const pages = [
    { path: "", frequency: "weekly" as const, priority: 1 },
    { path: "/blog", frequency: "daily" as const, priority: 0.9 },
    { path: "/projects", frequency: "monthly" as const, priority: 0.7 },
    { path: "/about", frequency: "monthly" as const, priority: 0.6 },
    { path: "/contact", frequency: "yearly" as const, priority: 0.5 },
    { path: "/tags", frequency: "weekly" as const, priority: 0.6 },
  ].map((item) => ({ url: `${siteConfig.url}${item.path}`, changeFrequency: item.frequency, priority: item.priority }));
  const postEntries = posts.map((post) => ({ url: `${siteConfig.url}/blog/${post.slug}`, lastModified: new Date(post.updated_at), changeFrequency: "monthly" as const, priority: 0.8 }));
  return [...pages, ...postEntries];
}
