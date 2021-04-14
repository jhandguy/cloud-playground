resource "helm_release" "gateway" {
  name      = "gateway"
  namespace = kubernetes_namespace.gateway.metadata.0.name
  chart     = "../gateway/helm"
  wait      = true

  values = [<<-EOF
    replicas: 1
    nodePort: ${var.node_port}
    configMap: ${kubernetes_config_map.gateway.metadata.0.name}
    secret: ${kubernetes_secret.gateway.metadata.0.name}
    image:
      secret: ${kubernetes_secret.gateway_image.metadata.0.name}
      registry: ${var.image_registry}
      repository: ${var.gateway_image_repository}
      tag: ${var.gateway_image_tag}
    ingressGateway:
      port: ${var.ingress_gateway_port}
    EOF
  ]
}