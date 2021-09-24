resource "helm_release" "s3" {
  name             = "s3"
  namespace        = "s3"
  chart            = "../s3/helm"
  create_namespace = true
  wait             = true
  version          = "1.0.0"

  values = [
    <<-EOF
    replicas: 1
    horizontalPodAutoscaler:
      minReplicas: 1
      maxReplicas: 2
      targetCPUUtilizationPercentage: 100
    nodePort: ${var.node_port}
    EOF
  ]
}
