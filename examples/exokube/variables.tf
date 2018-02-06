variable "key" {}
variable "secret" {}
variable "key_pair" {}

variable "zone" {
  default = "ch-dk-2"
}

variable "template" {
  default = "Linux Ubuntu 16.04 LTS 64-bit"
}

variable "ubuntu_flavor" {
  default = "xenial"
}
