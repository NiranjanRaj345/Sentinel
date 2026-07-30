import { DashboardShell } from "@/components/layout/DashboardShell";
import HistoryClient from "./HistoryClient";

export default function HistoryRoute() {
  return (
    <DashboardShell>
      <HistoryClient />
    </DashboardShell>
  );
}
