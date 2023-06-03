terraform {
  required_providers {
    bhyve = {
      source = "irth/bhyve"
    }
  }
}

provider "bhyve" {
  host = "bhyve-host"
}

data "bhyve_switches" "switches" {
}

output "switches" {
  value = data.bhyve_switches.switches
}

resource "bhyve_switch" "tf" {
  name    = "tf2"
  address = "10.21.38.0/24"
}

data "bhyve_config" "this" {}

output "config" {
  value = data.bhyve_config.this
}

resource "bhyve_iso" "mfsbsd" {
  image = true
  name = "mfsbsd_tf.iso"
  sha256sum = "9c461692692cf4f1de218439c16645e5c7e708fb8303d33eb20304c074bb1cb5"
  url = "https://mfsbsd.vx.sk/files/iso/13/amd64/mfsbsd-se-13.1-RELEASE-amd64.iso"
}
