resource "null_resource" "logs" {
  depends_on = [helm_release.cli]

  provisioner "local-exec" {
    command = <<-EOF
      kubectl logs job/cli -c cli -n cli
    EOF
  }
}
