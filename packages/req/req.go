package magic

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

type PreShop struct {
	Name string `json:"name"`

	Add string `json:"add"`

	Stars string `json:"stars"`

	Avatar string `json:"avatar"`

	Off string `json:"off"`

	Delivery string `json:"delivery"`
}

func Get_category(cat string) (preview []PreShop) {

	var results []PreShop
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return nil
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("shopinfo").C(cat)

	err = c.Find(nil).All(&results)

	if err != nil {
		return results
	}
	log.Print("\n category query failed:\n")
	log.Print(err)
	return nil
}
