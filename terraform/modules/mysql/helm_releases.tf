resource "helm_release" "mysql" {
  name             = "mysql"
  namespace        = "mysql"
  repository       = "https://charts.bitnami.com/bitnami"
  chart            = "mysql"
  create_namespace = true
  wait             = true
  version          = "9.6.0"

  values = [
    <<-EOF
    auth:
      database: ${var.database_name}
      username: ${var.user_name}
      password: ${random_password.user_password.result}
      rootPassword: ${random_password.root_password.result}
    primary:
      service:
        type: NodePort
        nodePorts:
          mysql: ${var.node_port}
      resources: null
    EOF
  ]
}
