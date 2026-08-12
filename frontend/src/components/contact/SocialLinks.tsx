import {
  AtSign,
  BriefcaseBusiness,
  Code2,
  Globe2,
  Mail,
  MessageCircle,
  Radio,
  Tv2,
  Video,
  type LucideIcon,
} from "lucide-react";

import type { SocialLink, SocialPlatform } from "@/types";

const icons: Record<SocialPlatform, LucideIcon> = {
  github: Code2,
  email: Mail,
  bilibili: Tv2,
  twitter: MessageCircle,
  linkedin: BriefcaseBusiness,
  youtube: Video,
  weibo: Radio,
  website: Globe2,
  other: AtSign,
};

export function SocialLinks({ links }: { links: SocialLink[] }) {
  if (!links.length) return null;

  return (
    <section aria-labelledby="social-links-title" className="space-y-4">
      <div className="text-center">
        <p className="eyebrow">Find me online</p>
        <h2 id="social-links-title" className="mt-2 text-2xl font-bold">社交账号</h2>
      </div>
      <div className="grid gap-3 sm:grid-cols-2">
        {links.map((link, index) => {
          const Icon = icons[link.platform] || AtSign;
          const external = link.href.startsWith("https://");
          return (
            <a
              key={`${link.platform}-${link.href}-${index}`}
              href={link.href}
              target={external ? "_blank" : undefined}
              rel={external ? "noopener noreferrer" : undefined}
              className="surface-card group flex items-center gap-4 p-4 transition hover:-translate-y-0.5 hover:border-primary/40 hover:shadow-md"
            >
              <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
                <Icon size={20} />
              </span>
              <span className="min-w-0">
                <span className="block text-sm font-semibold">{link.label}</span>
                <span className="block truncate text-xs text-muted-foreground">{link.href.replace(/^mailto:/, "")}</span>
              </span>
            </a>
          );
        })}
      </div>
    </section>
  );
}
