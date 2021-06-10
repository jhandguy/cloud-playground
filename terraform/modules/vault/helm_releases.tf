resource "helm_release" "vault" {
  name             = "vault"
  namespace        = "vault"
  repository       = "https://helm.releases.hashicorp.com"
  chart            = "vault"
  create_namespace = true
  wait             = true
  version          = "0.11.0" // TODO: version 0.12.0 is breaking

  values = [<<-EOF
    injector:
      metrics:
        enabled: true
    server:
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

  provisioner "local-exec" {
    command = <<-EOF
      kubectl wait --for=condition=ready --timeout=60s pod/vault-0 -n vault
      kubectl exec vault-0 -n vault -- vault auth enable kubernetes
      kubectl exec vault-0 -n vault -- sh -c 'vault write auth/kubernetes/config \
        issuer="https://kubernetes.default.svc.cluster.local" \
        token_reviewer_jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
        kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443" \
        kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt'
    EOF
  }
}