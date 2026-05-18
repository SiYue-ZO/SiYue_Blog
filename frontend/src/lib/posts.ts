import fs from "fs";
import path from "path";
import matter from "gray-matter";
import readingTime from "reading-time";

const postsDirectory = path.join(process.cwd(), "src/content/posts");

export interface PostFrontmatter {
  title: string;
  date: string;
  excerpt?: string;
  tags?: string[];
  coverImage?: string;
  author?: string;
  draft?: boolean;
}

export interface PostMeta extends PostFrontmatter {
  slug: string;
  readingTime: number;
}

export function getPostSlugs(): string[] {
  if (!fs.existsSync(postsDirectory)) {
    return [];
  }
  return fs
    .readdirSync(postsDirectory)
    .filter((file) => file.endsWith(".mdx"))
    .map((file) => file.replace(/\.mdx$/, ""));
}

export function getPostBySlug(slug: string): {
  meta: PostMeta;
  content: string;
} | null {
  const fullPath = path.join(postsDirectory, `${slug}.mdx`);
  if (!fs.existsSync(fullPath)) {
    return null;
  }

  const fileContents = fs.readFileSync(fullPath, "utf8");
  const { data, content } = matter(fileContents);
  const frontmatter = data as PostFrontmatter;

  const stats = readingTime(content);

  return {
    meta: {
      ...frontmatter,
      slug,
      readingTime: Math.ceil(stats.minutes),
    },
    content,
  };
}

export function getAllPosts(): PostMeta[] {
  const slugs = getPostSlugs();
  const posts = slugs
    .map((slug) => {
      const post = getPostBySlug(slug);
      if (!post) return null;
      return post.meta;
    })
    .filter((post): post is PostMeta => post !== null)
    .filter((post) => !post.draft)
    .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());

  return posts;
}

export function getPostsByTag(tag: string): PostMeta[] {
  return getAllPosts().filter(
    (post) =>
      post.tags?.some((t) => t.toLowerCase() === tag.toLowerCase())
  );
}

export function getAllTags(): { name: string; count: number }[] {
  const posts = getAllPosts();
  const tagMap = new Map<string, number>();

  posts.forEach((post) => {
    post.tags?.forEach((tag) => {
      tagMap.set(tag, (tagMap.get(tag) || 0) + 1);
    });
  });

  return Array.from(tagMap.entries())
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count);
}
