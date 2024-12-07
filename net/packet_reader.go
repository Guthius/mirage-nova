package net

type PacketReader struct {
	buffer []byte
}

// NewReader returns a new packet reader whose buffer holds the specified bytes.
func NewReader(bytes []byte) *PacketReader {
	return &PacketReader{
		buffer: bytes,
	}
}

// Read reads the specified number of bytes.
func (r *PacketReader) Read(count int) []byte {
	data := r.buffer[:count]
	r.buffer = r.buffer[count:]
	return data
}

// ReadByte reads a single byte.
func (r *PacketReader) ReadByte() byte {
	return r.Read(1)[0]
}

// ReadInteger reads a 16-bit integer.
func (r *PacketReader) ReadInteger() int {
	data := r.Read(2)
	return int(data[1])<<8 | int(data[0])
}

// ReadLong reads a 32-bit integer.
func (r *PacketReader) ReadLong() int {
	data := r.Read(4)
	return int(data[3])<<24 | int(data[2])<<16 | int(data[1])<<8 | int(data[0])
}

// ReadString reads a string,
func (r *PacketReader) ReadString() string {
	length := r.ReadInteger()
	return string(r.Read(length))
}

// Remaining returns the number of bytes left to be read.
func (r *PacketReader) Remaining() int {
	return len(r.buffer)
}
