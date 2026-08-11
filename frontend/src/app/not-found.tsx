import Link from "next/link";
import { Compass } from "lucide-react";

export default function NotFound() {
  return <div className="page-shell flex min-h-[65vh] flex-col items-center justify-center text-center"><Compass className="h-12 w-12 text-primary" /><p className="eyebrow mt-6">404</p><h1 className="mt-3 text-4xl font-bold">这里没有你要找的页面</h1><p className="mt-3 text-muted-foreground">链接可能已更新，或内容尚未发布。</p><Link href="/" className="button-primary mt-7">返回首页</Link></div>;
}
