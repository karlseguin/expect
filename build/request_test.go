package build

import (
	"compress/gzip"
	"io"
	"net/http"
	"testing"

	e "github.com/karlseguin/expect"
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

func (_ RequestTests) GzipsABodyOnSettingHeader() {
	req := Request().Body("abc123").Header("Content-Encoding", "gzip").Request
	assertGzipBody(req, "abc123")
}

func (_ RequestTests) GzipsABodyWhenHeaderAlreadySet() {
	req := Request().Header("Content-Encoding", "gzip").Body("over=9000!").Request
	assertGzipBody(req, "over=9000!")
}

func assertGzipBody(req *http.Request, expected string) {
	reader, _ := gzip.NewReader(req.Body)
	actual := make([]byte, 100)
	read := 0
	for {
		n, err := reader.Read(actual[read:])
		read += n
		if err == io.EOF {
			break
		}
	}
	e.Expect(string(actual[:read])).To.Equal(expected)
	reader.Close()
}
