export type DashboardOverview = {
  nodeName: string;
  version: string;
  status: DashboardStatus;
  uptimeMs: number;
  snapshot: MetricsSnapshot;
  activeAlerts: AlertSummary[];
  lastCollection: string;
};

export type DashboardStatus = "healthy" | "warning" | "critical" | "offline";

export type MetricsSnapshot = {
  metadata: {
    timestamp: string;
    collectionDurationMs: number;
    agent: {
      name: string;
      version: string;
      platform: string;
      architecture: string;
      goVersion: string;
    };
  };
  cpu: {
    usage_percent: number;
  };
  memory: {
    total_bytes: number;
    used_bytes: number;
    usage_percent: number;
  };
  disks: Array<{
    device: string;
    mountpoint: string;
    filesystem: string;
    total_bytes: number;
    used_bytes: number;
    free_bytes: number;
    usage_percent: number;
  }>;
  network: {
    hostname: string;
    interfaces: Array<{
      name: string;
      mac: string;
      addresses: string[];
    }>;
    io: {
      bytes_sent: number;
      bytes_received: number;
      packets_sent: number;
      packets_received: number;
    };
  };
  processes: Array<unknown>;
};

export type AlertSummary = {
  ruleId: string;
  ruleName: string;
  metric: string;
  value: number;
  threshold: number;
  severity: "info" | "warning" | "critical";
  triggeredAt: string;
};

export type HistoryPoint = {
  timestamp: string;
  cpuUsage: number;
  memoryUsage: number;
  diskUsage: number;
};

export type HistoryResponse = {
  period: string;
  points: HistoryPoint[];
};

export type CapabilityStatus = {
  capability: string;
  available: boolean;
  state: string;
  details?: string;
};

export type CapabilitiesResponse = {
  capabilities: CapabilityStatus[];
};
