"use client";

import { useEffect, useState } from "react";
import { ArrowDown, ArrowUp, ExternalLink, Loader2, RotateCcw, Save, Send } from "lucide-react";

import { api } from "@/lib/api";
import { homeSectionLabels } from "@/lib/home-page";
import type { HomePageConfig, HomePageContent, HomeSectionKey, SiteSettings } from "@/types";

type SectionMap = HomePageConfig["sections"];

export default function AdminHomePage() {
  const [content, setContent] = useState<HomePageContent | null>(null);
  const [settings, setSettings] = useState<SiteSettings | null>(null);
  const [config, setConfig] = useState<HomePageConfig | null>(null);
  const [selected, setSelected] = useState<HomeSectionKey>("hero");
  const [loading, setLoading] = useState(true);
  const [action, setAction] = useState<"save" | "publish" | "reset" | null>(null);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");

  useEffect(() => {
    Promise.all([api.getAdminHomePage(), api.getAdminSettings()])
      .then(([homeResponse, settingsResponse]) => { setContent(homeResponse.data); setConfig(homeResponse.data.draft_config); setSettings(settingsResponse.data); })
      .catch((cause) => setError(cause instanceof Error ? cause.message : "加载首页配置失败"))
      .finally(() => setLoading(false));
  }, []);

  function updateSection<K extends HomeSectionKey>(key: K, values: Partial<SectionMap[K]>) {
    setConfig((current) => current ? ({
      ...current,
      sections: { ...current.sections, [key]: { ...current.sections[key], ...values } },
    } as HomePageConfig) : current);
  }

  function move(key: HomeSectionKey, direction: -1 | 1) {
    setConfig((current) => {
      if (!current) return current;
      const index = current.order.indexOf(key);
      const target = index + direction;
      if (target < 0 || target >= current.order.length) return current;
      const order = [...current.order];
      [order[index], order[target]] = [order[target], order[index]];
      return { ...current, order };
    });
  }

  async function saveDraft(showNotice = true) {
    if (!config) return null;
    const normalized: HomePageConfig = {
      ...config,
      sections: {
        ...config.sections,
        intro: { ...config.sections.intro, paragraphs: config.sections.intro.paragraphs.map((item) => item.trim()).filter(Boolean) },
        tech_stack: { ...config.sections.tech_stack, items: config.sections.tech_stack.items.map((item) => item.trim()).filter(Boolean) },
      },
    };
    const response = await api.saveHomePageDraft(normalized);
    setContent(response.data);
    setConfig(response.data.draft_config);
    if (showNotice) setNotice("草稿已保存，公开首页尚未改变。");
    return response.data;
  }

  async function run(kind: "save" | "publish" | "reset") {
    if (!config || action) return;
    if (kind === "reset" && !window.confirm("放弃当前草稿，恢复为已发布版本？")) return;
    setAction(kind); setError(""); setNotice("");
    try {
      if (kind === "save") await saveDraft();
      if (kind === "publish") {
        await saveDraft(false);
        const response = await api.publishHomePage();
        setContent(response.data);
        setConfig(response.data.draft_config);
        setNotice("首页已发布，刷新公开首页即可看到更新。");
      }
      if (kind === "reset") {
        const response = await api.resetHomePageDraft();
        setContent(response.data);
        setConfig(response.data.draft_config);
        setNotice("草稿已恢复为当前发布版本。");
      }
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "操作失败");
    } finally {
      setAction(null);
    }
  }

  if (loading) return <div className="flex min-h-72 items-center justify-center"><Loader2 className="animate-spin text-primary" /><span className="sr-only">正在加载</span></div>;
  if (!config || !content) return <div className="empty-state"><h1>首页配置暂不可用</h1><p>{error || "请确认数据库迁移和后端服务已启动。"}</p></div>;

  const hero = config.sections.hero;
  const intro = config.sections.intro;
  const tech = config.sections.tech_stack;
  const projects = config.sections.featured_projects;
  const posts = config.sections.latest_posts;
  const contact = config.sections.contact_cta;

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div><p className="eyebrow">Homepage Studio</p><h1 className="mt-2 text-3xl font-bold">首页内容</h1><p className="mt-2 text-sm text-muted-foreground">编辑会先保存为草稿，点击发布后才会影响访客。</p></div>
        <a href="/" target="_blank" rel="noopener noreferrer" className="button-secondary"><ExternalLink size={16} />查看已发布首页</a>
      </div>

      {error && <p className="rounded-xl bg-destructive/10 p-4 text-sm text-destructive" role="alert">{error}</p>}
      {notice && <p className="rounded-xl bg-primary/10 p-4 text-sm text-primary" role="status">{notice}</p>}

      <section className="surface-card flex flex-wrap items-center justify-between gap-4 p-5">
        <div className="text-sm text-muted-foreground">版本 {content.config_version} · 上次发布 {new Date(content.published_at).toLocaleString("zh-CN")}</div>
        <div className="flex flex-wrap gap-2">
          <button type="button" className="button-secondary" disabled={Boolean(action)} onClick={() => run("reset")}><RotateCcw size={16} />{action === "reset" ? "恢复中" : "撤销草稿"}</button>
          <button type="button" className="button-secondary" disabled={Boolean(action)} onClick={() => run("save")}>
            {action === "save" ? <Loader2 size={16} className="animate-spin" /> : <Save size={16} />}保存草稿
          </button>
          <button type="button" className="button-primary" disabled={Boolean(action)} onClick={() => run("publish")}>
            {action === "publish" ? <Loader2 size={16} className="animate-spin" /> : <Send size={16} />}保存并发布
          </button>
        </div>
      </section>

      <div className="grid gap-6 lg:grid-cols-[280px_minmax(0,1fr)]">
        <aside className="surface-card h-fit p-3">
          <p className="px-2 pb-3 text-xs font-semibold uppercase tracking-wider text-muted-foreground">区块顺序</p>
          <div className="space-y-2">
            {config.order.map((key, index) => {
              const enabled = config.sections[key].enabled;
              return <div key={key} className={`rounded-xl border p-2 ${selected === key ? "border-primary bg-primary/5" : "border-border"}`}>
                <button type="button" onClick={() => setSelected(key)} className="flex w-full items-center justify-between gap-2 px-1 py-1 text-left text-sm font-medium"><span>{homeSectionLabels[key]}</span><span className={`text-xs ${enabled ? "text-primary" : "text-muted-foreground"}`}>{enabled ? "显示" : "隐藏"}</span></button>
                <div className="mt-2 flex gap-1">
                  <button type="button" className="rounded-md border p-1.5 disabled:opacity-30" disabled={index === 0} onClick={() => move(key, -1)} aria-label={`上移${homeSectionLabels[key]}`}><ArrowUp size={14} /></button>
                  <button type="button" className="rounded-md border p-1.5 disabled:opacity-30" disabled={index === config.order.length - 1} onClick={() => move(key, 1)} aria-label={`下移${homeSectionLabels[key]}`}><ArrowDown size={14} /></button>
                  <label className="ml-auto flex items-center gap-2 text-xs"><input type="checkbox" className="accent-primary" checked={enabled} onChange={(event) => updateSection(key, { enabled: event.target.checked })} />展示</label>
                </div>
              </div>;
            })}
          </div>
        </aside>

        <section className="surface-card p-6">
          <h2 className="text-xl font-semibold">{homeSectionLabels[selected]}</h2>
          <p className="mt-1 text-sm text-muted-foreground">字段长度和链接安全性会在服务端再次校验。</p>

          {selected === "hero" && <div className="mt-6 space-y-5">
            <Field label="标题前缀"><input className="input" maxLength={200} value={hero.title_prefix} onChange={(e) => updateSection("hero", { title_prefix: e.target.value })} /></Field>
            <Field label="强调文字" hint="留空时使用站点作者名称"><input className="input" maxLength={200} value={hero.highlight} onChange={(e) => updateSection("hero", { highlight: e.target.value })} /></Field>
            <Field label="副标题" hint="留空时使用站点作者简介"><textarea className="input" rows={3} maxLength={200} value={hero.subtitle} onChange={(e) => updateSection("hero", { subtitle: e.target.value })} /></Field>
            <ButtonFields title="主按钮" value={hero.primary_button} onChange={(primary_button) => updateSection("hero", { primary_button })} />
            <ButtonFields title="项目按钮" value={hero.secondary_button} onChange={(secondary_button) => updateSection("hero", { secondary_button })} />
          </div>}

          {selected === "intro" && <div className="mt-6 space-y-5">
            <Field label="标题"><input className="input" required maxLength={100} value={intro.title} onChange={(e) => updateSection("intro", { title: e.target.value })} /></Field>
            <Field label="头像地址" hint="留空使用站点默认头像；支持 /images/... 或 HTTPS"><input className="input" value={intro.image_url} onChange={(e) => updateSection("intro", { image_url: e.target.value })} /></Field>
            <Field label="正文段落" hint="每行作为一个段落，最多 6 段"><textarea className="input" rows={9} value={intro.paragraphs.join("\n")} onChange={(e) => updateSection("intro", { paragraphs: e.target.value.split("\n") })} /></Field>
          </div>}

          {selected === "tech_stack" && <div className="mt-6 space-y-5">
            <Field label="标题"><input className="input" maxLength={100} value={tech.title} onChange={(e) => updateSection("tech_stack", { title: e.target.value })} /></Field>
            <Field label="说明"><input className="input" maxLength={200} value={tech.subtitle} onChange={(e) => updateSection("tech_stack", { subtitle: e.target.value })} /></Field>
            <Field label="技术项目" hint="每行填写一个项目，顺序即展示顺序"><textarea className="input" rows={7} value={tech.items.join("\n")} onChange={(e) => updateSection("tech_stack", { items: e.target.value.split("\n") })} /></Field>
          </div>}

          {selected === "featured_projects" && <div className="mt-6 space-y-5">
            {settings?.project_page_visible === false && <p className="rounded-xl bg-amber-500/10 p-4 text-sm text-amber-700 dark:text-amber-300">项目页面当前已关闭，因此即使启用这个区块，公开首页也不会展示项目。</p>}
            <Field label="标题"><input className="input" maxLength={100} value={projects.title} onChange={(e) => updateSection("featured_projects", { title: e.target.value })} /></Field>
            <Field label="说明"><input className="input" maxLength={200} value={projects.subtitle} onChange={(e) => updateSection("featured_projects", { subtitle: e.target.value })} /></Field>
            <Field label="展示数量" hint="展示项目管理中标记为首页精选且公开的项目"><input type="number" min={1} max={12} className="input w-32" value={projects.limit} onChange={(e) => updateSection("featured_projects", { limit: Number(e.target.value) })} /></Field>
          </div>}

          {selected === "latest_posts" && <div className="mt-6 space-y-5">
            <Field label="标题"><input className="input" maxLength={100} value={posts.title} onChange={(e) => updateSection("latest_posts", { title: e.target.value })} /></Field>
            <Field label="说明"><input className="input" maxLength={200} value={posts.subtitle} onChange={(e) => updateSection("latest_posts", { subtitle: e.target.value })} /></Field>
            <div className="grid gap-4 sm:grid-cols-2">
              <Field label="内容来源"><select className="input" value={posts.source} onChange={(e) => updateSection("latest_posts", { source: e.target.value as "latest" | "pinned" })}><option value="latest">最新发布</option><option value="pinned">置顶文章</option></select></Field>
              <Field label="展示数量"><input type="number" min={1} max={12} className="input" value={posts.limit} onChange={(e) => updateSection("latest_posts", { limit: Number(e.target.value) })} /></Field>
            </div>
          </div>}

          {selected === "contact_cta" && <div className="mt-6 space-y-5">
            <Field label="标题"><input className="input" maxLength={100} value={contact.title} onChange={(e) => updateSection("contact_cta", { title: e.target.value })} /></Field>
            <Field label="说明"><textarea className="input" rows={4} maxLength={500} value={contact.description} onChange={(e) => updateSection("contact_cta", { description: e.target.value })} /></Field>
            <ButtonFields title="按钮" value={contact.button} onChange={(button) => updateSection("contact_cta", { button })} />
          </div>}
        </section>
      </div>
    </div>
  );
}

function Field({ label, hint, children }: { label: string; hint?: string; children: React.ReactNode }) {
  return <label className="block space-y-2"><span className="label">{label}</span>{children}{hint && <span className="block text-xs text-muted-foreground">{hint}</span>}</label>;
}

function ButtonFields({ title, value, onChange }: { title: string; value: { label: string; href: string }; onChange: (value: { label: string; href: string }) => void }) {
  return <fieldset className="rounded-xl border border-border p-4"><legend className="px-2 text-sm font-semibold">{title}</legend><div className="grid gap-4 sm:grid-cols-2"><Field label="文案"><input className="input" maxLength={40} value={value.label} onChange={(e) => onChange({ ...value, label: e.target.value })} /></Field><Field label="链接" hint="站内路径、HTTPS 或 mailto"><input className="input" maxLength={512} value={value.href} onChange={(e) => onChange({ ...value, href: e.target.value })} /></Field></div></fieldset>;
}
