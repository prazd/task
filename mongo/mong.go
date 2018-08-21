package mongo

import (
	"log"
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
	CONN    = LOCAL
)

func InvSin(id, password string) (string, string) {

	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	c := session.DB(IDBNAME).C(ICOL)
	colQuierier := bson.M{"id": id}
	var findR s.InvUser
	c.Find(colQuierier).One(&findR)
	if len(findR.Name) == 0 {
		return "", "not in db"
	}
	checkPass := comparePasswords(findR.Password, []byte(password))
	if checkPass != true {
		return "", "bad pass"
	}

	return findR.Name, "signIn"
}

func InvSup(id, name, number, password string) string {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	c := session.DB(IDBNAME).C(ICOL)
	var findR s.InvUser
	c.Find(bson.M{"id": id}).One(&findR)
	if len(findR.Name) != 0 {
		return "in db"
	}
	hashPass := hashAndSalt([]byte(password))
	err = c.Insert(&s.InvUser{Id: id, Name: name, Number: number, Password: hashPass})
	if err != nil {
		log.Fatal(err)
	}
	return "signUP"
}

func VolSin(number, password string) (string, string) {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	c := session.DB(VDBNAME).C(VCOL)
	colQuierier := bson.M{"number": number}
	var findR s.VolUser
	c.Find(colQuierier).One(&findR)
	if len(findR.Number) == 0 {
		return "", "not in db"
	}
	checkPass := comparePasswords(findR.Password, []byte(password))
	if checkPass != true {
		return "", "bad pass"
	}

	return findR.Name, "signIn"
}

func VolSup(name, number, password string) string {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	c := session.DB(VDBNAME).C(VCOL)
	var findR s.VolUser
	c.Find(bson.M{"number": number}).One(&findR)
	if len(findR.Name) != 0 {
		return "in db"
	}
	hashPass := hashAndSalt([]byte(password))
	err = c.Insert(&s.VolUser{Name: name, Number: number, Password: hashPass})
	if err != nil {
		log.Fatal(err)
	}
	return "signUP"
}

func VHelp(number, lat, long string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	c := session.DB(VDBNAME).C(VCOL)
	var findR s.VolUser
	colQuierier := bson.M{"number": number}
	c.Find(colQuierier).One(&findR)
	geo := [2]string{lat, long}

	if len(findR.Name) == 0 {
		return false
	} else {
		status := bson.M{"$set": bson.M{"canhelp": true, "geo": geo}}
		err = c.Update(colQuierier, status)
		if err != nil {
			log.Fatal(err)
		}
		return true
	}
}

func IHelp(id, lat, long string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
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
		status := bson.M{"$set": bson.M{"needhelp": true, "geo": geo}}
		err = c.Update(colQuierier, status)
		if err != nil {
			log.Fatal(err)
		}
		return true
	}
}

func VolEx(number string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	c := session.DB(VDBNAME).C(VCOL)
	var findR s.VolUser
	colQuierier := bson.M{"number": number}
	c.Find(colQuierier).One(&findR)
	if len(findR.Name) == 0 {
		return false
	} else {
		status := bson.M{"$set": bson.M{"canhelp": false}}
		err = c.Update(colQuierier, status)
		if err != nil {
			log.Fatal(err)
		}
		return true
	}
}

func InvEx(id string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	c := session.DB(IDBNAME).C(ICOL)
	var findR s.InvUser
	colQuierier := bson.M{"id": id}
	c.Find(colQuierier).One(&findR)
	if len(findR.Id) == 0 {
		return false
	} else {
		status := bson.M{"$set": bson.M{"needhelp": false}}
		err = c.Update(colQuierier, status)
		if err != nil {
			log.Fatal(err)
		}
		return true
	}
}

func GetGeoV() [][]string {
	var fGeo []s.VolUser
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := session.DB(VDBNAME).C(VCOL)
	err = c.Find(bson.M{}).All(&fGeo)
	if err != nil {
		log.Fatal(err)
	}

	var result [][]string
	for i, _ := range fGeo {
		sl := fGeo[i].Geo[:]
		sl = append(sl, strconv.FormatBool(fGeo[i].CanHelp))
		result = append(result, sl)
	}

	return result
}

func GetGeoI() [][]string {

	var fGeo []s.InvUser
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := session.DB(IDBNAME).C(ICOL)
	err = c.Find(bson.M{}).All(&fGeo)
	if err != nil {
		log.Fatal(err)
	}

	var result [][]string
	for i, _ := range fGeo {
		sl := fGeo[i].Geo[:]
		sl = append(sl, strconv.FormatBool(fGeo[i].NeedHelp))
		result = append(result, sl)
	}

	return result
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
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

// For bot

func QV() int {
	var allVol []s.VolUser
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := session.DB(VDBNAME).C(VCOL)
	err = c.Find(bson.M{}).All(&allVol)
	if err != nil {
		log.Fatal(err)
	}
	return len(allVol)

}

func QI() int {
	var allInv []s.InvUser
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := session.DB(IDBNAME).C(ICOL)
	err = c.Find(bson.M{}).All(&allInv)
	if err != nil {
		log.Fatal(err)
	}
	return len(allInv)

}

func SV() [][2]string {
	var allVol []s.VolUser
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := session.DB(VDBNAME).C(VCOL)
	err = c.Find(bson.M{}).All(&allVol)
	if err != nil {
		log.Fatal(err)
	}
	var result [][2]string
	for i, _ := range allVol {
		if allVol[i].CanHelp == true {
			count := [2]string{allVol[i].Name, allVol[i].Number}
			result = append(result, count)
		}
	}
	return result
}

func SI() [][2]string {
	var allInv []s.InvUser
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	c := session.DB(IDBNAME).C(ICOL)
	err = c.Find(bson.M{}).All(&allInv)
	if err != nil {
		log.Fatal(err)
	}
	var result [][2]string
	for i, _ := range allInv {
		if allInv[i].NeedHelp == true {
			count := [2]string{allInv[i].Id, allInv[i].Name}
			result = append(result, count)
		}
	}
	return result
}
