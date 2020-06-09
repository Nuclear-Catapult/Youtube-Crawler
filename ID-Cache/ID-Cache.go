package idcache

// #cgo LDFLAGS: -L${SRCDIR} -lcache
// #include "bst.h"
// #include "queue.h"
import "C"
import b64 "github.com/Nuclear-Catapult/NCcrawl-YT/ID-Cache/ytbase64"
import "fmt"

var queueCount uint64 = 0
var bstCount uint64 = 0

func Insert(str_id string) {
	int_id := C.ulong(b64.Decode64(str_id))
	if int(C.BST_insert(int_id)) != 0 {
		fmt.Println("Successful insertion")
		queueCount++
		bstCount++
//		if (bstCount - queueCount) % 100 == 0 {
//			fmt.Printf("Processed: %d, Waiting: %d, Total %d\n", bstCount - queueCount, queueCount, bstCount)
//		}
		C.enqueue(int_id)
	} else {
		fmt.Println("duplicate found")
	}
}

func Next() string {
	int_id := uint64(C.dequeue())
	if int_id == 0 {
		return ""
	}
	queueCount--
	return b64.Encode64(int_id)
}
