provider "aws" {
    region = var.region
}

locals {
    app_id = "${lower(var.app_name)}-${lower(var.app_env)}-${random_id.unique_suffix.hex}"
}

resource "random_id" "unique_suffix" {
    byte_length = 2
}

# resource "aws_iam_role" "lambda_exec" {
#     name_prefix = local.app_id
#     assume_role_policy = file("${path.cwd}/infra/policy/assume_role_policy.json")
#     inline_policy {
#       name = "access_dynamodb"
#       policy = file("${path.cwd}/infra/policy/lambda_access_dynamodb.json")
#     }
# }

# resource "aws_iam_policy_attachment" "role_attach" {
#     name       = "policy-${local.app_id}"
#     roles      = [aws_iam_role.lambda_exec.id]
#     count      = length(var.iam_policy_arn)
#     policy_arn = element(var.iam_policy_arn, count.index)
# }

data "archive_file" "lambda_zip" {
    type        = "zip"
    source_file = "${path.cwd}/bin/query"
    output_path = "${path.cwd}/bin/query.zip"
}

module "gw" {
    source = "./modules/aws/gw"
    query_api_invoke_arn = module.lambda_api.invoke_arn
    gw_env = var.app_env
}
module "lambda_api" {
    source = "./modules/aws/lambda"
    app_id = local.app_id
    query_api_path = data.archive_file.lambda_zip.output_path
    execution_arn = module.gw.execution_arn
    app_env = var.app_env
}
output "api_url" {
    value = module.gw.invoke_url
}
