import { DashboardShell } from "@/components/layout/DashboardShell";

export default function AlertsRoute() {
  return (
    <DashboardShell>
      <div className="space-y-4">
        <h1 className="text-2xl font-semibold text-white">Alerts</h1>
        <p className="text-sm text-slate-400">
          Active, resolved, and silenced alerts will be managed here.
        </p>
      </div>
    </DashboardShell>
  );
}
