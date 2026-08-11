/* eslint-disable @next/next/no-img-element */

import { Code2, ExternalLink, FolderGit2 } from "lucide-react";

import { getProjects } from "@/lib/content";
import { pageMetadata } from "@/lib/metadata";

export const metadata = pageMetadata("项目", "查看我的开源项目和作品", "/projects");
export const dynamic = "force-dynamic";

export default async function ProjectsPage() {
  const projects = await getProjects();
  return (
    <div className="page-shell py-12 md:py-16">
      <div className="space-y-8">
        <div className="space-y-3">
          <p className="eyebrow">Portfolio</p>
          <h1 className="text-4xl font-bold">项目</h1>
          <p className="text-lg text-muted-foreground">我参与开发的开源项目和个人作品</p>
        </div>
        {projects.length ? (
          <div className="grid gap-6 md:grid-cols-2">
            {projects.map((project) => (
              <article key={project.id} className="surface-card group p-6">
                {project.image_url && <img src={project.image_url} alt={`${project.title} 项目封面`} loading="lazy" referrerPolicy="no-referrer" className="mb-5 aspect-video w-full rounded-xl border border-border/60 object-cover" />}
                <div className="flex items-start justify-between gap-4">
                  <h2 className="text-xl font-semibold">{project.title}</h2>
                  {project.is_featured && <span className="badge">精选</span>}
                </div>
                <p className="mt-4 text-sm leading-6 text-muted-foreground">{project.description}</p>
                <div className="mt-5 flex flex-wrap gap-2">
                  {(project.tech_stack ?? []).map((tech) => <span key={tech} className="tag">{tech}</span>)}
                </div>
                <div className="mt-6 flex items-center gap-4">
                  {project.repo_url && <a href={project.repo_url} target="_blank" rel="noopener noreferrer" className="text-link"><Code2 size={16} />源码</a>}
                  {project.demo_url && <a href={project.demo_url} target="_blank" rel="noopener noreferrer" className="text-link"><ExternalLink size={16} />在线演示</a>}
                </div>
              </article>
            ))}
          </div>
        ) : (
          <div className="empty-state"><FolderGit2 className="h-8 w-8" /><h2>项目正在整理中</h2><p>稍后会在这里展示真实作品。</p></div>
        )}
      </div>
    </div>
  );
}
