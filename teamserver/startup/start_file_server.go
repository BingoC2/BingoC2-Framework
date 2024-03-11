package startup

import (
	"fmt"
	"net/http"
)

func StartFileServer(lhost string, lport string, binDir string) {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(binDir))
	mux.Handle("/bin/", http.StripPrefix("/bin", fileServer))

	fmt.Println("Starting bin file server on 0.0.0.0:4456")

	http.ListenAndServe("0.0.0.0:4456", mux)
}
