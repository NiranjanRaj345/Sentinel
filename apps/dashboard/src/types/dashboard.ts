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

export type OperationResult = {
  action: string;
  success: boolean;
  startedAt: string;
  finishedAt: string;
  message: string;
};

export type ActivityEvent = {
  id: string;
  type: string;
  severity: string;
  source: string;
  title: string;
  message: string;
  metadata?: Record<string, unknown>;
  createdAt: string;
};

export type EventsResponse = {
  events: ActivityEvent[];
};

export type RuleCondition = {
  field: string;
  operator: string;
  value: string;
};

export type Rule = {
  id: string;
  name: string;
  enabled: boolean;
  trigger: string;
  conditions: RuleCondition[];
  actions: string[];
};

export type RulesResponse = {
  rules: Rule[];
};

export type ResourceType = "remote_desktop" | "vpn" | "container_runtime" | "media_server" | "database" | "monitoring" | "application";

export type ResourceHealth = "healthy" | "degraded" | "unavailable" | "unknown";

export type Resource = {
  name: string;
  type: ResourceType;
  health: ResourceHealth;
  status: string;
  message?: string;
};

export type ResourcesResponse = {
  resources: Resource[];
};

export type ServiceItem = {
  name: string;
  status: string;
  action?: string;
  message?: string;
};

export type ServicesResponse = {
  services: ServiceItem[];
};

export type AutomationExecution = {
  id: string;
  ruleId: string;
  ruleName: string;
  action: string;
  success: boolean;
  message: string;
  createdAt: string;
};

export type AutomationExecutionsResponse = {
  executions: AutomationExecution[];
};

export type GuardianStatus = {
  status: "online" | "offline";
  firmware: string;
  uptime: number;
  powerButton: boolean;
  resetButton: boolean;
  powerLed: boolean;
  lastSeen: string;
};

export type GuardianActionResponse = {
  result: string;
};
