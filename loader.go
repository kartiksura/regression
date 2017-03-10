package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/sajari/regression"
)

var mdAdID map[string]float64
var mMedia map[string]float64
var mState map[string]float64
var r *regression.Regression

func load(digitalAdID string, media string, state string, req string, resp string) {
	var fdAdID float64
	var fMedia float64
	var fState float64
	ok := false

	reqs, err := strconv.ParseFloat(req, 64)
	if err != nil {
		fmt.Println("Error parsing req:", req)
	}
	resps, err := strconv.ParseFloat(resp, 64)
	if err != nil {
		fmt.Println("Error parsing req:", resp)
	}
	var fill float64 = 1.0
	if reqs > 0 {
		fill = resps / reqs
	}

	if fdAdID, ok = mdAdID[digitalAdID]; !ok {
		mdAdID[digitalAdID] = (float64)(len(mdAdID) + 1)
		fdAdID = mdAdID[digitalAdID]
	}
	if fMedia, ok = mMedia[media]; !ok {
		mMedia[media] = (float64)(len(mMedia) + 1)
		fMedia = mMedia[media]
	}
	if fState, ok = mState[state]; !ok {
		mState[state] = (float64)(len(mState) + 1)
		fState = mState[state]
	}
	r.Train(regression.DataPoint(fill, []float64{fdAdID, fMedia, fState}))

}

func readFile(fileName string) {
	csvfile, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// sanity check, display to standard output
	for i, each := range rawCSVdata {
		fmt.Println("loading:", each[0], each[1], each[2], each[3], each[4])
		if i > 0 && len(each) == 5 {
			load(each[0], each[1], each[2], each[3], each[4])
		}

	}
	r.Run()
	fmt.Printf("Regression formula:\n%v\n", r.Formula)
	fmt.Printf("Regression:\n%s\n", r)

}

//Predict ..
func Predict(digitalAdID string, media string, state string) float32 {
	var fdAdID float64
	var fMedia float64
	var fState float64

	ok := false
	if fdAdID, ok = mdAdID[digitalAdID]; !ok {
		return 1.0
	}
	if fMedia, ok = mMedia[media]; !ok {
		return 1.0
	}
	// if fState, ok = mState[state]; ok {
	// 	return 1.0
	// }

	prediction, err := r.Predict([]float64{fdAdID, fMedia, fState})
	if err != nil {
		fmt.Println(err)
		return 1.0
	}
	return float32(prediction)
}

func init() {
	mdAdID = make(map[string]float64)
	mMedia = make(map[string]float64)
	mState = make(map[string]float64)

	r = new(regression.Regression)
	r.SetObserved("Fill rate prediction")
	r.SetVar(0, "Digital Ad ID")
	r.SetVar(1, "Media ID")
	r.SetVar(2, "Geo State")

}
