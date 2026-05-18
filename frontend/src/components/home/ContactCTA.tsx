"use client";

import { motion } from "framer-motion";
import { Mail, ArrowRight } from "lucide-react";
import Link from "next/link";

export function ContactCTA() {
  return (
    <section className="mx-auto w-full max-w-7xl px-4 py-16 sm:px-6 lg:px-8 md:py-24">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5 }}
        className="relative overflow-hidden rounded-2xl border border-border/50 bg-card p-8 md:p-12 lg:p-16"
      >
        <div className="absolute inset-0 -z-10 bg-[radial-gradient(ellipse_at_center,_var(--tw-gradient-stops))] from-primary/10 via-transparent to-transparent" />

        <div className="mx-auto max-w-2xl text-center">
          <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
            <Mail className="h-8 w-8 text-primary" />
          </div>

          <h2 className="text-3xl font-bold md:text-4xl">
            有想法？聊聊吧
          </h2>

          <p className="mt-4 text-lg text-muted-foreground">
            无论是项目合作、技术交流，还是简单的问候，我都很乐意收到你的消息。
          </p>

          <Link
            href="/contact"
            className="mt-8 inline-flex h-11 items-center justify-center gap-2 rounded-md bg-primary px-8 text-sm font-medium text-primary-foreground shadow transition-colors hover:bg-primary/90"
          >
            联系我
            <ArrowRight size={16} />
          </Link>
        </div>
      </motion.div>
    </section>
  );
}
