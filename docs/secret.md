<!-- BEGIN_TF_DOCS -->

## Requirements

No requirements.
## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |
| <a name="provider_random"></a> [random](#provider\_random) | n/a |
## Inputs

No inputs.
## Modules

No modules.
## Resources

| Name | Type |
|------|------|
| [aws_secretsmanager_secret.secret](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/secretsmanager_secret) | resource |
| [aws_secretsmanager_secret_version.secret_val](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/secretsmanager_secret_version) | resource |
| [random_password.random_string](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password) | resource |
## Outputs

| Name | Description |
|------|-------------|
| <a name="output_sign_key"></a> [sign\_key](#output\_sign\_key) | The generated value for the sign key |
<!-- END_TF_DOCS -->