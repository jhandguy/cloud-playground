resource "helm_release" "gateway" {
  name      = "gateway"
  namespace = kubernetes_namespace.gateway.metadata.0.name
  chart     = "../gateway/helm"
  wait      = true

  values = [<<-EOF
    replicas: 1
    configMap: ${kubernetes_config_map.gateway.metadata.0.name}
    secret: ${kubernetes_secret.gateway.metadata.0.name}
    image:
      secret: ${kubernetes_secret.gateway_image.metadata.0.name}
      registry: ${var.image_registry}
    ingressGateway:
      port: ${var.ingress_gateway_port}
    deployments:
%{for name, node_port in var.node_ports~}
      ${name}:
        nodePort: ${node_port}
%{endfor~}
    EOF
  ]
}