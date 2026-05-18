"use client";

import { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import Link from "next/link";
import {
  LayoutDashboard,
  FileText,
  Tags,
  FolderOpen,
  LogOut,
  Menu,
  X,
} from "lucide-react";

const navItems = [
  {
    title: "Dashboard",
    href: "/admin/dashboard",
    icon: <LayoutDashboard size={18} />,
  },
  {
    title: "Posts",
    href: "/admin/posts",
    icon: <FileText size={18} />,
  },
  {
    title: "Tags",
    href: "/admin/tags",
    icon: <Tags size={18} />,
  },
  {
    title: "Projects",
    href: "/admin/projects",
    icon: <FolderOpen size={18} />,
  },
];

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const pathname = usePathname();
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem("access_token");
    if (!token && pathname !== "/admin/login") {
      router.push("/admin/login");
    } else {
      setIsAuthenticated(!!token);
    }
  }, [pathname, router]);

  const handleLogout = () => {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    router.push("/admin/login");
  };

  // Don't show admin layout for login page
  if (pathname === "/admin/login") {
    return <>{children}</>;
  }

  if (!isAuthenticated) {
    return null;
  }

  return (
    <div className="flex min-h-[calc(100vh-4rem)]">
      {/* Sidebar */}
      <aside
        className={`fixed inset-y-16 left-0 z-40 w-64 border-r border-border/40 bg-background transition-transform md:translate-x-0 ${
          sidebarOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <nav className="flex flex-col gap-2 p-4">
          {navItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              onClick={() => setSidebarOpen(false)}
              className={`flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors ${
                pathname === item.href
                  ? "bg-primary/10 text-primary"
                  : "text-muted-foreground hover:bg-accent hover:text-foreground"
              }`}
            >
              {item.icon}
              {item.title}
            </Link>
          ))}

          <button
            onClick={handleLogout}
            className="mt-auto flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-foreground transition-colors"
          >
            <LogOut size={18} />
            Logout
          </button>
        </nav>
      </aside>

      {/* Mobile menu button */}
      <button
        className="fixed bottom-4 right-4 z-50 rounded-full bg-primary p-3 text-primary-foreground shadow-lg md:hidden"
        onClick={() => setSidebarOpen(!sidebarOpen)}
      >
        {sidebarOpen ? <X size={20} /> : <Menu size={20} />}
      </button>

      {/* Main content */}
      <main className="flex-1 md:ml-64">
        <div className="container p-6">{children}</div>
      </main>
    </div>
  );
}
