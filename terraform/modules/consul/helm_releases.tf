resource "helm_release" "consul" {
  name             = "consul"
  namespace        = "consul"
  repository       = "https://helm.releases.hashicorp.com"
  chart            = "consul"
  create_namespace = true
  wait             = true
  version          = "1.1.0"

  values = [
    <<-EOF
    global:
      datacenter: consul
    connectInject:
      enabled: true
      resources: null
      replicas: 1
      cni:
        enabled: true
        resources: null
      transparentProxy:
        defaultEnabled: false
      sidecarProxy:
        resources: null
      initContainer:
        resources: null
    controller:
      resources: null
    server:
      replicas: 1
      resources: null
    ui:
      service:
        type: NodePort
        nodePort:
          http: ${var.node_port}
    ingressGateways:
      enabled: true
      defaults:
        resources: null
        initCopyConsulContainer:
          resources: null
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
