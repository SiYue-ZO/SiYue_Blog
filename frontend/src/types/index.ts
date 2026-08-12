export interface Post {
  id: number;
  title: string;
  slug: string;
  excerpt: string;
  content?: string;
  cover_image: string;
  status: "draft" | "published" | "archived";
  view_count: number;
  is_pinned: boolean;
  author_id: number;
  author?: User;
  tags?: Tag[];
  published_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface PublicPost {
  id: number;
  title: string;
  slug: string;
  excerpt: string;
  content?: string;
  cover_image: string;
  view_count: number;
  is_pinned: boolean;
  author?: Pick<User, "username" | "avatar_url">;
  tags: Tag[];
  published_at: string;
  updated_at: string;
}

export interface User {
  id: number;
  username: string;
  email: string;
  avatar_url: string;
  role: "admin" | "user";
  status: "active" | "inactive" | "banned";
  created_at: string;
  updated_at: string;
}

export interface Tag {
  id: number;
  name: string;
  slug: string;
  created_at: string;
  post_count?: number;
}

export interface Project {
  id: number;
  title: string;
  description: string;
  repo_url: string;
  demo_url: string;
  image_url: string;
  tech_stack: string[];
  is_featured: boolean;
  is_visible: boolean;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface SiteSettings {
  project_page_visible: boolean;
  social_links: SocialLink[] | null;
  footer: FooterConfig | null;
  updated_at: string;
}

export interface FooterConfig {
  enabled: boolean;
  brand: string;
  description: string;
  copyright: string;
  show_social_links: boolean;
  filing_icon_url: string;
  icp_number: string;
  icp_url: string;
  public_security_number: string;
  public_security_url: string;
  filing_information: string;
}

export type SocialPlatform = "github" | "email" | "bilibili" | "twitter" | "linkedin" | "youtube" | "weibo" | "website" | "other";

export interface SocialLink {
  platform: SocialPlatform;
  label: string;
  href: string;
}

export type HomeSectionKey = "hero" | "intro" | "tech_stack" | "featured_projects" | "latest_posts" | "contact_cta";

export interface HomeButton {
  label: string;
  href: string;
}

export interface HomePageConfig {
  schema_version: 1;
  order: HomeSectionKey[];
  sections: {
    hero: {
      enabled: boolean;
      title_prefix: string;
      highlight: string;
      subtitle: string;
      primary_button: HomeButton;
      secondary_button: HomeButton;
      show_social_links: boolean;
    };
    intro: {
      enabled: boolean;
      title: string;
      image_url: string;
      paragraphs: string[];
    };
    tech_stack: {
      enabled: boolean;
      title: string;
      subtitle: string;
      items: string[];
    };
    featured_projects: {
      enabled: boolean;
      title: string;
      subtitle: string;
      limit: number;
    };
    latest_posts: {
      enabled: boolean;
      title: string;
      subtitle: string;
      limit: number;
      source: "latest" | "pinned";
    };
    contact_cta: {
      enabled: boolean;
      title: string;
      description: string;
      button: HomeButton;
    };
  };
}

export interface HomePageContent {
  draft_config: HomePageConfig;
  published_config: HomePageConfig;
  config_version: number;
  updated_by?: number;
  published_at: string;
  updated_at: string;
}

export interface PaginatedResponse<T> {
  list: T[];
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export interface UserSession {
  id: number;
  username: string;
  role: "admin" | "user";
}

export interface DashboardStats {
  total_posts: number;
  published_posts: number;
  draft_posts: number;
  total_tags: number;
  total_projects: number;
  total_views: number;
  unread_messages: number;
}

export interface ContactMessage {
  id: number;
  name: string;
  email: string;
  message: string;
  status: "unread" | "read" | "archived";
  created_at: string;
  updated_at: string;
}

export interface Frontmatter {
  title: string;
  date: string;
  excerpt?: string;
  tags?: string[];
  coverImage?: string;
  author?: string;
  readingTime?: number;
}
