// For ECR Deployment
resource "aws_iam_user" "ecr_user" {
  name = "backend-ci-for-github-action"
}

resource "aws_iam_access_key" "ecr_user_key" {
  user = aws_iam_user.ecr_user.name
}

resource "aws_iam_user_policy" "ecr_deploy_policy" {
  name = "backend-ci-for-github-action"
  user = aws_iam_user.ecr_user.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:GetRepositoryPolicy",
          "ecr:DescribeRepositories",
          "ecr:ListImages",
          "ecr:DescribeImages",
          "ecr:BatchGetImage",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload",
          "ecr:PutImage",
          "ecr:CreateRepository",
        ]
        Effect   = "Allow"
        Resource = "*"
      }
    ]
  })
}
