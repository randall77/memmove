package memmove

func MemMove(dst, src *[2048]byte)
func MemMoveSSE2(dst, src *[2048]byte)
func MemMoveAVX(dst, src *[2048]byte)
