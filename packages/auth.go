package magic

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type user struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`

	Phone string `json:"phone"`

	Name string `json:"name"`

	Add string `json:"add"`

	X string `json:"x"`

	Y string `json:"y"`

	Rank string `json:"rank"`

	Level string `json:"level"`

	Pending string `json:"pending"`

	Avatar string `json:"avatar"`

	Log []string `json:"log"`
}

func Verify_phone(phone string) (verify_code string) {

	var result user

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return "-1"
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("userinfo").C("users")
	err = c.Find(bson.M{"phone": phone}).One(&result)

	if err != nil {

		inlog := []string{"null", "null"}

		log.Print(err)
		err = c.Insert(&user{Phone: phone, Name: "نام", Add: "آدرس", X: "0", Y: "0", Rank: "b", Level: "1", Pending: "null", Avatar: "avatar.jpg", Log: inlog})
		log.Print("\nnew user submited:" + phone + "\n")
		return "12345"
	} else {

		log.Print("\nduplicate user try to submit...\n")
		return "0"
	}

}
