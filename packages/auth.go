package magic

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type user struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`

	phone string `json:"phone"`

	name string `json:"name"`

	add string `json:"add"`

	x string `json:"x"`

	y string `json:"y"`

	rank string `json:"rank"`

	level string `json:"level"`

	pending string `json:"pending"`

	avatar string `json:"avatar"`

	log []string `json:"log"`
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
        var newUser user
		log.Print(err)
		err=c.Insert(&newUser{phone: phone, name: "نام", add:"آدرس", x: "0", y: "0", rank: "b", level: "1",pending:"null",avatar:"avatar.jpg",log[0]="0"})
		log.Print("\nnew user submited:" + phone + "\n")
		return "12345"
	} else {

		log.Print("\nduplicate user try to submit...\n")
		return "0"
	}

}
