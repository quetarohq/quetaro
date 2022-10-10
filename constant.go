package quetaro

const (
	JobIdKey = "_id"
)

// job status
const (
	JobStatusInvalid       = "invalid"
	JobStatusPending       = "pending"
	JobStatusInvokeFailure = "invoke_failure"
	JobStatusInvoked       = "invoked"
	JobStatusFailure       = "failure"
	JobStatusPass          = "pass"
)
