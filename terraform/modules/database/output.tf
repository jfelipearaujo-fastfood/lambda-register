output "db_host" {
  description = "The hostname of the database"
  value       = data.aws_db_instance.db.address
}

output "db_port" {
  description = "The port of the database"
  value       = data.aws_db_instance.db.port
}

output "db_name" {
  description = "The name of the database"
  value       = data.aws_db_instance.db.db_name
}

output "db_username" {
  description = "The username for the database"
  value       = data.aws_db_instance.db.master_username
}

output "db_pass" {
  description = "The password for the database"
  sensitive   = true
  value       = jsondecode(data.aws_secretsmanager_secret_version.db.secret_string)["password"]
}

output "security_group_id" {
  description = "The ID of the security group"
  value       = data.aws_security_group.db_security_group.id
}
