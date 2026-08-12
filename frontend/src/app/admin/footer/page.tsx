"use client";

import { useEffect, useMemo, useState } from "react";
import { ExternalLink, Loader2, RotateCcw, Save, ShieldCheck } from "lucide-react";

import { Footer } from "@/components/layout/Footer";
import { api } from "@/lib/api";
import { defaultFooterConfig } from "@/lib/constants";
import type { FooterConfig, SocialLink } from "@/types";

export default function AdminFooterPage() {
  const [config, setConfig] = useState<FooterConfig>(defaultFooterConfig);
  const [savedConfig, setSavedConfig] = useState<FooterConfig>(defaultFooterConfig);
  const [socialLinks, setSocialLinks] = useState<SocialLink[]>([]);
  const [loadState, setLoadState] = useState<"loading" | "loaded" | "error">("loading");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");

  useEffect(() => {
    loadSettings();
  }, []);

  function loadSettings() {
    setLoadState("loading");
    setError("");
    void api.getAdminSettings()
      .then((response) => {
        const footer = { ...defaultFooterConfig, ...(response.data.footer ?? {}) };
        setConfig(footer);
        setSavedConfig(footer);
        setSocialLinks(response.data.social_links ?? []);
        setLoadState("loaded");
      })
      .catch((cause) => {
        setLoadState("error");
        setError(cause instanceof Error ? cause.message : "加载页脚配置失败");
      });
  }

  const dirty = useMemo(() => JSON.stringify(config) !== JSON.stringify(savedConfig), [config, savedConfig]);

  useEffect(() => {
    function warnBeforeLeaving(event: BeforeUnloadEvent) {
      if (!dirty) return;
      event.preventDefault();
    }
    function confirmLinkNavigation(event: MouseEvent) {
      if (!dirty || event.defaultPrevented || event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
      const target = event.target as Element | null;
      const anchor = target?.closest("a");
      if (!anchor || anchor.target === "_blank" || anchor.href === window.location.href) return;
      if (!window.confirm("页脚还有未保存的修改，确定离开吗？")) event.preventDefault();
    }
    window.addEventListener("beforeunload", warnBeforeLeaving);
    document.addEventListener("click", confirmLinkNavigation, true);
    return () => {
      window.removeEventListener("beforeunload", warnBeforeLeaving);
      document.removeEventListener("click", confirmLinkNavigation, true);
    };
  }, [dirty]);

  function update(values: Partial<FooterConfig>) {
    setNotice("");
    setConfig((current) => ({ ...current, ...values }));
  }

  async function save(event: React.FormEvent) {
    event.preventDefault();
    setSaving(true); setError(""); setNotice("");
    try {
      const response = await api.updateSettings({ footer: config });
      const footer = { ...defaultFooterConfig, ...(response.data.footer ?? {}) };
      setConfig(footer);
      setSavedConfig(footer);
      setNotice("页脚设置已更新。");
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "保存失败");
    } finally {
      setSaving(false);
    }
  }

  if (loadState === "loading") return <div className="flex min-h-72 items-center justify-center"><Loader2 className="animate-spin text-primary" /><span className="sr-only">正在加载</span></div>;
  if (loadState === "error") return <div className="empty-state"><ShieldCheck size={28} /><h1>页脚配置加载失败</h1><p>{error || "请检查网络连接和后端服务。"}</p><button type="button" className="button-primary mt-5" onClick={loadSettings}>重新加载</button></div>;

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div><p className="eyebrow">Site footer</p><h1 className="mt-2 text-3xl font-bold">页脚与备案</h1><p className="mt-2 text-sm text-muted-foreground">统一管理全站底部的品牌、社交入口、版权和备案信息。</p></div>
        <div className="flex items-center gap-3">{dirty && <span className="text-xs font-medium text-amber-600 dark:text-amber-400" role="status">有未保存修改</span>}<a href="/" target="_blank" rel="noopener noreferrer" className="button-secondary"><ExternalLink size={16} />查看站点</a></div>
      </div>
      {error && <p className="rounded-xl bg-destructive/10 p-4 text-sm text-destructive" role="alert">{error}</p>}
      {notice && <p className="rounded-xl bg-primary/10 p-4 text-sm text-primary" role="status">{notice}</p>}
      <form onSubmit={save} className="surface-card space-y-6 p-5 md:p-6">
        <label className="flex items-center justify-between gap-4 rounded-xl border border-border p-4"><span><span className="block font-semibold">展示页脚</span><span className="mt-1 block text-xs text-muted-foreground">关闭后，所有公开页面都不再显示页脚。</span></span><input type="checkbox" checked={config.enabled} onChange={(event) => update({ enabled: event.target.checked })} className="h-5 w-5 accent-primary" /></label>

        <fieldset className="space-y-5 disabled:opacity-50" disabled={!config.enabled}>
          <label className="block space-y-2"><span className="label">品牌文字</span><input className="input" maxLength={80} value={config.brand} onChange={(event) => update({ brand: event.target.value })} placeholder="例如：SiYue Blog" /><span className="block text-xs text-muted-foreground">留空时不显示品牌文字。</span></label>
          <label className="block space-y-2"><span className="label">简介</span><textarea className="input min-h-28 resize-y" maxLength={300} value={config.description} onChange={(event) => update({ description: event.target.value })} placeholder="简短介绍这个站点" /><span className="block text-xs text-muted-foreground">最多 300 个字符，支持换行。</span></label>
          <label className="flex items-center justify-between gap-4 rounded-xl border border-border p-4"><span><span className="block font-semibold">展示社交账号</span><span className="mt-1 block text-xs text-muted-foreground">使用“联系页”中维护的社交账号和排列顺序。</span></span><input type="checkbox" checked={config.show_social_links} onChange={(event) => update({ show_social_links: event.target.checked })} className="h-5 w-5 accent-primary" /></label>
          <label className="block space-y-2"><span className="label">版权文案</span><input className="input" maxLength={200} value={config.copyright} onChange={(event) => update({ copyright: event.target.value })} placeholder="© {year} Your Name" /><span className="block text-xs text-muted-foreground"><code>{"{year}"}</code> 会自动替换为当前年份；留空时不展示版权栏。</span></label>

          <fieldset className="space-y-5 rounded-xl border border-border p-4 md:p-5">
            <legend className="px-2 text-sm font-semibold"><span className="inline-flex items-center gap-2"><ShieldCheck size={17} className="text-primary" />备案信息</span></legend>
            <p className="text-xs leading-5 text-muted-foreground">所有字段均可选，留空即不显示。链接请填写 HTTPS 地址，图标也可填写站内路径，例如 <code>/images/beian.png</code>。</p>
            <div className="grid gap-5 md:grid-cols-2">
              <label className="block space-y-2"><span className="label">备案图标</span><input className="input" maxLength={512} value={config.filing_icon_url} onChange={(event) => update({ filing_icon_url: event.target.value })} placeholder="/images/beian.png" /><span className="block text-xs text-muted-foreground">显示在“粤公网安备号”左侧。</span></label>
              <label className="block space-y-2"><span className="label">备案附加信息</span><textarea rows={2} className="input resize-y" maxLength={300} value={config.filing_information} onChange={(event) => update({ filing_information: event.target.value })} placeholder="例如：网站备案信息" /></label>
              <label className="block space-y-2"><span className="label">ICP备号</span><input className="input" maxLength={100} value={config.icp_number} onChange={(event) => update({ icp_number: event.target.value })} placeholder="粤ICP备12345678号" /></label>
              <label className="block space-y-2"><span className="label">ICP备案链接</span><input type="url" className="input" maxLength={512} value={config.icp_url} onChange={(event) => update({ icp_url: event.target.value })} placeholder="https://beian.miit.gov.cn/" /></label>
              <label className="block space-y-2"><span className="label">粤公网安备号</span><input className="input" maxLength={100} value={config.public_security_number} onChange={(event) => update({ public_security_number: event.target.value })} placeholder="粤公网安备 44000000000000号" /></label>
              <label className="block space-y-2"><span className="label">公安备案链接</span><input type="url" className="input" maxLength={512} value={config.public_security_url} onChange={(event) => update({ public_security_url: event.target.value })} placeholder="https://www.beian.gov.cn/..." /></label>
            </div>
          </fieldset>
        </fieldset>

        <div className="flex flex-wrap items-center justify-end gap-3"><button type="button" className="button-secondary" disabled={saving} onClick={() => { if (window.confirm("恢复为默认页脚内容？保存前仍可继续修改。")) { setNotice(""); setConfig(defaultFooterConfig); } }}><RotateCcw size={16} />恢复默认</button><button className="button-primary" disabled={saving || !dirty}>{saving ? <Loader2 size={16} className="animate-spin" /> : <Save size={16} />}{saving ? "保存中" : dirty ? "保存设置" : "已保存"}</button></div>
      </form>

      <section className="space-y-3" aria-labelledby="footer-preview-title">
        <div><h2 id="footer-preview-title" className="text-lg font-semibold">实时预览</h2><p className="mt-1 text-xs text-muted-foreground">预览仅用于确认内容和排列，实际宽度与背景跟随公开页面。</p></div>
        <div className="surface-card overflow-hidden"><Footer config={config} socialLinks={socialLinks} /></div>
      </section>
    </div>
  );
}
