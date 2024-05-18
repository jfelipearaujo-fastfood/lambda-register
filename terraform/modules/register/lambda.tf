data "aws_secretsmanager_secret" "db_url" {
  name = "db-customers-url-secret"
}

data "aws_secretsmanager_secret_version" "db_url_val" {
  secret_id = data.aws_secretsmanager_secret.db_url.id
}

resource "aws_lambda_function" "lambda_function" {
  function_name = "lambda_${var.lambda_name}"

  filename      = "./lambda.zip"
  role          = aws_iam_role.lambda_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  memory_size   = 128
  timeout       = 30

  environment {
    variables = {
      SIGN_KEY = var.sign_key
      DB_NAME  = "customers"
      DB_URL   = data.aws_secretsmanager_secret_version.db_url_val.secret_string
    }
  }

  source_code_hash = filebase64sha256("./lambda.zip")

  vpc_config {
    ipv6_allowed_for_dual_stack = false
    subnet_ids                  = data.aws_subnets.private_subnets.ids
    security_group_ids          = data.aws_security_groups.dbs_security_groups.ids
  }
}
