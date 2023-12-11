data "aws_ami" "ubuntu_2204" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners     = ["099720109477"] # Canonical (the company behind Ubuntu)
}


data "aws_iam_policy_document" "ecr_pull_access" {
  statement {
    actions = [
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage",
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetAuthorizationToken",
    ]
    resources = ["*"]
  }
}

data "aws_iam_policy_document" "s3_upload_policy" {
  statement {
    actions = [
      "s3:PutObject",
      "s3:PutObjectAcl",
    ]
    resources = ["${aws_s3_bucket.report.arn}/*"]
  }
}


resource "aws_iam_role" "ecr_pull_role" {
  name = "ecr-pull-role"

  # Remove jsondecode()
  assume_role_policy = data.aws_iam_policy_document.ec2_assume_role_policy.json
}

resource "aws_iam_role_policy" "ecr_pull_policy" {
  name   = "ecr-pull"
  role   = aws_iam_role.ecr_pull_role.id
  policy = data.aws_iam_policy_document.ecr_pull_access.json
}

resource "aws_iam_role_policy" "s3_push_policy" {
  name   = "S3UploadPolicy"
  role   = aws_iam_role.ecr_pull_role.id
  policy = data.aws_iam_policy_document.s3_upload_policy.json
}

data "aws_iam_policy_document" "ec2_assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}


resource "aws_instance" "web" {
  ami           = data.aws_ami.ubuntu_2204.id
  subnet_id   = aws_subnet.main.id
  # associate_public_ip_address = true
  security_groups = [aws_security_group.allow_tls.id]
  key_name = "logistics-system"
  instance_type = "t3.micro"
  user_data = file("./init-instance.sh")
    # 添加这一行以将角色关联到实例
  iam_instance_profile = aws_iam_instance_profile.ecr_pull_profile.name

}
resource "aws_eip" "lb" {
  instance = aws_instance.web.id
  domain   = "vpc"
}
resource "aws_iam_instance_profile" "ecr_pull_profile" {
  name = "your-instance-profile-name"
  role = aws_iam_role.ecr_pull_role.name
}

resource "aws_security_group" "allow_tls" {
  name        = "allow_tls"
  description = "Allow TLS inbound traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

}


output "eip" {
  value = aws_eip.lb.public_ip
  description = "The Elastic IP address associated with the instance"
}