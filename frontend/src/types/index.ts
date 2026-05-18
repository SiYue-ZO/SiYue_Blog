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
  published_at: string;
  created_at: string;
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

export interface Frontmatter {
  title: string;
  date: string;
  excerpt?: string;
  tags?: string[];
  coverImage?: string;
  author?: string;
  readingTime?: number;
}
