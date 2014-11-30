package expect

import (
	"github.com/karlseguin/expect/build"
	"github.com/karlseguin/expect/mock"
)

func MockConn() *mock.MockConn {
	return mock.Conn()
}

func RequestBuilder() *build.RequestBuilder {
	return build.Request()
}
