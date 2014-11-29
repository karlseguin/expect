package expect

import (
	"github.com/karlseguin/expect/mock"
	"github.com/karlseguin/expect/build"
)

func MockConn() *mock.MockConn {
	return mock.Conn()
}

func RequestBuilder() *build.RequestBuilder {
	return build.Request()
}
