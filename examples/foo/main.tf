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
