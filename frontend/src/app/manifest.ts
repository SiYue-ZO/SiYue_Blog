import type { MetadataRoute } from "next";
import { siteConfig } from "@/lib/constants";

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: siteConfig.name, short_name: "SiYue", description: siteConfig.description,
    start_url: "/", display: "standalone", background_color: "#151620", theme_color: "#786fd6",
    icons: [{ src: "/icon.svg", sizes: "any", type: "image/svg+xml" }],
  };
}
