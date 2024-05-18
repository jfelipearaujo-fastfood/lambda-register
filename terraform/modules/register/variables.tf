variable "lambda_name" {
  type        = string
  description = "The name of the lambda function"
}

variable "vpc_name" {
  type        = string
  description = "The name of the VPC"
}

variable "sign_key" {
  type        = string
  sensitive   = true
  description = "The sign key for the lambda function"
}
