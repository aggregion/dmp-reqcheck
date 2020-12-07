// +build 386 amd64

package reports

// #cgo CFLAGS: -g -Wall
// #include "cpuid_x86.h"
import "C"

func getSgxInfo() (*sgxInfo, error) {
	var eax, ebx, ecx, edx C.uint

	eax = 1
	ecx = 0
	C.native_cpuid(&eax, &ebx, &ecx, &edx)
	smx := (ecx >> 6) & 1

	eax = 7
	ecx = 0
	C.native_cpuid(&eax, &ebx, &ecx, &edx)
	flc := ecx & (1 << 30)
	available := (ebx >> 2) & 0x1

	eax = 0x12
	ecx = 0
	C.native_cpuid(&eax, &ebx, &ecx, &edx)
	sgx1 := eax & 0x1
	sgx2 := (eax >> 1) & 0x1

	return &sgxInfo{
		Smx:       smx > 0,
		Flc:       flc > 0,
		Available: available > 0,
		Version1:  sgx1 > 0,
		Version2:  sgx2 > 0,
	}, nil
}
