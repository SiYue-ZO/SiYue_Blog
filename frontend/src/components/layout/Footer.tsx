import Link from "next/link";
import { Code2, MessageCircle, Mail } from "lucide-react";

import { siteConfig } from "@/lib/constants";

export function Footer() {
  return (
    <footer className="border-t border-border/40 bg-background">
      <div className="mx-auto w-full max-w-7xl px-4 py-8 sm:px-6 lg:px-8 md:py-12">
        <div className="flex flex-col items-center justify-between gap-8 md:flex-row">
          <div className="space-y-2 text-center md:text-left">
            <Link href="/" className="inline-flex items-center space-x-2">
              <span className="text-xl font-bold">SiYue</span>
              <span className="text-xl font-bold text-primary">Blog</span>
            </Link>
            <p className="text-sm text-muted-foreground max-w-xs">
              {siteConfig.description}
            </p>
          </div>

          <div className="flex items-center space-x-6">
            {siteConfig.author.github && <a
              href={siteConfig.author.github}
              target="_blank"
              rel="noopener noreferrer"
              className="text-muted-foreground hover:text-primary transition-colors"
              aria-label="GitHub"
            >
              <Code2 size={20} />
            </a>}
            {siteConfig.author.twitter && <a
              href={siteConfig.author.twitter}
              target="_blank"
              rel="noopener noreferrer"
              className="text-muted-foreground hover:text-primary transition-colors"
              aria-label="Twitter"
            >
              <MessageCircle size={20} />
            </a>}
            <a
              href={`mailto:${siteConfig.author.email}`}
              className="text-muted-foreground hover:text-primary transition-colors"
              aria-label="Email"
            >
              <Mail size={20} />
            </a>
          </div>
        </div>

        <div className="mt-8 border-t border-border/40 pt-8 text-center text-sm text-muted-foreground">
          <p>
            &copy; {new Date().getFullYear()} {siteConfig.author.name}. All
            rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
}
