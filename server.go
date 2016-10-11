package main

/*
#cgo CFLAGS: -I./libfaced -I/usr/include
#cgo LDFLAGS: -L./libfaced -lfaced
#include "faced.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	path := C.CString("face.jpg")
	defer C.free(unsafe.Pointer(path))

	resp := C.GoString(C.FaceDetection(path))

	fmt.Println(resp)
}
