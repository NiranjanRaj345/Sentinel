"use client";

import { clsx } from "clsx";

type CardProps = {
  title?: string;
  description?: string;
  children: React.ReactNode;
  className?: string;
};

export function Card({ title, description, children, className }: CardProps) {
  return (
    <div
      className={clsx(
        "rounded-xl border border-white/10 bg-slate-900/60 p-5 shadow-sm",
        className
      )}
    >
      {(title || description) && (
        <div className="mb-4 space-y-1">
          {title && (
            <p className="text-sm font-semibold text-white">{title}</p>
          )}
          {description && (
            <p className="text-xs text-slate-400">{description}</p>
          )}
        </div>
      )}
      {children}
    </div>
  );
}
