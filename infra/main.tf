provider "aws" {
    region = var.region
}

resource "random_id" "unique_suffix" {
    byte_length = 2
}

data "archive_file" "lambda_zip" {
    type        = "zip"
    source_file = "bin/query"
    output_path = "bin/query.zip"
}

output "api_url" {
    value = aws_api_gateway_deployment.api_deployment.invoke_url
}