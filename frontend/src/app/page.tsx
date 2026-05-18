import { Hero } from "@/components/home/Hero";
import { IntroSection } from "@/components/home/IntroSection";
import { TechStack } from "@/components/home/TechStack";
import { FeaturedProjects } from "@/components/home/FeaturedProjects";
import { LatestPosts } from "@/components/home/LatestPosts";
import { ContactCTA } from "@/components/home/ContactCTA";
import { AnimatedSection } from "@/components/shared/AnimatedSection";
import { api } from "@/lib/api";
import { Post, Project } from "@/types";

async function getFeaturedProjects(): Promise<Project[]> {
  try {
    const response = await api.getFeaturedProjects();
    return (response.data as Project[]) || [];
  } catch {
    // Fallback data when API is not available
    return [
      {
        id: 1,
        title: "SiYue Blog",
        description: "A modern personal blog built with Next.js and Go",
        repo_url: "https://github.com/siyue/blog",
        demo_url: "https://siyue.dev",
        image_url: "",
        tech_stack: ["Next.js", "Go", "PostgreSQL"],
        is_featured: true,
        sort_order: 1,
        created_at: "",
        updated_at: "",
      },
      {
        id: 2,
        title: "Task Manager API",
        description: "A high-performance task management RESTful API",
        repo_url: "https://github.com/siyue/task-manager",
        demo_url: "",
        image_url: "",
        tech_stack: ["Go", "Gin", "PostgreSQL", "Redis"],
        is_featured: true,
        sort_order: 2,
        created_at: "",
        updated_at: "",
      },
      {
        id: 3,
        title: "React Component Library",
        description: "A reusable React component library",
        repo_url: "https://github.com/siyue/ui-components",
        demo_url: "https://ui.siyue.dev",
        image_url: "",
        tech_stack: ["React", "TypeScript", "Tailwind CSS"],
        is_featured: true,
        sort_order: 3,
        created_at: "",
        updated_at: "",
      },
    ];
  }
}

async function getLatestPosts(): Promise<Post[]> {
  try {
    const response = await api.getPosts({ page: 1, page_size: 3 });
    const data = response.data as { list: Post[] };
    return data.list || [];
  } catch {
    // Fallback data when API is not available
    return [
      {
        id: 1,
        title: "使用 Next.js 15 构建现代博客",
        slug: "building-blog-with-nextjs-15",
        excerpt: "探索 Next.js 15 的新特性，以及如何利用 App Router 构建一个高性能的博客系统。",
        cover_image: "",
        status: "published",
        view_count: 128,
        is_pinned: false,
        author_id: 1,
        tags: [
          { id: 1, name: "Next.js", slug: "nextjs", created_at: "" },
          { id: 2, name: "React", slug: "react", created_at: "" },
        ],
        published_at: "2026-05-10T00:00:00Z",
        created_at: "",
        updated_at: "",
      },
      {
        id: 2,
        title: "Go 语言实战：构建 RESTful API",
        slug: "go-restful-api",
        excerpt: "学习如何使用 Go 和 Gin 框架构建一个生产级别的 RESTful API。",
        cover_image: "",
        status: "published",
        view_count: 256,
        is_pinned: false,
        author_id: 1,
        tags: [
          { id: 3, name: "Go", slug: "golang", created_at: "" },
          { id: 4, name: "API", slug: "api", created_at: "" },
        ],
        published_at: "2026-05-05T00:00:00Z",
        created_at: "",
        updated_at: "",
      },
      {
        id: 3,
        title: "Docker 容器化最佳实践",
        slug: "docker-best-practices",
        excerpt: "深入了解 Docker 容器化的最佳实践，包括多阶段构建、安全性和性能优化。",
        cover_image: "",
        status: "published",
        view_count: 192,
        is_pinned: false,
        author_id: 1,
        tags: [
          { id: 5, name: "Docker", slug: "docker", created_at: "" },
          { id: 6, name: "DevOps", slug: "devops", created_at: "" },
        ],
        published_at: "2026-04-28T00:00:00Z",
        created_at: "",
        updated_at: "",
      },
    ];
  }
}

export default async function Home() {
  const [projects, posts] = await Promise.all([
    getFeaturedProjects(),
    getLatestPosts(),
  ]);

  return (
    <>
      <Hero />
      <AnimatedSection>
        <IntroSection />
      </AnimatedSection>
      <AnimatedSection delay={0.1}>
        <TechStack />
      </AnimatedSection>
      <AnimatedSection delay={0.2}>
        <FeaturedProjects projects={projects} />
      </AnimatedSection>
      <AnimatedSection delay={0.3}>
        <LatestPosts posts={posts} />
      </AnimatedSection>
      <AnimatedSection delay={0.4}>
        <ContactCTA />
      </AnimatedSection>
    </>
  );
}
