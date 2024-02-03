package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	key      = -1488
	dataSize = 256
	lockPos  = dataSize
	size     = dataSize + 8
)

func main() {
	if len(os.Args) == 0 {
		log.Fatalf("No args")
	}
	cmd := os.Args[1]
	switch cmd {
	case "write":
		if len(os.Args) != 4 {
			log.Fatalf("Use: %s write <msg> <pos>", os.Args[0])
		}
		pos, err := strconv.ParseUint(os.Args[3], 10, 64)
		if err != nil {
			log.Fatalf("Invalid pos: %v", err)
		}
		write(os.Args[2], pos)
	case "write-atomic":
		if len(os.Args) != 4 {
			log.Fatalf("Use: %s write <msg-max-len-8> <pos>", os.Args[0])
		}
		msg := []byte(os.Args[2])
		if len(msg) > 8 {
			log.Fatalf("Too long msg")
		}
		pos, err := strconv.ParseUint(os.Args[3], 10, 64)
		if err != nil {
			log.Fatalf("Invalid pos: %v", err)
		}
		writeAtomic(msg, pos)
	case "view":
		if len(os.Args) != 2 {
			log.Fatalf("Use: %s view", os.Args[0])
		}
		view()
	case "delete":
		if len(os.Args) != 2 {
			log.Fatalf("Use: %s delete", os.Args[0])
		}
		del()
	}
}

func write(msg string, pos uint64) {
	id, err := unix.SysvShmGet(key, size, 0)
	if err == unix.ENOENT {
		id, err = unix.SysvShmGet(key, size, 0600|unix.IPC_CREAT)
	}
	if err != nil {
		log.Fatalf("Get shared memory: %v", err)
	}

	mem, err := unix.SysvShmAttach(id, 0, 0)
	if err != nil {
		log.Fatalf("Attach to shared memory: %v", err)
	}

	copy(mem[pos:], msg)

	err = unix.SysvShmDetach(mem)
	if err != nil {
		log.Fatalf("Detach from shared memory: %v", err)
	}

	log.Print("Written")
}

func writeAtomic(msg []byte, pos uint64) {
	id, err := unix.SysvShmGet(key, size, 0)
	if err == unix.ENOENT {
		id, err = unix.SysvShmGet(key, size, 0600|unix.IPC_CREAT)
	}
	if err != nil {
		log.Fatalf("Get shared memory: %v", err)
	}

	mem, err := unix.SysvShmAttach(id, 0, 0)
	if err != nil {
		log.Fatalf("Attach to shared memory: %v", err)
	}

	lockUintPtr := (*uint64)(unsafe.Pointer(&mem[lockPos]))

	for i := 0; !atomic.CompareAndSwapUint64(lockUintPtr, 0, 1); i++ {
		if i > 0 && (i == 10 || i%100 == 0) {
			log.Printf("Spin locking... iteration #%d", i)
		}
	}

	log.Print("Lock acquired")

	posUintPtr := (*uint64)(unsafe.Pointer(&mem[pos]))
	oldMsgUint := *posUintPtr

	oldMsg := make([]byte, 8)
	binary.LittleEndian.PutUint64(oldMsg, oldMsgUint)
	newMsg := make([]byte, 8)
	for i := 0; i < len(msg); i++ {
		newMsg[i] = msg[i]
	}
	for i := len(msg); i < 8; i++ {
		newMsg[i] = oldMsg[i]
	}
	newMsgUint := binary.LittleEndian.Uint64(newMsg)

	*posUintPtr = newMsgUint

	atomic.StoreUint64(lockUintPtr, 0)

	log.Print("Lock released")

	err = unix.SysvShmDetach(mem)
	if err != nil {
		log.Fatalf("Detach from shared memory: %v", err)
	}

	log.Print("Written")
}

func view() {
	id, err := unix.SysvShmGet(key, size, 0600)
	if err != nil && err != unix.EEXIST {
		log.Fatalf("Get shared memory: %v", err)
	}

	mem, err := unix.SysvShmAttach(id, 0, 0)
	if err != nil {
		log.Fatalf("Attach to shared memory: %v", err)
	}

	cols := 64
	lines := dataSize / cols
	for l := 0; l < lines; l++ {
		start := l * cols
		end := start + cols
		fmt.Print("[")
		for c := start; c < end; c++ {
			char := mem[c]
			if char == 0 {
				char = '.'
			}
			fmt.Printf("%c", char)
		}
		fmt.Println("]")
	}
	lockMark := ' '
	for i := lockPos; i < len(mem); i++ {
		if mem[i] != 0 {
			lockMark = 'X'
		}
	}
	fmt.Printf("LOCK: [%c]\n", lockMark)

	err = unix.SysvShmDetach(mem)
	if err != nil {
		log.Fatalf("Detach from shared memory: %v", err)
	}
}

func del() {
	id, err := unix.SysvShmGet(key, size, 0600)
	if err != nil {
		log.Fatalf("Get shared memory: %v", err)
	}

	_, err = unix.SysvShmCtl(id, unix.IPC_RMID, nil)
	if err != nil {
		log.Fatalf("Delete shared memory: %v", err)
	}

	log.Print("Deleted")
}
