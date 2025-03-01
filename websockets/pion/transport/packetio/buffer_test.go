package packetio

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ahmad-khatib0/go/websockets/pion/transport/test"
	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	assert := assert.New(t)

	buffer := NewBuffer()
	packet := make([]byte, 4)

	// Write once
	n, err := buffer.Write([]byte{0, 1})
	assert.NoError(err)
	assert.Equal(2, n)

	// Read once
	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal([]byte{0, 1}, packet[:n])

	// Read deadline
	err = buffer.SetReadDeadline(time.Unix(0, 1))
	assert.NoError(err)
	n, err = buffer.Read(packet)
	var e net.Error
	if !errors.As(err, &e) || !e.Timeout() {
		t.Errorf("Unexpected error: %v", err)
	}
	assert.Equal(0, n)

	// Reset deadline
	err = buffer.SetReadDeadline(time.Time{})
	assert.NoError(err)

	// Write twice
	n, err = buffer.Write([]byte{2, 3, 4})
	assert.NoError(err)
	assert.Equal(3, n)

	n, err = buffer.Write([]byte{5, 6, 7})
	assert.NoError(err)
	assert.Equal(3, n)

	// Read twice
	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(3, n)
	assert.Equal([]byte{2, 3, 4}, packet[:n])

	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(3, n)
	assert.Equal([]byte{5, 6, 7}, packet[:n])

	// Write once prior to close.
	_, err = buffer.Write([]byte{3})
	assert.NoError(err)

	// Close
	err = buffer.Close()
	assert.NoError(err)

	// Future writes will error
	_, err = buffer.Write([]byte{4})
	assert.Error(err)

	// But we can read the remaining data.
	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(1, n)
	assert.Equal([]byte{3}, packet[:n])

	// Until EOF
	_, err = buffer.Read(packet)
	assert.Equal(io.EOF, err)
}

func testWraparound(t *testing.T, grow bool) {
	t.Helper()

	assert := assert.New(t)

	buffer := NewBuffer()
	err := buffer.grow()
	assert.NoError(err)

	buffer.head = len(buffer.data) - 13
	buffer.tail = buffer.head

	p1 := []byte{1, 2, 3}
	p2 := []byte{4, 5, 6}
	p3 := []byte{7, 8, 9}
	p4 := []byte{10, 11, 12}

	_, err = buffer.Write(p1)
	assert.NoError(err)
	_, err = buffer.Write(p2)
	assert.NoError(err)
	_, err = buffer.Write(p3)
	assert.NoError(err)

	packet := make([]byte, 10)

	n, err := buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(p1, packet[:n])

	if grow {
		err = buffer.grow()
		assert.NoError(err)
	}

	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(p2, packet[:n])

	_, err = buffer.Write(p4)
	assert.NoError(err)

	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(p3, packet[:n])
	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(p4, packet[:n])

	if !grow {
		assert.Equal(len(buffer.data), minSize)
	} else {
		assert.Equal(len(buffer.data), 2*minSize)
	}
}

func TestBufferWraparound(t *testing.T) {
	testWraparound(t, false)
}

func TestBufferWraparoundGrow(t *testing.T) {
	testWraparound(t, true)
}

func TestBufferAsync(t *testing.T) {
	assert := assert.New(t)

	buffer := NewBuffer()

	// Start up a goroutine to start a blocking read.
	done := make(chan struct{})
	go func() {
		packet := make([]byte, 4)

		n, err := buffer.Read(packet)
		assert.NoError(err)
		assert.Equal(2, n)
		assert.Equal([]byte{0, 1}, packet[:n])

		_, err = buffer.Read(packet)
		assert.Equal(io.EOF, err)

		close(done)
	}()

	// Wait for the reader to start reading.
	time.Sleep(time.Millisecond)

	// Write once
	n, err := buffer.Write([]byte{0, 1})
	assert.NoError(err)
	assert.Equal(2, n)

	// Wait for the reader to start reading again.
	time.Sleep(time.Millisecond)

	// Close will unblock the reader.
	err = buffer.Close()
	assert.NoError(err)

	<-done
}

func TestBufferLimitCount(t *testing.T) {
	assert := assert.New(t)

	buffer := NewBuffer()
	buffer.SetLimitCount(2)

	assert.Equal(0, buffer.Count())

	// Write twice
	n, err := buffer.Write([]byte{0, 1})
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal(1, buffer.Count())

	n, err = buffer.Write([]byte{2, 3})
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal(2, buffer.Count())

	// Over capacity
	_, err = buffer.Write([]byte{4, 5})
	assert.Equal(ErrFull, err)
	assert.Equal(2, buffer.Count())

	// Read once
	packet := make([]byte, 4)
	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal([]byte{0, 1}, packet[:n])
	assert.Equal(1, buffer.Count())

	// Write once
	n, err = buffer.Write([]byte{6, 7})
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal(2, buffer.Count())

	// Over capacity
	_, err = buffer.Write([]byte{8, 9})
	assert.Equal(ErrFull, err)
	assert.Equal(2, buffer.Count())

	// Read twice
	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal([]byte{2, 3}, packet[:n])
	assert.Equal(1, buffer.Count())

	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal([]byte{6, 7}, packet[:n])
	assert.Equal(0, buffer.Count())

	// Nothing left.
	err = buffer.Close()
	assert.NoError(err)
}

func TestBufferLimitSize(t *testing.T) {
	assert := assert.New(t)

	buffer := NewBuffer()
	buffer.SetLimitSize(11)

	assert.Equal(0, buffer.Size())

	// Write twice
	n, err := buffer.Write([]byte{0, 1})
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal(4, buffer.Size())

	n, err = buffer.Write([]byte{2, 3})
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal(8, buffer.Size())

	// Over capacity
	_, err = buffer.Write([]byte{4, 5})
	assert.Equal(ErrFull, err)
	assert.Equal(8, buffer.Size())

	// Cheeky write at exact size.
	n, err = buffer.Write([]byte{6})
	assert.NoError(err)
	assert.Equal(1, n)
	assert.Equal(11, buffer.Size())

	// Read once
	packet := make([]byte, 4)
	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal([]byte{0, 1}, packet[:n])
	assert.Equal(7, buffer.Size())

	// Write once
	n, err = buffer.Write([]byte{7, 8})
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal(11, buffer.Size())

	// Over capacity
	_, err = buffer.Write([]byte{9, 10})
	assert.Equal(ErrFull, err)
	assert.Equal(11, buffer.Size())

	// Read everything
	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal([]byte{2, 3}, packet[:n])
	assert.Equal(7, buffer.Size())

	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(1, n)
	assert.Equal([]byte{6}, packet[:n])
	assert.Equal(4, buffer.Size())

	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(2, n)
	assert.Equal([]byte{7, 8}, packet[:n])
	assert.Equal(0, buffer.Size())

	// Nothing left.
	err = buffer.Close()
	assert.NoError(err)
}

func TestBufferLimitSizes(t *testing.T) {
	if sizeHardLimit {
		t.Skip("skipping since packetioSizeHardLimit is enabled")
	}
	sizes := []int{
		128 * 1024,
		1024 * 1024,
		8 * 1024 * 1024,
		0, // default
	}
	const headerSize = 2
	const packetSize = 0x8000

	for _, size := range sizes {
		size := size
		name := "default"
		if size > 0 {
			name = fmt.Sprintf("%dkBytes", size/1024)
		}

		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			buffer := NewBuffer()
			if size == 0 {
				size = maxSize
			} else {
				buffer.SetLimitSize(size + headerSize)
			}
			now := time.Now()
			assert.NoError(buffer.SetReadDeadline(now.Add(5 * time.Second))) // Set deadline to avoid test deadlock

			nPackets := size / (packetSize + headerSize)

			for range nPackets {
				_, err := buffer.Write(make([]byte, packetSize))
				assert.NoError(err)
			}

			// Next write is expected to be errored.
			_, err := buffer.Write(make([]byte, packetSize))
			assert.Error(err, ErrFull)

			packet := make([]byte, size)
			for range nPackets {
				n, err := buffer.Read(packet)
				assert.NoError(err)
				assert.Equal(packetSize, n)
				if err != nil {
					t.FailNow()
				}
			}
		})
	}
}

func TestBufferMisc(t *testing.T) {
	assert := assert.New(t)

	buffer := NewBuffer()

	// Write once
	n, err := buffer.Write([]byte{0, 1, 2, 3})
	assert.NoError(err)
	assert.Equal(4, n)

	// Try to read with a short buffer
	packet := make([]byte, 3)
	_, err = buffer.Read(packet)
	assert.Equal(io.ErrShortBuffer, err)

	// Close
	err = buffer.Close()
	assert.NoError(err)

	// Make sure you can Close twice
	err = buffer.Close()
	assert.NoError(err)
}

func TestBufferAlloc(t *testing.T) {
	packet := make([]byte, 1024)

	test := func(f func(count int) func(), count int, maxVal float64) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()

			// AllocsPerRun returns the average number of allocations during calls to f
			// runs the operation f(count) 3 times and returns the average number of memory
			// allocations per run.
			allocs := testing.AllocsPerRun(3, f(count))
			if allocs > maxVal {
				t.Errorf("count=%v, max=%v, got %v",
					count, maxVal, allocs,
				)
			}
		}
	}

	writer := func(count int) func() {
		return func() {
			buffer := NewBuffer()
			for range count {
				_, err := buffer.Write(packet)
				if err != nil {
					t.Errorf("Write: %v", err)

					break
				}
			}
		}
	}
	// Example: Allocations for 100 Writes
	// Let’s simulate the buffer growth and allocations for 100 writes, assuming
	// each write is 1024 bytes (as in the test).
	//  1- Initial State:
	//      Buffer size: 2048 bytes (minSize).
	//      Can store 2 packets (2048 / 1024 = 2).
	//  2- First Growth:
	//      After 2 writes, the buffer is full.
	//      New size: 4096 bytes (double the current size).
	//      Allocation: 1.
	//  3- Second Growth:
	//      After 4 writes, the buffer is full.
	//      New size: 8192 bytes.
	//      Allocation: 2.
	//  4- Third Growth:
	//      After 8 writes, the buffer is full.
	//      New size: 16384 bytes.
	//      Allocation: 3.
	//  5- Fourth Growth:
	//      After 16 writes, the buffer is full.
	//      New size: 32768 bytes.
	//      Allocation: 4.
	//  6- Fifth Growth:
	//      After 32 writes, the buffer is full.
	//      New size: 65536 bytes.
	//      Allocation: 5.
	//  7- Sixth Growth:
	//      After 64 writes, the buffer is full.
	//      New size: 131072 bytes (128 KB, reaching cutoffSize).
	//      Allocation: 6.
	//  8- Seventh Growth:
	//      After 128 writes, the buffer is full.
	//      New size: 163840 bytes (25% increase).
	//      Allocation: 7.
	//  However, since we’re only doing 100 writes, the buffer grows 6 times (up to 131072 bytes).
	// Each growth results in 1 allocation, so we have 6 allocations from growth.
	//  9- Additional Allocations:
	//     Storing packet length (2 bytes per write): No additional allocations.
	//     Copying packet data: No additional allocations.
	//     Mutex operations and notifications: No additional allocations.
	// 10- Total Allocations:
	//     Growth-related allocations: 6.
	//     Other overhead: ~5 (e.g., internal bookkeeping, circular buffer management).
	//     Total: ~11 allocations for 100 writes.
	t.Run("100 writes", test(writer, 100, 11))
	t.Run("200 writes", test(writer, 200, 14))
	t.Run("400 writes", test(writer, 400, 17))
	t.Run("1000 writes", test(writer, 1000, 21))

	wr := func(count int) func() {
		return func() {
			buffer := NewBuffer()
			for range count {
				_, err := buffer.Write(packet)
				if err != nil {
					t.Fatalf("Write: %v", err)
				}
				_, err = buffer.Read(packet)
				if err != nil {
					t.Fatalf("Read: %v", err)
				}
			}
		}
	}

	t.Run("100 writes and reads", test(wr, 100, 5))
	t.Run("1000 writes and reads", test(wr, 1000, 5))
	t.Run("10000 writes and reads", test(wr, 10000, 5))
}

func benchmarkBufferWR(b *testing.B, size int64, write bool, grow int) { // nolint:unparam
	b.Helper()
	buffer := NewBuffer()
	packet := make([]byte, size)

	// Grow the buffer first
	pad := make([]byte, 1022)
	for buffer.Size() < grow {
		_, err := buffer.Write(pad)
		if err != nil {
			b.Fatalf("Write: %v", err)
		}
	}
	for buffer.Size() > 0 {
		_, err := buffer.Read(pad)
		if err != nil {
			b.Fatalf("Write: %v", err)
		}
	}

	if write {
		_, err := buffer.Write(packet)
		if err != nil {
			b.Fatalf("Write: %v", err)
		}
	}

	b.SetBytes(size)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := buffer.Write(packet)
		if err != nil {
			b.Fatalf("Write: %v", err)
		}
		_, err = buffer.Read(packet)
		if err != nil {
			b.Fatalf("Read: %v", err)
		}
	}
}

// 1- Simulates a typical usage scenario where the buffer is frequently drained (empty).
// 2- Measures the performance of the buffer when it is often empty, which is a common
//
//	case in real-world applications.
func BenchmarkBufferWR14(b *testing.B) {
	benchmarkBufferWR(b, 14, false, 128000)
}

func BenchmarkBufferWR140(b *testing.B) {
	benchmarkBufferWR(b, 140, false, 128000)
}

func BenchmarkBufferWR1400(b *testing.B) {
	benchmarkBufferWR(b, 1400, false, 128000)
}

// 1- Simulates a stress scenario where the buffer is always full, forcing wraparound behavior.
// 2- Measures the performance of the buffer under continuous load, ensuring that the wraparound
//
//	logic is efficient and doesn’t introduce significant overhead.
func BenchmarkBufferWWR14(b *testing.B) {
	benchmarkBufferWR(b, 14, true, 128000)
}

func BenchmarkBufferWWR140(b *testing.B) {
	benchmarkBufferWR(b, 140, true, 128000)
}

func BenchmarkBufferWWR1400(b *testing.B) {
	benchmarkBufferWR(b, 1400, true, 128000)
}

func benchmarkBuffer(b *testing.B, size int64) {
	b.Helper()

	buffer := NewBuffer()
	b.SetBytes(size)

	done := make(chan struct{})
	go func() {
		packet := make([]byte, size)

		for {
			_, err := buffer.Read(packet)
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				b.Error(err)

				break
			}
		}

		close(done)
	}()

	packet := make([]byte, size)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var err error
		for {
			_, err = buffer.Write(packet)
			if !errors.Is(err, ErrFull) {
				break
			}
			time.Sleep(time.Microsecond)
		}
		if err != nil {
			b.Fatal(err)
		}
	}

	err := buffer.Close()
	if err != nil {
		b.Fatal(err)
	}

	<-done
}

func BenchmarkBuffer14(b *testing.B) {
	benchmarkBuffer(b, 14)
}

func BenchmarkBuffer140(b *testing.B) {
	benchmarkBuffer(b, 140)
}

func BenchmarkBuffer1400(b *testing.B) {
	benchmarkBuffer(b, 1400)
}

func TestBufferConcurrentRead(t *testing.T) {
	assert := assert.New(t)

	buffer := NewBuffer()
	packet := make([]byte, 4)

	// Write twice
	n, err := buffer.Write([]byte{2, 3, 4})
	assert.NoError(err)
	assert.Equal(3, n)

	n, err = buffer.Write([]byte{5, 6, 7})
	assert.NoError(err)
	assert.Equal(3, n)

	// Read twice
	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(3, n)
	assert.Equal([]byte{2, 3, 4}, packet[:n])

	n, err = buffer.Read(packet)
	assert.NoError(err)
	assert.Equal(3, n)
	assert.Equal([]byte{5, 6, 7}, packet[:n])

	errCh := make(chan error, 2)
	readIntoErr := func() {
		packet := make([]byte, 4)
		_, readErr := buffer.Read(packet)
		errCh <- readErr
	}
	go readIntoErr()
	go readIntoErr()

	// Close
	err = buffer.Close()
	assert.NoError(err)

	err = <-errCh
	assert.Equal(io.EOF, err)
	err = <-errCh
	assert.Equal(io.EOF, err)
}

func TestBufferConcurrentReadWrite(t *testing.T) {
	defer test.TimeOut(time.Second * 5).Stop()

	assert := assert.New(t)

	buffer := NewBuffer()

	numPkts := 1000
	var numRead uint64
	allRead := make(chan struct{})

	readPkts := func(count int) {
		packet := make([]byte, 4)
		for range count {
			_, readErr := buffer.Read(packet)
			if readErr != nil {
				return
			}
			if atomic.AddUint64(&numRead, 1) == uint64(numPkts) { //nolint:gosec
				close(allRead)

				return
			}
		}
	}
	go readPkts(numPkts)
	go readPkts(numPkts / 100)

	for range numPkts {
		_, writeErr := buffer.Write([]byte{2, 3, 4})
		assert.NoError(writeErr)
	}

	<-allRead
	assert.NoError(buffer.Close())
}
