provider "template" {
  version = "~> 1.0"
}

provider "exoscale" {
  version = "~> 0.9.19"
  key = "${var.key}"
  secret = "${var.secret}"

  compute_endpoint = "https://ppapi.exoscale.ch/compute"
}

resource "exoscale_compute" "exokube" {
  display_name = "exokube"
  size = "Medium"
  disk_size = 50
  zone = "${var.zone}"
  template = "${var.template}"
  key_pair = "${var.key_pair}"

  security_groups = [
    "${exoscale_security_group.exokube.name}",
  ]

  user_data = "${data.template_cloudinit_config.exokube.rendered}"
}
