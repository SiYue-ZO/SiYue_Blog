export type MarkdownCommand =
  | "heading-2"
  | "heading-3"
  | "bold"
  | "italic"
  | "inline-code"
  | "link"
  | "image"
  | "quote"
  | "bullet-list"
  | "ordered-list"
  | "task-list"
  | "code-block"
  | "table"
  | "divider";

export interface MarkdownEdit {
  value: string;
  selectionStart: number;
  selectionEnd: number;
}

export interface MarkdownHeading {
  level: 2 | 3;
  text: string;
  position: number;
}

function replaceSelection(value: string, start: number, end: number, replacement: string, selectedStart = 0, selectedLength = replacement.length): MarkdownEdit {
  return {
    value: `${value.slice(0, start)}${replacement}${value.slice(end)}`,
    selectionStart: start + selectedStart,
    selectionEnd: start + selectedStart + selectedLength,
  };
}

function wrapSelection(value: string, start: number, end: number, before: string, after: string, placeholder: string): MarkdownEdit {
  const selected = value.slice(start, end) || placeholder;
  return replaceSelection(value, start, end, `${before}${selected}${after}`, before.length, selected.length);
}

function prefixLines(value: string, start: number, end: number, prefix: (index: number) => string): MarkdownEdit {
  const lineStart = value.lastIndexOf("\n", Math.max(0, start - 1)) + 1;
  const nextLineBreak = value.indexOf("\n", end);
  const lineEnd = nextLineBreak === -1 ? value.length : nextLineBreak;
  const selectedLines = value.slice(lineStart, lineEnd);
  const replacement = selectedLines
    .split("\n")
    .map((line, index) => `${prefix(index)}${line}`)
    .join("\n");

  return replaceSelection(value, lineStart, lineEnd, replacement, start - lineStart + prefix(0).length, Math.max(0, end - start));
}

function heading(value: string, start: number, end: number, level: 2 | 3): MarkdownEdit {
  const lineStart = value.lastIndexOf("\n", Math.max(0, start - 1)) + 1;
  const nextLineBreak = value.indexOf("\n", end);
  const lineEnd = nextLineBreak === -1 ? value.length : nextLineBreak;
  const line = value.slice(lineStart, lineEnd);
  const content = line.replace(/^#{1,6}\s+/, "") || "小标题";
  const prefix = `${"#".repeat(level)} `;
  return replaceSelection(value, lineStart, lineEnd, `${prefix}${content}`, prefix.length, content.length);
}

export function applyMarkdownEdit(value: string, start: number, end: number, command: MarkdownCommand): MarkdownEdit {
  switch (command) {
    case "heading-2":
      return heading(value, start, end, 2);
    case "heading-3":
      return heading(value, start, end, 3);
    case "bold":
      return wrapSelection(value, start, end, "**", "**", "重点内容");
    case "italic":
      return wrapSelection(value, start, end, "_", "_", "强调内容");
    case "inline-code":
      return wrapSelection(value, start, end, "`", "`", "code");
    case "link":
      return wrapSelection(value, start, end, "[", "](https://)", "链接文字");
    case "image":
      return wrapSelection(value, start, end, "![", "](https://)", "图片说明");
    case "quote":
      return prefixLines(value, start, end, () => "> ");
    case "bullet-list":
      return prefixLines(value, start, end, () => "- ");
    case "ordered-list":
      return prefixLines(value, start, end, (index) => `${index + 1}. `);
    case "task-list":
      return prefixLines(value, start, end, () => "- [ ] ");
    case "code-block": {
      const selected = value.slice(start, end) || "在这里输入代码";
      return replaceSelection(value, start, end, `\`\`\`text\n${selected}\n\`\`\``, 8, selected.length);
    }
    case "table": {
      const table = "| 知识点 | 说明 |\n| --- | --- |\n| 示例 | 内容 |";
      return replaceSelection(value, start, end, table, 2, 3);
    }
    case "divider":
      return replaceSelection(value, start, end, "\n\n---\n\n", 7, 0);
  }
}

export function extractMarkdownHeadings(content: string): MarkdownHeading[] {
  const headings: MarkdownHeading[] = [];
  const pattern = /^(#{2,3})\s+(.+?)\s*#*\s*$/gm;
  let match: RegExpExecArray | null;

  while ((match = pattern.exec(content)) !== null) {
    headings.push({
      level: match[1].length as 2 | 3,
      text: match[2].replace(/[*_`~[\]]/g, "").trim(),
      position: match.index,
    });
  }

  return headings;
}
