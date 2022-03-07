resource "null_resource" "cluster" {
  triggers = {
    cluster_name = var.cluster_name
  }

  provisioner "local-exec" {
    command = <<-EOF
      kind create cluster --name=${self.triggers.cluster_name} --image=kindest/node:${var.node_image} --config=- <<CFG
        apiVersion: kind.x-k8s.io/v1alpha4
        kind: Cluster
        nodes:
          - role: control-plane
            extraPortMappings:
%{for node_port in random_shuffle.node_ports.result~}
              - containerPort: ${node_port}
                hostPort: ${node_port}
%{endfor~}
      CFG
    EOF
  }

  provisioner "local-exec" {
    when = destroy

    command = <<-EOF
      kind delete cluster --name=${self.triggers.cluster_name}
    EOF
  }
}
