import { AlertCircle, AlertTriangle, Info, Lightbulb } from "lucide-react";
import { ReactNode } from "react";

type CalloutType = "info" | "warning" | "danger" | "tip";

interface CalloutProps {
  type?: CalloutType;
  title?: string;
  children: ReactNode;
}

const calloutConfig: Record<
  CalloutType,
  { icon: ReactNode; className: string; defaultTitle: string }
> = {
  info: {
    icon: <Info size={18} />,
    className:
      "border-blue-500/30 bg-blue-500/10 text-blue-400",
    defaultTitle: "提示",
  },
  warning: {
    icon: <AlertTriangle size={18} />,
    className:
      "border-yellow-500/30 bg-yellow-500/10 text-yellow-400",
    defaultTitle: "警告",
  },
  danger: {
    icon: <AlertCircle size={18} />,
    className:
      "border-red-500/30 bg-red-500/10 text-red-400",
    defaultTitle: "危险",
  },
  tip: {
    icon: <Lightbulb size={18} />,
    className:
      "border-green-500/30 bg-green-500/10 text-green-400",
    defaultTitle: "技巧",
  },
};

export function Callout({ type = "info", title, children }: CalloutProps) {
  const config = calloutConfig[type];

  return (
    <div
      className={`my-6 rounded-lg border-l-4 p-4 ${config.className}`}
    >
      <div className="flex items-center gap-2 font-semibold mb-2">
        {config.icon}
        <span>{title || config.defaultTitle}</span>
      </div>
      <div className="text-sm leading-relaxed opacity-90">{children}</div>
    </div>
  );
}
