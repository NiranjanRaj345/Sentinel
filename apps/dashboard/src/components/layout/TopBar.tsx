"use client";

import { Shield } from "lucide-react";

export function TopBar() {
  return (
    <header className="flex items-center justify-between border-b border-white/10 bg-slate-950/80 px-6 py-3 backdrop-blur">
      <div className="flex items-center gap-3">
        <Shield className="h-5 w-5 text-emerald-400" />
        <div>
          <p className="text-sm font-semibold text-white">Mission Dashboard</p>
          <p className="text-xs text-slate-400">
            Sentinel Node Agent — Personal Infrastructure
          </p>
        </div>
      </div>
    </header>
  );
}
