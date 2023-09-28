resource "aws_api_gateway_rest_api" "test_project" {
    name = local.app_id
}

resource "aws_api_gateway_method" "root" {
    rest_api_id = aws_api_gateway_rest_api.test_project.id
    resource_id = aws_api_gateway_rest_api.test_project.root_resource_id
    http_method = "ANY"
    authorization = "NONE"
}

resource "aws_api_gateway_resource" "query" {
    parent_id = aws_api_gateway_rest_api.test_project.root_resource_id
    rest_api_id = aws_api_gateway_rest_api.test_project.id
    path_part = "query"
}

resource "aws_api_gateway_method" "query_method" {
    rest_api_id = aws_api_gateway_rest_api.test_project.id
    resource_id = aws_api_gateway_resource.query.id
    http_method = "POST"
    authorization = "NONE"
}

resource "aws_api_gateway_method" "query_option_method" {
    rest_api_id = aws_api_gateway_rest_api.test_project.id
    resource_id = aws_api_gateway_resource.query.id
    http_method = "OPTIONS"
    authorization = "NONE"
}

resource "aws_api_gateway_method_response" "cors_method_response_200" {
    rest_api_id   = aws_api_gateway_rest_api.test_project.id
    resource_id   = aws_api_gateway_resource.query.id
    http_method   = aws_api_gateway_method.query_option_method.http_method
    status_code   = "200"
    response_parameters = {
        "method.response.header.Access-Control-Allow-Origin" = true
    }
    depends_on = [aws_api_gateway_method.query_option_method]
}

resource "aws_api_gateway_integration" "integration" {
    rest_api_id             = aws_api_gateway_rest_api.test_project.id
    resource_id             = aws_api_gateway_method.query_method.resource_id
    http_method             = aws_api_gateway_method.query_method.http_method
    integration_http_method = "POST"
    type                    = "AWS_PROXY"
    uri                     = aws_lambda_function.query_trade.invoke_arn
}

resource "aws_api_gateway_integration" "options_integration" {
    rest_api_id   = aws_api_gateway_rest_api.test_project.id
    resource_id   = aws_api_gateway_resource.query.id
    http_method   = aws_api_gateway_method.query_option_method.http_method
    type          = "MOCK"
    depends_on = [aws_api_gateway_method.query_option_method]
}

resource "aws_api_gateway_integration" "integration_root" {
  rest_api_id             = aws_api_gateway_rest_api.test_project.id
  resource_id             = aws_api_gateway_method.root.resource_id
  http_method             = aws_api_gateway_method.root.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.query_trade.invoke_arn
}


resource "aws_api_gateway_deployment" "api_deployment" {
    rest_api_id = aws_api_gateway_rest_api.test_project.id
    stage_name = "Dev"
    depends_on = [ 
        aws_api_gateway_integration.integration,
        aws_api_gateway_integration.integration_root,
        aws_api_gateway_integration.options_integration
    ]
}

resource "aws_lambda_permission" "apigw" {
    statement_id  = "AllowAPIGatewayInvoke"
    action        = "lambda:InvokeFunction"
    function_name = "${aws_lambda_function.query_trade.function_name}"
    principal     = "apigateway.amazonaws.com"
    source_arn = "${aws_api_gateway_rest_api.test_project.execution_arn}/*/*"
}