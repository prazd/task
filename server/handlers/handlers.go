package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/prazd/task/server/mongo"
	s "github.com/prazd/task/sett"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func InvSignUp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	var inv s.InvUser
	err = json.Unmarshal(body, &inv)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}

	up := mongo.InvSup(inv.Id, inv.Name, inv.Phone, inv.Password)

	resp := struct {
		Resp string `json:"resp"`
	}{up}

	js, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	w.Write(js)
}

func InvSignIn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	var inv s.InvUser
	err = json.Unmarshal(body, &inv)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	sIn, name, phone := mongo.InvSin(inv.Id, inv.Password)

	respFull := struct {
		Resp  string `json:"resp"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}{sIn, name, phone}
	respExcp := struct {
		Resp string `json:"resp"`
	}{sIn}
	if len(name) == 0 && len(phone) == 0 {
		js, err := json.Marshal(respExcp)
		if err != nil {
			log.Println(err)
		}
		w.Write(js)
	} else {
		js, err := json.Marshal(respFull)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "Error 500")
		}
		w.Write(js)
	}

}

func VolSignUp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}

	var vol s.VolUser
	err = json.Unmarshal(body, &vol)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}

	up := mongo.VolSup(vol.Name, vol.Phone, vol.Password)

	resp := struct {
		Resp string `json:"resp"`
	}{up}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
		fmt.Fprint(w, "Error 500")
		return
	}
	w.Write(js)

}

func VolSignIn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	var vol s.VolUser
	err = json.Unmarshal(body, &vol)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	sIn, name, phone := mongo.VolSin(vol.Phone, vol.Password)

	respFull := struct {
		Resp  string `json:"resp"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}{sIn, name, phone}

	respExcp := struct {
		Resp string `json:"resp"`
	}{sIn}
	if len(name) == 0 && len(phone) == 0 {
		js, err := json.Marshal(respExcp)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "Error 500")
			return
		}
		w.Write(js)
	} else {
		js, err := json.Marshal(respFull)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "Error 500")
			return
		}
		w.Write(js)
	}

}

func VolHelp(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
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
		fmt.Fprint(w, "Error 500")
		return
	}

	help := mongo.VHelp(ch.Phone, ch.Latitude, ch.Longitude)
	resp := struct {
		Resp string `json:"resp"`
	}{strconv.FormatBool(help)}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
		fmt.Fprint(w, "Error 500")
		return
	}
	w.Write(js)
}

func InvHelp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
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
		fmt.Fprint(w, "Error 500")
		return
	}

	help := mongo.IHelp(nh.Id, nh.Latitude, nh.Longitude)
	resp := struct {
		Resp string `json:"resp"`
	}{strconv.Itoa(help)}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
		fmt.Fprint(w, "Error 500")
		return
	}
	w.Write(js)
}

func VolExit(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}

	var vol s.VolUser
	err = json.Unmarshal(body, &vol)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}

	help := mongo.VolEx(vol.Phone)
	resp := struct {
		Resp string `json:"resp"`
	}{strconv.FormatBool(help)}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
		fmt.Fprint(w, "Error 500")
		return
	}
	w.Write(js)
}

func InvExit(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}

	var inv s.InvUser
	err = json.Unmarshal(body, &inv)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}

	help := mongo.InvEx(inv.Id)
	resp := struct {
		Resp string `json:"resp"`
	}{strconv.FormatBool(help)}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
		fmt.Fprint(w, "Error 500")
		return
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
		fmt.Fprint(w, "Error 500")
		return
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
		fmt.Fprint(w, "Error 500")
		return
	}
	w.Write(js)
}

func GetVRev(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	var vol s.VolUser
	err = json.Unmarshal(body, &vol)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}

	good, bad := mongo.GetVolReviews(vol.Phone)

	resp := struct {
		Good int `json:"goodrev"`
		Bad  int `json:"badreviews"`
	}{good, bad}

	js, excp := json.Marshal(resp)
	if excp != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	w.Write(js)
}

// GP
func VolGP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
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
		fmt.Fprint(w, "Error 500")
		return
	}

	help := mongo.VGP(ch.Phone, ch.Latitude, ch.Longitude)
	resp := struct {
		Resp string `json:"resp"`
	}{strconv.FormatBool(help)}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
		fmt.Fprint(w, "Error 500")
		return
	}
	w.Write(js)
}

func InvGP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}

	type setGeo struct {
		Id        string `json:"id"`
		Longitude string `json:"longitude"`
		Latitude  string `json:"latitude"`
	}

	var ch setGeo

	err = json.Unmarshal(body, &ch)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}

	help := mongo.IGP(ch.Id, ch.Latitude, ch.Longitude)
	resp := struct {
		Resp string `json:"resp"`
	}{strconv.FormatBool(help)}

	js, bad := json.Marshal(resp)
	if bad != nil {
		log.Println(bad)
		fmt.Fprint(w, "Error 500")
		return
	}
	w.Write(js)
}

//

func VolHelpInv(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
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
		fmt.Fprint(w, "Error 500")
		return
	}
	info, name, phone, geo := mongo.VolGetInv(getH.Phone, getH.Conid)

	respEx := struct {
		Resp string `json:"resp"`
	}{info}

	resp := struct {
		Resp  string    `json:"resp"`
		Name  string    `json:"name"`
		Phone string    `json:"phone"`
		Geo   [2]string `json:"geo"`
	}{info, name, phone, geo}

	if len(name) == 0 {
		js, bad := json.Marshal(respEx)
		if bad != nil {
			log.Println(err)
			fmt.Fprint(w, "Error 500")
			return
		}
		w.Write(js)
	} else {
		js, bad := json.Marshal(resp)
		if bad != nil {
			log.Println(err)
			fmt.Fprint(w, "Error 500")
			return
		}
		w.Write(js)
	}

}

func IStop(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}

	type StopHelp struct {
		Conid  string `json:"conid"`
		Phone  string `json:"phone"`
		Review string `json:"review"`
	}

	var stop StopHelp
	err = json.Unmarshal(body, &stop)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	if len(stop.Review) == 0 {
		js, err := json.Marshal(bson.M{"resp": "not set review"})
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "Error 500")
			return
		}
		w.Write(js)
	} else {
		sh := mongo.InvStopHelp(stop.Conid, stop.Phone, stop.Review)
		js, err := json.Marshal(bson.M{"resp": strconv.FormatBool(sh)})
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "Error 500")
			return
		}
		w.Write(js)
	}
}

func HelperInfo(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	type InvInfo struct {
		Id string `json:"id"`
	}
	var i InvInfo
	err = json.Unmarshal(body, &i)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	info, helper := mongo.Helper(i.Id)
	if info == "false" {
		js, err := json.Marshal(bson.M{"resp": "bad"})
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "Error 500")
			return
		}
		w.Write(js)
	} else if info == "true" {
		js, err := json.Marshal(bson.M{"resp": "true", "name": helper.Name, "phone": helper.Phone, "longitude": helper.Geo[1], "latitude": helper.Geo[0]})
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "Error 500")
			return
		}
		w.Write(js)
	}
}

func VolRen(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	type VolPhone struct {
		Phone string `json:"phone"`
	}
	var v VolPhone
	err = json.Unmarshal(body, &v)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	stop := mongo.VolStop(v.Phone)
	js, err := json.Marshal(bson.M{"resp": strconv.FormatBool(stop)})
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	w.Write(js)
}

func HelperGeo(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	type InvId struct {
		Id string `json:"id"`
	}
	var iId InvId
	err = json.Unmarshal(body, &iId)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "Error 500")
		return
	}
	info, geo := mongo.HGeo(iId.Id)
	if info == false {
		js, err := json.Marshal(bson.M{"resp": "bad"})
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "Error 500")
			return
		}
		w.Write(js)
	} else {
		js, err := json.Marshal(bson.M{"resp": "nice", "longitude": geo[1], "latitude": geo[0]})
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "Error 500")
			return
		}
		w.Write(js)
	}
}
