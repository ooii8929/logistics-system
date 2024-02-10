terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
    # profile = "alvin"
    profile = "backyard"
    region = "ap-northeast-1"
    default_tags {
        tags = {
            Owner = "AlvinLin"
            Terraform = "true"
        }
    }
}

