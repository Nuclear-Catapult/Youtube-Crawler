package idcache

// #cgo LDFLAGS: -L${SRCDIR} -lcache
// #include "bst.h"
// #include "queue.h"
import "C"
import "fmt"

var queueCount uint64 = 0
var bstCount uint64 = 0

func Insert(id uint64) {
	c_id := C.ulong(id)
	if uint64(C.BST_insert(c_id)) != 0 {
		queueCount++
		bstCount++
		C.enqueue(c_id)
	}
}

func TryAgainLater(id uint64) {
	queueCount++
	C.enqueue(C.ulong(id))
}

func Next() (uint64) {
	queueCount--
	return uint64(C.dequeue())
}

func status() {
	fmt.Printf("Processed: %d, Waiting: %d, Total %d\n", bstCount - queueCount, queueCount, bstCount)
}
