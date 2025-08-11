package main

type OurByteBuffer struct {
	bytes        []byte
	readPosition int
}

func NewBufferString(s string) OurByteBuffer {
	return OurByteBuffer{
		bytes:        []byte(s),
		readPosition: 0,
	}
}

func (b *OurByteBuffer) Bytes() []byte {
	return b.bytes
}

func (b *OurByteBuffer) Write(bytes []byte) (int, error) {
	b.bytes = append(b.bytes, bytes...)
	return len(bytes), nil
}

func (b *OurByteBuffer) Read(to []byte) (int, error) {
	remainingBytes := len(b.bytes) - b.readPosition
	bytesToRead := len(to)
	if remainingBytes < bytesToRead {
		bytesToRead = remainingBytes
	}
	copy(to, b.bytes[b.readPosition:b.readPosition+bytesToRead])
	b.readPosition += bytesToRead
	return bytesToRead, nil
}
