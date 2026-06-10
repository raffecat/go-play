package main

/*
#cgo LDFLAGS: -L../rust-play/target/debug -lplay_lib
#include "../rust-play/target/rust-play.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

func main() {
	SendBuffer()
	SendString()
	BorrowBuffer()
	result := C.add_numbers(10, 20)
	fmt.Println("Result from Rust:", result)
	runtime.GC()
}

func SendBuffer() {
	buf := []byte{1, 2, 3, 4, 5}
	// Pass the pointer to the first element and the length.
	// Pinning: this memory is safely pinned for the duration of this call.
	C.process_buffer((*C.uint8_t)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
}

func SendString() {
	msg := "Hello from Go!"
	// Pass the pointer to the string data and the length.
	// Pinning: this memory is safely pinned for the duration of this call.
	C.print_go_string((*C.uint8_t)(unsafe.Pointer(unsafe.StringData(msg))), C.size_t(len(msg)))
}

type RustBuf struct {
	ptr []byte
	raw *C.uint8_t
	len C.size_t
}

func BorrowBuffer() *RustBuf {
	var length C.size_t

	// 1. Get the pointer from Rust
	ptr := C.get_rust_buffer((*C.size_t)(unsafe.Pointer(&length)))

	// 2. Wrap the pointer in a Go slice (No copy happens here!)
	// Note: In Go 1.21+, use unsafe.Slice(ptr, length)
	buf := &RustBuf{ptr: unsafe.Slice((*byte)(ptr), int(length)), raw: ptr, len: length}

	// 3. Attach the "destructor"
	// CAUTION: does not fire unless RustBuf is actully returned.
	runtime.SetFinalizer(buf, func(r *RustBuf) {
		fmt.Println("Finalizer: Freeing Rust string memory")
		C.free_rust_buffer(r.raw, r.len)
	})

	fmt.Printf("Go received: %s\n", string(buf.ptr))

	return buf
}
