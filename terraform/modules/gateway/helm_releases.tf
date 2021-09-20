resource "helm_release" "gateway" {
  name             = "gateway"
  namespace        = "gateway"
  chart            = "../gateway/helm"
  create_namespace = true
  wait             = true
  version          = "1.0.0"

  values = [<<-EOF
    replicas: 1
    horizontalPodAutoscaler:
      minReplicas: 1
      maxReplicas: 1
      targetCPUUtilizationPercentage: 50
    ingressGateway:
      port: ${var.ingress_gateway_port}
    services:
%{for name, node_port in var.node_ports~}
      ${name}:
        nodePort: ${node_port}
%{endfor~}
    EOF
  ]
}