package idcache

// #cgo LDFLAGS: -L${SRCDIR} -lcache -pthread
// #include "hash_table.h"
// #include "queue.h"
import "C"
import "fmt"
import "log"

var queueCount int64 = 0
var tableCount int64 = 0

func Insert(id int64) bool {
	c_id := C.long(id)
	if int64(C.table_insert(c_id)) != 0 {
		C.enqueue(c_id)
		queueCount++
		tableCount++
		return true
	}
	return false
}

func TryAgainLater(id int64) {
	queueCount++
	C.enqueue(C.long(id))
}

func Next() int64 {
	status := int64(C.dequeue())
	if status == 0 {
		fmt.Println("Warning: Attempting to dequeue empty queue.")
	} else {
		queueCount--
	}
	return status
}

func Key_Insert(id int64) {
	c_id := C.long(id)
	if int64(C.table_insert(c_id)) == 0 {
		log.Fatal("Fatal: Loading the primary keys resulted in duplicate values.")
	}
	tableCount++
}

func Status() {
	fmt.Printf("IDs: Processed %d, Waiting %d, Total Cached %d\n",
		tableCount-queueCount, queueCount, tableCount)
}

func QueueCount() int64 {
	return queueCount
}
