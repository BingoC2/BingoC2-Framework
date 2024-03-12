package startup

import (
	"fmt"
	"net/http"
)

func StartBinServer(lhost string) {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./bin"))
	mux.Handle("/bin/", http.StripPrefix("/bin", fileServer))

	fmt.Println("Starting bin file server on 0.0.0.0:4456")

	http.ListenAndServe("0.0.0.0:4456", mux)
}

func StartLootServer(lhost string) {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./loot"))
	mux.Handle("/loot/", http.StripPrefix("/loot", fileServer))

	fmt.Println("Starting loot file server on 0.0.0.0:4457")

	http.ListenAndServe("0.0.0.0:4457", mux)
}

func StartFileServer(lhost string) {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./files"))
	mux.Handle("/files/", http.StripPrefix("/files", fileServer))

	fmt.Println("Starting loot file server on 0.0.0.0:4458")

	http.ListenAndServe("0.0.0.0:4458", mux)
}
