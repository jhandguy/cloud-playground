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

          path = "/${random_pet.s3_uri_prefix.id}"
        }
      }
    }
  }

  wait_for_load_balancer = true
}

resource "kubernetes_ingress" "dynamo" {
  metadata {
    name      = helm_release.dynamo.name
    namespace = helm_release.dynamo.namespace
  }

  spec {
    rule {
      http {
        path {
          backend {
            service_name = data.kubernetes_service.dynamo.metadata.0.name
            service_port = data.kubernetes_service.dynamo.spec.0.port.0.port
          }

          path = "/${random_pet.dynamo_uri_prefix.id}"
        }
      }
    }
  }

  wait_for_load_balancer = true
}