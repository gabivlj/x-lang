package http

import (
	"encoding/json"
	"net/http"
	"xlang/runtime"

	"github.com/julienschmidt/httprouter"
)

type requestRun struct {
	Code string `json:"code"`
}

// RunServer runs the http server of the xlang runtime
func RunServer() {
	router := httprouter.New()
	router.OPTIONS("/api/v1", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Write([]byte("nice"))
	})
	router.PUT("/api/v1", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		var body requestRun
		json.NewDecoder(r.Body).Decode(&body)
		output := runtime.Parse(body.Code)
		output.Print()
		res := map[string]interface{}{"data": output, "status": 200}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		bytes, _ := json.Marshal(&res)

		w.Write(bytes)
	})

	router.ServeFiles("/*filepath", http.Dir("/"))

	http.ListenAndServe(":8080", router)
}
