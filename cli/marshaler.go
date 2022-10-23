package cli

import "io"

type marshaler interface {
	marshal(interface{}, io.Writer) error
}

type rawMarshaler struct {
	m marshaler
}

func (m *rawMarshaler) marshal(v interface{}, w io.Writer) error {
	if s, ok := v.(string); ok {
		_, err := w.Write([]byte(s))
		return err
	}
	return m.m.marshal(v, w)
}
