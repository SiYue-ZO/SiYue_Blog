import { CodeBlock } from "./CodeBlock";
import { Callout } from "./Callout";

export const mdxComponents = {
  pre: ({ children, ...props }: React.HTMLAttributes<HTMLPreElement>) => {
    const codeElement = children as React.ReactElement<{
      className?: string;
      children?: string;
    }>;
    if (
      codeElement?.type === "code" &&
      typeof codeElement.props.children === "string"
    ) {
      return (
        <CodeBlock
          className={codeElement.props.className}
          filename={undefined}
        >
          {codeElement.props.children}
        </CodeBlock>
      );
    }
    return <pre {...props}>{children}</pre>;
  },
  Callout,
  // Add more custom components as needed
  h2: ({ children, ...props }: React.HTMLAttributes<HTMLHeadingElement>) => (
    <h2 className="text-2xl font-bold mt-12 mb-4" {...props}>
      {children}
    </h2>
  ),
  h3: ({ children, ...props }: React.HTMLAttributes<HTMLHeadingElement>) => (
    <h3 className="text-xl font-semibold mt-8 mb-3" {...props}>
      {children}
    </h3>
  ),
  p: ({ children, ...props }: React.HTMLAttributes<HTMLParagraphElement>) => (
    <p className="leading-7 mb-4" {...props}>
      {children}
    </p>
  ),
  ul: ({ children, ...props }: React.HTMLAttributes<HTMLUListElement>) => (
    <ul className="list-disc list-inside space-y-2 mb-4" {...props}>
      {children}
    </ul>
  ),
  ol: ({ children, ...props }: React.HTMLAttributes<HTMLOListElement>) => (
    <ol className="list-decimal list-inside space-y-2 mb-4" {...props}>
      {children}
    </ol>
  ),
  li: ({ children, ...props }: React.HTMLAttributes<HTMLLIElement>) => (
    <li className="ml-4" {...props}>
      {children}
    </li>
  ),
  blockquote: ({
    children,
    ...props
  }: React.HTMLAttributes<HTMLQuoteElement>) => (
    <blockquote
      className="border-l-4 border-primary/30 pl-4 italic my-4 text-muted-foreground"
      {...props}
    >
      {children}
    </blockquote>
  ),
  a: ({ children, href, ...props }: React.AnchorHTMLAttributes<HTMLAnchorElement>) => (
    <a
      href={href}
      className="text-primary underline underline-offset-4 hover:text-primary/80 transition-colors"
      target={href?.startsWith("http") ? "_blank" : undefined}
      rel={href?.startsWith("http") ? "noopener noreferrer" : undefined}
      {...props}
    >
      {children}
    </a>
  ),
  img: ({ src, alt, ...props }: React.ImgHTMLAttributes<HTMLImageElement>) => (
    // eslint-disable-next-line @next/next/no-img-element
    <img
      src={src}
      alt={alt || ""}
      className="rounded-lg my-6 max-w-full"
      loading="lazy"
      {...props}
    />
  ),
  hr: () => <hr className="my-8 border-border/50" />,
  table: ({ children, ...props }: React.HTMLAttributes<HTMLTableElement>) => (
    <div className="overflow-x-auto my-6">
      <table className="w-full border-collapse" {...props}>
        {children}
      </table>
    </div>
  ),
  th: ({ children, ...props }: React.HTMLAttributes<HTMLTableCellElement>) => (
    <th
      className="border border-border/50 px-4 py-2 text-left font-semibold bg-muted/50"
      {...props}
    >
      {children}
    </th>
  ),
  td: ({ children, ...props }: React.HTMLAttributes<HTMLTableCellElement>) => (
    <td className="border border-border/50 px-4 py-2" {...props}>
      {children}
    </td>
  ),
};
