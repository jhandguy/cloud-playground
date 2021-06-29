resource "helm_release" "gateway" {
  name             = "gateway"
  namespace        = "gateway"
  chart            = "../gateway/helm"
  create_namespace = true
  wait             = true

  values = [<<-EOF
    replicas: 1
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