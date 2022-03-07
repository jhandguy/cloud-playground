resource "helm_release" "consul" {
  name             = "consul"
  namespace        = "consul"
  repository       = "https://helm.releases.hashicorp.com"
  chart            = "consul"
  create_namespace = true
  wait             = true
  version          = "0.40.0"

  values = [
    <<-EOF
    global:
      datacenter: consul
    connectInject:
      enabled: true
      resources:
        limits:
          cpu: 25m
          memory: 50Mi
        requests:
          cpu: 25m
          memory: 50Mi
      replicas: 1
      transparentProxy:
        defaultEnabled: false
      sidecarProxy:
        resources:
          limits:
            cpu: 25m
            memory: 50Mi
          requests:
            cpu: 25m
            memory: 50Mi
      initContainer:
        resources:
          limits:
            cpu: 25m
            memory: 50Mi
          requests:
            cpu: 25m
            memory: 50Mi
    controller:
      enabled: true
      limits:
        cpu: 25m
        memory: 50Mi
      requests:
        cpu: 25m
        memory: 50Mi
    client:
      enabled: true
      resources:
        requests:
          cpu: 25m
          memory: 50Mi
        limits:
          cpu: 25m
          memory: 50Mi
    server:
      replicas: 1
      resources:
        requests:
          cpu: 25m
          memory: 50Mi
        limits:
          cpu: 25m
          memory: 50Mi
    ui:
      service:
        type: NodePort
        nodePort:
          http: ${var.node_port}
    ingressGateways:
      enabled: true
      resources:
        requests:
          cpu: 25m
          memory: 50Mi
        limits:
          cpu: 25m
          memory: 50Mi
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
