resource "helm_release" "gateway" {
  name             = "gateway"
  namespace        = "gateway"
  chart            = "../gateway/helm"
  create_namespace = true
  wait             = true

  values = [<<-EOF
    replicas: 1
    image:
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