resource "helm_release" "sql" {
  name             = "sql-${var.feature}"
  namespace        = "sql-${var.feature}"
  chart            = "../../../sql/helm"
  create_namespace = true
  wait             = true
  version          = "1.0.0"

  values = [
    <<-EOF
    replicas: ${var.replicas}
    nodePorts:
      http: ${var.node_ports.0}
      metrics: ${var.node_ports.1}
    ingress:
      host: ${var.ingress_host}
%{if var.rate_limit_requests > 0}
      rateLimitRequests: ${var.rate_limit_requests}
%{endif}
    prometheus:
      enabled: ${var.prometheus_enabled}
      groupName: ${var.feature == "mysql" ? "MySQL" : title(var.feature)}
    horizontalPodAutoscaler:
      minReplicas: ${var.replicas}
      maxReplicas: ${var.replicas * 2}
    EOF
  ]

  dynamic "set_sensitive" {
    for_each = var.secrets

    content {
      name  = "secrets.${set_sensitive.key}"
      value = base64encode(set_sensitive.value)
    }
  }
}
