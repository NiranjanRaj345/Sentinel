"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { clsx } from "clsx";
import {
  LayoutDashboard,
  Activity,
  History,
  Bell,
  Settings,
  PanelLeftClose,
  PanelLeftOpen,
  ListTree,
  FileText,
} from "lucide-react";
import { useUIStore } from "@/stores/ui";

const navigation = [
  { name: "Overview", href: "/", icon: LayoutDashboard },
  { name: "Monitoring", href: "/monitoring", icon: Activity },
  { name: "History", href: "/history", icon: History },
  { name: "Activity", href: "/activity", icon: ListTree },
  { name: "Alerts", href: "/alerts", icon: Bell },
  { name: "Rules", href: "/rules", icon: FileText },
  { name: "Settings", href: "/settings", icon: Settings },
];

export function Sidebar() {
  const pathname = usePathname();
  const sidebarOpen = useUIStore((state) => state.sidebarOpen);
  const toggleSidebar = useUIStore((state) => state.toggleSidebar);

  return (
    <aside
      className={clsx(
        "flex flex-col border-r border-white/10 bg-slate-950 transition-all duration-200",
        sidebarOpen ? "w-64" : "w-16"
      )}
    >
      <div className="flex items-center justify-between px-4 py-5">
        {sidebarOpen ? (
          <span className="text-sm font-semibold tracking-wide text-white">
            SENTINEL
          </span>
        ) : (
          <span className="text-sm font-semibold tracking-wide text-white">
            S
          </span>
        )}
        <button
          onClick={toggleSidebar}
          className="rounded-md p-1 text-slate-400 hover:bg-white/10 hover:text-white"
        >
          {sidebarOpen ? (
            <PanelLeftClose className="h-4 w-4" />
          ) : (
            <PanelLeftOpen className="h-4 w-4" />
          )}
        </button>
      </div>

      <nav className="flex-1 space-y-1 px-2">
        {navigation.map((item) => {
          const isActive =
            item.href === "/"
              ? pathname === "/"
              : pathname.startsWith(item.href);

          return (
            <Link
              key={item.name}
              href={item.href}
              className={clsx(
                "flex items-center rounded-md px-3 py-2 text-sm font-medium transition-colors",
                isActive
                  ? "bg-white/10 text-white"
                  : "text-slate-400 hover:bg-white/5 hover:text-white"
              )}
            >
              <item.icon className="h-4 w-4 shrink-0" />
              {sidebarOpen && <span className="ml-3">{item.name}</span>}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}
