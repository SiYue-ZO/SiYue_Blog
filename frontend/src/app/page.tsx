import { Hero } from "@/components/home/Hero";
import { IntroSection } from "@/components/home/IntroSection";
import { TechStack } from "@/components/home/TechStack";
import { FeaturedProjects } from "@/components/home/FeaturedProjects";
import { LatestPosts } from "@/components/home/LatestPosts";
import { ContactCTA } from "@/components/home/ContactCTA";
import { AnimatedSection } from "@/components/shared/AnimatedSection";
import { getPosts, getProjects } from "@/lib/content";
import { pageMetadata } from "@/lib/metadata";
import { siteConfig } from "@/lib/constants";

export const dynamic = "force-dynamic";
export const metadata = pageMetadata(siteConfig.name, siteConfig.description, "/");

export default async function Home() {
  const [projects, posts] = await Promise.all([
    getProjects(true),
    getPosts({ page: 1, page_size: 3 }),
  ]);

  return (
    <>
      <Hero />
      <AnimatedSection><IntroSection /></AnimatedSection>
      <AnimatedSection delay={0.1}><TechStack /></AnimatedSection>
      <AnimatedSection delay={0.2}><FeaturedProjects projects={projects} /></AnimatedSection>
      <AnimatedSection delay={0.3}><LatestPosts posts={posts.list} /></AnimatedSection>
      <AnimatedSection delay={0.4}><ContactCTA /></AnimatedSection>
    </>
  );
}
