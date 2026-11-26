"use client";

import { useEffect, useRef, useState } from "react";
import { Eye } from "lucide-react";

import { api } from "@/lib/api";

export function PostViewCounter({ slug, initialCount, enabled }: {
  slug: string;
  initialCount: number;
  enabled: boolean;
}) {
  const [count, setCount] = useState(initialCount);
  const pending = useRef<ReturnType<typeof api.recordPostView> | null>(null);

  useEffect(() => {
    if (!enabled) return;
    let active = true;
    // Reuse the same request during React Strict Mode's effect replay.
    // A fresh mount (including a page refresh) creates a new reading event.
    pending.current ??= api.recordPostView(slug, crypto.randomUUID());
    void pending.current.then((response) => {
      if (active) setCount(response.data.view_count);
    }).catch(() => {
      // Keep the server-rendered count if reporting fails; reading still works.
    });
    return () => { active = false; };
  }, [enabled, slug]);

  return (
    <span className="inline-flex items-center gap-1.5">
      <Eye size={14} />{count} 次阅读
    </span>
  );
}
