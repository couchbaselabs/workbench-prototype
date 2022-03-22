export type CloudClusterHealth = {
  status: string;
  health?: string;
  bucket_stats?: any;
  node_stats?: any;
}

export type CloudCluster = {
  id: string;
  cloudId: string;
  name: string;
  nodes: number;
  projectId: string;
  services: string;
  tenantId: string;
  expanded?: boolean;
};

export type CloudClusterPage = {
  cursor: any;
  data: CloudCluster[];
}

export type NodeSummary = {
  node_uuid: string;
  version: string;
  host: string;
  status: string;
  cluster_membership: string;
  services: string[];

  expanded?: boolean;
};

export type Cluster = {
  uuid: string;
  name: string;
  alias?: string;
  enterprise: boolean;
  nodes_summary: NodeSummary[];
  heart_beat_issue?: number;
  last_update: string;
};

export type ClusterAddRequest = {
  host: string;
  user: string;
  password: string;
};

export type BucketSummary = {
  name: string;
  storage_backend: string;
  quota: number;
  quota_used: number;
  num_replicas: number;
  items: number;
  bucket_type: string;

  expanded?: boolean;
};

export type StatusResults = {
  cluster: string;
  node?: string;
  bucket?: string;
  log_file?: string;
  result: {
    name: string;
    remediation?: string;
    status: string;
    time: string;
    version: number;
    value?: any;
  };

  expanded?: boolean;
};

export type ClusterStatusResults = {
  nodes_summary: NodeSummary[];
  buckets_summary: BucketSummary[];
  last_update: string;
  name: string;
  uuid: string;
};

export type Alert = {
  message: string;
  type: string;
  timeout: number;
  timeoutFn?: any;
};

export type Address = {
  ip: string;
  port: number;
}

export function heartBeatMessage(issue: number | undefined): string {
  switch (issue) {
    case undefined:
      return '';
    case 0:
      return 'none';
    case 1:
      return (
        'The user and password given are no longer valid or do not have the required permissions.' +
        'Please update them.'
      );
    case 2:
      return (
        'Could not establish connection with the cluster during last heartbeat. Please check the cluster is ' +
        'still online.'
      );
    case 3:
      return 'The given host no longer points to the same cluster. The cluster UUID has changed.';
    default:
      return `Unknown heart beat issue - ${issue}.`;
  }
}
