module "database" {
  source = "./modules/database"

  db_name = "fastfood"
}

module "secret" {
  source = "./modules/secret"

  lambda_name = "register"
}

module "register" {
  source = "./modules/register"

  lambda_name = "register"

  sign_key    = module.secret.sign_key
  db_host     = module.database.db_host
  db_port     = module.database.db_port
  db_name     = module.database.db_name
  db_username = module.database.db_username
  db_password = module.database.db_pass

  private_subnets   = var.private_subnets
  security_group_id = module.database.security_group_id

  depends_on = [
    module.secret
  ]
}
