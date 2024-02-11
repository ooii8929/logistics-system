module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "20.2.1"

  cluster_name                   = "${local.project_name}-eks"
  cluster_endpoint_public_access = true

  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
    }
  }

  vpc_id                   = module.vpc.vpc_id
  subnet_ids               = module.vpc.private_subnets
  control_plane_subnet_ids = module.vpc.intra_subnets

  eks_managed_node_groups = {
    ascode-cluster-wg = {
      min_size     = 1
      max_size     = 3
      desired_size = 3

      instance_types = ["t3.small"]
      capacity_type  = "ON_DEMAND"
    }
  }
}

# resource "aws_iam_role" "admin" {
#   name = "EKSAdminIAMRole"
#   assume_role_policy = data.aws_iam_policy_document.admin.json
# }


# data "aws_iam_policy_document" "admin" {
#   statement {
#     actions = ["sts:AssumeRoleWithWebIdentity"]
#     effect = "Allow"
#     principals {
#       type = "Federated"
#       identifiers = [
#         module.eks.oidc_provider_arn,
#       ]
#     }
#   }
# }

# resource "aws_security_group_rule" "eks_cluster_allow_bastion_sg" {
#   security_group_id = module.eks.cluster_security_group_id

#   type        = "ingress"
#   from_port   = 443
#   to_port     = 443
#   protocol    = "tcp"
#   source_security_group_id = aws_security_group.bastion_sg.id
# }

# For public
resource "aws_security_group_rule" "eks_cluster_api_open" {
  security_group_id = module.eks.cluster_security_group_id

  type        = "ingress"
  from_port   = 443
  to_port     = 443
  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
}


resource "aws_ecr_repository" "this" {
  name = "${local.project_name}-backend"
    lifecycle {
    prevent_destroy = false
  }
}

output "ecr_registry_id" {
  value = aws_ecr_repository.this.registry_id
}

output "eks_id" {
  value =  module.eks.cluster_id
}

output "eks_cluster_endpoint" {
  description = "The endpoint for the EKS cluster"
  value       = module.eks.cluster_endpoint
}
