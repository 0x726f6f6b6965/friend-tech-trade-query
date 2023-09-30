module "api_gateway" {
    source  = "terraform-aws-modules/apigateway-v2/aws"
    version = "2.2.2"
    protocol_type = "HTTP"
    name = "Friend.Tech_API"
    cors_configuration = {
        allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
        allow_methods = ["*"]
        allow_origins = ["*"]
    }
    # body = file("${path.module}/docs/openapi.yaml")

    create_api_domain_name = false

    integrations = {
        "POST /query" = {
            lambda_arn             = var.query_api_invoke_arn
            payload_format_version = "2.0"
            timeout_milliseconds   = 12000
        }
    }
}

# module "api_gateway_plan" {
#     source = "value"
# }

# resource "aws_api_gateway_resource" "query_resource" {
#     parent_id = module.api_gateway.aws_api_gateway_rest_api_root_resource_id
#     rest_api_id = module.api_gateway.aws_api_gateway_rest_api_id
#     path_part = "query"
# }

# resource "aws_api_gateway_method" "query_method" {
#     rest_api_id = module.api_gateway.aws_api_gateway_rest_api_id
#     resource_id = aws_api_gateway_resource.query_resource.id
#     http_method = "POST"
#     authorization = "NONE"
# }

# resource "aws_api_gateway_integration" "integration" {
#     rest_api_id             = module.api_gateway.aws_api_gateway_rest_api_id
#     resource_id             = aws_api_gateway_method.query_method.resource_id
#     http_method             = aws_api_gateway_method.query_method.http_method
#     integration_http_method = "POST"
#     type                    = "AWS_PROXY"
#     uri                     = var.query_api_invoke_arn
# }

# resource "aws_api_gateway_deployment" "api_deployment" {
#     rest_api_id = module.api_gateway.aws_api_gateway_rest_api_id
#     stage_name = var.gw_env
#     depends_on = [ 
#         aws_api_gateway_integration.integration
#     ]
# }

output "execution_arn" {
    value = "${module.api_gateway.apigatewayv2_api_execution_arn}/*/*"
}

output "invoke_url" {
    value = module.api_gateway.default_apigatewayv2_stage_invoke_url
}