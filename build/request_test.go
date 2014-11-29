package build

import (
	e "github.com/karlseguin/expect"
	"testing"
)

type RequestTests struct{}

func Test_Request(t *testing.T) {
	e.Expectify(new(RequestTests), t)
}

func (_ RequestTests) SetsProtocol() {
	req := Request().Proto(2, 4).Request
	e.Expect(req.Proto).To.Equal("HTTP/2.4")
	e.Expect(req.ProtoMajor).To.Equal(2)
	e.Expect(req.ProtoMinor).To.Equal(4)
}

func (_ RequestTests) SetsURL() {
	req := Request().URLString("http://openmymind.net/test?a=1").Request
	e.Expect(req.URL.Host).To.Equal("openmymind.net")
	e.Expect(req.URL.Path).To.Equal("/test")
	e.Expect(req.URL.RawQuery).To.Equal("a=1")
}

func (_ RequestTests) SetsPath() {
	req := Request().Path("/hello").Request
	e.Expect(req.URL.Host).To.Equal("local.test")
	e.Expect(req.URL.Path).To.Equal("/hello")
	e.Expect(req.URL.RawQuery).To.Equal("")
}

func (_ RequestTests) SetsRawQuery() {
	req := Request().RawQuery("a=1&b=2").Request
	e.Expect(req.URL.Query()["a"][0]).To.Equal("1")
	e.Expect(req.URL.Query()["b"][0]).To.Equal("2")
	e.Expect(req.URL.RawQuery).To.Equal("a=1&b=2")
}

func (_ RequestTests) SetsHost() {
	req := Request().Host("openmymind.io").Request
	e.Expect(req.URL.Host).To.Equal("openmymind.io")
	e.Expect(req.URL.Path).To.Equal("/spice")
	e.Expect(req.Host).To.Equal("openmymind.io")
}
