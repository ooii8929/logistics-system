resource "aws_s3_bucket" "report" {
  bucket = "alvin-report"

  # for delete
  force_destroy = true
}

resource "aws_s3_bucket_public_access_block" "report" {
  bucket = aws_s3_bucket.report.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}
