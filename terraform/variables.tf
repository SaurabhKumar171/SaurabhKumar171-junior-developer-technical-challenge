variable "project_id" {
  description = "rick-morty-project-442813"
  type        = string
}

variable "region" {
  description = "GCP Region"
  default     = "asia-south2"
}

variable "zone" {
  description = "GCP Zone"
  default     = "asia-south2-c"
}

variable "disk_size" {
  description = "Size of the MongoDB disk in GB"
  type        = number
  default     = 10
}

variable "machine_type" {
  description = "GCP machine type for MongoDB"
  type        = string
  default     = "e2-medium"
}
