resource "aws_secretsmanager_secret" "secret" {
  name = "lambda_${var.lambda_name}_sign_key"
}

resource "random_password" "random_string" {
  length  = 20
  special = false
}

resource "aws_secretsmanager_secret_version" "secret_val" {
  secret_id     = aws_secretsmanager_secret.secret.id
  secret_string = random_password.random_string.result
}
