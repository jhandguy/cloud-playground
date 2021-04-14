resource "helm_release" "s3" {
  name      = "s3"
  namespace = kubernetes_namespace.s3.metadata.0.name
  chart     = "../s3/helm"
  wait      = true

  values = [<<-EOF
    replicas: 1
    nodePort: ${var.node_port}
    configMap: ${kubernetes_config_map.s3.metadata.0.name}
    secret: ${kubernetes_secret.s3.metadata.0.name}
    image:
      secret: ${kubernetes_secret.s3_image.metadata.0.name}
      registry: ${var.image_registry}
      repository: ${var.s3_image_repository}
      tag: ${var.s3_image_tag}
    EOF
  ]
}