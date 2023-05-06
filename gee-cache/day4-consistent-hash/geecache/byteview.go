package geecache

type ByteView struct {
	b []byte
}

func (r ByteView) Len() int {
	return len(r.b)
}

func (r ByteView) ByteSlice() []byte {
	return cloneBytes(r.b)
}

func (r ByteView) String() string {
	return string(r.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
