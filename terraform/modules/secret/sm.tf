resource "aws_secretsmanager_secret" "secret" {
  name = "lambda_sign_key"

  recovery_window_in_days = 0 # means no recovery window
}

resource "random_password" "random_string" {
  length  = 20
  special = false
}

resource "aws_secretsmanager_secret_version" "secret_val" {
  secret_id     = aws_secretsmanager_secret.secret.id
  secret_string = random_password.random_string.result
}
