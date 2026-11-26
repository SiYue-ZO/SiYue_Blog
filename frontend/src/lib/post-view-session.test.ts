import { afterEach, describe, expect, it, vi } from "vitest";

import { getPostViewSessionEvent } from "./post-view-session";

function session() {
  const values = new Map<string, string>();
  const storage = {
    getItem: (key: string) => values.get(key) ?? null,
    setItem: (key: string, value: string) => { values.set(key, value); },
  };
  vi.stubGlobal("window", { sessionStorage: storage });
  return { values, storage };
}

describe("article view session", () => {
  afterEach(() => vi.unstubAllGlobals());

  it("reuses the event after a reload or a return visit", async () => {
    session();
    const first = getPostViewSessionEvent(12);
    expect(first).toMatch(/^[0-9a-f-]{36}$/);
    getPostViewSessionEvent(13);
    // Reload the module too: identity must come from storage, not module state.
    vi.resetModules();
    const reloaded = await import("./post-view-session");
    expect(reloaded.getPostViewSessionEvent(12)).toBe(first);
    expect(getPostViewSessionEvent(12)).toBe(first);
  });

  it("uses independent events for other articles and new sessions", () => {
    session();
    const first = getPostViewSessionEvent(12);
    expect(getPostViewSessionEvent(13)).not.toBe(first);
    session();
    expect(getPostViewSessionEvent(12)).not.toBe(first);
  });

  it("does not count when session storage is blocked or cannot be written", () => {
    vi.stubGlobal("window", { get sessionStorage() { throw new Error("blocked"); } });
    expect(getPostViewSessionEvent(12)).toBeNull();
    const { storage } = session();
    storage.setItem = () => { throw new Error("quota exceeded"); };
    expect(getPostViewSessionEvent(12)).toBeNull();
  });

  it("replaces an invalid stored event and retains the valid replacement", () => {
    const { values } = session();
    values.set("siyue:post-view:v1:12", "invalid");
    const event = getPostViewSessionEvent(12);
    expect(event).not.toBeNull();
    expect(event).not.toBe("invalid");
    expect(getPostViewSessionEvent(12)).toBe(event);
  });

  it("skips server rendering and local fallback articles", () => {
    vi.stubGlobal("window", undefined);
    expect(getPostViewSessionEvent(12)).toBeNull();
    session();
    expect(getPostViewSessionEvent(-1)).toBeNull();
    expect(getPostViewSessionEvent(0)).toBeNull();
  });
});
