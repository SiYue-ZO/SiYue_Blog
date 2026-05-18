import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { ThemeProvider } from "next-themes";

import { Header } from "@/components/layout/Header";
import { Footer } from "@/components/layout/Footer";

import "./globals.css";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-sans",
});

export const metadata: Metadata = {
  title: {
    default: "SiYue Blog",
    template: "%s | SiYue Blog",
  },
  description: "A modern personal blog and portfolio built with Next.js and Go",
  keywords: ["blog", "portfolio", "developer", "next.js", "golang"],
  authors: [{ name: "SiYue" }],
  openGraph: {
    type: "website",
    locale: "zh_CN",
    siteName: "SiYue Blog",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-CN" suppressHydrationWarning>
      <body className={`${inter.variable} font-sans antialiased`}>
        <ThemeProvider
          attribute="class"
          defaultTheme="dark"
          enableSystem
          disableTransitionOnChange
        >
          <div className="relative flex min-h-screen flex-col">
            <Header />
            <main className="flex-1">{children}</main>
            <Footer />
          </div>
        </ThemeProvider>
      </body>
    </html>
  );
}
