package mongo

import (
	"log"

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

func InvSin(id, password string) string {

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
		return "not in db"
	}
	checkPass := comparePasswords(findR.Password, []byte(password))
	if checkPass != true {
		return "bad pass"
	}
	status := bson.M{"$set": bson.M{"needhelp": true}}
	err = c.Update(colQuierier, status)
	if err != nil {
		log.Fatal(err)
	}
	return findR.Name
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

func VolSin(number, password string) string {
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
		return "not in db"
	}
	checkPass := comparePasswords(findR.Password, []byte(password))
	if checkPass != true {
		return "bad pass"
	}
	status := bson.M{"$set": bson.M{"canhelp": true}}
	err = c.Update(colQuierier, status)
	if err != nil {
		log.Fatal(err)
	}

	return findR.Name
}

func VolSup(name, number, password string, geo [2]string) string {
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
	err = c.Insert(&s.VolUser{Name: name, Number: number, Geo: geo, Password: hashPass})
	if err != nil {
		log.Fatal(err)
	}
	return "signUP"
}

func GetGeoV() [][2]string {
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

	var result [][2]string
	for i, _ := range fGeo {
		result = append(result, fGeo[i].Geo)
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
