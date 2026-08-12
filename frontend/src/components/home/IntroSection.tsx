"use client";

/* eslint-disable @next/next/no-img-element */

import { motion } from "framer-motion";
import Image from "next/image";

import { siteConfig } from "@/lib/constants";
import type { HomePageConfig } from "@/types";

export function IntroSection({ config }: { config: HomePageConfig["sections"]["intro"] }) {
	const imageURL = config.image_url || siteConfig.author.avatar;
  return (
    <section className="mx-auto w-full max-w-7xl px-4 py-16 sm:px-6 lg:px-8 md:py-24">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5 }}
        className="flex flex-col items-center gap-8 md:flex-row md:gap-12"
      >
        <div className="relative h-48 w-48 shrink-0 overflow-hidden rounded-full border-4 border-primary/20">
          {imageURL.startsWith("/") ? <Image src={imageURL} alt={siteConfig.author.name} fill className="object-cover" priority /> : <img src={imageURL} alt={siteConfig.author.name} className="h-full w-full object-cover" referrerPolicy="no-referrer" />}
        </div>

        <div className="space-y-4 text-center md:text-left">
          <h2 className="text-3xl font-bold">{config.title}</h2>
          {config.paragraphs.map((paragraph, index) => <p key={index} className="text-lg text-muted-foreground leading-relaxed">{paragraph}</p>)}
        </div>
      </motion.div>
    </section>
  );
}
