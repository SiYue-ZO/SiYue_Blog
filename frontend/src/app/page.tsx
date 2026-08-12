import { Hero } from "@/components/home/Hero";
import { IntroSection } from "@/components/home/IntroSection";
import { TechStack } from "@/components/home/TechStack";
import { FeaturedProjects } from "@/components/home/FeaturedProjects";
import { LatestPosts } from "@/components/home/LatestPosts";
import { ContactCTA } from "@/components/home/ContactCTA";
import { AnimatedSection } from "@/components/shared/AnimatedSection";
import { getFeaturedPosts, getHomePageConfig, getPosts, getProjects, getSiteSettings } from "@/lib/content";
import { pageMetadata } from "@/lib/metadata";
import { siteConfig } from "@/lib/constants";
import { Fragment } from "react";

export const dynamic = "force-dynamic";
export const metadata = pageMetadata(siteConfig.name, siteConfig.description, "/");

export default async function Home() {
  const [settings, config] = await Promise.all([
    getSiteSettings(),
    getHomePageConfig(),
  ]);
  const postsSection = config.sections.latest_posts;
  const projectsSection = config.sections.featured_projects;
  const [posts, projects] = await Promise.all([
    postsSection.enabled
      ? postsSection.source === "pinned"
        ? getFeaturedPosts(postsSection.limit)
        : getPosts({ page: 1, page_size: postsSection.limit }).then((response) => response.list)
      : Promise.resolve([]),
    projectsSection.enabled && settings.project_page_visible
      ? getProjects(true).then((items) => items.slice(0, projectsSection.limit))
      : Promise.resolve([]),
  ]);

  const sections = {
    hero: config.sections.hero.enabled ? <Hero config={config.sections.hero} projectPageVisible={settings.project_page_visible} /> : null,
    intro: config.sections.intro.enabled ? <AnimatedSection><IntroSection config={config.sections.intro} /></AnimatedSection> : null,
    tech_stack: config.sections.tech_stack.enabled ? <AnimatedSection><TechStack config={config.sections.tech_stack} /></AnimatedSection> : null,
    featured_projects: projectsSection.enabled && settings.project_page_visible ? <AnimatedSection><FeaturedProjects projects={projects} config={projectsSection} /></AnimatedSection> : null,
    latest_posts: postsSection.enabled ? <AnimatedSection><LatestPosts posts={posts} config={postsSection} /></AnimatedSection> : null,
    contact_cta: config.sections.contact_cta.enabled ? <AnimatedSection><ContactCTA config={config.sections.contact_cta} /></AnimatedSection> : null,
  };

  return (
    <>
      {config.order.map((key) => <Fragment key={key}>{sections[key]}</Fragment>)}
    </>
  );
}
