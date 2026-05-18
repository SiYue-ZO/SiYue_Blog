"use client";

import { useEffect, useState } from "react";
import {
  FileText,
  Tags,
  FolderOpen,
  Eye,
  TrendingUp,
} from "lucide-react";

import { api } from "@/lib/api";

interface Stats {
  totalPosts: number;
  totalTags: number;
  totalProjects: number;
  totalViews: number;
}

export default function AdminDashboardPage() {
  const [stats, setStats] = useState<Stats>({
    totalPosts: 0,
    totalTags: 0,
    totalProjects: 0,
    totalViews: 0,
  });

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [postsRes, tagsRes, projectsRes] = await Promise.all([
          api.getPosts({ page: 1, page_size: 1 }),
          api.getTags(),
          api.getProjects(),
        ]);

        const postsData = postsRes.data as { total: number };
        const tagsData = tagsRes.data as unknown[];
        const projectsData = projectsRes.data as unknown[];

        setStats({
          totalPosts: postsData.total || 0,
          totalTags: tagsData?.length || 0,
          totalProjects: projectsData?.length || 0,
          totalViews: 0, // TODO: Implement view tracking
        });
      } catch {
        // Use fallback data
        setStats({
          totalPosts: 3,
          totalTags: 10,
          totalProjects: 4,
          totalViews: 576,
        });
      }
    };

    fetchStats();
  }, []);

  const statCards = [
    {
      title: "Total Posts",
      value: stats.totalPosts,
      icon: <FileText className="h-6 w-6" />,
      color: "text-blue-500",
      bgColor: "bg-blue-500/10",
    },
    {
      title: "Total Tags",
      value: stats.totalTags,
      icon: <Tags className="h-6 w-6" />,
      color: "text-green-500",
      bgColor: "bg-green-500/10",
    },
    {
      title: "Total Projects",
      value: stats.totalProjects,
      icon: <FolderOpen className="h-6 w-6" />,
      color: "text-purple-500",
      bgColor: "bg-purple-500/10",
    },
    {
      title: "Total Views",
      value: stats.totalViews,
      icon: <Eye className="h-6 w-6" />,
      color: "text-orange-500",
      bgColor: "bg-orange-500/10",
    },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground">
          Welcome to the admin panel
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {statCards.map((card) => (
          <div
            key={card.title}
            className="rounded-lg border border-border/50 bg-card p-6"
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-muted-foreground">{card.title}</p>
                <p className="mt-2 text-3xl font-bold">{card.value}</p>
              </div>
              <div
                className={`flex h-12 w-12 items-center justify-center rounded-lg ${card.bgColor} ${card.color}`}
              >
                {card.icon}
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Quick Actions */}
      <div className="rounded-lg border border-border/50 bg-card p-6">
        <h2 className="text-xl font-semibold mb-4">Quick Actions</h2>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <a
            href="/admin/posts"
            className="flex items-center gap-3 rounded-md border border-border/50 p-4 transition-colors hover:bg-accent"
          >
            <FileText className="h-5 w-5 text-primary" />
            <div>
              <p className="font-medium">Manage Posts</p>
              <p className="text-sm text-muted-foreground">
                Create, edit, or delete posts
              </p>
            </div>
          </a>

          <a
            href="/admin/tags"
            className="flex items-center gap-3 rounded-md border border-border/50 p-4 transition-colors hover:bg-accent"
          >
            <Tags className="h-5 w-5 text-primary" />
            <div>
              <p className="font-medium">Manage Tags</p>
              <p className="text-sm text-muted-foreground">
                Organize your content with tags
              </p>
            </div>
          </a>

          <a
            href="/admin/projects"
            className="flex items-center gap-3 rounded-md border border-border/50 p-4 transition-colors hover:bg-accent"
          >
            <FolderOpen className="h-5 w-5 text-primary" />
            <div>
              <p className="font-medium">Manage Projects</p>
              <p className="text-sm text-muted-foreground">
                Showcase your work
              </p>
            </div>
          </a>
        </div>
      </div>

      {/* Recent Activity */}
      <div className="rounded-lg border border-border/50 bg-card p-6">
        <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
          <TrendingUp className="h-5 w-5" />
          Recent Activity
        </h2>
        <div className="space-y-4">
          <div className="flex items-center gap-4 rounded-md bg-muted/50 p-4">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
              <FileText className="h-5 w-5 text-primary" />
            </div>
            <div>
              <p className="font-medium">New post published</p>
              <p className="text-sm text-muted-foreground">
                &quot;使用 Next.js 15 构建现代博客&quot; was published
              </p>
            </div>
          </div>
          <div className="flex items-center gap-4 rounded-md bg-muted/50 p-4">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-green-500/10">
              <Tags className="h-5 w-5 text-green-500" />
            </div>
            <div>
              <p className="font-medium">New tag created</p>
              <p className="text-sm text-muted-foreground">
                &quot;Next.js&quot; tag was added
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
