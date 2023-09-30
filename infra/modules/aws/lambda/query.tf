# data "aws_caller_identity" "current" {}

module "lambda_query" {
    source  = "terraform-aws-modules/lambda/aws"
    version = "~> 6.0"
    function_name = var.app_id
    description = "Friend.Tech query trade record lambda function"
    handler = "query"
    runtime = "go1.x"
    publish = true
    create_package         = false
    local_existing_package = var.query_api_path
    attach_policy_jsons = true
    policy_jsons = [
        # file("${path.module}/assume_role_policy.json"),
        file("${path.module}/../../../policy/lambda_access_dynamodb.json")
    ]
    number_of_policy_jsons = 1

    allowed_triggers = {
        APIGatewayAny = {
            service    = "apigateway"
            source_arn = var.execution_arn
        }
    }
}


output "invoke_arn" {
    value = module.lambda_query.lambda_function_invoke_arn
}