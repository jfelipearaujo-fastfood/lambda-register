data "aws_db_instance" "db" {
  db_instance_identifier = var.db_name
}

data "aws_secretsmanager_secret" "db" {
  arn = data.aws_db_instance.db.master_user_secret[0].secret_arn

  depends_on = [
    data.aws_db_instance.db
  ]
}

data "aws_secretsmanager_secret_version" "db" {
  secret_id = data.aws_secretsmanager_secret.db.id
}
