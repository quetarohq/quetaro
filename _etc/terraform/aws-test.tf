########################################################################
# Lambda
########################################################################

data "archive_file" "qtr_job_test" {
  type        = "zip"
  output_path = "qtr_job_test.zip"
  source_file = "qtr_job_test.js"
}

resource "aws_lambda_function" "qtr_job_test" {
  function_name    = "qtr-job-test"
  runtime          = "nodejs18.x"
  role             = aws_iam_role.qtr_job.arn
  handler          = "qtr_job_test.handler"
  filename         = data.archive_file.qtr_job_test.output_path
  source_code_hash = data.archive_file.qtr_job_test.output_base64sha256
}

resource "aws_lambda_function_event_invoke_config" "qtr_job_test" {
  function_name                = aws_lambda_function.qtr_job_test.function_name
  maximum_event_age_in_seconds = 60
  maximum_retry_attempts       = 0

  destination_config {
    on_success {
      destination = aws_sqs_queue.qtr_outlet_success_test.arn
    }
    on_failure {
      destination = aws_sqs_queue.qtr_outlet_failure_test.arn
    }
  }
}

########################################################################
# SQS
########################################################################

resource "aws_sqs_queue" "qtr_intake_test" {
  name = "qtr-intake-test"
}

resource "aws_sqs_queue" "qtr_outlet_success_test" {
  name = "qtr-outlet-success-test"
}

resource "aws_sqs_queue" "qtr_outlet_failure_test" {
  name = "qtr-outlet-failure-test"
}
