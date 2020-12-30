resource "kubernetes_ingress" "jenkins" {
  metadata {
    name      = helm_release.jenkins.name
    namespace = helm_release.jenkins.namespace
  }

  spec {
    rule {
      http {
        path {
          backend {
            service_name = data.kubernetes_service.jenkins.metadata.0.name
            service_port = data.kubernetes_service.jenkins.spec.0.port.0.port
          }

          path = local.jenkins_uri_prefix
        }
      }
    }
  }

  wait_for_load_balancer = true
}

resource "kubernetes_ingress" "localstack" {
  metadata {
    name      = helm_release.localstack.name
    namespace = helm_release.localstack.namespace
  }

  spec {
    rule {
      http {
        path {
          backend {
            service_name = data.kubernetes_service.localstack.metadata.0.name
            service_port = data.kubernetes_service.localstack.spec.0.port.0.port
          }

          path = "/"
        }
      }
    }
  }

  wait_for_load_balancer = true
}

resource "kubernetes_ingress" "s3" {
  metadata {
    name      = helm_release.s3.name
    namespace = helm_release.s3.namespace
  }

  spec {
    rule {
      http {
        path {
          backend {
            service_name = data.kubernetes_service.s3.metadata.0.name
            service_port = data.kubernetes_service.s3.spec.0.port.0.port
          }

          path = local.s3_uri_prefix
        }
      }
    }
  }

  wait_for_load_balancer = true
}