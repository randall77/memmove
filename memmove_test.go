package memmove

import (
	"testing"
	"unsafe"
)

var x [2048]byte
var y [2048]byte

var z struct {
	pad1 [64]byte
	data [2048]byte
	pad2 [64]byte
}

func TestMemMove(t *testing.T) {
	testMemMove(t, MemMove)
}
func TestMemMoveSSE2(t *testing.T) {
	testMemMove(t, MemMoveSSE2)
}
func TestMemMoveAVX(t *testing.T) {
	testMemMove(t, MemMoveAVX)
}
func testMemMove(t *testing.T, f func(dst, src *[2048]byte)) {
	for i := range z.pad1 {
		z.pad1[i] = 0
	}
	for i := range z.pad2 {
		z.pad2[i] = 0
	}
	for i := range z.data {
		z.data[i] = 0
	}
	n := byte(1)
	for i := range x {
		x[i] = n
		n *= 3
	}

	f(&z.data, &x)

	for i := range z.pad1 {
		if z.pad1[i] != 0 {
			t.Fatalf("overwrite in prepad %d", i)
		}
	}
	for i := range z.pad2 {
		if z.pad2[i] != 0 {
			t.Fatalf("overwrite in postpad %d", i)
		}
	}
	n = byte(1)
	for i := range z.data {
		if z.data[i] != n {
			t.Fatalf("bad copy @ %d", i)
		}
		n *= 3
	}
}

func BenchmarkMemMove(b *testing.B) {
	benchMemMove(b, MemMove)
}

func BenchmarkMemMoveSSE2(b *testing.B) {
	benchMemMove(b, MemMoveSSE2)
}

func BenchmarkMemMoveAVX(b *testing.B) {
	benchMemMove(b, MemMoveAVX)
}

func benchMemMove(b *testing.B, f func(dst, src *[2048]byte)) {
	b.SetBytes(int64(unsafe.Sizeof(x)))
	for i := 0; i < b.N; i++ {
		f(&x, &y)
	}
}

// Results:
// 3.7 GHz Quad-Core Intel Xeon E5 (mac pro)
// BenchmarkMemMove-8    	20000000	        77.5 ns/op	26413.79 MB/s
// BenchmarkMemMoveSSE2-8	50000000	        37.1 ns/op	55234.76 MB/s
// BenchmarkMemMoveAVX-8 	50000000	        37.1 ns/op	55229.94 MB/s
// (prefetch doesn't help, aligned reads/writes don't help)

// Intel(R) Core(TM) i7-4600U CPU @ 2.10GHz (HP laptop)
// BenchmarkMemMove-4    	20000000	        85.1 ns/op	24056.98 MB/s
// BenchmarkMemMoveSSE2-4	30000000	        44.4 ns/op	46176.00 MB/s
// BenchmarkMemMoveAVX-4 	50000000	        24.3 ns/op	84248.36 MB/s

// Intel(R) Xeon(R) CPU E5-1650 0 @ 3.20GHz (work desktop)
// BenchmarkMemMove-12     20000000                84.1 ns/op      24357.54 MB/s
// BenchmarkMemMoveSSE2-12 30000000                37.5 ns/op      54555.20 MB/s
// BenchmarkMemMoveAVX-12  30000000                38.2 ns/op      53670.40 MB/s

// Intel(R) Xeon(R) CPU E5-2676 v3 @ 2.40GHz  (amazon ec2)
// BenchmarkMemMove-2     20000000       100 ns/op   20440.44 MB/s
// BenchmarkMemMoveSSE2-2 30000000        52.5 ns/op 38990.39 MB/s
// BenchmarkMemMoveAVX-2  50000000        29.0 ns/op 70503.32 MB/s

// AMD Opteron(tm) Processor 4122 (dreamhost)
// BenchmarkMemMove-4		 5000000	       387 ns/op	5285.22 MB/s
// BenchmarkMemMoveSSE2-4	 3000000	       427 ns/op	4787.95 MB/s
// no AVX instructions
