resource "helm_release" "postgresql" {
  name             = "postgresql"
  namespace        = "postgresql"
  repository       = "https://charts.bitnami.com/bitnami"
  chart            = "postgresql"
  create_namespace = true
  wait             = true
  version          = "12.1.6"

  values = [
    <<-EOF
    auth:
      database: ${var.database_name}
      username: ${var.user_name}
      password: ${random_password.user_password.result}
      postgresPassword: ${random_password.postgres_password.result}
    primary:
      service:
        type: NodePort
        nodePorts:
          postgresql: ${var.node_port}
      resources: null
    EOF
  ]
}
