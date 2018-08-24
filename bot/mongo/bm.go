package mongo

import (
	"log"
	"strconv"

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

type InvUser struct {
	Id       string
	Name     string
	Number   string
	Password string
	NeedHelp bool
}

type VolUser struct {
	Name        string
	Number      string
	Geo         [2]string
	Password    string
	CanHelp     bool
	GoodReviews int
	BadReviews  int
}

func QV() int {
	var allVol []VolUser
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	c := session.DB(VDBNAME).C(VCOL)
	err = c.Find(bson.M{}).All(&allVol)
	if err != nil {
		log.Println(err)
	}
	return len(allVol)

}

func QI() int {
	var allInv []InvUser
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	c := session.DB(IDBNAME).C(ICOL)
	err = c.Find(bson.M{}).All(&allInv)
	if err != nil {
		log.Println(err)
	}
	return len(allInv)

}

func SV() [][4]string {
	var allVol []VolUser
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	c := session.DB(VDBNAME).C(VCOL)
	err = c.Find(bson.M{}).All(&allVol)
	if err != nil {
		log.Println(err)
	}
	var result [][4]string
	for i, _ := range allVol {
		if allVol[i].CanHelp == true {
			count := [4]string{allVol[i].Name, allVol[i].Number, strconv.Itoa(allVol[i].GoodReviews), strconv.Itoa(allVol[i].BadReviews)}
			result = append(result, count)
		}
	}
	return result
}

func SI() [][2]string {
	var allInv []InvUser
	session, err := mgo.Dial(CONN)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	c := session.DB(IDBNAME).C(ICOL)
	err = c.Find(bson.M{}).All(&allInv)
	if err != nil {
		log.Println(err)
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
