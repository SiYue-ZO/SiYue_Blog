import type { Metadata } from "next";
import {
  Code,
  Database,
  Globe,
  Server,
  Terminal,
} from "lucide-react";

export const metadata: Metadata = {
  title: "关于",
  description: "了解我的技能和经历",
};

const skills = [
  {
    category: "前端",
    icon: <Code className="h-5 w-5" />,
    items: ["React", "Next.js", "TypeScript", "Tailwind CSS", "Vue.js"],
  },
  {
    category: "后端",
    icon: <Server className="h-5 w-5" />,
    items: ["Go", "Node.js", "Python", "PostgreSQL", "Redis"],
  },
  {
    category: "DevOps",
    icon: <Terminal className="h-5 w-5" />,
    items: ["Docker", "Kubernetes", "CI/CD", "GitHub Actions", "Linux"],
  },
  {
    category: "其他",
    icon: <Globe className="h-5 w-5" />,
    items: ["Git", "REST API", "GraphQL", "Microservices", "Agile"],
  },
];

const timeline = [
  {
    year: "2026 - 至今",
    title: "全栈开发工程师",
    company: "某科技公司",
    description: "负责核心业务系统的架构设计和开发，使用 Go 和 React 技术栈。",
  },
  {
    year: "2024 - 2026",
    title: "后端开发工程师",
    company: "某互联网公司",
    description: "参与微服务架构的改造和优化，提升了系统性能和可维护性。",
  },
  {
    year: "2022 - 2024",
    title: "初级开发工程师",
    company: "某创业公司",
    description: "从零开始学习 Web 开发，参与多个项目的开发和维护。",
  },
];

export default function AboutPage() {
  return (
    <div className="mx-auto w-full max-w-7xl px-4 py-12 sm:px-6 lg:px-8 md:py-16">
      <div className="mx-auto max-w-3xl space-y-16">
        {/* Header */}
        <div className="space-y-4">
          <h1 className="text-4xl font-bold">关于我</h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            我是一名热爱技术的全栈开发者，专注于构建高质量的 Web 应用。
            我喜欢探索新技术，解决复杂问题，并通过代码创造有价值的产品。
          </p>
        </div>

        {/* Skills */}
        <section className="space-y-8">
          <h2 className="text-2xl font-bold">技能</h2>
          <div className="grid gap-6 md:grid-cols-2">
            {skills.map((skill) => (
              <div
                key={skill.category}
                className="rounded-lg border border-border/50 bg-card p-6"
              >
                <div className="mb-4 flex items-center gap-3">
                  <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary">
                    {skill.icon}
                  </div>
                  <h3 className="text-lg font-semibold">{skill.category}</h3>
                </div>
                <div className="flex flex-wrap gap-2">
                  {skill.items.map((item) => (
                    <span
                      key={item}
                      className="inline-flex items-center rounded-md bg-secondary px-2.5 py-1 text-sm text-secondary-foreground"
                    >
                      {item}
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Timeline */}
        <section className="space-y-8">
          <h2 className="text-2xl font-bold">经历</h2>
          <div className="space-y-8">
            {timeline.map((item, index) => (
              <div
                key={index}
                className="relative pl-8 before:absolute before:left-0 before:top-0 before:h-full before:w-px before:bg-border"
              >
                <div className="absolute left-0 top-0 flex h-6 w-6 -translate-x-1/2 items-center justify-center rounded-full border-2 border-primary bg-background">
                  <div className="h-2 w-2 rounded-full bg-primary" />
                </div>
                <div className="space-y-2">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="text-sm font-medium text-primary">
                      {item.year}
                    </span>
                    <h3 className="text-lg font-semibold">{item.title}</h3>
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {item.company}
                  </p>
                  <p className="text-muted-foreground">{item.description}</p>
                </div>
              </div>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
}
