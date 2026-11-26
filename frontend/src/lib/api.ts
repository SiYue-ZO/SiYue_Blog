import type {
  ApiResponse,
  ContactMessage,
  DashboardStats,
  HomePageConfig,
  HomePageContent,
  PaginatedResponse,
  Post,
  PublicPost,
  Project,
  SiteSettings,
  Tag,
  UserSession,
} from "@/types";

type FetchOptions = RequestInit & {
  next?: { revalidate?: number; tags?: string[] };
};

export class ApiError extends Error {
  constructor(
    message: string,
    public readonly status: number,
    public readonly code?: number,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

function publicBaseUrl() {
  if (typeof window !== "undefined") return "/api/v1";
  return process.env.API_INTERNAL_URL || "http://localhost:8080/api/v1";
}

function adminBaseUrl() {
  if (typeof window !== "undefined") return "/api/admin";
  if (process.env.API_INTERNAL_ADMIN_URL) return process.env.API_INTERNAL_ADMIN_URL;
  return publicBaseUrl().replace(/\/api\/v1\/?$/, "/api/admin");
}

async function request<T>(baseUrl: string, endpoint: string, options: FetchOptions = {}): Promise<ApiResponse<T>> {
  const timeout = AbortSignal.timeout(10_000);
  const signal = options.signal ? AbortSignal.any([options.signal, timeout]) : timeout;
  const headers = new Headers(options.headers);
  if (options.body && !headers.has("Content-Type")) headers.set("Content-Type", "application/json");

  const response = await fetch(`${baseUrl}${endpoint}`, {
    ...options,
    headers,
    signal,
    credentials: "include",
  });
  const payload = (await response.json().catch(() => null)) as ApiResponse<T> | null;
  if (!response.ok) {
    throw new ApiError(payload?.message || `请求失败 (${response.status})`, response.status, payload?.code);
  }
  if (!payload) throw new ApiError("服务器返回了无效响应", response.status);
  return payload;
}

async function publicRequest<T>(endpoint: string, options: FetchOptions = {}) {
  return request<T>(publicBaseUrl(), endpoint, options);
}

let refreshPromise: Promise<ApiResponse<unknown>> | null = null;

function refreshSession() {
  if (!refreshPromise) {
    refreshPromise = publicRequest("/auth/refresh", { method: "POST" }).finally(() => {
      refreshPromise = null;
    });
  }
  return refreshPromise;
}

async function adminRequest<T>(endpoint: string, options: FetchOptions = {}, retry = true): Promise<ApiResponse<T>> {
  try {
    return await request<T>(adminBaseUrl(), endpoint, { cache: "no-store", ...options });
  } catch (error) {
    if (retry && error instanceof ApiError && error.status === 401 && typeof window !== "undefined") {
      await refreshSession();
      return adminRequest<T>(endpoint, options, false);
    }
    throw error;
  }
}

async function sessionRequest<T>(endpoint: string, retry = true): Promise<ApiResponse<T>> {
  try {
    return await publicRequest<T>(endpoint, { cache: "no-store" });
  } catch (error) {
    if (retry && error instanceof ApiError && error.status === 401 && typeof window !== "undefined") {
      await refreshSession();
      return sessionRequest<T>(endpoint, false);
    }
    throw error;
  }
}

function queryString(params: Record<string, string | number | undefined>) {
  const query = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== "") query.set(key, String(value));
  });
  const value = query.toString();
  return value ? `?${value}` : "";
}

export const api = {
  health: () => publicRequest<{ status: string; version: string }>("/health"),
  getPosts: (params: { page?: number; page_size?: number; tag?: string; search?: string } = {}) =>
    publicRequest<PaginatedResponse<PublicPost>>(`/posts${queryString(params)}`, { cache: "no-store" }),
  getPostBySlug: (slug: string) =>
    publicRequest<PublicPost>(`/posts/${encodeURIComponent(slug)}`, { cache: "no-store" }),
  recordPostView: (slug: string, eventId: string) =>
    publicRequest<{ view_count: number }>(`/posts/${encodeURIComponent(slug)}/views`, {
      method: "POST", body: JSON.stringify({ event_id: eventId }), cache: "no-store",
    }),
  getFeaturedPosts: (limit = 3) => publicRequest<PublicPost[]>(`/posts/featured${queryString({ limit })}`, { cache: "no-store" }),
  getTags: () => publicRequest<Tag[]>("/tags", { cache: "no-store" }),
  getSettings: () => publicRequest<SiteSettings>("/settings", { cache: "no-store" }),
  getHomePage: () => publicRequest<HomePageConfig>("/home", { cache: "no-store" }),
  getProjects: () => publicRequest<Project[]>("/projects", { cache: "no-store" }),
  getFeaturedProjects: () => publicRequest<Project[]>("/projects/featured", { cache: "no-store" }),
  submitContact: (data: { name: string; email: string; message: string; website?: string }) =>
    publicRequest<{ received: boolean }>("/contact", { method: "POST", body: JSON.stringify(data) }),

  login: (email: string, password: string) =>
    publicRequest<{ expires_in: number; user: UserSession }>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),
  me: () => sessionRequest<UserSession>("/auth/me"),
  logout: () => publicRequest<{ logged_out: boolean }>("/auth/logout", { method: "POST" }),

  getStats: () => adminRequest<DashboardStats>("/stats"),
  getAdminSettings: () => adminRequest<SiteSettings>("/settings"),
  updateSettings: (data: Partial<Pick<SiteSettings, "project_page_visible" | "social_links" | "footer">>) =>
    adminRequest<SiteSettings>("/settings", { method: "PATCH", body: JSON.stringify(data) }),
  getAdminHomePage: () => adminRequest<HomePageContent>("/home"),
  saveHomePageDraft: (data: HomePageConfig) =>
    adminRequest<HomePageContent>("/home/draft", { method: "PUT", body: JSON.stringify(data) }),
  publishHomePage: () => adminRequest<HomePageContent>("/home/publish", { method: "POST" }),
  resetHomePageDraft: () => adminRequest<HomePageContent>("/home/reset", { method: "POST" }),
  getAdminPosts: (params: { page?: number; page_size?: number; status?: string; search?: string } = {}) =>
    adminRequest<PaginatedResponse<Post>>(`/posts${queryString(params)}`),
  getAdminPost: (id: number) => adminRequest<Post>(`/posts/${id}`),
  createPost: (data: Partial<Post> & { tag_ids: number[] }) =>
    adminRequest<Post>("/posts", { method: "POST", body: JSON.stringify(data) }),
  updatePost: (id: number, data: Partial<Post> & { tag_ids?: number[] }) =>
    adminRequest<Post>(`/posts/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  deletePost: (id: number) => adminRequest<{ deleted: boolean }>(`/posts/${id}`, { method: "DELETE" }),
  getAdminTags: () => adminRequest<Tag[]>("/tags"),
  createTag: (data: { name: string; slug: string }) =>
    adminRequest<Tag>("/tags", { method: "POST", body: JSON.stringify(data) }),
  deleteTag: (id: number) => adminRequest<{ deleted: boolean }>(`/tags/${id}`, { method: "DELETE" }),
  getAdminProjects: () => adminRequest<Project[]>("/projects"),
  createProject: (data: Omit<Project, "id" | "created_at" | "updated_at">) =>
    adminRequest<Project>("/projects", { method: "POST", body: JSON.stringify(data) }),
  updateProject: (id: number, data: Partial<Project>) =>
    adminRequest<Project>(`/projects/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  deleteProject: (id: number) => adminRequest<{ deleted: boolean }>(`/projects/${id}`, { method: "DELETE" }),
  getMessages: (params: { page?: number; page_size?: number; status?: string } = {}) =>
    adminRequest<PaginatedResponse<ContactMessage>>(`/messages${queryString(params)}`),
  updateMessage: (id: number, status: ContactMessage["status"]) =>
    adminRequest<{ updated: boolean }>(`/messages/${id}`, { method: "PATCH", body: JSON.stringify({ status }) }),
  deleteMessage: (id: number) => adminRequest<{ deleted: boolean }>(`/messages/${id}`, { method: "DELETE" }),
  changePassword: (currentPassword: string, newPassword: string) =>
    adminRequest<{ changed: boolean; login_required: boolean }>("/password", {
      method: "POST",
      body: JSON.stringify({ current_password: currentPassword, new_password: newPassword }),
    }),
};
