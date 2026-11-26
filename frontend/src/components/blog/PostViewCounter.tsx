"use client";

import { useEffect, useRef, useState } from "react";
import { Eye } from "lucide-react";

import { api } from "@/lib/api";
import { getPostViewSessionEvent } from "@/lib/post-view-session";

export function PostViewCounter({ postId, slug, initialCount }: {
  postId: number;
  slug: string;
  initialCount: number;
}) {
  const [count, setCount] = useState(initialCount);
  const pending = useRef<ReturnType<typeof api.recordPostView> | null>(null);

  useEffect(() => {
    const eventId = getPostViewSessionEvent(postId);
    if (!eventId) return;
    let active = true;
    // Reuse the same request during React Strict Mode's effect replay.
    // Reloads and return visits reuse the ID stored in this tab's session.
    pending.current ??= api.recordPostView(slug, eventId);
    void pending.current.then((response) => {
      if (active) setCount(response.data.view_count);
    }).catch(() => {
      // Keep the server-rendered count if reporting fails; reading still works.
    });
    return () => { active = false; };
  }, [postId, slug]);

  return (
    <span className="inline-flex items-center gap-1.5">
      <Eye size={14} />{count} 次阅读
    </span>
  );
}
