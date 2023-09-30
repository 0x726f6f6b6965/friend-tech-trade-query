variable "region" {
    default = "ap-northeast-1"
}

variable "prefix" {
    default = "test-project"
}

variable "app_name" {
    description = "Query Friend.Tech trading record"
    default     = "test-project-api"
}

variable "app_env" {
    description = "Application environment tag"
    default     = "dev"
}

variable "rest_api_query_path" {
    default     = "query"
    type        = string
}

variable "rest_api_domain_name" {
    default     = "test-project.com"
    description = "Domain name of the API Gateway REST API for self-signed TLS certificate"
    type        = string
}

variable "rest_api_name" {
    default     = "api-gateway-rest-api-openapi-test-project"
    description = "Name of the API Gateway REST API (can be used to trigger redeployments)"
    type        = string
}

# variable "iam_policy_arn" {
#     description = "IAM Policy to be attached to role"
#     type        = list(string)
#     default = [
#         "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
#     ]
# }