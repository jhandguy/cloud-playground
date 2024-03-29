resource "helm_release" "mimir" {
  name             = "mimir"
  namespace        = "mimir"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "mimir-distributed"
  create_namespace = true
  wait             = true
  version          = "5.1.3"

  values = [
    <<-EOF
    alertmanager:
      enabled: false
    compactor:
      resources:
        limits:
          cpu: 0
          memory: 0
        requests:
          cpu: 0
          memory: 0
    distributor:
      resources:
        limits:
          cpu: 0
          memory: 0
        requests:
          cpu: 0
          memory: 0
    gateway:
      enabledNonEnterprise: true
      service:
        legacyPort: null
    ingester:
      replicas: 1
      resources:
        limits:
          cpu: 0
          memory: 0
        requests:
          cpu: 0
          memory: 0
      zoneAwareReplication:
        enabled: false
    mimir:
      structuredConfig:
%{if var.localstack_enabled}
        common:
          storage:
            backend: s3
            s3:
              endpoint: ${var.aws_s3_cluster_endpoint}
              region: ${var.aws_region}
              secret_access_key: ${var.aws_secret_access_key}
              access_key_id: ${var.aws_access_key_id}
              insecure: true
        blocks_storage:
          s3:
            bucket_name: ${var.aws_s3_bucket}
%{endif}
        ingester:
          ring:
            replication_factor: 1
    minio:
      enabled: ${!var.localstack_enabled}
      resources:
        limits:
          cpu: 0
          memory: 0
        requests:
          cpu: 0
          memory: 0
    nginx:
      enabled: false
    overrides_exporter:
      enabled: false
    querier:
      replicas: 1
      resources:
        limits:
          cpu: 0
          memory: 0
        requests:
          cpu: 0
          memory: 0
    query_frontend:
      replicas: 1
      resources:
        limits:
          cpu: 0
          memory: 0
        requests:
          cpu: 0
          memory: 0
    query_scheduler:
      enabled: false
    rollout_operator:
      enabled: false
    ruler:
      enabled: false
    serviceMonitor:
      enabled: ${var.prometheus_enabled}
      labels:
        release: prometheus
    store_gateway:
      zoneAwareReplication:
        enabled: false
      resources:
        limits:
          cpu: 0
          memory: 0
        requests:
          cpu: 0
          memory: 0
    EOF
  ]
}
