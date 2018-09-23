package mongo

import (
	"log"
	"os"
	"strconv"

	s "../sett"

	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	LOCAL   = "localhost:27017"
	DOCKER  = "mongodb://mongo:27017"
	IDBNAME = "inv"
	VDBNAME = "vol"
	ICOL    = "invalids"
	VCOL    = "volonters"
)

var CONN = os.Getenv("CONN")

func InvSin(id, password string) (string, string, string) {

	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	c := session.DB(IDBNAME).C(ICOL)
	colQuierier := bson.M{"id": id}
	var findR s.InvUser
	c.Find(colQuierier).One(&findR)
	if len(findR.Name) == 0 {
		return "not in db", "", ""
	}

	checkPass := comparePasswords(findR.Password, []byte(password))
	if checkPass == false {
		return "bad pass", "", ""
	} else {
		online := bson.M{"$set": bson.M{"online": true}}
		err = c.Update(colQuierier, online)
		if err != nil {
			log.Println(err)
		}
		return "signIn", findR.Name, findR.Phone
	}

}

func InvSup(id, name, phone, password string) string {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	if len(id) == 0 {
		return "empty id"
	} else if len(name) == 0 {
		return "empty name"
	} else if len(phone) == 0 {
		return "empty phone"
	} else if len(password) == 0 {
		return "empty password"
	}
	c := session.DB(IDBNAME).C(ICOL)
	var findR s.InvUser
	c.Find(bson.M{"id": id}).One(&findR)

	if len(findR.Name) != 0 {
		return "in db"
	}

	hashPass := hashAndSalt([]byte(password))
	err = c.Insert(&s.InvUser{Id: id, Name: name, Phone: phone, Password: hashPass, State: 0})
	if err != nil {
		log.Println(err)
	}
	return "signUP"
}

func VolSin(phone, password string) (string, string, string) {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	c := session.DB(VDBNAME).C(VCOL)
	colQuierier := bson.M{"phone": phone}
	var findR s.VolUser
	c.Find(colQuierier).One(&findR)
	if len(findR.Phone) == 0 {
		return "not in db", "", ""
	}
	checkPass := comparePasswords(findR.Password, []byte(password))
	if checkPass == false {
		return "bad pass", "", ""
	} else {
		online := bson.M{"$set": bson.M{"online": true}}
		err = c.Update(colQuierier, online)
		if err != nil {
			log.Println(err)
		}
		return "signIn", findR.Name, findR.Phone
	}
}

func VolSup(name, phone, password string) string {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	if len(phone) == 0 {
		return "empty phone"
	} else if len(name) == 0 {
		return "empty name"
	} else if len(password) == 0 {
		return "empty password"
	}
	c := session.DB(VDBNAME).C(VCOL)
	var findR s.VolUser
	c.Find(bson.M{"phone": phone}).One(&findR)
	if len(findR.Name) != 0 {
		return "in db"
	}
	hashPass := hashAndSalt([]byte(password))
	err = c.Insert(&s.VolUser{Name: name, Phone: phone, Password: hashPass, State: 0})
	if err != nil {
		log.Println(err)
	}
	return "signUP"

}

func VHelp(phone, lat, long string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	c := session.DB(VDBNAME).C(VCOL)
	var findR s.VolUser
	colQuierier := bson.M{"phone": phone}
	c.Find(colQuierier).One(&findR)
	geo := [2]string{lat, long}

	if len(findR.Name) == 0 {
		return false
	} else if findR.Online == false {
		return false
	} else {
		status := bson.M{"$set": bson.M{"state": 1, "geo": geo}}
		err = c.Update(colQuierier, status)
		if err != nil {
			log.Println(err)
		}
		return true
	}
}

// GP

func VGP(phone, lat, long string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	c := session.DB(VDBNAME).C(VCOL)
	var findR s.VolUser
	colQuierier := bson.M{"phone": phone}
	c.Find(colQuierier).One(&findR)
	geo := [2]string{lat, long}

	if len(findR.Name) == 0 {
		return false
	} else if findR.Online == false {
		return false
	} else {
		status := bson.M{"$set": bson.M{"geo": geo}}
		err = c.Update(colQuierier, status)
		if err != nil {
			log.Println(err)
		}
		return true
	}
}

func IGP(id, lat, long string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	c := session.DB(IDBNAME).C(ICOL)
	var findR s.InvUser
	colQuierier := bson.M{"id": id}
	c.Find(colQuierier).One(&findR)
	geo := [2]string{lat, long}

	if len(findR.Name) == 0 {
		return false
	} else if findR.Online == false {
		return false
	} else {
		status := bson.M{"$set": bson.M{"geo": geo}}
		err = c.Update(colQuierier, status)
		if err != nil {
			log.Println(err)
		}
		return true
	}
}

//

func IHelp(id, lat, long string) int {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	c := session.DB(IDBNAME).C(ICOL)
	var findR s.InvUser
	var readyInv []s.InvUser

	c.Find(bson.M{"state": 1}).All(&readyInv)
	conID := len(readyInv) + 1
	colQuierier := bson.M{"id": id}
	c.Find(colQuierier).One(&findR)
	geo := [2]string{lat, long}

	if len(findR.Name) == 0 {
		return -1
	} else if findR.Online == false {
		return -1
	} else if findR.State == 1 {
		return -1
	} else {
		status := bson.M{"$set": bson.M{"state": 1, "geo": geo, "conid": conID}}
		err = c.Update(colQuierier, status)
		if err != nil {
			log.Println(err)
		}
		return conID
	}
}

func VolEx(phone string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	c := session.DB(VDBNAME).C(VCOL)
	var findR s.VolUser
	colQuierier := bson.M{"phone": phone}
	c.Find(colQuierier).One(&findR)
	if len(findR.Name) == 0 {
		return false
	} else {
		status := bson.M{"$set": bson.M{"state": 0, "conid": -1, "online": false}}
		err = c.Update(colQuierier, status)
		if err != nil {
			log.Println(err)
		}
		return true
	}
}

func InvEx(id string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	c := session.DB(IDBNAME).C(ICOL)
	var findR s.InvUser
	colQuierier := bson.M{"id": id}
	c.Find(colQuierier).One(&findR)
	if len(findR.Id) == 0 {
		return false
	} else {
		status := bson.M{"$set": bson.M{"state": 0, "conid": 0, "online": false}}
		err = c.Update(colQuierier, status)
		if err != nil {
			log.Println(err)
		}
		return true
	}
}

func GetGeoV() [][]string {
	var fGeo []s.VolUser
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	c := session.DB(VDBNAME).C(VCOL)
	err = c.Find(bson.M{}).All(&fGeo)
	if err != nil {
		log.Println(err)
	}

	var result [][]string
	for i, _ := range fGeo {
		sl := fGeo[i].Geo[:]
		sl = append(sl, strconv.Itoa(fGeo[i].State), fGeo[i].Name, fGeo[i].Phone)
		result = append(result, sl)
	}
	return result
}

func GetGeoI() [][]string {

	var fGeo []s.InvUser
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	c := session.DB(IDBNAME).C(ICOL)
	err = c.Find(bson.M{}).All(&fGeo)
	if err != nil {
		log.Println(err)
	}

	var result [][]string
	for i, _ := range fGeo {
		sl := fGeo[i].Geo[:]
		sl = append(sl, strconv.Itoa(fGeo[i].State), fGeo[i].Name, fGeo[i].Id, fGeo[i].Phone)
		result = append(result, sl)
	}

	return result
}

func GetVolReviews(phone string) (int, int) {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.
			Println(err)
	}
	defer session.Close()
	var vol s.VolUser
	c := session.DB(VDBNAME).C(VCOL)
	err = c.Find(bson.M{"phone": phone}).One(&vol)
	if err != nil {
		log.Println(err)
	}
	return vol.GoodReviews, vol.BadReviews
}

func ChangeVReview(id, phone, review string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	c := session.DB(VDBNAME).C(VCOL)
	// i := session.DB(VDBNAME).C(VCOL)
	var vol s.VolUser

	colQuierier := bson.M{"phone": phone}
	err = c.Find(colQuierier).One(&vol)
	if err != nil {
		log.Println(err)
	}
	if len(vol.Phone) == 0 {
		return false
	} else {
		if review == "bad" {
			rev := bson.M{"$set": bson.M{"badreviews": vol.BadReviews + 1}}
			err = c.Update(colQuierier, rev)
			if err != nil {
				log.Println(err)
			}
		} else if review == "good" {
			rev := bson.M{"$set": bson.M{"goodreviews": vol.GoodReviews + 1}}
			err = c.Update(colQuierier, rev)
			if err != nil {
				log.Println(err)
			}
		} else {
			return false
		}
	}
	return true
}

func VolGetInv(phone, conid string) (string, string, string, [2]string) {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	var vol s.VolUser
	v := session.DB(VDBNAME).C(VCOL)
	v.Find(bson.M{"phone": phone}).One(&vol)
	if len(vol.Name) == 0 {
		return "vol not found", "", "", [2]string{"", ""}
	}
	checkVol := VolCheck(&vol)
	if checkVol == false {
		return "vol not ready", "", "", [2]string{"", ""}
	}

	i := session.DB(IDBNAME).C(ICOL)
	var inv s.InvUser
	invId, err := strconv.Atoi(conid)
	if err != nil {
		log.Println()
		return "bad conid", "", "", [2]string{"", ""}
	}
	err = i.Find(bson.M{"conid": invId}).One(&inv)
	if err != nil {
		log.Println(err)
		return "bad inv find", "", "", [2]string{"", ""}
	}
	if len(inv.Name) == 0 {
		return "user not found", "", "", [2]string{"", ""}

	} else if len(inv.Helper) != 0 || len(vol.InTrouble) != 0 {
		return "busy", "", "", [2]string{"", ""}
	} else {
		iStateAndHelper := bson.M{"$set": bson.M{"state": 2, "helper": phone}}
		vStateAndInTrouble := bson.M{"$set": bson.M{"state": 2, "introuble": inv.Id}}
		err = i.Update(bson.M{"conid": invId}, iStateAndHelper)
		if err != nil {
			log.Println(err)
			return "bad inv set", "", "", [2]string{"", ""}
		}
		err = v.Update(bson.M{"phone": phone}, vStateAndInTrouble)

		if err != nil {
			log.Println(err)
			return "bad vol set", "", "", [2]string{"", ""}
		}
		return "nice", inv.Name, inv.Phone, inv.Geo
	}
}

func InvStopHelp(conid, phone string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	var vol s.VolUser
	var inv s.InvUser

	v := session.DB(VDBNAME).C(VCOL)

	i := session.DB(IDBNAME).C(ICOL)
	invId, err := strconv.Atoi(conid)
	if err != nil {
		log.Println(err)
	}
	err = v.Find(bson.M{"phone": phone}).One(&vol)
	if err != nil {
		log.Println(err)
	}
	err = i.Find(bson.M{"conid": invId}).One(&inv)
	if err != nil {
		log.Println(err)
	}

	if len(vol.Name) == 0 || len(inv.Name) == 0 {
		log.Println(vol.Name, inv.Name)
		return false
	} else if vol.Online == false || inv.Online == false {
		log.Println(vol.Online, inv.Online)
		return false
	} else if vol.State != 2 || inv.State != 2 {
		log.Println(vol.State, inv.State)
		return false
	} else {
		stop := bson.M{"$set": bson.M{"state": 0, "introuble": ""}}

		istop := bson.M{"$set": bson.M{"state": 0, "conid": 0, "helper": ""}}

		v.Update(bson.M{"phone": phone}, stop)
		i.Update(bson.M{"conid": invId}, istop)

		return true
	}
}

// Get helper for inv
func Helper(id string) (string, s.VolUser) {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	var inv s.InvUser
	var vol s.VolUser

	i := session.DB(IDBNAME).C(ICOL)
	v := session.DB(VDBNAME).C(VCOL)

	err = i.Find(bson.M{"id": id}).One(&inv)
	if err != nil {
		log.Println(err)
	}
	if len(inv.Helper) == 0 || inv.State != 2 {
		return "false", vol
	} else {
		err = v.Find(bson.M{"phone": inv.Helper}).One(&vol)
		if err != nil {
			log.Println(err)
		}
		if len(vol.Name) == 0 {
			return "false", vol
		} else {
			return "true", vol
		}
	}

}

// If ren vol :
func VolStop(phone string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	var inv s.InvUser
	var vol s.VolUser

	i := session.DB(IDBNAME).C(ICOL)
	v := session.DB(VDBNAME).C(VCOL)

	err = v.Find(bson.M{"phone": phone}).One(&vol)
	if err != nil {
		log.Println(err)
	}
	if len(vol.Name) == 0 {
		return false
	} else if len(vol.InTrouble) == 0 {
		return false
	} else if vol.Online == false {
		return false
	} else {
		err = i.Find(bson.M{"id": vol.InTrouble}).One(&inv)
		if err != nil {
			log.Println(err)
		}
		volInfo := bson.M{"$set": bson.M{"state": 1, "introuble": ""}}
		invInfo := bson.M{"$set": bson.M{"state": 1, "helper": ""}}
		err = v.Update(bson.M{"phone": phone}, volInfo)
		if err != nil {
			log.Println(err)
		}
		err = i.Update(bson.M{"id": inv.Id}, invInfo)
		if err != nil {
			log.Println(err)
		}
		return true
	}

}

func HGeo(id string) (bool, [2]string) {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	var inv s.InvUser
	var vol s.VolUser

	i := session.DB(IDBNAME).C(ICOL)
	v := session.DB(VDBNAME).C(VCOL)

	err = i.Find(bson.M{"id": id}).One(&inv)
	if err != nil {
		log.Println(err)
	}
	if len(inv.Name) == 0 {
		return false, [2]string{"", ""}
	} else if len(inv.Helper) == 0 {
		return false, [2]string{"", ""}
	} else if inv.Online == false {
		return false, [2]string{"", ""}
	} else {
		err = v.Find(bson.M{"phone": inv.Helper}).One(&vol)
		if err != nil {
			log.Println(err)
		}
		if len(vol.Name) == 0 {
			return false, [2]string{"", ""}
		} else if len(vol.InTrouble) == 0 {
			return false, [2]string{"", ""}
		} else if vol.Online == false {
			return false, [2]string{"", ""}
		}
		return true, vol.Geo
	}
}

func VolCheck(vol *s.VolUser) bool {
	if vol.State == 0 {
		return false
	} else {
		return true
	}
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}
	return true
}
