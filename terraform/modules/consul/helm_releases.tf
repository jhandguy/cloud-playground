resource "helm_release" "consul" {
  name             = "consul"
  namespace        = "consul"
  repository       = "https://helm.releases.hashicorp.com"
  chart            = "consul"
  create_namespace = true
  wait             = true

  values = [<<-EOF
    global:
      datacenter: consul
      metrics:
        enabled: true
    connectInject:
      enabled: true
    controller:
      enabled: true
    client:
      enabled: true
    server:
      replicas: 1
    ui:
      service:
        type: NodePort
        nodePort:
          http: ${var.node_port}
    ingressGateways:
      enabled: true
      defaults:
        replicas: 1
        service:
          type: NodePort
      gateways:
%{for name, node_port in var.node_ports~}
        - name: ${name}
          service:
            ports:
              - port: ${var.ingress_gateway_port}
                nodePort: ${node_port}
%{endfor~}
    EOF
  ]
}