import type { Metadata } from "next";
import { ExternalLink, Github } from "lucide-react";

import { api } from "@/lib/api";
import { Project } from "@/types";

export const metadata: Metadata = {
  title: "项目",
  description: "查看我的开源项目和作品",
};

async function getProjects(): Promise<Project[]> {
  try {
    const response = await api.getProjects();
    return (response.data as Project[]) || [];
  } catch {
    // Fallback data when API is not available
    return [
      {
        id: 1,
        title: "SiYue Blog",
        description: "一个现代化的个人博客系统，使用 Next.js 15 和 Go 构建，支持 MDX、暗黑模式、SEO 优化。",
        repo_url: "https://github.com/siyue/blog",
        demo_url: "https://siyue.dev",
        image_url: "",
        tech_stack: ["Next.js", "Go", "PostgreSQL", "Docker"],
        is_featured: true,
        sort_order: 1,
        created_at: "",
        updated_at: "",
      },
      {
        id: 2,
        title: "Task Manager API",
        description: "一个高性能的任务管理 RESTful API，使用 Go 和 Gin 框架构建，支持 JWT 认证。",
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
        description: "一套可复用的 React 组件库，基于 Tailwind CSS 和 Radix UI 构建。",
        repo_url: "https://github.com/siyue/ui-components",
        demo_url: "https://ui.siyue.dev",
        image_url: "",
        tech_stack: ["React", "TypeScript", "Tailwind CSS", "Radix UI"],
        is_featured: false,
        sort_order: 3,
        created_at: "",
        updated_at: "",
      },
      {
        id: 4,
        title: "CLI Tool",
        description: "一个实用的命令行工具，用于快速生成项目模板和自动化开发任务。",
        repo_url: "https://github.com/siyue/cli-tool",
        demo_url: "",
        image_url: "",
        tech_stack: ["Go", "Cobra"],
        is_featured: false,
        sort_order: 4,
        created_at: "",
        updated_at: "",
      },
    ];
  }
}

export default async function ProjectsPage() {
  const projects = await getProjects();

  return (
    <div className="mx-auto w-full max-w-7xl px-4 py-12 sm:px-6 lg:px-8 md:py-16">
      <div className="space-y-8">
        <div className="space-y-4">
          <h1 className="text-4xl font-bold">项目</h1>
          <p className="text-lg text-muted-foreground">
            我参与开发的开源项目和个人作品
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2">
          {projects.map((project) => (
            <div
              key={project.id}
              className="group relative overflow-hidden rounded-lg border border-border/50 bg-card p-6 transition-all hover:border-primary/50 hover:shadow-lg"
            >
              <div className="space-y-4">
                <div className="flex items-start justify-between">
                  <h2 className="text-xl font-semibold">{project.title}</h2>
                  {project.is_featured && (
                    <span className="inline-flex items-center rounded-md bg-primary/10 px-2 py-1 text-xs font-medium text-primary">
                      精选
                    </span>
                  )}
                </div>

                <p className="text-sm text-muted-foreground">
                  {project.description}
                </p>

                <div className="flex flex-wrap gap-2">
                  {project.tech_stack.map((tech) => (
                    <span
                      key={tech}
                      className="inline-flex items-center rounded-md bg-secondary px-2 py-1 text-xs font-medium text-secondary-foreground"
                    >
                      {tech}
                    </span>
                  ))}
                </div>

                <div className="flex items-center gap-3 pt-2">
                  {project.repo_url && (
                    <a
                      href={project.repo_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
                    >
                      <Github size={16} />
                      源码
                    </a>
                  )}
                  {project.demo_url && (
                    <a
                      href={project.demo_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
                    >
                      <ExternalLink size={16} />
                      在线演示
                    </a>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
