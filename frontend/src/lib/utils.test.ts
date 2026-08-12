import { describe, expect, it } from "vitest";

import { cn, estimateReadingTime, slugify } from "./utils";

describe("slugify", () => {
  it("keeps useful Unicode characters and normalizes separators", () => {
    expect(slugify("  Next.js 16：升级 指南  ")).toBe("next-js-16-升级-指南");
  });

  it("collapses repeated separators", () => {
    expect(slugify("hello___world -- again")).toBe("hello-world-again");
  });
});

describe("estimateReadingTime", () => {
  it("returns zero for empty content", () => {
    expect(estimateReadingTime("   ")).toBe(0);
  });

  it("rounds a partial minute up", () => {
    expect(estimateReadingTime("one two three")).toBe(1);
  });

  it("counts Chinese characters instead of treating the note as one word", () => {
    expect(estimateReadingTime("学".repeat(401))).toBe(2);
  });
});

describe("cn", () => {
  it("resolves conflicting Tailwind classes", () => {
    expect(cn("px-2", false && "hidden", "px-4")).toBe("px-4");
  });
});
