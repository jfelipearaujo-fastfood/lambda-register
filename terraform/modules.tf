module "database" {
  source = "./modules/database"

  db_name = "customers"
}

module "secret" {
  source = "./modules/secret"
}

module "register" {
  source = "./modules/register"

  lambda_name = "register"
  vpc_name    = var.vpc_name

  sign_key = module.secret.sign_key

  depends_on = [
    module.secret
  ]
}
