package bingo_errors

import "errors"

var (
	ErrCmdNotSupported  = errors.New("command is not currently supported")
	ErrTcpPortInUse     = errors.New("tcp port specified is unavailable")
	ErrNameInUse        = errors.New("name already in use")
	ErrStartingListener = errors.New("error starting listener")
	ErrInvalidRHOST     = errors.New("must specify rhost when listener listens on all interfaces")
	ErrInvalidListener  = errors.New("invalid listener")
	ErrInvalidOS        = errors.New("invalid os specified")
	ErrInvalidArch      = errors.New("invalid arch specified")
	ErrDeadSession      = errors.New("session has expired")
	ErrInvalidPath      = errors.New("invalid path specified")
)
