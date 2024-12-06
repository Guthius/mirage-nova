package packet

type Reader struct {
	buffer []byte
}

func NewReader(buffer []byte) *Reader {
	return &Reader{
		buffer: buffer,
	}
}

func (r *Reader) Read(length int) []byte {
	data := r.buffer[:length]
	r.buffer = r.buffer[length:]
	return data
}

func (r *Reader) ReadByte() byte {
	return r.Read(1)[0]
}

func (r *Reader) ReadInteger() int {
	data := r.Read(2)
	return int(data[1])<<8 | int(data[0])
}

func (r *Reader) ReadLong() int {
	data := r.Read(4)
	return int(data[3])<<24 | int(data[2])<<16 | int(data[1])<<8 | int(data[0])
}

func (r *Reader) ReadString() string {
	length := r.ReadInteger()
	return string(r.Read(length))
}

func (r *Reader) Remaining() int {
	return len(r.buffer)
}
