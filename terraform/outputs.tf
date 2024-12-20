output "instance_name" {
  value = google_compute_instance.mongodb_vm.name
}

output "instance_ip" {
  value = google_compute_instance.mongodb_vm.network_interface[0].access_config[0].nat_ip
}
