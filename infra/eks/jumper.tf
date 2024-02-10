resource "aws_iam_role" "bastion" {
  name = "eks-bastion-role"
  

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"

        Principal =  {
          Service = "eks.amazonaws.com"
        }

        
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "eks_describe" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.bastion.name
}

# SSM
resource "aws_iam_role_policy_attachment" "bastion_ssm_policy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM"
  role       = aws_iam_role.bastion.name
}


resource "aws_iam_instance_profile" "bastion" {
  name = "eks-bastion-instance-profile"
  # role = aws_iam_role.bastion.name
  role = aws_iam_role.bastion.name
}

resource "aws_iam_policy" "custom_eks_policy" {
  name        = "CustomEKSPolicy"
  description = "A custom policy for describing EKS clusters"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "eks:*"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "configmaps:*"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "custom_eks_policy" {
  policy_arn = aws_iam_policy.custom_eks_policy.arn
  role       = aws_iam_role.bastion.name
}


resource "aws_instance" "bastion" {
  ami           = "ami-0c056d433176c20ec"            # Replace with the desired Amazon Linux 2 AMI ID for your region
  instance_type = "t2.micro"

  key_name          = "sre-alvin"       # Set the key pair name which includes the public key for SSH access
  vpc_security_group_ids = [aws_security_group.bastion_sg.id]
  subnet_id             = module.vpc.private_subnets[0]

  iam_instance_profile = aws_iam_instance_profile.bastion.name

  user_data = <<-EOF
              #!/bin/bash
              sudo curl --silent --location -o /usr/local/bin/kubectl https://s3.us-west-2.amazonaws.com/amazon-eks/1.21.5/2022-01-21/bin/linux/amd64/kubectl
              sudo chmod +x /usr/local/bin/kubectl
              curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
              sudo mv /tmp/eksctl /usr/local/bin
              aws eks update-kubeconfig --region ap-northeast-1 --name alvin-dev-practice-eks
              EOF
  tags = {
    Name = "${local.project_name}-bastion"
  }
}



output "bastion_instance_id" {
  description = "The instance ID of the bastion host"
  value       = aws_instance.bastion.id
}


resource "aws_security_group" "bastion_sg" {
  name        = "${local.project_name}-bastion-sg"
  description = "Allow SSH access to bastion host"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Change this to restrict access to specific IP addresses or CIDR ranges
  }
    ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Change this to restrict access to specific IP addresses or CIDR ranges
  }


egress {
    from_port   = 0             # Allow all outbound traffic
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "EKS-Bastion-SG"
  }
}


### Public Bastion

# resource "aws_eip" "bastion" {
#   instance = aws_instance.bastion.id

#   tags = {
#     Name = "${local.project_name}-bastion-eip"
#   }
# }

# output "bastion_elastic_ip" {
#   value       = aws_eip.bastion.public_ip
#   description = "Elastic IP of the bastion host"
# }
