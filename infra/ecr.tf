resource "aws_ecr_repository" "application" {
  name = "logistics-track"

  # for delete
  force_delete = true
}

resource "aws_ecr_repository" "database" {
  name = "logistics-database"

  # for delete
  force_delete = true
}

output "application_repository_url" {
  description = "The URL of the created ECR repository"
  value       = aws_ecr_repository.application.repository_url
}

output "database_repository_url" {
  description = "The URL of the created DB repository"
  value       = aws_ecr_repository.database.repository_url
}