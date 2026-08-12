export const siteConfig = {
  name: process.env.NEXT_PUBLIC_SITE_NAME || "SiYue Blog",
  description: process.env.NEXT_PUBLIC_SITE_DESCRIPTION || "记录技术实践与持续思考的个人博客",
  url: process.env.SITE_URL || process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000",
  ogImage: "/images/avatar.jpg",
  author: {
    name: process.env.NEXT_PUBLIC_AUTHOR_NAME || "SiYue",
    bio: process.env.NEXT_PUBLIC_AUTHOR_BIO || "专注于构建可靠、清晰且有温度的软件产品。",
    avatar: "/images/rana.png",
    email: process.env.NEXT_PUBLIC_AUTHOR_EMAIL || "contact@siyue.dev",
    location: process.env.NEXT_PUBLIC_AUTHOR_LOCATION || "中国",
    github: process.env.NEXT_PUBLIC_AUTHOR_GITHUB || "https://github.com/SiYue-ZO",
    twitter: process.env.NEXT_PUBLIC_AUTHOR_TWITTER || "",
  },
  nav: [
    { title: "首页", href: "/" },
    { title: "博客", href: "/blog" },
    { title: "项目", href: "/projects" },
    { title: "关于", href: "/about" },
    { title: "联系", href: "/contact" },
  ],
  techStack: [
    { name: "TypeScript", icon: "typescript" },
    { name: "React", icon: "react" },
    { name: "Next.js", icon: "nextjs" },
    { name: "Go", icon: "go" },
    { name: "PostgreSQL", icon: "postgresql" },
    { name: "Docker", icon: "docker" },
    { name: "Tailwind CSS", icon: "tailwindcss" },
    { name: "Git", icon: "git" },
  ],
};

export const defaultSocialLinks = [
  { platform: "github" as const, label: "GitHub", href: siteConfig.author.github },
  { platform: "email" as const, label: "邮箱", href: `mailto:${siteConfig.author.email}` },
];

export const defaultFooterConfig = {
  enabled: true,
  brand: siteConfig.name,
  description: siteConfig.description,
  copyright: `© {year} ${siteConfig.author.name}. All rights reserved.`,
  show_social_links: true,
  filing_icon_url: "",
  icp_number: "",
  icp_url: "https://beian.miit.gov.cn/",
  public_security_number: "",
  public_security_url: "",
  filing_information: "",
};
