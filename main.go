// package main

// import (
// 	"juicer/cmd"
// 	"log"
// )

// func main() {
// 	if err := cmd.Run(); err != nil {
// 		log.Fatalf("failed to run chess server")
// 	}
// }

package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("lol 123")); err != nil {
			fmt.Printf("errored")
			return
		}
	})

	http.HandleFunc("/api/v1/health/alive", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("lol 123")); err != nil {
			fmt.Printf("errored")
			return
		}
	})

	if err := http.ListenAndServe(":1337", nil); err != nil {
		fmt.Printf("Err: %+v\n", err)
	}
}
