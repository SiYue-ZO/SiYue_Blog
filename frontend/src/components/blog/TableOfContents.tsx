"use client";

import { useEffect, useState } from "react";

interface Heading {
  id: string;
  text: string;
  level: number;
}

interface TableOfContentsProps {
  content: string;
}

export function TableOfContents({ content }: TableOfContentsProps) {
  const [activeId, setActiveId] = useState<string>("");
  const [headings, setHeadings] = useState<Heading[]>([]);

  useEffect(() => {
    const frame = requestAnimationFrame(() => {
      const elements = Array.from(document.querySelectorAll<HTMLElement>(".article-content h2[id], .article-content h3[id]"));
      setHeadings(elements.map((element) => ({
        id: element.id,
        text: element.textContent?.trim() || "",
        level: Number(element.tagName.slice(1)),
      })).filter((heading) => heading.text));
    });
    return () => cancelAnimationFrame(frame);
  }, [content]);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setActiveId(entry.target.id);
          }
        });
      },
      { rootMargin: "-80px 0px -80% 0px" }
    );

    headings.forEach(({ id }) => {
      const element = document.getElementById(id);
      if (element) {
        observer.observe(element);
      }
    });

    return () => observer.disconnect();
  }, [headings]);

  if (headings.length === 0) {
    return null;
  }

  return (
    <nav className="hidden xl:block sticky top-24 max-h-[calc(100vh-6rem)] overflow-auto">
      <h4 className="text-sm font-semibold mb-4">目录</h4>
      <ul className="space-y-2 text-sm">
        {headings.map((heading) => (
          <li
            key={heading.id}
            style={{ paddingLeft: `${(heading.level - 2) * 1}rem` }}
          >
            <a
              href={`#${heading.id}`}
              className={`block py-1 transition-colors hover:text-foreground ${
                activeId === heading.id
                  ? "text-primary font-medium"
                  : "text-muted-foreground"
              }`}
            >
              {heading.text}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  );
}
