import { describe, expect, it } from "vitest";

import { applyMarkdownEdit, extractMarkdownHeadings } from "./markdown-editor";

describe("applyMarkdownEdit", () => {
  it("wraps the selected text and keeps it selected", () => {
    expect(applyMarkdownEdit("理解闭包", 2, 4, "bold")).toEqual({
      value: "理解**闭包**",
      selectionStart: 4,
      selectionEnd: 6,
    });
  });

  it("adds task markers to every selected line", () => {
    expect(applyMarkdownEdit("阅读\n练习", 0, 5, "task-list").value).toBe("- [ ] 阅读\n- [ ] 练习");
  });

  it("replaces an existing heading level", () => {
    expect(applyMarkdownEdit("### 原标题", 4, 7, "heading-2").value).toBe("## 原标题");
  });
});

describe("extractMarkdownHeadings", () => {
  it("extracts second- and third-level headings with source positions", () => {
    expect(extractMarkdownHeadings("# 标题\n\n## 基础\n正文\n### `示例`")).toEqual([
      { level: 2, text: "基础", position: 6 },
      { level: 3, text: "示例", position: 15 },
    ]);
  });
});
