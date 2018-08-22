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

func GetVolReviews(number string) (int, int) {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.
			Fatal(err)
	}
	defer session.Close()
	var vol s.VolUser
	c := session.DB(VDBNAME).C(VCOL)
	err = c.Find(bson.M{"number": number}).One(&vol)
	if err != nil {
		log.Fatal(err)
	}
	return vol.GoodReviews, vol.BadReviews
}

func ChangeVReview(number, review string) bool {
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	c := session.DB(VDBNAME).C(VCOL)
	var vol s.VolUser
	colQuierier := bson.M{"number": number}
	err = c.Find(colQuierier).One(&vol)
	if err != nil {
		log.Fatal(err)
	}
	if len(vol.Number) == 0 {
		return false
	} else {
		if review == "bad" {
			rev := bson.M{"$set": bson.M{"badreviews": vol.BadReviews + 1}}
			err = c.Update(colQuierier, rev)
			if err != nil {
				log.Fatal(err)
			}
		} else if review == "good" {
			rev := bson.M{"$set": bson.M{"goodreviews": vol.GoodReviews + 1}}
			err = c.Update(colQuierier, rev)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			return false
		}
	}
	return true
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
