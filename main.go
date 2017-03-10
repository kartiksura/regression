package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
)

func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func handler(w http.ResponseWriter, r *http.Request) {
	digitalAdID := r.URL.Query().Get("ad_id")
	media := r.URL.Query().Get("media")
	geoState := r.URL.Query().Get("geoState")

	fmt.Fprintln(w, Predict(digitalAdID, media, geoState))
}

func main() {

	var port, path string
	flag.StringVar(&port, "p", "8000", "specify port to use.  defaults to 8000.")
	flag.StringVar(&path, "f", "input.csv", "file to be loaded. defaults to input.csv")

	flag.Parse()

	fmt.Println("starting on port = ", port)
	fmt.Println("reading file = ", path)
	readFile(path)

	http.HandleFunc("/sr", handler)
	http.ListenAndServe(":"+port, nil)
}
