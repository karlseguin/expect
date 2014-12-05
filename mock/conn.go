package mock

import (
	"net"
	"time"
)

type MockConn struct {
	t         int
	err       error
	block     bool
	readIndex int
	first     bool
	readings  [][]byte
	Closed    bool
	Written   [][]byte
}

// Creates a new mock object which satisfies the net.Conn interface
func Conn() *MockConn {
	return &MockConn{
		first:    true,
		Written:  make([][]byte, 0, 4),
		readings: make([][]byte, 0, 4),
	}
}

func (c *MockConn) Read(out []byte) (int, error) {
	if c.first {
		c.first = false
		if c.err != nil {
			return 0, c.err
		}
	}

	if c.readIndex == len(c.readings) {
		if c.block == false {
			return 0, nil
		}
		time.Sleep(time.Minute)
	}
	bytes := c.readings[c.readIndex]
	end := len(bytes)
	if end > len(out) {
		end = len(out)
	}
	copy(out, bytes[:end])
	length := len(bytes[:end])
	c.readings[c.readIndex] = bytes[end:]
	if len(bytes[end:]) == 0 {
		c.readIndex++
	} else {
		c.readings[c.readIndex] = bytes[end:]
	}
	return length, nil
}

func (c *MockConn) Write(b []byte) (int, error) {
	if c.err != nil {
		return 0, c.err
	}
	c.Written = append(c.Written, b)
	return len(b), nil
}

func (c *MockConn) Close() error {
	c.Closed = true
	return nil
}

func (c *MockConn) LocalAddr() net.Addr {
	return nil
}

func (c *MockConn) RemoteAddr() net.Addr {
	return nil
}

func (c *MockConn) SetDeadline(t time.Time) error {
	c.SetReadDeadline(t)
	c.SetWriteDeadline(t)
	return nil
}

func (c *MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

// noop
func (c *MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (c *MockConn) Error(err error) *MockConn {
	c.err = err
	return c
}

// Define what each call to Read() will read
func (c *MockConn) Reading(data ...[]byte) *MockConn {
	for _, d := range data {
		c.readings = append(c.readings, d)
	}
	return c
}

func (c *MockConn) Block() *MockConn {
	c.block = true
	return c
}

// Returns whatever data hasn't been read yet
func (c *MockConn) Drain() []byte {
	left := 0
	for i := c.readIndex; i < len(c.readings); i++ {
		left += len(c.readings[i])
	}
	data := make([]byte, 0, left)
	for i := c.readIndex; i < len(c.readings); i++ {
		data = append(data, c.readings[i]...)
	}
	return data
}

type TimeoutError struct {
	error
}

func (t *TimeoutError) Timeout() bool {
	return true
}

func (t *TimeoutError) Temporary() bool {
	return true
}
