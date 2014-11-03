package expect

var (
	FailureHandler = &FailurePostHandler{}
	SuccessHandler = &SuccessPostHandler{}
)

type PostHandler interface {
	Message(format string, args ...interface{})
}

type FailurePostHandler struct {
}

func (h *FailurePostHandler) Message(format string, args ...interface{}) {
	runner.ErrorMessage(format, args...)
}

type SuccessPostHandler struct {
}

func (h *SuccessPostHandler) Message(format string, args ...interface{}) {

}
