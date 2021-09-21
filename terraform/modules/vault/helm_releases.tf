resource "helm_release" "vault" {
  name             = "vault"
  namespace        = "vault"
  repository       = "https://helm.releases.hashicorp.com"
  chart            = "vault"
  create_namespace = true
  wait             = true
  version          = "0.16.0"

  values = [<<-EOF
    injector:
      enabled: false
    server:
      updateStrategyType: RollingUpdate
      dev:
        enabled: true
        devRootToken: "${random_password.root_token.result}"
      ha:
        enabled: true
        replicas: 1
    ui:
      enabled: true
      serviceType: NodePort
      serviceNodePort: ${var.node_port}
    csi:
      enabled: true
    EOF
  ]
}