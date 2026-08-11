import type { Metadata } from "next";

import { siteConfig } from "@/lib/constants";

export function pageMetadata(title: string, description: string, path: string): Metadata {
  return {
    title,
    description,
    alternates: { canonical: path },
    openGraph: { title, description, type: "website", url: path, images: [siteConfig.ogImage] },
    twitter: { card: "summary_large_image", title, description, images: [siteConfig.ogImage] },
  };
}
