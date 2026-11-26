const uuidV4 = /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

// Persist the event before sending it so reloads and failed-request retries use
// the same database deduplication key throughout this tab's session.
export function getPostViewSessionEvent(postId: number): string | null {
  if (postId <= 0 || typeof window === "undefined") return null;
  try {
    const key = `siyue:post-view:v1:${postId}`;
    const existing = window.sessionStorage.getItem(key);
    if (existing && uuidV4.test(existing)) return existing;
    const eventId = crypto.randomUUID();
    window.sessionStorage.setItem(key, eventId);
    return eventId;
  } catch {
    // Without persistent session storage, skip reporting rather than count
    // every refresh as another reader. The article remains available.
    return null;
  }
}
