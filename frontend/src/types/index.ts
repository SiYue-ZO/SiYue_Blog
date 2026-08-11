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
