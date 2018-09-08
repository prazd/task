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
	if checkPass != true {
		return "bad pass", "", ""
	}

	return "signIn", findR.Name, findR.Phone
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

func VolSin(phone, password string) (string, string) {
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
		return "not in db", ""
	}
	checkPass := comparePasswords(findR.Password, []byte(password))
	if checkPass != true {
		return "bad pass", ""
	}

	return "signIn", findR.Name
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
	} else {
		status := bson.M{"$set": bson.M{"geo": geo}}
		err = c.Update(colQuierier, status)
		if err != nil {
			log.Println(err)
		}
		return true
	}
}

func IHelp(id, lat, long string) bool {
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
	} else {
		status := bson.M{"$set": bson.M{"state": 1, "geo": geo}}
		err = c.Update(colQuierier, status)
		if err != nil {
			log.Println(err)
		}
		return true
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
		status := bson.M{"$set": bson.M{"canhelp": false, "state": 0}}
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
		status := bson.M{"$set": bson.M{"needhelp": false, "state": 0}}
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
		sl = append(sl, strconv.Itoa(fGeo[i].State), fGeo[i].Name, fGeo[i].Phone)
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

func ChangeVReview(phone, review string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()
	c := session.DB(VDBNAME).C(VCOL)
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

func FindHelp(invId, volPhone string) (string, string, string) {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	volc := session.DB(VDBNAME).C(VCOL)
	invc := session.DB(IDBNAME).C(ICOL)

	// Vol Busy
	var vol s.VolUser
	vColQuierier := bson.M{"phone": volPhone}
	err = volc.Find(vColQuierier).One(&vol)
	if err != nil {
		log.Println(err)
	}
	if len(vol.Phone) == 0 {
		return "bad", "not found", "..."
	} else {
		vBusy := bson.M{"$set": bson.M{"state": 2}}
		err = volc.Update(vColQuierier, vBusy)
		if err != nil {
			log.Println(err)
		}
	}

	// Inv Busy
	var inv s.InvUser
	iColQuierier := bson.M{"id": invId}
	err = invc.Find(iColQuierier).One(&inv)
	if err != nil {
		log.Println(err)
	}
	if len(inv.Name) == 0 {
		return "bad", "...", "not found"
	} else {
		iBusy := bson.M{"$set": bson.M{"state": 2}}
		err = invc.Update(iColQuierier, iBusy)
		if err != nil {
			log.Println(err)
		}
	}
	return "nice", "busy", "busy"
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
