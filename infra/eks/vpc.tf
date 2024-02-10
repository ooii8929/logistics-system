module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "5.1.2"

  name = local.project_name
  cidr = "10.0.0.0/16"

  azs             = ["ap-northeast-1a", "ap-northeast-1c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24"]
  intra_subnets   = ["10.0.122.0/24", "10.0.123.0/24"]

  enable_nat_gateway = true

  public_subnet_tags = {
    "kubernetes.io/role/elb" = 1
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = 1
  }

  # public_subnet_tags = {
  #   Terraform                  = "true"
  #   Environment                = "dev"
  #   "kubernetes.io/cluster/${local.project_name}-eks" = "shared"
  #   "kubernetes.io/role/elb"                      = "1"
  #   MapPublicIpOnLaunch                           = "true"
  # }
  # enable_nat_gateway = true
  enable_vpn_gateway = true

  enable_dns_hostnames = true
  enable_dns_support = true



  tags = {
    Terraform = "true"
    Environment = "dev"
  }
}


output "private_subnets_from_vpc" {
  value = module.vpc.private_subnets
}