package cgi2echo

import (
	"bufio"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/cgi"
	"os"
)

type CGI2Echo struct {
	echo *echo.Echo
}

func New() *CGI2Echo {
	return &CGI2Echo{echo.New()}
}

func (e *CGI2Echo) Echo() *echo.Echo {
	return e.echo
}

func (e *CGI2Echo) Serve() error {
	req, err := cgi.Request()
	if err != nil {
		return err
	}
	rw := &response{
		header: make(http.Header),
		bufw:   bufio.NewWriter(os.Stdout),
	}
	e.echo.ServeHTTP(rw, req)
	return rw.bufw.Flush()
}

type response struct {
	header         http.Header
	code           int
	wroteHeader    bool
	wroteCGIHeader bool
	bufw           *bufio.Writer
}

func (r *response) Flush() {
	r.bufw.Flush()
}

func (r *response) Header() http.Header {
	return r.header
}

func (r *response) Write(p []byte) (n int, err error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	if !r.wroteCGIHeader {
		r.writeCGIHeader(p)
	}
	return r.bufw.Write(p)
}

func (r *response) WriteHeader(code int) {
	if r.wroteHeader {
		return
	}
	r.wroteHeader = true
	r.code = code
}

// writeCGIHeader finalizes the header sent to the client and writes it to the output.
// p is not written by writeHeader, but is the first chunk of the body
// that will be written. It is sniffed for a Content-Type if none is
// set explicitly.
func (r *response) writeCGIHeader(p []byte) {
	if r.wroteCGIHeader {
		return
	}
	r.wroteCGIHeader = true
	fmt.Fprintf(r.bufw, "Status: %d %s\r\n", r.code, http.StatusText(r.code))
	if _, hasType := r.header["Content-Type"]; !hasType {
		r.header.Set("Content-Type", http.DetectContentType(p))
	}
	r.header.Write(r.bufw)
	r.bufw.WriteString("\r\n")
	r.bufw.Flush()
}
