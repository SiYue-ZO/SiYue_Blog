import type { Metadata, Viewport } from "next";
import { ThemeProvider } from "next-themes";

import { Header } from "@/components/layout/Header";
import { Footer } from "@/components/layout/Footer";
import { siteConfig } from "@/lib/constants";
import { getSiteSettings } from "@/lib/content";

import "./globals.css";

export const metadata: Metadata = {
  metadataBase: new URL(siteConfig.url),
  title: { default: siteConfig.name, template: `%s | ${siteConfig.name}` },
  description: siteConfig.description,
  applicationName: siteConfig.name,
  authors: [{ name: siteConfig.author.name, url: siteConfig.author.github }],
  creator: siteConfig.author.name,
  alternates: { types: { "application/rss+xml": "/rss.xml" } },
  openGraph: { type: "website", locale: "zh_CN", siteName: siteConfig.name, title: siteConfig.name, description: siteConfig.description, images: [siteConfig.ogImage] },
  twitter: { card: "summary_large_image", title: siteConfig.name, description: siteConfig.description, images: [siteConfig.ogImage] },
  robots: { index: true, follow: true },
  manifest: "/manifest.webmanifest",
};

export const viewport: Viewport = { colorScheme: "light dark", themeColor: [{ media: "(prefers-color-scheme: light)", color: "#f8f9fc" }, { media: "(prefers-color-scheme: dark)", color: "#151620" }] };

export default async function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  const settings = await getSiteSettings();
  return (
    <html lang="zh-CN" suppressHydrationWarning data-scroll-behavior="smooth">
      <body>
        <ThemeProvider attribute="class" defaultTheme="system" enableSystem disableTransitionOnChange>
          <a href="#main-content" className="fixed left-4 top-4 z-[100] -translate-y-24 rounded-lg bg-primary px-4 py-2 text-primary-foreground transition-transform focus:translate-y-0">跳到主要内容</a>
          <div className="relative flex min-h-screen flex-col">
            <Header projectPageVisible={settings.project_page_visible} />
            <main id="main-content" className="flex-1">{children}</main>
            <Footer config={settings.footer} socialLinks={settings.social_links} />
          </div>
        </ThemeProvider>
      </body>
    </html>
  );
}
