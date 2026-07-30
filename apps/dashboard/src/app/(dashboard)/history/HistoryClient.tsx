"use client";

import React from "react";
import { useDashboardHistory } from "@/services/dashboard";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { Card } from "@/components/ui/Card";
import { RefreshCw } from "lucide-react";

type Period = "1h" | "24h" | "7d";

const PERIODS: { value: Period; label: string }[] = [
  { value: "1h", label: "Last Hour" },
  { value: "24h", label: "Last 24 Hours" },
  { value: "7d", label: "Last 7 Days" },
];

function formatTimestamp(value: string) {
  const date = new Date(value);
  return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

function SkeletonChart() {
  return (
    <div className="h-64 w-full animate-pulse rounded-lg bg-white/5" />
  );
}

function Timeline({ events }: { events: any[] }) {
  if (events.length === 0) {
    return (
      <Card title="Infrastructure Timeline" description="Recent operational events.">
        <p className="text-sm text-slate-400">No events recorded in this period.</p>
      </Card>
    );
  }

  return (
    <Card title="Infrastructure Timeline" description="Recent operational events.">
      <div className="space-y-4">
        {events.map((event, index) => (
          <div key={index} className="flex gap-4">
            <div className="text-xs text-slate-400 w-16 shrink-0">
              {formatTimestamp(event.timestamp)}
            </div>
            <div className="flex-1 border-l border-white/10 pl-4">
              <p className="text-sm font-medium text-white">{event.title}</p>
              {event.description && (
                <p className="text-xs text-slate-400">{event.description}</p>
              )}
            </div>
          </div>
        ))}
      </div>
    </Card>
  );
}

export default function HistoryClient() {
  const [period, setPeriod] = React.useState<Period>("1h");
  const {
    data,
    isLoading,
    isError,
    refetch,
    isFetching,
  } = useDashboardHistory(period);

  const chartData = React.useMemo(() => {
    if (!data?.points) return [];
    return data.points.map((point) => ({
      time: formatTimestamp(point.timestamp),
      cpu: point.cpuUsage,
      memory: point.memoryUsage,
      disk: point.diskUsage,
    }));
  }, [data]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold text-white">History</h1>
          <p className="text-sm text-slate-400">
            Operational history for your infrastructure.
          </p>
        </div>
        <div className="flex gap-2">
          {PERIODS.map((p) => (
            <button
              key={p.value}
              onClick={() => setPeriod(p.value)}
              className={`rounded-lg px-3 py-1.5 text-sm font-medium transition-colors ${
                period === p.value
                  ? "bg-white/10 text-white"
                  : "text-slate-400 hover:text-white"
              }`}
            >
              {p.label}
            </button>
          ))}
        </div>
      </div>

      {isLoading ? (
        <div className="space-y-4">
          <SkeletonChart />
          <SkeletonChart />
          <SkeletonChart />
        </div>
      ) : isError || !data ? (
        <Card title="Unable to load history" description="Something went wrong while fetching operational history.">
          <div className="flex items-center gap-3">
            <button
              onClick={() => refetch()}
              className="flex items-center gap-2 rounded-lg bg-white/10 px-3 py-1.5 text-sm font-medium text-white hover:bg-white/20"
            >
              <RefreshCw className="h-4 w-4" />
              Retry
            </button>
          </div>
        </Card>
      ) : chartData.length === 0 ? (
        <Card title="No historical data available" description="History will appear here once metrics have been collected.">
          <p className="text-sm text-slate-400">
            Start collecting metrics or adjust the time window.
          </p>
        </Card>
      ) : (
        <div className="space-y-6">
          <Card title="CPU Usage" description="Processor utilization over time.">
            <div className="h-64 w-full">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={chartData}>
                  <XAxis
                    dataKey="time"
                    stroke="#475569"
                    tick={{ fill: "#94a3b8", fontSize: 12 }}
                  />
                  <YAxis
                    stroke="#475569"
                    tick={{ fill: "#94a3b8", fontSize: 12 }}
                    domain={[0, 100]}
                  />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "#0f172a",
                      borderColor: "#1e293b",
                      borderRadius: "0.5rem",
                    }}
                    itemStyle={{ color: "#f8fafc" }}
                  />
                  <Line
                    type="monotone"
                    dataKey="cpu"
                    stroke="#38bdf8"
                    strokeWidth={2}
                    dot={false}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </Card>

          <Card title="Memory Usage" description="RAM utilization over time.">
            <div className="h-64 w-full">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={chartData}>
                  <XAxis
                    dataKey="time"
                    stroke="#475569"
                    tick={{ fill: "#94a3b8", fontSize: 12 }}
                  />
                  <YAxis
                    stroke="#475569"
                    tick={{ fill: "#94a3b8", fontSize: 12 }}
                    domain={[0, 100]}
                  />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "#0f172a",
                      borderColor: "#1e293b",
                      borderRadius: "0.5rem",
                    }}
                    itemStyle={{ color: "#f8fafc" }}
                  />
                  <Line
                    type="monotone"
                    dataKey="memory"
                    stroke="#a78bfa"
                    strokeWidth={2}
                    dot={false}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </Card>

          <Card title="Disk Usage" description="Disk utilization over time.">
            <div className="h-64 w-full">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={chartData}>
                  <XAxis
                    dataKey="time"
                    stroke="#475569"
                    tick={{ fill: "#94a3b8", fontSize: 12 }}
                  />
                  <YAxis
                    stroke="#475569"
                    tick={{ fill: "#94a3b8", fontSize: 12 }}
                    domain={[0, 100]}
                  />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: "#0f172a",
                      borderColor: "#1e293b",
                      borderRadius: "0.5rem",
                    }}
                    itemStyle={{ color: "#f8fafc" }}
                  />
                  <Line
                    type="monotone"
                    dataKey="disk"
                    stroke="#34d399"
                    strokeWidth={2}
                    dot={false}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </Card>

          <Timeline events={[]} />
        </div>
      )}
    </div>
  );
}
