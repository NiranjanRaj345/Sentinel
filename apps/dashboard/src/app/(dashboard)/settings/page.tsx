"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";

export default function SettingsRoute() {
  const [token, setToken] = React.useState(
    typeof window !== "undefined" ? localStorage.getItem("sentinel_token") || "" : ""
  );

  const handleSave = () => {
    if (token.trim()) {
      localStorage.setItem("sentinel_token", token.trim());
    } else {
      localStorage.removeItem("sentinel_token");
    }
  };

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">Settings</h1>
          <p className="text-sm text-slate-400">
            Node configuration, alert rules, and theme settings will live here.
          </p>
        </div>

        <div className="rounded-lg border border-white/10 bg-white/5 p-4">
          <h2 className="text-sm font-semibold text-white">API Token</h2>
          <p className="mt-1 text-xs text-slate-400">
            Optional. Set a Bearer token to authenticate dashboard requests.
          </p>
          <input
            type="text"
            value={token}
            onChange={(e) => setToken(e.target.value)}
            placeholder="sentinel_..."
            className="mt-3 w-full rounded-lg border border-white/10 bg-slate-950 px-3 py-2 text-sm text-white placeholder:text-slate-500 focus:border-white/20 focus:outline-none"
          />
          <button
            onClick={handleSave}
            className="mt-3 rounded-lg bg-white/10 px-4 py-2 text-sm font-medium text-white hover:bg-white/20"
          >
            Save Token
          </button>
        </div>
      </div>
    </DashboardShell>
  );
}
