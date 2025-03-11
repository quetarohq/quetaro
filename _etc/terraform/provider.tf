variable "localstack_host" {
  type    = string
  default = "localhost"
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.90.1"
    }
  }
}

provider "aws" {
  region                      = terraform.workspace == "aws" ? null : "us-east-1"
  access_key                  = terraform.workspace == "aws" ? null : "mock_access_key"
  secret_key                  = terraform.workspace == "aws" ? null : "mock_secret_key"
  skip_credentials_validation = terraform.workspace != "aws"
  skip_metadata_api_check     = terraform.workspace != "aws"
  skip_requesting_account_id  = terraform.workspace != "aws"

  endpoints {
    iam    = terraform.workspace == "aws" ? null : "http://${var.localstack_host}:4566"
    lambda = terraform.workspace == "aws" ? null : "http://${var.localstack_host}:4566"
    sqs    = terraform.workspace == "aws" ? null : "http://${var.localstack_host}:4566"
    logs   = terraform.workspace == "aws" ? null : "http://${var.localstack_host}:4566"
  }
}
