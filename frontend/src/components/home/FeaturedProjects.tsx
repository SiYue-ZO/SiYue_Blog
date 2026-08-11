"use client";

/* eslint-disable @next/next/no-img-element */

import { motion } from "framer-motion";
import { Code2, ExternalLink } from "lucide-react";
import Link from "next/link";

import { Project } from "@/types";

interface FeaturedProjectsProps {
  projects: Project[];
}

export function FeaturedProjects({ projects }: FeaturedProjectsProps) {
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
            <h2 className="text-3xl font-bold">精选项目</h2>
            <p className="mt-2 text-muted-foreground">
              我最近参与的一些项目
            </p>
          </div>
          <Link
            href="/projects"
            className="text-sm text-primary hover:underline"
          >
            查看全部 →
          </Link>
        </div>

        {projects.length ? <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {projects.map((project, index) => (
            <motion.div
              key={project.id}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.3, delay: index * 0.1 }}
              whileHover={{ y: -4 }}
              className="group relative overflow-hidden rounded-lg border border-border/50 bg-card p-6 transition-all hover:border-primary/50 hover:shadow-lg"
            >
              {project.image_url && <img src={project.image_url} alt={`${project.title} 项目封面`} loading="lazy" referrerPolicy="no-referrer" className="mb-5 aspect-video w-full rounded-xl border border-border/60 object-cover" />}
              <div className="space-y-3">
                <h3 className="text-xl font-semibold">{project.title}</h3>
                <p className="text-sm text-muted-foreground line-clamp-3">
                  {project.description}
                </p>
                <div className="flex flex-wrap gap-2">
                  {(project.tech_stack ?? []).map((tech) => (
                    <span
                      key={tech}
                      className="inline-flex items-center rounded-md bg-primary/10 px-2 py-1 text-xs font-medium text-primary"
                    >
                      {tech}
                    </span>
                  ))}
                </div>
              </div>

              <div className="mt-4 flex items-center gap-3">
                {project.repo_url && (
                  <a
                    href={project.repo_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
                  >
                    <Code2 size={16} />
                    源码
                  </a>
                )}
                {project.demo_url && (
                  <a
                    href={project.demo_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
                  >
                    <ExternalLink size={16} />
                    演示
                  </a>
                )}
              </div>
            </motion.div>
          ))}
        </div> : <div className="empty-state"><h3>项目正在整理中</h3><p>真实作品会从管理后台发布到这里。</p></div>}
      </motion.div>
    </section>
  );
}
