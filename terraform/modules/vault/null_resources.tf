resource "null_resource" "auth" {
  depends_on = [helm_release.vault]

  provisioner "local-exec" {
    command = <<-EOF
      kubectl exec vault-0 -n vault -- vault auth enable kubernetes
      kubectl exec vault-0 -n vault -- sh -c 'vault write auth/kubernetes/config \
        issuer="https://kubernetes.default.svc.cluster.local" \
        token_reviewer_jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
        kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443" \
        kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt'
    EOF
  }
}

resource "null_resource" "secrets" {
  depends_on = [null_resource.auth]
  for_each   = var.secrets

  provisioner "local-exec" {
    command = <<-EOF
      kubectl exec vault-0 -n vault -- vault kv put secret/${each.key}%{for key, value in each.value} ${key}="${value}"%{endfor}
      kubectl exec vault-0 -n vault -- sh -c 'vault policy write ${each.key} - <<EOF
      path "secret/data/${each.key}" {
        capabilities = ["read"]
      }
      EOF'
      kubectl exec vault-0 -n vault -- vault write auth/kubernetes/role/${each.key} \
        bound_service_account_names=${each.key} \
        bound_service_account_namespaces=${each.key} \
        policies=${each.key} \
        ttl=24h
    EOF
  }
}