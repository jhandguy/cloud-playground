resource "grafana_dashboard" "s3" {
  depends_on = [helm_release.prometheus]
  for_each   = toset(var.grafana_dashboards)

  config_json = file("${path.module}/dashboards/${each.key}.json")
}