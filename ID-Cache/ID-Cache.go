package idcache

// #cgo LDFLAGS: -L${SRCDIR} -lcache
// #include "bst.h"
// #include "queue.h"
import "C"
import "fmt"

var queueCount uint64 = 0
var bstCount uint64 = 0

func Insert(id uint64) bool {
	c_id := C.ulong(id)
	if uint64(C.BST_insert(c_id)) != 0 {
		C.enqueue(c_id)
		queueCount++
		bstCount++
		return true
	}
	return false
}

func TryAgainLater(id uint64) {
	queueCount++
	C.enqueue(C.ulong(id))
}

func Next() (uint64) {
	status := uint64(C.dequeue())
	if status == 0 {
		fmt.Println("Warning: Attempting to dequeue empty queue.")
	} else {
		queueCount--
	}
	return status
}

func Status() {
	fmt.Printf("IDs: Processed %d, Waiting %d, Total %d\n",
	bstCount - queueCount, queueCount, bstCount)
}

func QueueCount() uint64 {
	return queueCount
}
