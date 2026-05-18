import { ApiResponse } from "@/types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

async function fetchApi<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const url = `${API_BASE_URL}${endpoint}`;

  const defaultHeaders: HeadersInit = {
    "Content-Type": "application/json",
  };

  const response = await fetch(url, {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.message || `API error: ${response.status}`);
  }

  return response.json();
}

export const api = {
  // Health
  health: () => fetchApi<{ status: string; version: string }>("/health"),

  // Posts
  getPosts: (params?: {
    page?: number;
    page_size?: number;
    tag?: string;
    search?: string;
  }) => {
    const searchParams = new URLSearchParams();
    if (params?.page) searchParams.set("page", String(params.page));
    if (params?.page_size) searchParams.set("page_size", String(params.page_size));
    if (params?.tag) searchParams.set("tag", params.tag);
    if (params?.search) searchParams.set("search", params.search);
    const query = searchParams.toString();
    return fetchApi(`/posts${query ? `?${query}` : ""}`);
  },

  getPostBySlug: (slug: string) => fetchApi(`/posts/${slug}`),

  getFeaturedPosts: () => fetchApi("/posts/featured"),

  // Tags
  getTags: () => fetchApi("/tags"),

  // Projects
  getProjects: () => fetchApi("/projects"),
  getFeaturedProjects: () => fetchApi("/projects/featured"),

  // Auth
  login: (email: string, password: string) =>
    fetchApi("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    }),

  register: (username: string, email: string, password: string) =>
    fetchApi("/auth/register", {
      method: "POST",
      body: JSON.stringify({ username, email, password }),
    }),

  refreshToken: (refreshToken: string) =>
    fetchApi("/auth/refresh", {
      method: "POST",
      body: JSON.stringify({ refresh_token: refreshToken }),
    }),

  // Contact
  submitContact: (data: { name: string; email: string; message: string }) =>
    fetchApi("/contact", {
      method: "POST",
      body: JSON.stringify(data),
    }),
};
