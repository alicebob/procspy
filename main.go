package main

/*

#cgo linux CFLAGS: -DLINUXV=313007 -DGLIBCV=218 -DHASIPv6 -D_FILE_OFFSET_BITS=64 -D_LARGEFILE64_SOURCE -DHAS_STRFTIME -DLSOF_VSTR="3.13.7" -O
#cgo darwin CFLAGS: -DHASIPv6 -DHASUTMPX -DDARWINV=1100 -DHAS_STRFTIME -DLSOF_VSTR="13.1.0" -O
#cgo LDFLAGS: -L ./ -llsof


#include "lsof.h"

void
lsof_init()
{
    if (!(Namech = (char *)malloc(MAXPATHLEN + 1))) {
        (void) fprintf(stderr, "%s: no space for name buffer\n", Pn);
        Exit(1);
    }
	Namechl = (size_t)(MAXPATHLEN + 1);

	Fnet = 1;
	FnetTy = 0;
	Selflags = 0;
	Selflags |= SELNET;
	Selinet = 1;
	Selall = 0;
	Fhost = 0;
	Fport = 0;
	initialize();
	(void) hashSfile();
}
*/
import "C"

import (
	"fmt"
	"net"
	"time"
	"unsafe"
)

type ConnProc struct {
	Protocol  string
	LocalAddr net.IP
	LocalPort int
	PID       int
	Name      string
}

func init() {
	C.lsof_init()
}

func lsof() []ConnProc {
	C.gather_proc_info()

	res := []ConnProc{}
	var p C.struct_lproc
	var i int
	for i = 0; i < int(C.Nlproc); i++ {
		// Lproc is a pointer to NULL, Go can't make it an array, it seems.
		ptr := uintptr(unsafe.Pointer(C.Lproc)) + unsafe.Sizeof(p)*uintptr(i)
		myp := (*C.struct_lproc)(unsafe.Pointer(ptr))
		if myp.pss > 0 {
			for lf := myp.file; lf != nil; lf = lf.next {
				// Only use files with 2 (a)ddress(f)amilies.
				// Listens will only have one.
				if lf.li[0].af == 0 || lf.li[1].af == 0 {
					continue
				}
				var address net.IP
				if lf.li[0].af == C.AF_INET {
					address = net.IP(C.GoBytes(unsafe.Pointer(&lf.li[0].ia[0]), 4))
				}
				if lf.li[0].af == C.AF_INET6 {
					address = net.IP(C.GoBytes(unsafe.Pointer(&lf.li[0].ia[0]), 4*4))
				}
				res = append(res, ConnProc{
					Name:      C.GoString(myp.cmd),
					PID:       int(myp.pid),
					Protocol:  C.GoString(&lf.iproto[0]),
					LocalPort: int(lf.li[0].p),
					LocalAddr: address,
				})
			}
		}
		C.free_lproc(myp)
	}
	C.Nlproc = 0 // Needed to prevent a memory leak.
	return res
}

func main() {
	// HASNCACHE NcacheReload = 1
	// C.print_lsof()
	for {
		fmt.Printf("Go Go:\n")
		for _, p := range lsof() {
			fmt.Printf(" - %v\n", p)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
