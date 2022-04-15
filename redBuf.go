package sc

import (
	// "fmt"
	"sync"
	"time"

	g "github.com/Yessir-kim/sessionController/var"
)

type rebuffer struct {
	mutex sync.Mutex
	rob []byte
	queue [][]byte
	mark []bool
	idx int
	seq int
	total int
}

func (b *rebuffer) store(buf []byte, seq int, total int) bool {

		b.mutex.Lock()

		b.total = total

		if b.seq != seq {
			idx := seq - (b.seq + 1)
			// fmt.Printf("\tstore the data in queue[%d]\n", idx)

			b.mark[idx] = true
			b.queue[idx] = append([]byte(nil), buf...)
			// copy(b.queue[idx], buf[:len(buf)])

			b.mutex.Unlock()

			return true
		}

		temp := b.rob[:b.idx]
		b.rob = append(temp, buf...)
		b.idx += len(buf)
		b.seq++

		// dequeue
		for i := 0; i < g.SLICE_SIZE; i++ {
			if b.mark[i] { // true
				temp := b.rob[:b.idx]
				b.rob = append(temp, b.queue[i]...)
				b.idx += len(b.queue[i])
				b.seq++
			} else { // false
				q := make([][]byte, g.QUEUE_SIZE)
				for i := 0; i < g.QUEUE_SIZE; i++ {
					q[i] = make([]byte, g.PAYLOAD_SIZE)
				}
				copy(q, b.queue[i + 1:])
				b.queue = nil // free
				b.queue = q

				m := make([]bool, g.QUEUE_SIZE)
				copy(m, b.mark[i + 1:])
				b.mark = nil // free
				b.mark = m

				break
			}

		}

		// fmt.Printf("Successful storing rob[:%d]\n", b.idx)

		b.mutex.Unlock()

		return true
}

func (b *rebuffer) read(buf []byte) (int) {

		// The remaining data that has not been received from the sender exists  
		for b.total != b.idx {} // need to fix to make more efficient 

		b.mutex.Lock()

		size := 0

		if b.idx < len(buf) { // less than buffer 
			size = b.idx
			copy(buf, b.rob[:b.idx])
			b.idx = 0
		} else { // greater than buffer
			size = len(buf)
			copy(buf, b.rob[:len(buf)])
			b.idx -= len(buf) // maybe zero... 
		}

		if b.idx == 0 {
			// To keep the underlying array, slice the 'slice' to zero length.
			b.rob = b.rob[:0]
		} else {
			temp := make([]byte, g.PAYLOAD_SIZE * g.SLICE_SIZE)
			copy(temp, b.rob[len(buf):])
			b.rob = nil // free
			b.rob = temp
		}
		/*
		fmt.Printf("Successful reading\n")
		fmt.Printf("\tRemaining reorder buf size : %d\n", cap(b.rob))
		fmt.Printf("\tIdx: %d\n", b.idx)
		*/
		b.mutex.Unlock()

		return size
}

func New() *rebuffer {
		rob := rebuffer{
			rob: make([]byte, g.PAYLOAD_SIZE * g.SLICE_SIZE),
			queue: make([][]byte, g.QUEUE_SIZE),
			mark: make([]bool, g.QUEUE_SIZE),
			idx: 0,
			seq: 0,
			total: 0,
		}

		for i := 0; i < g.QUEUE_SIZE; i++ {
			rob.queue[i] = make([]byte, g.PAYLOAD_SIZE)
		}

		return &rob
}
