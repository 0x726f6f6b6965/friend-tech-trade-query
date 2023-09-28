resource "aws_lambda_function" "query_trade" {
    filename = data.archive_file.lambda_zip.output_path
    function_name = local.app_id
    handler = "query"
    source_code_hash = base64sha256(data.archive_file.lambda_zip.output_path)
    runtime = "go1.x"
    role = "${aws_iam_role.lambda_exec.arn}"
}

resource "aws_iam_role" "lambda_exec" {
    name_prefix = local.app_id
    assume_role_policy = file("./policy/assume_role_policy.json")
    inline_policy {
      name = "access_dynamodb"
      policy = file("./policy/lambda_access_dynamodb.json")
    }
}

resource "aws_iam_policy_attachment" "role_attach" {
    name       = "policy-${local.app_id}"
    roles      = [aws_iam_role.lambda_exec.id]
    count      = length(var.iam_policy_arn)
    policy_arn = element(var.iam_policy_arn, count.index)
}