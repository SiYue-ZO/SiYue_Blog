import { Clock, Mail, MapPin } from "lucide-react";

import { ContactForm } from "@/components/contact/ContactForm";
import { siteConfig } from "@/lib/constants";
import { pageMetadata } from "@/lib/metadata";

export const metadata = pageMetadata("联系", "与我取得联系", "/contact");

export default function ContactPage() {
  const cards = [
    { icon: Mail, title: "邮箱", value: siteConfig.author.email },
    { icon: MapPin, title: "位置", value: siteConfig.author.location },
    { icon: Clock, title: "响应时间", value: "通常在 48 小时内" },
  ];
  return (
    <div className="page-shell py-12 md:py-16">
      <div className="mx-auto max-w-3xl space-y-10">
        <div className="space-y-3 text-center"><p className="eyebrow">Contact</p><h1 className="text-4xl font-bold">一起聊聊</h1><p className="text-lg text-muted-foreground">项目合作、技术交流，或只是打个招呼，都很欢迎。</p></div>
        <div className="grid gap-4 sm:grid-cols-3">{cards.map(({ icon: Icon, title, value }) => <div key={title} className="surface-card p-5 text-center"><Icon className="mx-auto h-5 w-5 text-primary" /><h2 className="mt-3 text-sm font-semibold">{title}</h2><p className="mt-1 break-words text-xs text-muted-foreground">{value}</p></div>)}</div>
        <ContactForm />
      </div>
    </div>
  );
}
