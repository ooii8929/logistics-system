locals {
  project_name = "${var.project}-${terraform.workspace}-${var.sub_project}"
  
  intra_subnets   = ["10.123.5.0/24", "10.123.6.0/24"]

}

variable "region" {
  description = "aws region"
  default     = "ap-northeast-1"
  type        = string
}

variable "project" {
  description = "name of the project"
  default     = "alvin"
  type        = string
}

variable "sub_project" {
  description = "name of the subproject"
  default     = "practice"
  type        = string
}


