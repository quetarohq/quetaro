########################################################################
# IAM
########################################################################

data "aws_iam_policy_document" "qtr_job_assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "qtr_job_sqs" {
  statement {
    actions = [
      "sqs:SendMessage",
    ]
    resources = [
      aws_sqs_queue.qtr_outlet_success.arn,
      aws_sqs_queue.qtr_outlet_failure.arn,
    ]
  }
}

data "aws_iam_policy_document" "qtr_job_logs" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = [
      "*",
    ]
  }
}

resource "aws_iam_role" "qtr_job" {
  name_prefix        = "qtr-job-"
  assume_role_policy = data.aws_iam_policy_document.qtr_job_assume_role_policy.json

  inline_policy {
    name   = "sqs"
    policy = data.aws_iam_policy_document.qtr_job_sqs.json
  }

  inline_policy {
    name   = "logs"
    policy = data.aws_iam_policy_document.qtr_job_logs.json
  }
}

########################################################################
# Lambda
########################################################################

data "archive_file" "qtr_job" {
  type        = "zip"
  output_path = "qtr_job.zip"
  source_file = "qtr_job.js"
}

resource "aws_lambda_function" "qtr_job" {
  function_name    = "qtr-job"
  runtime          = "nodejs18.x"
  role             = aws_iam_role.qtr_job.arn
  handler          = "qtr_job.handler"
  filename         = data.archive_file.qtr_job.output_path
  source_code_hash = data.archive_file.qtr_job.output_base64sha256
}

resource "aws_lambda_function_event_invoke_config" "qtr_job" {
  function_name                = aws_lambda_function.qtr_job.function_name
  maximum_event_age_in_seconds = 60
  maximum_retry_attempts       = 0

  destination_config {
    on_success {
      destination = aws_sqs_queue.qtr_outlet_success.arn
    }
    on_failure {
      destination = aws_sqs_queue.qtr_outlet_failure.arn
    }
  }
}

resource "aws_cloudwatch_log_group" "qtr_jop" {
  name = "/aws/lambda/qtr-job"
}

########################################################################
# SQS
########################################################################

resource "aws_sqs_queue" "qtr_intake" {
  name = "qtr-intake"
}

resource "aws_sqs_queue" "qtr_outlet_success" {
  name = "qtr-outlet-success"
}

resource "aws_sqs_queue" "qtr_outlet_failure" {
  name = "qtr-outlet-failure"
}
