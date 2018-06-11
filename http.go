package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlers() {

	routerPath("/api", httpGet, func(w http.ResponseWriter, r *http.Request) {
		log.Print("in")
		// composite response body
		var res = map[string]string{"result": "succ", "name": "test"}
		response, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

}
