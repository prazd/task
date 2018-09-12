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

	http.HandleFunc("/inv/up", PostOnly(InvSignUp))     // Sign up (inv)
	http.HandleFunc("/inv/in", PostOnly(InvSignIn))     // Sign in (inv), get inv info
	http.HandleFunc("/vol/up", PostOnly(VolSignUp))     // Sign up (vol)
	http.HandleFunc("/vol/in", PostOnly(VolSignIn))     // Sign in (vol), get vol info
	http.HandleFunc("/vol/ex", PostOnly(VolExit))       // Exit (vol)
	http.HandleFunc("/inv/ex", PostOnly(InvExit))       // Exit (inv)
	http.HandleFunc("/vol/geolist", GetOnly(VGeoList))  // Get geolocation and user info (vol)
	http.HandleFunc("/inv/geolist", GetOnly(IGeoList))  // Get geolocation and user info (inb)
	http.HandleFunc("/vol/getrev", PostOnly(GetVRev))   // Get reviews about vol
	http.HandleFunc("/vol/chrev", PostOnly(ChangeVRev)) // evaluate vol
	http.HandleFunc("/vol/ch", PostOnly(VolHelp))       // (Vol) canHelp - set state(1) and geolocation
	http.HandleFunc("/inv/nh", PostOnly(InvHelp))       // (Inv) needHelp - set state(1) and geolocation
	http.HandleFunc("/findhelp", PostOnly(FHelp))       // ...
	http.HandleFunc("/vol/gp", PostOnly(VolGP))         // Set Vol(geo)
	http.HandleFunc("/vol/help", PostOnly(VolHelpInv))  // Set state(2) vol and in, get vol info
	http.HandleFunc("/inv/stophelp", PostOnly(IStop))   // Set vol and inv state(0)

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

	up := mongo.InvSup(inv.Id, inv.Name, inv.Phone, inv.Password)

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
	sIn, name, phone := mongo.InvSin(inv.Id, inv.Password)

	resp := struct {
		Resp  string `json:"resp"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}{sIn, name, phone}

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

	up := mongo.VolSup(vol.Name, vol.Phone, vol.Password)

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
	sIn, name, phone := mongo.VolSin(vol.Phone, vol.Password)

	resp := struct {
		Resp  string `json:"resp"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}{sIn, name, phone}

	if len(name) == 0 && len(phone) == 0 {
		js, err := json.Marshal(resp.Resp)
		if err != nil {
			log.Println(err)
		}
		w.Write(js)
	} else {
		js, bad := json.Marshal(resp)
		if bad != nil {
			log.Println(bad)
		}
		w.Write(js)
	}
}

func VolHelp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	type canHelp struct {
		Phone     string
		Longitude string
		Latitude  string
	}

	var ch canHelp

	err = json.Unmarshal(body, &ch)
	if err != nil {
		log.Println(err)
	}

	help := mongo.VHelp(ch.Phone, ch.Latitude, ch.Longitude)
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
	}{strconv.Itoa(help)}

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

	help := mongo.VolEx(vol.Phone)
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

	good, bad := mongo.GetVolReviews(vol.Phone)

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
		ID     string `json:"id"`
		Phone  string `json:"phone"`
		Review string `json:"review"`
	}
	var req Req
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println(err)
	}
	rev := mongo.ChangeVReview(req.ID, req.Phone, req.Review)
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
		Id    string
		Phone string
	}
	var help Help
	err = json.Unmarshal(body, &help)
	if err != nil {
		log.Println(err)
	}
	resp, rVol, rInv := mongo.FindHelp(help.Id, help.Phone)

	js, exc := json.Marshal(bson.M{"resp": resp, "inv": rInv, "vol": rVol})
	if exc != nil {
		log.Println(err)
	}
	w.Write(js)
}

// GP
func VolGP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	type setGeo struct {
		Phone     string
		Longitude string
		Latitude  string
	}

	var ch setGeo

	err = json.Unmarshal(body, &ch)
	if err != nil {
		log.Println(err)
	}

	help := mongo.VGP(ch.Phone, ch.Latitude, ch.Longitude)
	resp := struct {
		Resp string `json:"resp"`
	}{strconv.FormatBool(help)}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
	}
	w.Write(js)
}

func VolHelpInv(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	// var InvID string
	type Help struct {
		Conid string `json:"conid"`
		Phone string `json:"phone"`
	}

	var getH Help
	err = json.Unmarshal(body, &getH)
	if err != nil {
		log.Println(err)
	}
	name, phone, geo := mongo.VolGetInv(getH.Phone, getH.Conid)

	resp := struct {
		Name  string    `json:"name"`
		Phone string    `json:"phone"`
		Geo   [2]string `json:"geo"`
	}{name, phone, geo}
	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(err)
	}
	w.Write(js)
}

func IStop(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	type StopHelp struct {
		Conid string `json:"conid"`
		Phone string `json:"phone"`
	}

	var stop StopHelp
	err = json.Unmarshal(body, &stop)
	if err != nil {
		log.Println(err)
	}
	sh := mongo.InvStopHelp(stop.Conid, stop.Phone)
	js, bad := json.Marshal(strconv.FormatBool(sh))
	if bad != nil {
		log.Println(err)
	}
	w.Write(js)
}

type handler func(w http.ResponseWriter, r *http.Request)

func PostOnly(h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h(w, r)
			return
		}
		http.Error(w, "post only", http.StatusMethodNotAllowed)
	}
}

func GetOnly(h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h(w, r)
			return
		}
		http.Error(w, "get only", http.StatusMethodNotAllowed)
	}
}
