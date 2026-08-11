import { afterEach, describe, expect, it, vi } from "vitest";

import { api } from "./api";

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
