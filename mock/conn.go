package mock

type Conn struct {
	t         int
	err       error
	block     bool
	readIndex int
	first     bool
	readings  [][]byte
	Closed    bool
	Written   [][]byte
}

func Conn() *Conn {
	return &Conn{
		first:    true,
		Written:  make([][]byte, 0, 4),
		readings: make([][]byte, 0, 4),
	}
}

func (c *Conn) Read(out []byte) (int, error) {
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

func (c *Conn) Write(b []byte) (int, error) {
	if c.err != nil {
		return 0, c.err
	}
	c.Written = append(c.Written, b)
	return len(b), nil
}

func (c *Conn) Close() error {
	c.Closed = true
	return nil
}

func (c *Conn) LocalAddr() net.Addr {
	return nil
}

func (c *Conn) RemoteAddr() net.Addr {
	return nil
}

func (c *Conn) SetDeadline(t time.Time) error {
	c.SetReadDeadline(t)
	c.SetWriteDeadline(t)
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (c *Conn) Error(err error) *Conn {
	c.err = err
	return c
}

func (c *Conn) Reading(data ...[][]byte) *Conn {
	for _, d := range data {
		c.readings = append(c.readings, d)
	}
	return c
}

func (c *Conn) Block() *Conn {
	c.block = true
	return c
}
