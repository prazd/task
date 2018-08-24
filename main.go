package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"./mongo"
	s "./sett"
	"gopkg.in/mgo.v2/bson"
)

func main() {

	http.HandleFunc("/inv/up", InvSignUp)
	http.HandleFunc("/inv/in", InvSignIn)
	http.HandleFunc("/vol/up", VolSignUp)
	http.HandleFunc("/vol/in", VolSignIn)
	http.HandleFunc("/vol/ex", VolExit)
	http.HandleFunc("/inv/ex", InvExit)
	http.HandleFunc("/vol/geolist", VGeoList)
	http.HandleFunc("/inv/geolist", IGeoList)
	http.HandleFunc("/vol/getrev", GetVRev)
	http.HandleFunc("/vol/chrev", ChangeVRev)
	http.HandleFunc("/vol/ch", VolHelp)
	http.HandleFunc("/inv/nh", InvHelp)
	http.HandleFunc("/findhelp", FHelp)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func InvSignUp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var inv s.InvUser
	err = json.Unmarshal(body, &inv)
	if err != nil {
		log.Println(err)
	}

	up := mongo.InvSup(inv.Id, inv.Name, inv.Number, inv.Password)

	resp := struct {
		Resp string `json:"resp"`
	}{up}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
	}
	w.Write(js)
}

func InvSignIn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var inv s.InvUser
	err = json.Unmarshal(body, &inv)
	if err != nil {
		log.Println(err)
	}
	sIn, name, number := mongo.InvSin(inv.Id, inv.Password)

	resp := struct {
		Resp   string `json:"resp"`
		Name   string `json:"name"`
		Number string `json:"number"`
	}{sIn, name, number}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
	}
	w.Write(js)
}

func VolSignUp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	var vol s.VolUser
	err = json.Unmarshal(body, &vol)
	if err != nil {
		log.Println(err)
	}

	up := mongo.VolSup(vol.Name, vol.Number, vol.Password)

	resp := struct {
		Resp string `json:"resp"`
	}{up}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
	}
	w.Write(js)

}

func VolSignIn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var vol s.VolUser
	err = json.Unmarshal(body, &vol)
	if err != nil {
		log.Println(err)
	}
	sIn, name := mongo.VolSin(vol.Number, vol.Password)

	resp := struct {
		Resp string `json:"resp"`
		Name string `json:"name"`
	}{sIn, name}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
	}
	w.Write(js)

}

func VolHelp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	type canHelp struct {
		Number    string
		Longitude string
		Latitude  string
	}

	var ch canHelp

	err = json.Unmarshal(body, &ch)
	if err != nil {
		log.Println(err)
	}

	help := mongo.VHelp(ch.Number, ch.Latitude, ch.Longitude)
	resp := struct {
		Resp string `json:"resp"`
	}{strconv.FormatBool(help)}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
	}
	w.Write(js)
}

func InvHelp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	type needHelp struct {
		Id        string
		Longitude string
		Latitude  string
	}

	var nh needHelp

	err = json.Unmarshal(body, &nh)
	if err != nil {
		log.Println(err)
	}

	help := mongo.IHelp(nh.Id, nh.Latitude, nh.Longitude)
	resp := struct {
		Resp string `json:"resp"`
	}{strconv.FormatBool(help)}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
	}
	w.Write(js)
}

func VolExit(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	var vol s.VolUser
	err = json.Unmarshal(body, &vol)
	if err != nil {
		log.Println(err)
	}

	help := mongo.VolEx(vol.Number)
	resp := struct {
		Resp string `json:"resp"`
	}{strconv.FormatBool(help)}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
	}
	w.Write(js)
}

func InvExit(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	var inv s.InvUser
	err = json.Unmarshal(body, &inv)
	if err != nil {
		log.Println(err)
	}

	help := mongo.InvEx(inv.Id)
	resp := struct {
		Resp string `json:"resp"`
	}{strconv.FormatBool(help)}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
	}
	w.Write(js)
}

func VGeoList(w http.ResponseWriter, r *http.Request) {
	geolist := mongo.GetGeoV()
	resp := struct {
		Resp [][]string `json:"resp"`
	}{Resp: geolist}
	js, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	w.Write(js)
}

func IGeoList(w http.ResponseWriter, r *http.Request) {
	geolist := mongo.GetGeoI()
	resp := struct {
		Resp [][]string `json:"resp"`
	}{Resp: geolist}
	js, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	w.Write(js)
}

func GetVRev(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var vol s.VolUser
	err = json.Unmarshal(body, &vol)
	if err != nil {
		log.Println(err)
	}

	good, bad := mongo.GetVolReviews(vol.Number)

	resp := struct {
		Good int `json:"goodrev"`
		Bad  int `json:"badreviews"`
	}{good, bad}

	js, excp := json.Marshal(resp)
	if excp != nil {
		log.Println(err)
	}
	w.Write(js)
}

func ChangeVRev(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	type Req struct {
		Number string `json:"number"`
		Review string `json:"review"`
	}
	var req Req
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println(err)
	}
	rev := mongo.ChangeVReview(req.Number, req.Review)
	js, excp := json.Marshal(bson.M{"resp": rev})
	if excp != nil {
		log.Println(err)
	}
	w.Write(js)
}

func FHelp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	type Help struct {
		Id     string
		Number string
	}
	var help Help
	err = json.Unmarshal(body, &help)
	if err != nil {
		log.Println(err)
	}
	resp, rVol, rInv := mongo.FindHelp(help.Id, help.Number)

	js, exc := json.Marshal(bson.M{"resp": resp, "inv": rInv, "vol": rVol})
	if exc != nil {
		log.Println(err)
	}
	w.Write(js)
}
