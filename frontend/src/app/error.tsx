"use client";

import { useEffect } from "react";
import { AlertTriangle } from "lucide-react";

export default function ErrorPage({ error, reset }: { error: Error & { digest?: string }; reset: () => void }) {
  useEffect(() => { console.error(error); }, [error]);
  return <div className="page-shell flex min-h-[65vh] flex-col items-center justify-center text-center"><AlertTriangle className="h-12 w-12 text-destructive" /><h1 className="mt-5 text-3xl font-bold">页面暂时无法加载</h1><p className="mt-3 text-muted-foreground">请稍后重试；如果问题持续存在，请联系站点管理员。</p><button onClick={reset} className="button-primary mt-7">重新加载</button></div>;
}
