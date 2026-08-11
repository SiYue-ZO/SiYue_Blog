import React, { isValidElement } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeSlug from "rehype-slug";
import rehypeAutolinkHeadings from "rehype-autolink-headings";
import rehypeSanitize from "rehype-sanitize";

import { CodeBlock } from "@/components/mdx/CodeBlock";

export function MarkdownContent({ content }: { content: string }) {
  return (
    <div className="article-content">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[
          rehypeSanitize,
          [rehypeSlug, { prefix: "user-content-" }],
          [rehypeAutolinkHeadings, { behavior: "wrap" }],
        ]}
        components={{
          pre({ children }) {
            if (isValidElement(children)) {
              const props = children.props as { className?: string; children?: React.ReactNode };
              return <CodeBlock className={props.className}>{String(props.children ?? "").replace(/\n$/, "")}</CodeBlock>;
            }
            return <pre>{children}</pre>;
          },
          a({ href, children, node, ...props }) {
            void node;
            const external = href?.startsWith("http");
            return <a href={href} target={external ? "_blank" : undefined} rel={external ? "noopener noreferrer" : undefined} {...props}>{children}</a>;
          },
          img({ alt, node, ...props }) {
            void node;
            // Markdown image dimensions are unknown at authoring time.
            // eslint-disable-next-line @next/next/no-img-element
            return <img alt={alt || ""} loading="lazy" {...props} />;
          },
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
}
