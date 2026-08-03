"use client";

import React from "react";
import { useToastStore } from "@/stores/toast";
import { X } from "lucide-react";

const styles = {
  success: "border-emerald-500/30 bg-emerald-500/10 text-emerald-200",
  error: "border-red-500/30 bg-red-500/10 text-red-200",
  info: "border-white/10 bg-white/5 text-slate-200",
};

export function Toast() {
  const { toasts, removeToast } = useToastStore();

  if (toasts.length === 0) {
    return null;
  }

  return (
    <div className="fixed right-4 top-4 z-50 flex max-w-sm flex-col gap-2">
      {toasts.map((toast) => (
        <div
          key={toast.id}
          className={`flex items-center justify-between rounded-lg border px-4 py-3 text-sm shadow-xl ${styles[toast.type]}`}
        >
          <span>{toast.message}</span>
          <button
            onClick={() => removeToast(toast.id)}
            className="ml-3 rounded-md p-1 opacity-70 hover:opacity-100"
          >
            <X className="h-3 w-3" />
          </button>
        </div>
      ))}
    </div>
  );
}
