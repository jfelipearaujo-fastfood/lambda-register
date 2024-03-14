output "sign_key" {
  description = "The generated value for the sign key"
  sensitive   = true
  value       = random_password.random_string.result
}
