package main

import (
	"log"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const (
	PIPE_ACCESS_INBOUND  = 0x00000001
	PIPE_ACCESS_OUTBOUND = 0x00000002
	PIPE_ACCESS_DUPLEX   = 0x00000003

	PIPE_WAIT                  = 0x00000000
	PIPE_NOWAIT                = 0x00000001
	PIPE_READMODE_BYTE         = 0x00000000
	PIPE_READMODE_MESSAGE      = 0x00000002
	PIPE_TYPE_BYTE             = 0x00000000
	PIPE_TYPE_MESSAGE          = 0x00000004
	PIPE_ACCEPT_REMOTE_CLIENTS = 0x00000000
	PIPE_REJECT_REMOTE_CLIENTS = 0x00000008

	PIPE_UNLIMITED_INSTANCES = 255

	NMPWAIT_WAIT_FOREVER     = 0xffffffff
	NMPWAIT_NOWAIT           = 0x00000001
	NMPWAIT_USE_DEFAULT_WAIT = 0x00000000
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	handle := syscall.NewLazyDLL("Kernel32.dll")
	if handle == nil {
		log.Println("handle")
		return
	}
	pfnCreateNamePipe := handle.NewProc("CreateNamedPipeA")
	if pfnCreateNamePipe == nil {
		log.Println("pfnCreateNamePipe")
		return
	}

	pipeName := []byte(`\\.\pipe\mynamedpipe`)
	hpipe, _, err := pfnCreateNamePipe.Call(uintptr(unsafe.Pointer(&pipeName[0])),
		PIPE_ACCESS_DUPLEX,
		PIPE_TYPE_MESSAGE|PIPE_READMODE_MESSAGE|PIPE_WAIT,
		PIPE_UNLIMITED_INSTANCES,
		4096, 4096, NMPWAIT_WAIT_FOREVER, uintptr(0))
	// 成功
	if err != nil && !strings.Contains(err.Error(), "successfully") {
		log.Println("CreateNamePipe")
		return
	}

	count := 1
	w := &sync.WaitGroup{}
	w.Add(count)
	for i := 0; i < count; i++ {
		go writeRoutine(w)
	}

	//接收显示
	buf := make([]byte, 4096)
	var done uint32 = 0
	for {
		err = syscall.ReadFile(syscall.Handle(hpipe), buf, &done, nil)
		if err != nil {
			log.Println("ReadFile", err)

		}
		log.Println(string(buf[:done]), done)
	}

	w.Wait()
}

func writeRoutine(w *sync.WaitGroup) {
	defer w.Done()

	path, _ := syscall.UTF16PtrFromString(`\\.\pipe\mynamedpipe`)
	var hpipe syscall.Handle
	var err error = nil
	for {
		hpipe, err = syscall.CreateFile(path, syscall.GENERIC_WRITE,
			syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_READ, nil, syscall.OPEN_EXISTING, 0, 0)
		if err != nil {
			log.Println("CreateFile", err)
			time.Sleep(time.Second)
			continue
		}
		log.Println(hpipe)
		break
	}

	for i := 0; i < 100000; i++ {
		var done uint32 = 0
		err = syscall.WriteFile(hpipe, []byte("hello"), &done, nil)
		if err != nil {
			log.Println("WriteFile", err)
		}
		time.Sleep(time.Second)
	}

}
