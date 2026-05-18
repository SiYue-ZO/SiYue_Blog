"use client";

import { motion } from "framer-motion";
import Image from "next/image";

import { siteConfig } from "@/lib/constants";

export function IntroSection() {
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
          <Image
            src={siteConfig.author.avatar}
            alt={siteConfig.author.name}
            fill
            className="object-cover"
            priority
          />
        </div>

        <div className="space-y-4 text-center md:text-left">
          <h2 className="text-3xl font-bold">关于我</h2>
          <p className="text-lg text-muted-foreground leading-relaxed">
            我是一名热爱技术的全栈开发者，专注于构建现代化的 Web 应用。
            我喜欢探索新技术，分享开发经验，并通过代码创造有价值的产品。
          </p>
          <p className="text-lg text-muted-foreground leading-relaxed">
            这个博客是我记录技术思考、分享项目经验的地方。
            希望我的文章能对你有所帮助。
          </p>
        </div>
      </motion.div>
    </section>
  );
}
