package build

import (
	. "github.com/karlseguin/expect"
	"testing"
)

type RequestTests struct{}

func Test_Request(t *testing.T) {
	Expectify(new(RequestTests), t)
}

func (r *RequestTests) SetsProtocol() {
	req := Request().Proto(2, 4).Request
	Expect(req.Proto).To.Equal("HTTP/2.4")
	Expect(req.ProtoMajor).To.Equal(2)
	Expect(req.ProtoMinor).To.Equal(4)
}

func (r *RequestTests) SetsURL() {
	req := Request().URLString("http://openmymind.net/test?a=1").Request
	Expect(req.URL.Host).To.Equal("openmymind.net")
	Expect(req.URL.Path).To.Equal("/test")
	Expect(req.URL.RawQuery).To.Equal("a=1")
}

func (r *RequestTests) SetsPath() {
	req := Request().Path("/hello").Request
	Expect(req.URL.Host).To.Equal("local.test")
	Expect(req.URL.Path).To.Equal("/hello")
	Expect(req.URL.RawQuery).To.Equal("")
}

func (r *RequestTests) SetsHost() {
	req := Request().Host("openmymind.io").Request
	Expect(req.URL.Host).To.Equal("openmymind.io")
	Expect(req.URL.Path).To.Equal("/spice")
	Expect(req.Host).To.Equal("openmymind.io")
}
