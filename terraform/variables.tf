variable "project_id" {
  description = "project_id"
  type        = string
}

variable "gcp_credentials_path" {
  description = "gcp_credentials_path"
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

variable "source_ip" {
  description = "source_ip" // instead can put you local ip for security or put (0.0.0.0/0) to allow all
  type        = string
}