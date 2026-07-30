import { DashboardShell } from "@/components/layout/DashboardShell";

export default function SettingsRoute() {
  return (
    <DashboardShell>
      <div className="space-y-4">
        <h1 className="text-2xl font-semibold text-white">Settings</h1>
        <p className="text-sm text-slate-400">
          Node configuration, alert rules, and theme settings will live here.
        </p>
      </div>
    </DashboardShell>
  );
}
