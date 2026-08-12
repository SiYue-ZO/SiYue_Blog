import type { HomePageConfig, HomeSectionKey } from "@/types";

export const homeSectionLabels: Record<HomeSectionKey, string> = {
  hero: "首屏介绍",
  intro: "关于我",
  tech_stack: "技术栈",
  featured_projects: "精选项目",
  latest_posts: "文章列表",
  contact_cta: "联系引导",
};

export const defaultHomePageConfig: HomePageConfig = {
  schema_version: 1,
  order: ["hero", "intro", "tech_stack", "featured_projects", "latest_posts", "contact_cta"],
  sections: {
    hero: {
      enabled: true,
      title_prefix: "你好，我是",
      highlight: "",
      subtitle: "",
      primary_button: { label: "阅读博客", href: "/blog" },
      secondary_button: { label: "查看项目", href: "/projects" },
      show_social_links: true,
    },
    intro: {
      enabled: true,
      title: "关于我",
      image_url: "",
      paragraphs: [
        "我是一名热爱技术的全栈开发者，专注于构建现代化的 Web 应用。我喜欢探索新技术，分享开发经验，并通过代码创造有价值的产品。",
        "这个博客是我记录技术思考、分享项目经验的地方。希望我的文章能对你有所帮助。",
      ],
    },
    tech_stack: {
      enabled: true,
      title: "技术栈",
      subtitle: "我日常使用的技术和工具",
      items: ["TypeScript", "React", "Next.js", "Go", "PostgreSQL", "Docker", "Tailwind CSS", "Git"],
    },
    featured_projects: { enabled: true, title: "精选项目", subtitle: "我最近参与的一些项目", limit: 12 },
    latest_posts: { enabled: true, title: "最新文章", subtitle: "分享我的技术思考和经验", limit: 3, source: "latest" },
    contact_cta: {
      enabled: true,
      title: "有想法？聊聊吧",
      description: "无论是项目合作、技术交流，还是简单的问候，我都很乐意收到你的消息。",
      button: { label: "联系我", href: "/contact" },
    },
  },
};
