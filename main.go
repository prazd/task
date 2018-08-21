package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"./mongo"
	s "./sett"
)

func main() {
	http.HandleFunc("/inv/up", InvSignUp)
	http.HandleFunc("/inv/in", InvSignIn)
	http.HandleFunc("/vol/up", VolSignUp)
	http.HandleFunc("/vol/in", VolSignIn)
	http.HandleFunc("/vol/exit", VolEx)
	http.HandleFunc("/inv/exit", InvEx)
	http.HandleFunc("/geolist", GeoList)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func InvSignUp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var inv s.InvUser
	err = json.Unmarshal(body, &inv)
	if err != nil {
		log.Fatal(err)
	}

	up := mongo.InvSup(inv.Id, inv.Name, inv.Number, inv.Password)

	resp := struct {
		Resp string `json:"resp"`
	}{up}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Fatal(bad)
	}
	w.Write(js)
}

func InvSignIn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var inv s.InvUser
	err = json.Unmarshal(body, &inv)
	if err != nil {
		log.Fatal(err)
	}
	name, sIn := mongo.InvSin(inv.Id, inv.Password)

	resp := struct {
		Resp string `json:"resp"`
		Name string `json:"name"`
	}{sIn, name}
	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Fatal(bad)
	}
	w.Write(js)
}

func VolSignUp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var vol s.VolUser
	err = json.Unmarshal(body, &vol)
	if err != nil {
		log.Fatal(err)
	}

	up := mongo.VolSup(vol.Name, vol.Number, vol.Password, vol.Geo)

	resp := struct {
		Resp string `json:"resp"`
	}{up}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Fatal(bad)
	}
	w.Write(js)

}

func VolSignIn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var vol s.VolUser
	err = json.Unmarshal(body, &vol)
	if err != nil {
		log.Fatal(err)
	}
	name, sIn := mongo.VolSin(vol.Number, vol.Password)

	resp := struct {
		Resp string `json:"resp"`
		Name string `json:"name"`
	}{sIn, name}
	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Fatal(bad)
	}
	w.Write(js)

}

func VolEx(w http.ResponseWriter, r *http.Request) {

}

func InvEx(w http.ResponseWriter, r *http.Request) {

}

func GeoList(w http.ResponseWriter, r *http.Request) {
	geolist := mongo.GetGeoV()
	resp := struct {
		Resp [][]string `json:"resp"`
	}{Resp: geolist}
	js, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(js)
}
