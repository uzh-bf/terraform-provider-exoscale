data "template_file" "exokube" {
  template = "${file("cloud-config.yaml")}"

  vars {
    fqdn = "exokube"
    ubuntu = "${var.ubuntu_flavor}"
  }
}

data "template_cloudinit_config" "exokube" {
  part {
    filename = "init.cfg"
    content_type = "text/cloud-config"
    content = "${data.template_file.exokube.rendered}"
  }
}
