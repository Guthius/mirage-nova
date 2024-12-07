package net

type PacketWriter struct {
	buffer []byte
}

// NewWriter returns a new packet Writer whose buffer has the default size.
func NewWriter() *PacketWriter {
	return &PacketWriter{
		buffer: make([]byte, 0, 128),
	}
}

// NewBufferSize returns a new packet Writer whose buffer as at least the specified size.
func NewWriterSize(size int) *PacketWriter {
	if size < 32 {
		size = 32
	}
	return &PacketWriter{
		buffer: make([]byte, size),
	}
}

// Write writes the contents of bytes into the packet.
func (w *PacketWriter) Write(bytes []byte) {
	w.buffer = append(w.buffer, bytes...)
}

// WriteByte writes a single byte. The returned error is always nil.
func (w *PacketWriter) WriteByte(value byte) error {
	w.Write([]byte{value})
	return nil
}

// WriteInteger writes a 16-bit integer.
func (w *PacketWriter) WriteInteger(value int) {
	w.Write([]byte{byte(value), byte(value >> 8)})
}

// WriteLong writes a 32-bit integer.
func (w *PacketWriter) WriteLong(value int) {
	w.Write([]byte{byte(value), byte(value >> 8), byte(value >> 16), byte(value >> 24)})
}

// WriteString writes a string.
func (w *PacketWriter) WriteString(value string) {
	w.WriteInteger(len(value))
	w.Write([]byte(value))
}

func (w *PacketWriter) Bytes() []byte {
	return w.buffer
}
