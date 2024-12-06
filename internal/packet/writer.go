package packet

type Writer struct {
	buffer []byte
}

func NewWriter() *Writer {
	return &Writer{
		buffer: make([]byte, 0, 128),
	}
}

func (w *Writer) Write(data []byte) {
	w.buffer = append(w.buffer, data...)
}

func (w *Writer) WriteInteger(value int) {
	w.Write([]byte{byte(value), byte(value >> 8)})
}

func (w *Writer) WriteLong(value int) {
	w.Write([]byte{byte(value), byte(value >> 8), byte(value >> 16), byte(value >> 24)})
}

func (w *Writer) WriteString(value string) {
	w.WriteInteger(len(value))
	w.Write([]byte(value))
}

func (w *Writer) Bytes() []byte {
	return w.buffer
}
