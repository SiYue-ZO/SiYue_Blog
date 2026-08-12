"use client";

import { motion, type Variants } from "framer-motion";

import type { HomePageConfig } from "@/types";

const iconVariants: Variants = {
  hidden: { opacity: 0, scale: 0.8 },
  visible: (i: number) => ({
    opacity: 1,
    scale: 1,
    transition: {
      delay: i * 0.1,
      duration: 0.3,
      ease: "easeOut",
    },
  }),
};

export function TechStack({ config }: { config: HomePageConfig["sections"]["tech_stack"] }) {
  return (
    <section className="mx-auto w-full max-w-7xl px-4 py-16 sm:px-6 lg:px-8 md:py-24">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5 }}
        className="space-y-8"
      >
        <div className="text-center">
          <h2 className="text-3xl font-bold">{config.title}</h2>
          <p className="mt-2 text-muted-foreground">
            {config.subtitle}
          </p>
        </div>

        <motion.div
          className="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-4"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true }}
        >
          {config.items.map((tech, index) => (
            <motion.div
              key={`${tech}-${index}`}
              custom={index}
              variants={iconVariants}
              whileHover={{ scale: 1.05 }}
              className="flex flex-col items-center gap-3 rounded-lg border border-border/50 bg-card p-6 transition-colors hover:border-primary/50 hover:bg-accent"
            >
              <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <span className="text-2xl font-bold">
                  {tech.charAt(0)}
                </span>
              </div>
              <span className="text-sm font-medium">{tech}</span>
            </motion.div>
          ))}
        </motion.div>
      </motion.div>
    </section>
  );
}
