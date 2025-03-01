// Package packetio provides packet buffer
package packetio

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/ahmad-khatib0/go/websockets/pion/transport/deadline"
)

var errPacketTooBig = errors.New("packet too big")

// BufferPacketType allow the Buffer to know which packet protocol is writing.
type BufferPacketType int

const (
	// RTPBufferPacket indicates the Buffer that is handling RTP packets.
	RTPBufferPacket BufferPacketType = 1
	// RTCPBufferPacket indicates the Buffer that is handling RTCP packets.
	RTCPBufferPacket BufferPacketType = 2
)

// Buffer allows writing packets to an intermediate buffer, which can then
// be read form. This is verify similar to bytes.Buffer but avoids combining
// multiple writes into a single read.
type Buffer struct {
	mutex sync.Mutex

	// this is a circular buffer.  If head <= tail, then the useful data is in the
	// interval [head, tail[.  If tail < head, then the useful data is the union of
	// [head, len[ and [0, tail[. In order to avoid ambiguity when head = tail, we
	// always leave an unused byte in the buffer.
	data         []byte
	head, tail   int
	notify       chan struct{}
	closed       bool
	count        int
	limitSize    int
	limitCount   int
	readDeadlien *deadline.Deadline
}

const (
	minSize    = 2048
	cutoffSize = 128 * 1024
	maxSize    = 4 * 1024 * 1024
)

// NewBuffer creates a new Buffer.
func NewBuffer() *Buffer {
	return &Buffer{
		notify:       make(chan struct{}, 1),
		readDeadlien: deadline.New(),
	}
}

// available returns true if the buffer is large enough to fit a packet
// of the given size, taking overhead into account.
func (b *Buffer) available(size int) bool {
	available := b.head - b.tail
	if available <= 0 {
		available += len(b.data)
	}

	// we interpret head=tail as empty, so always keep a byte free
	if size+2+1 > available {
		return false
	}

	return true
}

// grow increases the size of the buffer.  If it returns nil, then the
// buffer has been grown.  It returns ErrFull if hits a limit.
//
// # Purpose of the grow Function
//
// 1- Expand the Buffer:
//
//   - When the buffer is full (i.e., head and tail meet, leaving no space for new data),
//
//     the grow function is called to increase the buffer's capacity.
//
// 2- Handle Limits:
//
//   - The function respects any size limits (limitSize) and ensures the buffer does not
//
//     grow beyond a maximum size (maxSize).
//
// 3- Copy Data:
//
//   - The existing data in the buffer is copied into the new, larger slice while
//
//     maintaining the circular buffer's structure.
//
// Example Scenario
// Suppose the buffer has the following state:
//
// data = [1, 2, 3, 4, 0, 0, 0, 0] (size = 8)
//
// head = 2, tail = 4 (data is [3, 4])
//
// cutoffSize = 16, minSize = 8, maxSize = 64, limitSize = 0
//
// When grow is called:
// 1- The new size is calculated as 2 * 8 = 16 (since len(b.data) < cutoffSize).
//
// 2- A new slice of size 16 is allocated.
//
// 3- The data [3, 4] is copied into the new slice.
//
// 4- The buffer's state is updated:
//
//	. head = 0, tail = 2
//
//	. data = [3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0].
func (b *Buffer) grow() error {
	var newSize int

	// If the current buffer size (len(b.data)) is less than a predefined
	// cutoffSize, the buffer size is doubled.
	if len(b.data) < cutoffSize {
		newSize = 2 * len(b.data)

		// If the buffer size is larger than or equal to cutoffSize, the buffer size
		// is increased by 25% (to avoid excessive memory usage for large buffers).
	} else {
		newSize = 5 * len(b.data) / 4
	}

	// Ensures the new size is at least minSize. This prevents the buffer from becoming too small.
	if newSize < minSize {
		newSize = minSize
	}

	// If there is no explicit size limit (b.limitSize <= 0) or if a hard limit is
	// enforced (sizeHardLimit), the new size is capped(will be) at maxSize.
	if (b.limitSize <= 0 || sizeHardLimit) && newSize > maxSize {
		newSize = maxSize
	}

	// one byte slack: If an explicit size limit (b.limitSize) is set, the new size is
	// capped at b.limitSize + 1. The +1 ensures there is always one unused byte in the
	// buffer (to avoid ambiguity when head = tail).
	if b.limitSize > 0 && newSize > b.limitSize+1 {
		newSize = b.limitSize + 1
	}

	// If the new size is not larger than the current size, the function returns ErrFull.
	// This indicates that the buffer cannot grow further (e.g., due to size limits).
	if newSize <= len(b.data) {
		return ErrFull
	}

	newData := make([]byte, newSize)
	var n int

	if b.head <= b.tail {
		// data was contiguous: copied directly into the new slice.
		n = copy(newData, b.data[b.head:b.tail])
	} else {
		// data was discontinuous
		n = copy(newData, b.data[b.head:])      // From head to the end of the slice.
		n += copy(newData[n:], b.data[:b.tail]) // From the start of the slice to tail
	}

	// The head is reset to 0 (start of the new slice).
	b.head = 0
	b.tail = n
	b.data = newData

	return nil
}

// Write appends a copy of the packet data to the buffer.
// Returns ErrFull if the packet doesn't fit.
//
// Note that the packet size is limited to 65536 bytes since v0.11.0 due
// to the internal data structure.
func (b *Buffer) Write(packet []byte) (int, error) { //nolint:cyclop
	// 65536 bytes
	if len(packet) >= 0x10000 {
		return 0, errPacketTooBig
	}

	b.mutex.Lock()
	if b.closed {
		b.mutex.Unlock()

		return 0, io.ErrClosedPipe
	}

	if (b.limitCount > 0 && b.count >= b.limitCount) || (b.limitSize > 0 && b.size()+2+len(packet) > b.limitSize) {
		b.mutex.Unlock()

		return 0, ErrFull
	}

	// grow the buffer until the packet fits
	for !b.available(len(packet)) {
		err := b.grow()
		if err != nil {
			b.mutex.Unlock()

			return 0, err
		}
	}

	// Store the length of the packet: Why Not Store the Length Once?
	//
	// Storing the length as two bytes (high byte and low byte) is necessary because:
	// - 16-bit Length:
	//   - The packet length is a 16-bit value, which requires 2 bytes to store.
	//   - Storing it as a single byte would limit the packet size to 255 bytes,
	//     which is insufficient for many use cases.
	// - Serialization: Serializing the length into two bytes allows the receiver to easily
	//   reconstruct the original 16-bit length value when reading the packet.
	// - Compatibility:
	//   This approach is consistent with how many protocols (e.g., TCP/IP, UDP)
	//   handle length fields in their headers.
	b.data[b.tail] = uint8(len(packet) >> 8) //nolint:gosec
	b.tail++
	if b.tail >= len(b.data) {
		b.tail = 0
	}

	b.data[b.tail] = uint8(len(packet)) //nolint:gosec
	b.tail++
	if b.tail >= len(b.data) {
		b.tail = 0
	}

	// store the packet
	n := copy(b.data[b.tail:], packet)
	b.tail += n
	if b.tail >= len(b.data) {
		// we reached the end, wrap around
		m := copy(b.data, packet[n:]) // copied to the beginning of the buffer
		b.tail = m                    //  the new position after wrapping around.
	}
	b.count++

	select {
	case b.notify <- struct{}{}:
	default:
	}

	b.mutex.Unlock()

	return len(packet), nil
}

// Read populates the given byte slice, returning the number of bytes read.
//
// Blocks until data is available or the buffer is closed. Returns io.ErrShortBuffer
//
// is the packet is too small to copy the Write. Returns io.EOF if the buffer is closed.
func (b *Buffer) Read(packet []byte) (n int, err error) { //nolint:cyclop
	// Return immediately if the deadline is already exceeded.
	select {
	case <-b.readDeadlien.Done():
		return 0, &netError{ErrTimeout, true, true}
	default:
	}

	for {
		b.mutex.Lock()

		if b.head != b.tail { //nolint:nestif
			// decode the packet size
			n1 := b.data[b.head]
			b.head++
			if b.head >= len(b.data) {
				b.head = 0
			}

			n2 := b.data[b.head]
			b.head++
			if b.head >= len(b.data) {
				b.head = 0
			}
			count := int((uint16(n1) << 8) | uint16(n2))

			// determine the number of bytes we'll actually copy
			copied := min(count, len(packet))

			// copy the data
			if b.head+copied < len(b.data) {
				// from current head position to the requested position (the len of the packet)
				copy(packet, b.data[b.head:b.head+copied])
			} else {
				k := copy(packet, b.data[b.head:])
				copy(packet[k:], b.data[:copied-k])
			}

			// advance head, discarding any data that wasn't copied
			b.head += count
			if b.head >= len(b.data) {
				b.head -= len(b.data)
			}

			if b.head == b.tail {
				// the buffer is empty, reset to beginning
				// in order to improve cache locality.
				b.head = 0
				b.tail = 0
			}

			b.count--
			b.mutex.Unlock()

			if copied < count {
				return copied, io.ErrShortBuffer
			}

			return copied, nil
		}

		if b.closed {
			b.mutex.Unlock()

			return 0, io.EOF
		}
		b.mutex.Unlock()

		select {
		case <-b.readDeadlien.Done():
			return 0, &netError{ErrTimeout, true, true}
		case <-b.notify:
		}
	}
}

// Close the buffer, unblocking any pending reads.
// Data in the buffer can still be read, Read will return io.EOF only when empty.
func (b *Buffer) Close() (err error) {
	b.mutex.Lock()

	if b.closed {
		b.mutex.Unlock()

		return nil
	}

	b.closed = true
	close(b.notify)
	b.mutex.Unlock()

	return nil
}

// Count returns the number of packets in the buffer.
func (b *Buffer) Count() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.count
}

// SetLimitCount controls the maximum number of packets that can be buffered.
// Causes Write to return ErrFull when this limit is reached.
// A zero value will disable this limit.
func (b *Buffer) SetLimitCount(limit int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.limitCount = limit
}

// Size returns the total byte size of packets in the buffer, including
// a small amount of administrative overhead.
func (b *Buffer) Size() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.size()
}

func (b *Buffer) size() int {
	size := b.tail - b.head
	if size < 0 {
		size += len(b.data)
	}

	return size
}

// SetLimitSize controls the maximum number of bytes that can be buffered.
// Causes Write to return ErrFull when this limit is reached.
// A zero value means 4MB since v0.11.0.
//
// User can set packetioSizeHardLimit build tag to enable 4MB hard limit.
// When packetioSizeHardLimit build tag is set, SetLimitSize exceeding
// the hard limit will be silently discarded.
func (b *Buffer) SetLimitSize(limit int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.limitSize = limit
}

// SetReadDeadline sets the deadline for the Read operation.
// Setting to zero means no deadline.
func (b *Buffer) SetReadDeadline(t time.Time) error {
	b.readDeadlien.Set(t)

	return nil
}
