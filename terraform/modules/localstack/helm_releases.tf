resource "helm_release" "localstack" {
  name             = "localstack"
  namespace        = "localstack"
  repository       = "https://localstack.github.io/helm-charts"
  chart            = "localstack"
  create_namespace = true
  wait             = true
  version          = "0.5.5"

  values = [
    <<-EOF
    startServices: s3,dynamodb
    image:
      tag: 0.14.5
    service:
      edgeService:
        nodePort: ${var.node_port}
      externalServicePorts:
        start: 0
        end: 0
%{if var.consul_enabled}
    podAnnotations:
      'consul.hashicorp.com/connect-inject': "true"
      'consul.hashicorp.com/connect-service': "localstack"
      'consul.hashicorp.com/connect-service-port': "edge"
%{endif}
    persistence:
      enabled: true
    EOF
  ]
}
