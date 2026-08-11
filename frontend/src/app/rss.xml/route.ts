import { getPosts } from "@/lib/content";
import { siteConfig } from "@/lib/constants";

function escapeXml(value: string) {
  return value.replace(/[<>&'\"]/g, (character) => ({ "<": "&lt;", ">": "&gt;", "&": "&amp;", "'": "&apos;", '"': "&quot;" })[character] || character);
}

export async function GET() {
  const posts = await getPosts({ page: 1, page_size: 50 });
  const items = posts.list.map((post) => `
    <item>
      <title>${escapeXml(post.title)}</title>
      <link>${siteConfig.url}/blog/${encodeURIComponent(post.slug)}</link>
      <guid isPermaLink="true">${siteConfig.url}/blog/${encodeURIComponent(post.slug)}</guid>
      <description>${escapeXml(post.excerpt)}</description>
      <pubDate>${new Date(post.published_at).toUTCString()}</pubDate>
    </item>`).join("");
  const xml = `<?xml version="1.0" encoding="UTF-8" ?>
  <rss version="2.0"><channel>
    <title>${escapeXml(siteConfig.name)}</title>
    <link>${siteConfig.url}</link>
    <description>${escapeXml(siteConfig.description)}</description>
    <language>zh-CN</language>${items}
  </channel></rss>`;
  return new Response(xml, { headers: { "Content-Type": "application/rss+xml; charset=utf-8", "Cache-Control": "public, max-age=300, stale-while-revalidate=3600" } });
}
