package hg

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type FILE_RENAME_INFO struct {
	Union struct {
		ReplaceIfExists bool
		Flags           uint32
	}
	RootDirectory  windows.Handle
	FileNameLength uint32
	FileName       [1]uint16
}

type FILE_DELETE_INFO struct {
	DeleteFile bool
}

func openHandle(pwPath *uint16) (windows.Handle, error) {
	handle, err := windows.CreateFile(
		pwPath,
		windows.DELETE,
		0,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)

	if err != nil {
		return 0, err
	}

	return handle, nil
}

func renameHandle(hHandle windows.Handle, newName string) error {
	var fRename FILE_RENAME_INFO
	DS_STREAM_RENAME, err := windows.UTF16FromString(newName)
	if err != nil {
		return err
	}

	lpwStream := &DS_STREAM_RENAME
	fRename.FileNameLength = uint32(unsafe.Sizeof(lpwStream))

	windows.NewLazyDLL("kernel32.dll").NewProc("RtlCopyMemory").Call(
		uintptr(unsafe.Pointer(&fRename.FileName[0])),
		uintptr(unsafe.Pointer(lpwStream)),
		unsafe.Sizeof(lpwStream),
	)

	err = windows.SetFileInformationByHandle(
		hHandle,
		windows.FileRenameInfo,
		(*byte)(unsafe.Pointer(&fRename)),
		uint32(unsafe.Sizeof(fRename)+unsafe.Sizeof(lpwStream)),
	)

	return err
}

func disposeProcess(hHandle windows.Handle) error {
	var fDel FILE_DELETE_INFO
	fDel.DeleteFile = true

	err := windows.SetFileInformationByHandle(
		hHandle,
		windows.FileDispositionInfo,
		(*byte)(unsafe.Pointer(&fDel)),
		uint32(unsafe.Sizeof(fDel)),
	)

	return err
}

func Suicide() error {
	var wcPath [windows.MAX_PATH + 1]uint16
	var hCurrentProcess windows.Handle

	// get module file name
	_, err := windows.GetModuleFileName(0, &wcPath[0], windows.MAX_PATH)
	if err != nil {
		return err
	}

	// open handle on current process
	hCurrentProcess, err = openHandle(&wcPath[0])
	if err != nil || hCurrentProcess == windows.InvalidHandle {
		return err
	}

	// rename handle of current process to random string of 6 characters
	newHandleName := RandomStr(6)
	err = renameHandle(hCurrentProcess, newHandleName)
	if err != nil {
		return nil
	}

	// close handle
	windows.CloseHandle(hCurrentProcess)

	// reopen handle
	hCurrentProcess, err = openHandle(&wcPath[0])
	if err != nil || hCurrentProcess == windows.InvalidHandle {
		windows.CloseHandle(hCurrentProcess)
		return err
	}

	// despose of process
	err = disposeProcess(hCurrentProcess)
	if err != nil {
		return err
	}

	return nil
}
