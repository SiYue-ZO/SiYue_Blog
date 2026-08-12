/* eslint-disable @next/next/no-img-element */

import Link from "next/link";
import { AtSign, BriefcaseBusiness, Code2, Globe2, Mail, MessageCircle, Play, Radio } from "lucide-react";

import { defaultFooterConfig, defaultSocialLinks } from "@/lib/constants";
import type { FooterConfig, SocialLink, SocialPlatform } from "@/types";

const socialIcons = {
  github: Code2,
  email: Mail,
  bilibili: Play,
  twitter: MessageCircle,
  linkedin: BriefcaseBusiness,
  youtube: Play,
  weibo: Radio,
  website: Globe2,
  other: AtSign,
} satisfies Record<SocialPlatform, typeof Code2>;

export function Footer({ config = defaultFooterConfig, socialLinks = defaultSocialLinks }: { config?: FooterConfig | null; socialLinks?: SocialLink[] | null }) {
  // Merge rather than replace so JSON saved before new footer fields were
  // introduced keeps working after deployment.
  const footer = { ...defaultFooterConfig, ...(config ?? {}) };
  const links = socialLinks ?? defaultSocialLinks;
  if (!footer.enabled) return null;

  const hasFilingInformation = Boolean(
    footer.filing_icon_url || footer.icp_number || footer.public_security_number || footer.filing_information,
  );

  return (
    <footer className="border-t border-border/40 bg-background">
      <div className="mx-auto w-full max-w-7xl px-4 py-8 sm:px-6 lg:px-8 md:py-12">
        <div className="flex flex-col items-center justify-between gap-8 md:flex-row">
          <div className="space-y-2 text-center md:text-left">
            {footer.brand && <Link href="/" className="inline-flex items-center space-x-2"><span className="text-xl font-bold">{footer.brand}</span></Link>}
            {footer.description && <p className="max-w-xs whitespace-pre-line text-sm text-muted-foreground">{footer.description}</p>}
          </div>

          {footer.show_social_links && Boolean(links.length) && <div className="flex flex-wrap items-center justify-center gap-6">
            {links.map((link, index) => {
              const Icon = socialIcons[link.platform] ?? AtSign;
              return <a key={`${link.platform}-${link.href}-${index}`} href={link.href} target={link.href.startsWith("mailto:") ? undefined : "_blank"} rel={link.href.startsWith("mailto:") ? undefined : "noopener noreferrer"} className="text-muted-foreground transition-colors hover:text-primary" aria-label={link.label} title={link.label}><Icon size={20} /></a>;
            })}
          </div>}
        </div>

        {(footer.copyright || hasFilingInformation) && <div className="mt-8 space-y-3 border-t border-border/40 pt-8 text-center text-sm text-muted-foreground">
          {footer.copyright && <p>{footer.copyright.replaceAll("{year}", String(new Date().getFullYear()))}</p>}
          {hasFilingInformation && <div className="flex flex-wrap items-center justify-center gap-x-4 gap-y-2 text-xs sm:text-sm">
            {footer.icp_number && <FilingItem href={footer.icp_url}>{footer.icp_number}</FilingItem>}
            {footer.public_security_number && <FilingItem href={footer.public_security_url}>
              {footer.filing_icon_url && <img src={footer.filing_icon_url} alt="" width={18} height={18} referrerPolicy="no-referrer" className="h-[18px] w-[18px] object-contain" />}
              <span>{footer.public_security_number}</span>
            </FilingItem>}
            {!footer.public_security_number && footer.filing_icon_url && <img src={footer.filing_icon_url} alt="备案图标" width={18} height={18} referrerPolicy="no-referrer" className="h-[18px] w-[18px] object-contain" />}
            {footer.filing_information && <span className="whitespace-pre-line">{footer.filing_information}</span>}
          </div>}
        </div>}
      </div>
    </footer>
  );
}

function FilingItem({ href, children }: { href?: string; children: React.ReactNode }) {
  const className = "inline-flex items-center gap-1.5 transition-colors hover:text-primary";
  if (!href) return <span className="inline-flex items-center gap-1.5">{children}</span>;
  return <a href={href} target="_blank" rel="noopener noreferrer" className={className}>{children}</a>;
}
