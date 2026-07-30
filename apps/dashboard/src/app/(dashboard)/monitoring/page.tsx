import { DashboardShell } from "@/components/layout/DashboardShell";

export default function MonitoringRoute() {
  return (
    <DashboardShell>
      <div className="space-y-4">
        <h1 className="text-2xl font-semibold text-white">Monitoring</h1>
        <p className="text-sm text-slate-400">
          Live CPU, memory, disk, network, and process telemetry arriving in a
          future sprint.
        </p>
      </div>
    </DashboardShell>
  );
}
