resource "helm_release" "dynamo" {
  name      = "dynamo"
  namespace = kubernetes_namespace.dynamo.metadata.0.name
  chart     = "../dynamo/helm"
  wait      = true

  values = [<<-EOF
    replicas: 1
    nodePort: ${var.node_port}
    configMap: ${kubernetes_config_map.dynamo.metadata.0.name}
    secret: ${kubernetes_secret.dynamo.metadata.0.name}
    image:
      secret: ${kubernetes_secret.dynamo_image.metadata.0.name}
      registry: ${var.image_registry}
    EOF
  ]
}