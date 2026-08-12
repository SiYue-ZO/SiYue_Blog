import "server-only";

import { cache } from "react";

import { api, ApiError } from "@/lib/api";
import {
  getAllPosts as getLocalPosts,
  getAllTags as getLocalTags,
  getPostBySlug as getLocalPostBySlug,
} from "@/lib/posts";
import type { Project, PublicPost, SiteSettings, Tag } from "@/types";
import type { HomePageConfig } from "@/types";
import { defaultHomePageConfig } from "@/lib/home-page";
import { defaultFooterConfig, defaultSocialLinks } from "@/lib/constants";

function localPost(slug: string, index = 0): PublicPost | null {
  const post = getLocalPostBySlug(slug);
  if (!post || post.meta.draft) return null;
  return {
    id: -(index + 1),
    title: post.meta.title,
    slug: post.meta.slug,
    excerpt: post.meta.excerpt || "",
    content: post.content,
    cover_image: post.meta.coverImage || "",
    view_count: 0,
    is_pinned: false,
    tags: (post.meta.tags || []).map((name, tagIndex) => ({
      id: -(tagIndex + 1), name, slug: name.toLowerCase(), created_at: post.meta.date, post_count: 1,
    })),
    published_at: post.meta.date,
    updated_at: post.meta.date,
  };
}

function allowLocalFallback() {
  // Docker production always sets API_INTERNAL_URL. Keeping this fallback only
  // for native development and build-time generation prevents deleted database
  // content from being resurrected by a source-controlled Markdown file.
  return !process.env.API_INTERNAL_URL;
}

export async function getPosts(params: { page?: number; page_size?: number; tag?: string; search?: string } = {}) {
  try {
    return (await api.getPosts(params)).data;
  } catch (error) {
    if (!allowLocalFallback()) throw error;
    const search = params.search?.toLowerCase().trim();
    const tag = params.tag?.toLowerCase();
    const list = getLocalPosts()
      .map((post, index) => localPost(post.slug, index))
      .filter((post): post is PublicPost => Boolean(post))
      .filter((post) => !tag || post.tags?.some((item) => item.slug === tag || item.name.toLowerCase() === tag))
      .filter((post) => !search || `${post.title} ${post.excerpt}`.toLowerCase().includes(search));
    return { list, page: 1, page_size: list.length || 10, total: list.length, total_pages: list.length ? 1 : 0 };
  }
}

export async function getPostBySlug(slug: string) {
  try {
    return (await api.getPostBySlug(slug)).data;
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) return null;
    if (!allowLocalFallback()) throw error;
    return localPost(slug);
  }
}

export async function getTags(): Promise<Tag[]> {
  try {
    return (await api.getTags()).data;
  } catch (error) {
    if (!allowLocalFallback()) throw error;
    return getLocalTags().map((tag, index) => ({
      id: -(index + 1), name: tag.name, slug: tag.name.toLowerCase(), post_count: tag.count, created_at: "",
    }));
  }
}

export async function getProjects(featured = false): Promise<Project[]> {
  try {
    return (await (featured ? api.getFeaturedProjects() : api.getProjects())).data;
  } catch (error) {
    if (!allowLocalFallback()) throw error;
    return [];
  }
}

export async function getFeaturedPosts(limit = 3): Promise<PublicPost[]> {
  try {
    return (await api.getFeaturedPosts(limit)).data;
  } catch (error) {
    if (!allowLocalFallback()) throw error;
    return (await getPosts({ page: 1, page_size: limit })).list.filter((post) => post.is_pinned);
  }
}

export const getSiteSettings = cache(async (): Promise<SiteSettings> => {
  try {
    return (await api.getSettings()).data;
  } catch (error) {
    if (!allowLocalFallback()) throw error;
    return { project_page_visible: true, social_links: defaultSocialLinks, footer: defaultFooterConfig, updated_at: "" };
  }
});

export const getHomePageConfig = cache(async (): Promise<HomePageConfig> => {
  try {
    return (await api.getHomePage()).data;
  } catch (error) {
    if (!allowLocalFallback()) throw error;
    return defaultHomePageConfig;
  }
});
