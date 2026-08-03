"use client";

import React from "react";
import { DashboardShell } from "@/components/layout/DashboardShell";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

function Section({ title, description, children }: { title: string; description: string; children: React.ReactNode }) {
  return (
    <div className="rounded-lg border border-white/10 bg-white/5 p-4">
      <h2 className="text-sm font-semibold text-white">{title}</h2>
      <p className="mt-1 text-xs text-slate-400">{description}</p>
      <div className="mt-4 space-y-3">{children}</div>
    </div>
  );
}

function Field({ label, value, onChange, placeholder, type = "text" }: { label: string; value: string; onChange: (value: string) => void; placeholder?: string; type?: string }) {
  return (
    <div>
      <label className="text-xs font-medium text-slate-300">{label}</label>
      <input
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className="mt-1 w-full rounded-lg border border-white/10 bg-slate-950 px-3 py-2 text-sm text-white placeholder:text-slate-500 focus:border-white/20 focus:outline-none"
      />
    </div>
  );
}

export default function SettingsRoute() {
  const queryClient = useQueryClient();
  const [token, setToken] = React.useState(
    typeof window !== "undefined" ? localStorage.getItem("sentinel_token") || "" : ""
  );
  const [telegramBotToken, setTelegramBotToken] = React.useState("");
  const [telegramChatID, setTelegramChatID] = React.useState("");

  const { data: capabilities } = useQuery({
    queryKey: ["dashboard", "capabilities"],
    queryFn: async () => {
      const base = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080";
      const res = await fetch(`${base}/dashboard/capabilities`, { cache: "no-store", headers: { Accept: "application/json" } });
      if (!res.ok) throw new Error("failed");
      return res.json();
    },
    staleTime: 60_000,
  });

  const telegramCapability = capabilities?.capabilities?.find((c: any) => c.capability === "notifications" || c.capability === "telegram");

  const testMutation = useMutation({
    mutationFn: async () => {
      const base = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080";
      const headers: HeadersInit = { "Content-Type": "application/json", Accept: "application/json" };
      const storedToken = typeof window !== "undefined" ? localStorage.getItem("sentinel_token") : "";
      if (storedToken) headers.Authorization = `Bearer ${storedToken}`;

      const res = await fetch(`${base}/notifications/test`, {
        method: "POST",
        headers,
        body: JSON.stringify({ provider: "telegram" }),
      });
      if (!res.ok) throw new Error("test failed");
      return res.json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
    },
  });

  const handleSaveToken = () => {
    if (token.trim()) {
      localStorage.setItem("sentinel_token", token.trim());
    } else {
      localStorage.removeItem("sentinel_token");
    }
    alert("Token saved.");
  };

  const handleSaveTelegram = () => {
    alert("Telegram settings saved. Update config.yaml to persist bot_token and chat_id.");
  };

  return (
    <DashboardShell>
      <div className="space-y-6">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">Settings</h1>
          <p className="text-sm text-slate-400">
            API token, notification channels, and integration settings.
          </p>
        </div>

        <Section title="API Token" description="Optional. Set a Bearer token to authenticate dashboard requests.">
          <Field label="Token" value={token} onChange={setToken} placeholder="sentinel_..." />
          <button
            onClick={handleSaveToken}
            className="rounded-lg bg-white/10 px-4 py-2 text-sm font-medium text-white hover:bg-white/20"
          >
            Save Token
          </button>
        </Section>

        <Section title="Telegram" description="Send notifications to a Telegram chat.">
          <Field label="Bot Token" value={telegramBotToken} onChange={setTelegramBotToken} placeholder="123456:ABC-DEF..." />
          <Field label="Chat ID" value={telegramChatID} onChange={setTelegramChatID} placeholder="-1001234567890" />
          <div className="flex items-center gap-3">
            <button
              onClick={handleSaveTelegram}
              className="rounded-lg bg-white/10 px-4 py-2 text-sm font-medium text-white hover:bg-white/20"
            >
              Save Telegram Settings
            </button>
            <button
              onClick={() => testMutation.mutate()}
              disabled={testMutation.isPending}
              className="rounded-lg bg-emerald-500/10 px-4 py-2 text-sm font-medium text-emerald-300 hover:bg-emerald-500/20 disabled:opacity-50"
            >
              {testMutation.isPending ? "Sending..." : "Send Test Notification"}
            </button>
          </div>
          {telegramCapability && (
            <p className="text-xs text-slate-400">
              Status: {telegramCapability.available ? "Available" : "Unavailable"} — {telegramCapability.state}
            </p>
          )}
          {testMutation.isError && (
            <p className="text-xs text-red-300">Test failed. Check Telegram configuration.</p>
          )}
          {testMutation.isSuccess && (
            <p className="text-xs text-emerald-300">Test notification sent.</p>
          )}
        </Section>
      </div>
    </DashboardShell>
  );
}
