locals {
  cluster_name = "minimal-gce-private.example.com"
  project      = "testproject"
  region       = "us-test1"
}

output "cluster_name" {
  value = "minimal-gce-private.example.com"
}

output "project" {
  value = "testproject"
}

output "region" {
  value = "us-test1"
}

provider "google" {
  region = "us-test1"
}

resource "google_compute_disk" "d1-etcd-events-minimal-gce-private-example-com" {
  labels = {
    "k8s-io-cluster-name" = "minimal-gce-private-example-com"
    "k8s-io-etcd-events"  = "1-2f1"
    "k8s-io-role-master"  = "master"
  }
  name = "d1-etcd-events-minimal-gce-private-example-com"
  size = 20
  type = "pd-ssd"
  zone = "us-test1-a"
}

resource "google_compute_disk" "d1-etcd-main-minimal-gce-private-example-com" {
  labels = {
    "k8s-io-cluster-name" = "minimal-gce-private-example-com"
    "k8s-io-etcd-main"    = "1-2f1"
    "k8s-io-role-master"  = "master"
  }
  name = "d1-etcd-main-minimal-gce-private-example-com"
  size = 20
  type = "pd-ssd"
  zone = "us-test1-a"
}

resource "google_compute_firewall" "cidr-to-master-minimal-gce-private-example-com" {
  allow {
    ports    = ["443"]
    protocol = "tcp"
  }
  allow {
    ports    = ["4194"]
    protocol = "tcp"
  }
  name          = "cidr-to-master-minimal-gce-private-example-com"
  network       = google_compute_network.default.name
  source_ranges = ["100.64.0.0/10"]
  target_tags   = ["minimal-gce-private-example-com-k8s-io-role-master"]
}

resource "google_compute_firewall" "cidr-to-node-minimal-gce-private-example-com" {
  allow {
    protocol = "tcp"
  }
  allow {
    protocol = "udp"
  }
  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "esp"
  }
  allow {
    protocol = "ah"
  }
  allow {
    protocol = "sctp"
  }
  name          = "cidr-to-node-minimal-gce-private-example-com"
  network       = google_compute_network.default.name
  source_ranges = ["100.64.0.0/10"]
  target_tags   = ["minimal-gce-private-example-com-k8s-io-role-node"]
}

resource "google_compute_firewall" "kubernetes-master-https-minimal-gce-private-example-com" {
  allow {
    ports    = ["443"]
    protocol = "tcp"
  }
  name          = "kubernetes-master-https-minimal-gce-private-example-com"
  network       = google_compute_network.default.name
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["minimal-gce-private-example-com-k8s-io-role-master"]
}

resource "google_compute_firewall" "master-to-master-minimal-gce-private-example-com" {
  allow {
    protocol = "tcp"
  }
  allow {
    protocol = "udp"
  }
  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "esp"
  }
  allow {
    protocol = "ah"
  }
  allow {
    protocol = "sctp"
  }
  name        = "master-to-master-minimal-gce-private-example-com"
  network     = google_compute_network.default.name
  source_tags = ["minimal-gce-private-example-com-k8s-io-role-master"]
  target_tags = ["minimal-gce-private-example-com-k8s-io-role-master"]
}

resource "google_compute_firewall" "master-to-node-minimal-gce-private-example-com" {
  allow {
    protocol = "tcp"
  }
  allow {
    protocol = "udp"
  }
  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "esp"
  }
  allow {
    protocol = "ah"
  }
  allow {
    protocol = "sctp"
  }
  name        = "master-to-node-minimal-gce-private-example-com"
  network     = google_compute_network.default.name
  source_tags = ["minimal-gce-private-example-com-k8s-io-role-master"]
  target_tags = ["minimal-gce-private-example-com-k8s-io-role-node"]
}

resource "google_compute_firewall" "node-to-master-minimal-gce-private-example-com" {
  allow {
    ports    = ["443"]
    protocol = "tcp"
  }
  allow {
    ports    = ["4194"]
    protocol = "tcp"
  }
  name        = "node-to-master-minimal-gce-private-example-com"
  network     = google_compute_network.default.name
  source_tags = ["minimal-gce-private-example-com-k8s-io-role-node"]
  target_tags = ["minimal-gce-private-example-com-k8s-io-role-master"]
}

resource "google_compute_firewall" "node-to-node-minimal-gce-private-example-com" {
  allow {
    protocol = "tcp"
  }
  allow {
    protocol = "udp"
  }
  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "esp"
  }
  allow {
    protocol = "ah"
  }
  allow {
    protocol = "sctp"
  }
  name        = "node-to-node-minimal-gce-private-example-com"
  network     = google_compute_network.default.name
  source_tags = ["minimal-gce-private-example-com-k8s-io-role-node"]
  target_tags = ["minimal-gce-private-example-com-k8s-io-role-node"]
}

resource "google_compute_firewall" "nodeport-external-to-node-minimal-gce-private-example-com" {
  allow {
    ports    = ["30000-32767"]
    protocol = "tcp"
  }
  allow {
    ports    = ["30000-32767"]
    protocol = "udp"
  }
  name        = "nodeport-external-to-node-minimal-gce-private-example-com"
  network     = google_compute_network.default.name
  source_tags = ["minimal-gce-private-example-com-k8s-io-role-node"]
  target_tags = ["minimal-gce-private-example-com-k8s-io-role-node"]
}

resource "google_compute_firewall" "ssh-external-to-master-minimal-gce-private-example-com" {
  allow {
    ports    = ["22"]
    protocol = "tcp"
  }
  name          = "ssh-external-to-master-minimal-gce-private-example-com"
  network       = google_compute_network.default.name
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["minimal-gce-private-example-com-k8s-io-role-master"]
}

resource "google_compute_firewall" "ssh-external-to-node-minimal-gce-private-example-com" {
  allow {
    ports    = ["22"]
    protocol = "tcp"
  }
  name          = "ssh-external-to-node-minimal-gce-private-example-com"
  network       = google_compute_network.default.name
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["minimal-gce-private-example-com-k8s-io-role-node"]
}

resource "google_compute_instance_group_manager" "a-master-us-test1-a-minimal-gce-private-example-com" {
  base_instance_name = "master-us-test1-a"
  name               = "a-master-us-test1-a-minimal-gce-private-example-com"
  target_size        = 1
  version {
    instance_template = google_compute_instance_template.master-us-test1-a-minimal-gce-private-example-com.self_link
  }
  zone = "us-test1-a"
}

resource "google_compute_instance_group_manager" "a-nodes-minimal-gce-private-example-com" {
  base_instance_name = "nodes"
  name               = "a-nodes-minimal-gce-private-example-com"
  target_size        = 2
  version {
    instance_template = google_compute_instance_template.nodes-minimal-gce-private-example-com.self_link
  }
  zone = "us-test1-a"
}

resource "google_compute_instance_template" "master-us-test1-a-minimal-gce-private-example-com" {
  can_ip_forward = true
  disk {
    auto_delete  = true
    boot         = true
    device_name  = "persistent-disks-0"
    disk_name    = ""
    disk_size_gb = 64
    disk_type    = "pd-standard"
    interface    = ""
    mode         = "READ_WRITE"
    source       = ""
    source_image = "https://www.googleapis.com/compute/v1/projects/cos-cloud/global/images/cos-stable-57-9202-64-0"
    type         = "PERSISTENT"
  }
  machine_type = "n1-standard-1"
  metadata = {
    "cluster-name"                    = "minimal-gce-private.example.com"
    "kops-k8s-io-instance-group-name" = "master-us-test1-a"
    "ssh-keys"                        = file("${path.module}/data/google_compute_instance_template_master-us-test1-a-minimal-gce-private-example-com_metadata_ssh-keys")
    "startup-script"                  = file("${path.module}/data/google_compute_instance_template_master-us-test1-a-minimal-gce-private-example-com_metadata_startup-script")
  }
  name_prefix = "master-us-test1-a-minimal-asf34c-"
  network_interface {
    network = google_compute_network.default.name
  }
  scheduling {
    automatic_restart   = true
    on_host_maintenance = "MIGRATE"
    preemptible         = false
  }
  service_account {
    email  = "default"
    scopes = ["https://www.googleapis.com/auth/compute", "https://www.googleapis.com/auth/monitoring", "https://www.googleapis.com/auth/logging.write", "https://www.googleapis.com/auth/devstorage.read_write", "https://www.googleapis.com/auth/ndev.clouddns.readwrite"]
  }
  tags = ["minimal-gce-private-example-com-k8s-io-role-master"]
}

resource "google_compute_instance_template" "nodes-minimal-gce-private-example-com" {
  can_ip_forward = true
  disk {
    auto_delete  = true
    boot         = true
    device_name  = "persistent-disks-0"
    disk_name    = ""
    disk_size_gb = 128
    disk_type    = "pd-standard"
    interface    = ""
    mode         = "READ_WRITE"
    source       = ""
    source_image = "https://www.googleapis.com/compute/v1/projects/cos-cloud/global/images/cos-stable-57-9202-64-0"
    type         = "PERSISTENT"
  }
  machine_type = "n1-standard-2"
  metadata = {
    "cluster-name"                    = "minimal-gce-private.example.com"
    "kops-k8s-io-instance-group-name" = "nodes"
    "ssh-keys"                        = file("${path.module}/data/google_compute_instance_template_nodes-minimal-gce-private-example-com_metadata_ssh-keys")
    "startup-script"                  = file("${path.module}/data/google_compute_instance_template_nodes-minimal-gce-private-example-com_metadata_startup-script")
  }
  name_prefix = "nodes-minimal-gce-private-4aopo5-"
  network_interface {
    network = google_compute_network.default.name
  }
  scheduling {
    automatic_restart   = true
    on_host_maintenance = "MIGRATE"
    preemptible         = false
  }
  service_account {
    email  = "default"
    scopes = ["https://www.googleapis.com/auth/compute", "https://www.googleapis.com/auth/monitoring", "https://www.googleapis.com/auth/logging.write", "https://www.googleapis.com/auth/devstorage.read_only"]
  }
  tags = ["minimal-gce-private-example-com-k8s-io-role-node"]
}

resource "google_compute_network" "default" {
  auto_create_subnetworks = true
  name                    = "default"
}

resource "google_compute_router" "nat-minimal-gce-private-example-com" {
  name    = "nat-minimal-gce-private-example-com"
  network = "https://www.googleapis.com/compute/v1/projects/testproject/global/networks/default"
}

resource "google_compute_router_nat" "nat-minimal-gce-private-example-com" {
  name                               = "nat-minimal-gce-private-example-com"
  nat_ip_allocate_option             = "AUTO_ONLY"
  region                             = "us-test1"
  router                             = "nat-minimal-gce-private-example-com"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
}

terraform {
  required_version = ">= 0.12.26"
  required_providers {
    google = {
      "source"  = "hashicorp/google"
      "version" = ">= 2.19.0"
    }
  }
}
