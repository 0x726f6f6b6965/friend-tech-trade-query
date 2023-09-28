locals {
  app_id = "${lower(var.app_name)}-${lower(var.app_env)}-${random_id.unique_suffix.hex}"
}