package main

import (
	"encoding/json"
	"net/http"
)

func handlers() {

	routerPath("/api", httpGet, func(w http.ResponseWriter, r *http.Request) {
		// composite response body
		var res = map[string]string{"result": "succ", "name": "test"}
		response, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(path))))

}
