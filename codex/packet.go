package codex

import "io"

type Packet struct {
	closer io.Closer
	reader io.Reader
	writer io.Writer
}

func NewPacket(rw io.ReadWriter) (Codex, error) {
	codec := &Packet{
		reader: rw,
		writer: rw,
	}
	codec.closer, _ = rw.(io.Closer)
	return codec, nil
}

func (c *Packet) Receive(msg interface{}) error {
	_, err := msg.(io.ReaderFrom).ReadFrom(c.reader)
	return err
}

func (c *Packet) Send(msg interface{}) error {
	_, err := msg.(io.WriterTo).WriteTo(c.writer)
	return err
}

func (c *Packet) Close() error {
	if c.closer != nil {
		return c.closer.Close()
	}
	return nil
}
