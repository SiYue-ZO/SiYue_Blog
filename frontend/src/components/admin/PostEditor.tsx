"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import {
  ArrowLeft,
  Bold,
  Check,
  Clock3,
  Code2,
  Columns2,
  Eye,
  FileText,
  Archive,
  Heading2,
  Heading3,
  ImageIcon,
  Italic,
  Link2,
  List,
  ListOrdered,
  ListTodo,
  Loader2,
  Minus,
  PanelLeft,
  Quote,
  Save,
  Table2,
} from "lucide-react";

import { MarkdownContent } from "@/components/blog/MarkdownContent";
import { api } from "@/lib/api";
import { applyMarkdownEdit, extractMarkdownHeadings, type MarkdownCommand } from "@/lib/markdown-editor";
import { cn, estimateReadingTime, slugify } from "@/lib/utils";
import type { Post, Tag } from "@/types";

type FormState = {
  title: string;
  slug: string;
  excerpt: string;
  content: string;
  cover_image: string;
  status: Post["status"];
  is_pinned: boolean;
  tag_ids: number[];
};

type EditorMode = "write" | "split" | "preview";
type SaveState = "idle" | "local" | "saving" | "saved" | "error";

type LocalDraft = {
  form: FormState;
  savedAt: string;
};

const emptyForm: FormState = {
  title: "",
  slug: "",
  excerpt: "",
  content: "",
  cover_image: "",
  status: "draft",
  is_pinned: false,
  tag_ids: [],
};

const noteTemplate = `## 学习目标

- [ ] 我希望理解什么？

## 核心概念

用自己的话解释概念，并记录一个例子。

## 实践与验证

\`\`\`text
在这里记录代码、命令或实验结果
\`\`\`

## 疑问与延伸

- 仍未解决的问题：
- 下一步：
`;

const toolbarItems: Array<{
  command: MarkdownCommand;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
}> = [
  { command: "heading-2", label: "二级标题", icon: Heading2 },
  { command: "heading-3", label: "三级标题", icon: Heading3 },
  { command: "bold", label: "粗体", icon: Bold },
  { command: "italic", label: "斜体", icon: Italic },
  { command: "inline-code", label: "行内代码", icon: Code2 },
  { command: "link", label: "链接", icon: Link2 },
  { command: "image", label: "图片", icon: ImageIcon },
  { command: "quote", label: "引用", icon: Quote },
  { command: "bullet-list", label: "无序列表", icon: List },
  { command: "ordered-list", label: "有序列表", icon: ListOrdered },
  { command: "task-list", label: "任务清单", icon: ListTodo },
  { command: "code-block", label: "代码块", icon: Code2 },
  { command: "table", label: "表格", icon: Table2 },
  { command: "divider", label: "分隔线", icon: Minus },
];

function serializeForm(form: FormState) {
  return JSON.stringify(form);
}

function isLocalDraft(value: unknown): value is LocalDraft {
  if (!value || typeof value !== "object") return false;
  const candidate = value as Partial<LocalDraft>;
  const draft = candidate.form;
  return typeof candidate.savedAt === "string"
    && Boolean(draft)
    && typeof draft?.title === "string"
    && typeof draft.slug === "string"
    && typeof draft.excerpt === "string"
    && typeof draft.content === "string"
    && typeof draft.cover_image === "string"
    && ["draft", "published", "archived"].includes(draft.status)
    && typeof draft.is_pinned === "boolean"
    && Array.isArray(draft.tag_ids);
}

export function PostEditor({ postId }: { postId?: number }) {
  const router = useRouter();
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const savingRef = useRef(false);
  const lastAutoSaveAttemptRef = useRef("");
  const [currentId, setCurrentId] = useState(postId);
  const [form, setForm] = useState<FormState>(emptyForm);
  const [tags, setTags] = useState<Tag[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [hydrated, setHydrated] = useState(false);
  const [error, setError] = useState("");
  const [slugEdited, setSlugEdited] = useState(Boolean(postId));
  const [mode, setMode] = useState<EditorMode>("split");
  const [saveState, setSaveState] = useState<SaveState>("idle");
  const [savedAt, setSavedAt] = useState<Date | null>(null);
  const [serverStatus, setServerStatus] = useState<Post["status"]>("draft");
  const [serverSnapshot, setServerSnapshot] = useState(serializeForm(emptyForm));

  const storageKey = `siyue-post-draft:${currentId ?? "new"}`;
  const dirty = useMemo(() => serializeForm(form) !== serverSnapshot, [form, serverSnapshot]);
  const headings = useMemo(() => extractMarkdownHeadings(form.content), [form.content]);
  const readingTime = useMemo(() => estimateReadingTime(form.content), [form.content]);

  useEffect(() => {
    let active = true;
    Promise.all([api.getAdminTags(), postId ? api.getAdminPost(postId) : Promise.resolve(null)])
      .then(([tagResponse, postResponse]) => {
        if (!active) return;
        setTags(tagResponse.data);
        const serverForm: FormState = postResponse
          ? {
              title: postResponse.data.title,
              slug: postResponse.data.slug,
              excerpt: postResponse.data.excerpt || "",
              content: postResponse.data.content || "",
              cover_image: postResponse.data.cover_image || "",
              status: postResponse.data.status,
              is_pinned: postResponse.data.is_pinned,
              tag_ids: postResponse.data.tags?.map((tag) => tag.id) || [],
            }
          : emptyForm;

        setForm(serverForm);
        setServerStatus(serverForm.status);
        setServerSnapshot(serializeForm(serverForm));

        const localValue = window.localStorage.getItem(`siyue-post-draft:${postId ?? "new"}`);
        if (localValue) {
          try {
            const localDraft: unknown = JSON.parse(localValue);
            if (isLocalDraft(localDraft) && serializeForm(localDraft.form) !== serializeForm(serverForm)) {
              const restore = window.confirm(`发现 ${new Date(localDraft.savedAt).toLocaleString("zh-CN")} 保存的本地内容，是否恢复？`);
              if (restore) {
                setForm(localDraft.form);
                setSaveState("local");
                setSavedAt(new Date(localDraft.savedAt));
              } else {
                window.localStorage.removeItem(`siyue-post-draft:${postId ?? "new"}`);
              }
            }
          } catch {
            window.localStorage.removeItem(`siyue-post-draft:${postId ?? "new"}`);
          }
        }
      })
      .catch((cause) => setError(cause instanceof Error ? cause.message : "加载失败"))
      .finally(() => {
        if (active) {
          setLoading(false);
          setHydrated(true);
        }
      });
    return () => {
      active = false;
    };
  }, [postId]);

  useEffect(() => {
    if (!hydrated || !dirty) return;
    const timer = window.setTimeout(() => {
      const now = new Date();
      window.localStorage.setItem(storageKey, JSON.stringify({ form, savedAt: now.toISOString() } satisfies LocalDraft));
      setSavedAt(now);
      if (!savingRef.current) setSaveState("local");
    }, 600);
    return () => window.clearTimeout(timer);
  }, [dirty, form, hydrated, storageKey]);

  useEffect(() => {
    if (!dirty) return;
    const warn = (event: BeforeUnloadEvent) => event.preventDefault();
    const guardLinks = (event: MouseEvent) => {
      const link = (event.target as HTMLElement).closest("a");
      if (link && link.href && new URL(link.href, window.location.href).origin === window.location.origin) {
        if (!window.confirm("当前修改尚未保存到服务器，确定离开吗？本地备份仍会保留。")) {
          event.preventDefault();
          event.stopPropagation();
        }
      }
    };
    window.addEventListener("beforeunload", warn);
    document.addEventListener("click", guardLinks, true);
    return () => {
      window.removeEventListener("beforeunload", warn);
      document.removeEventListener("click", guardLinks, true);
    };
  }, [dirty]);

  const set = useCallback(<K extends keyof FormState>(key: K, value: FormState[K]) => {
    setForm((current) => ({ ...current, [key]: value }));
  }, []);

  function updateTitle(title: string) {
    setForm((current) => ({ ...current, title, slug: slugEdited ? current.slug : slugify(title) }));
  }

  function toggleTag(id: number) {
    setForm((current) => ({
      ...current,
      tag_ids: current.tag_ids.includes(id) ? current.tag_ids.filter((item) => item !== id) : [...current.tag_ids, id],
    }));
  }

  const savePost = useCallback(
    async (status: Post["status"], silent = false) => {
      if (savingRef.current) return;
      if (!form.title.trim() || !form.slug.trim()) {
        if (!silent) setError("请先填写标题和链接别名。未完成的内容已保存在当前浏览器中。");
        return;
      }
      if (status === "published" && (!form.excerpt.trim() || !form.content.trim())) {
        setError("发布前请填写摘要和正文。");
        return;
      }

      savingRef.current = true;
      setSaving(true);
      setSaveState("saving");
      setError("");
      const nextForm = { ...form, status };
      const oldStorageKey = storageKey;
      try {
        if (currentId) {
          await api.updatePost(currentId, nextForm);
        } else {
          const response = await api.createPost(nextForm);
          setCurrentId(response.data.id);
          router.replace(`/admin/posts/${response.data.id}`);
        }
        const now = new Date();
        setForm(nextForm);
        setServerStatus(status);
        setServerSnapshot(serializeForm(nextForm));
        setSavedAt(now);
        setSaveState("saved");
        window.localStorage.removeItem(oldStorageKey);
      } catch (cause) {
        setSaveState("error");
        setError(cause instanceof Error ? cause.message : "保存失败");
      } finally {
        savingRef.current = false;
        setSaving(false);
      }
    },
    [currentId, form, router, storageKey],
  );

  const saveLocalDraft = useCallback(() => {
    const now = new Date();
    window.localStorage.setItem(storageKey, JSON.stringify({ form, savedAt: now.toISOString() } satisfies LocalDraft));
    setSavedAt(now);
    setSaveState("local");
  }, [form, storageKey]);

  const saveWork = useCallback(() => {
    if (serverStatus === "published") saveLocalDraft();
    else void savePost(serverStatus === "archived" ? "archived" : "draft");
  }, [saveLocalDraft, savePost, serverStatus]);

  useEffect(() => {
    if (!hydrated || !currentId || serverStatus !== "draft" || form.status !== "draft" || !dirty || saving) return;
    if (!form.title.trim() || !form.slug.trim()) return;
    const snapshot = serializeForm(form);
    if (snapshot === lastAutoSaveAttemptRef.current) return;
    const timer = window.setTimeout(() => {
      lastAutoSaveAttemptRef.current = snapshot;
      void savePost("draft", true);
    }, 2500);
    return () => window.clearTimeout(timer);
  }, [currentId, dirty, form, hydrated, savePost, saving, serverStatus]);

  useEffect(() => {
    const shortcut = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "s") {
        event.preventDefault();
        saveWork();
      }
    };
    window.addEventListener("keydown", shortcut);
    return () => window.removeEventListener("keydown", shortcut);
  }, [saveWork]);

  function runCommand(command: MarkdownCommand) {
    const textarea = textareaRef.current;
    const start = textarea?.selectionStart ?? form.content.length;
    const end = textarea?.selectionEnd ?? start;
    const edit = applyMarkdownEdit(form.content, start, end, command);
    set("content", edit.value);
    requestAnimationFrame(() => {
      textareaRef.current?.focus();
      textareaRef.current?.setSelectionRange(edit.selectionStart, edit.selectionEnd);
    });
  }

  function jumpToHeading(position: number) {
    setMode((current) => (current === "preview" ? "split" : current));
    requestAnimationFrame(() => {
      textareaRef.current?.focus();
      textareaRef.current?.setSelectionRange(position, position);
      const before = form.content.slice(0, position);
      textareaRef.current?.scrollTo({ top: Math.max(0, before.split("\n").length * 24 - 120), behavior: "smooth" });
    });
  }

  function insertTemplate() {
    if (form.content.trim() && !window.confirm("正文中已有内容，确定在末尾追加学习笔记模板吗？")) return;
    set("content", form.content.trim() ? `${form.content.trim()}\n\n${noteTemplate}` : noteTemplate);
  }

  const saveLabel = saveState === "saving"
    ? "保存中…"
    : saveState === "saved"
      ? "已保存到服务器"
      : saveState === "local"
        ? "已备份到浏览器"
        : saveState === "error"
          ? "保存失败"
          : dirty
            ? "有未保存修改"
            : "内容已同步";

  if (loading) {
    return <div className="flex min-h-80 items-center justify-center"><Loader2 className="animate-spin text-primary" /></div>;
  }

  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
        saveWork();
      }}
      className="space-y-5"
    >
      <header className="sticky top-16 z-20 -mx-2 flex flex-wrap items-center justify-between gap-3 border-b border-border/80 bg-background/95 px-2 py-3 backdrop-blur">
        <div className="min-w-0">
          <Link href="/admin/posts" className="text-link"><ArrowLeft size={16} />返回文章</Link>
          <div className="mt-1 flex items-center gap-2 text-xs text-muted-foreground" aria-live="polite">
            {saveState === "saving" ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : saveState === "saved" ? <Check className="h-3.5 w-3.5 text-emerald-600" /> : <Clock3 className="h-3.5 w-3.5" />}
            <span>{saveLabel}{savedAt ? ` · ${savedAt.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" })}` : ""}</span>
          </div>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <button type="button" className="button-secondary" disabled={saving || !dirty} onClick={saveWork}>
            {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
            {serverStatus === "published" ? "保存编辑稿" : serverStatus === "archived" ? "保存归档" : "保存草稿"}
          </button>
          <button type="button" className="button-primary" disabled={saving} onClick={() => void savePost("published")}>
            {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : <Eye className="h-4 w-4" />}
            {serverStatus === "published" ? "更新发布" : "发布笔记"}
          </button>
        </div>
      </header>

      {error && <p className="rounded-xl bg-destructive/10 p-4 text-sm text-destructive" role="alert">{error}</p>}

      <section className="surface-card space-y-4 p-5">
        <div className="space-y-2">
          <label className="sr-only" htmlFor="title">标题</label>
          <input id="title" className="w-full bg-transparent text-2xl font-bold tracking-tight outline-none placeholder:text-muted-foreground/55 md:text-3xl" required maxLength={255} value={form.title} onChange={(event) => updateTitle(event.target.value)} placeholder="这篇学习笔记讲什么？" />
        </div>
        <div className="flex flex-wrap items-center gap-x-5 gap-y-2 border-t border-border/70 pt-3 text-xs text-muted-foreground">
          <span>{form.content.length.toLocaleString()} 字符</span>
          <span>约 {readingTime} 分钟阅读</span>
          <span>{headings.length} 个章节</span>
          {!currentId && <button type="button" className="font-medium text-primary hover:underline" onClick={insertTemplate}>使用学习笔记模板</button>}
        </div>
      </section>

      <section className="surface-card overflow-hidden">
        <div className="flex flex-wrap items-center justify-between gap-3 border-b border-border bg-muted/30 p-2.5">
          <div className="flex flex-wrap items-center gap-1" role="toolbar" aria-label="Markdown 格式工具">
            {toolbarItems.map(({ command, label, icon: Icon }) => (
              <button key={command} type="button" title={label} aria-label={label} onClick={() => runCommand(command)} className="rounded-lg p-2 text-muted-foreground transition hover:bg-accent hover:text-foreground">
                <Icon className="h-4 w-4" />
              </button>
            ))}
          </div>
          <div className="flex rounded-lg border border-border bg-background p-1" aria-label="编辑器显示模式">
            {([
              ["write", "编辑", PanelLeft],
              ["split", "分栏", Columns2],
              ["preview", "预览", Eye],
            ] as const).map(([value, label, Icon]) => (
              <button key={value} type="button" onClick={() => setMode(value)} aria-pressed={mode === value} className={cn("inline-flex items-center gap-1.5 rounded-md px-2.5 py-1.5 text-xs font-medium transition", mode === value ? "bg-primary text-primary-foreground" : "text-muted-foreground hover:text-foreground")}>
                <Icon className="h-3.5 w-3.5" />{label}
              </button>
            ))}
          </div>
        </div>

        <div className={cn("grid min-h-[38rem]", mode === "split" && "lg:grid-cols-2 lg:divide-x lg:divide-border")}>
          {mode !== "preview" && (
            <div className="min-w-0">
              <label className="sr-only" htmlFor="content">正文（Markdown）</label>
              <textarea ref={textareaRef} id="content" className="h-full min-h-[38rem] w-full resize-y bg-transparent p-5 font-mono text-sm leading-6 outline-none placeholder:text-muted-foreground/60" value={form.content} onChange={(event) => set("content", event.target.value)} placeholder={"## 小标题\n\n开始记录今天学到的内容…"} spellCheck={false} />
            </div>
          )}
          {mode !== "write" && (
            <div className="min-w-0 overflow-auto bg-background p-5 lg:max-h-[calc(100vh-11rem)]">
              {form.content.trim() ? <MarkdownContent content={form.content} /> : <div className="empty-state h-full"><Eye className="h-6 w-6" /><h2>预览会显示在这里</h2><p>输入 Markdown，查看最终文章效果。</p></div>}
            </div>
          )}
        </div>
      </section>

      <div className="grid gap-5 xl:grid-cols-[minmax(0,1fr)_20rem]">
        <details className="surface-card p-5" open>
          <summary className="cursor-pointer font-semibold">文章设置</summary>
          <div className="mt-5 grid gap-5 md:grid-cols-2">
            <div className="space-y-2">
              <label className="label" htmlFor="slug">链接别名</label>
              <div className="flex items-center rounded-xl border border-input bg-background/80 pl-4 focus-within:border-primary focus-within:ring-3 focus-within:ring-primary/15">
                <span className="text-sm text-muted-foreground">/blog/</span>
                <input id="slug" className="min-w-0 flex-1 bg-transparent px-1 py-3 text-sm outline-none" required maxLength={255} value={form.slug} onChange={(event) => { setSlugEdited(true); set("slug", slugify(event.target.value)); }} />
              </div>
            </div>
            <div className="space-y-2">
              <label className="label" htmlFor="cover-image">封面地址</label>
              <input id="cover-image" type="url" className="input" value={form.cover_image} onChange={(event) => set("cover_image", event.target.value)} placeholder="https://…" />
            </div>
            <div className="space-y-2 md:col-span-2">
              <div className="flex items-center justify-between"><label className="label" htmlFor="excerpt">摘要</label><span className="text-xs text-muted-foreground">{form.excerpt.length}/500</span></div>
              <textarea id="excerpt" className="input" rows={3} maxLength={500} value={form.excerpt} onChange={(event) => set("excerpt", event.target.value)} placeholder="发布时必填；草稿阶段可以稍后补充" />
            </div>
          </div>
        </details>

        <aside className="space-y-5">
          <section className="surface-card p-5">
            <h2 className="font-semibold">笔记大纲</h2>
            {headings.length ? (
              <nav className="mt-3 max-h-64 space-y-1 overflow-auto" aria-label="笔记大纲">
                {headings.map((heading, index) => (
                  <button key={`${heading.position}-${index}`} type="button" onClick={() => jumpToHeading(heading.position)} className={cn("block w-full truncate rounded-lg px-2 py-1.5 text-left text-sm text-muted-foreground hover:bg-accent hover:text-foreground", heading.level === 3 && "pl-6 text-xs")} title={heading.text}>
                    {heading.text}
                  </button>
                ))}
              </nav>
            ) : <p className="mt-3 text-sm leading-6 text-muted-foreground">使用 H2、H3 标题后，大纲会自动生成。</p>}
          </section>

          <section className="surface-card space-y-4 p-5">
            <h2 className="font-semibold">发布设置</h2>
            <div className="flex items-center justify-between gap-3 text-sm"><span className="text-muted-foreground">当前状态</span><span className="badge">{{ draft: "草稿", published: "已发布", archived: "已归档" }[serverStatus]}</span></div>
            <label className="flex cursor-pointer items-center gap-3 text-sm"><input type="checkbox" checked={form.is_pinned} onChange={(event) => set("is_pinned", event.target.checked)} className="h-4 w-4 accent-primary" />设为精选文章</label>
            {serverStatus === "published" && <p className="rounded-lg bg-amber-500/10 p-3 text-xs leading-5 text-amber-700 dark:text-amber-300">已发布文章不会自动保存到服务器；点击“更新发布”后才会影响线上内容。</p>}
            {currentId && serverStatus !== "archived" && <button type="button" className="button-secondary w-full" disabled={saving} onClick={() => { if (window.confirm("归档后文章将不再公开显示，确定继续吗？")) void savePost("archived"); }}><Archive className="h-4 w-4" />归档文章</button>}
          </section>

          <section className="surface-card p-5">
            <h2 className="font-semibold">标签</h2>
            <div className="mt-3 max-h-64 space-y-2 overflow-auto">
              {tags.length ? tags.map((tag) => <label key={tag.id} className="flex cursor-pointer items-center gap-3 rounded-lg px-2 py-1.5 text-sm hover:bg-accent"><input type="checkbox" checked={form.tag_ids.includes(tag.id)} onChange={() => toggleTag(tag.id)} className="h-4 w-4 accent-primary" />{tag.name}</label>) : <p className="text-sm text-muted-foreground">请先创建标签。</p>}
            </div>
          </section>
        </aside>
      </div>

      <p className="flex items-center gap-2 text-xs text-muted-foreground"><FileText className="h-3.5 w-3.5" />草稿可以暂时不填写摘要和正文；发布前会进行完整检查。按 Ctrl/Cmd + S 快速保存。</p>
    </form>
  );
}
