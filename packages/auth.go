package magic

import (
	"log"

	mgo "mgo.v2"
	"mgo.v2/bson"
)

func Verify_phone(phone string) (verify_code int32) {

	var result string
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Fatal("\n!!!!-- DB connection error:")
		log.Fatal(err)
		log.Fatal("\n")
		return -1
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("userinfo").C("users")
	err = c.Find(bson.M{"phone": phone}).One(&result)

	if err != nil {

		log.Fatal(err)
		log.Fatal("\nnew user submited:" + phone + "\n")
		return 12345
	} else {

		log.Fatal("\nduplicate user try to submit...\n")
		return 0
	}

}
