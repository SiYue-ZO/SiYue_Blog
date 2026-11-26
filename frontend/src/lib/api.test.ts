import { afterEach, describe, expect, it, vi } from "vitest";

import { api } from "./api";

describe("post view reporting", () => {
  afterEach(() => vi.unstubAllGlobals());

  it("posts an event separately from fetching article content", async () => {
    vi.stubGlobal("window", {});
    const fetchMock = vi.fn<typeof fetch>(async () => jsonResponse(200, {
      code: 0, message: "ok", data: { view_count: 7 },
    }));
    vi.stubGlobal("fetch", fetchMock);

    await api.getPostBySlug("hello 世界");
    const result = await api.recordPostView("hello 世界", "12345678-1234-4234-8234-123456789012");

    expect(fetchMock.mock.calls[0]?.[0]).toBe("/api/v1/posts/hello%20%E4%B8%96%E7%95%8C");
    expect(fetchMock.mock.calls[1]?.[0]).toBe("/api/v1/posts/hello%20%E4%B8%96%E7%95%8C/views");
    expect(fetchMock.mock.calls[1]?.[1]).toMatchObject({
      method: "POST", cache: "no-store",
      body: JSON.stringify({ event_id: "12345678-1234-4234-8234-123456789012" }),
    });
    expect(result.data.view_count).toBe(7);
  });
});

function jsonResponse(status: number, data: unknown) {
  return new Response(JSON.stringify(data), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

describe("API session refresh", () => {
  afterEach(() => vi.unstubAllGlobals());

  it("shares one refresh request across concurrent admin 401 responses", async () => {
    vi.stubGlobal("window", {});
    let adminAttempts = 0;
    let refreshAttempts = 0;
    const fetchMock = vi.fn(async (input: string | URL | Request) => {
      const url = String(input);
      if (url === "/api/v1/auth/refresh") {
        refreshAttempts += 1;
        return jsonResponse(200, { code: 0, message: "ok", data: { expires_in: 900 } });
      }
      if (url === "/api/admin/stats") {
        adminAttempts += 1;
        if (adminAttempts <= 2) {
          return jsonResponse(401, { code: 2000, message: "expired", data: null });
        }
        return jsonResponse(200, { code: 0, message: "ok", data: {} });
      }
      throw new Error(`unexpected request: ${url}`);
    });
    vi.stubGlobal("fetch", fetchMock);

    await Promise.all([api.getStats(), api.getStats()]);

    expect(refreshAttempts).toBe(1);
    expect(adminAttempts).toBe(4);
  });
});
