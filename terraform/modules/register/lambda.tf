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
      # DB_HOST     = module.rds_proxy.proxy_endpoint
      SIGN_KEY = var.sign_key
      DB_HOST  = var.db_host
      DB_PORT  = var.db_port
      DB_NAME  = var.db_name
      DB_USER  = var.db_username
      DB_PASS  = var.db_password
    }
  }

  source_code_hash = filebase64sha256("./lambda.zip")

  vpc_config {
    ipv6_allowed_for_dual_stack = false
    subnet_ids                  = var.private_subnets
    security_group_ids          = [var.security_group_id]
  }
}
