package magic

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"mgo.v2/bson"
)

//////////////////////preview of each shop/////////////////////////
type PreShop struct {
	Name string `json:"name"`

	Add string `json:"add"`

	Phone string `json:"phone"`

	Stars int64 `json:"star"`

	Avatar string `json:"avatar"`

	Off string `json:"off"`

	Delivery int64 `json:"delivery"`
}

////////////////////////////////////////////////////////////////////////

/////////////////////////////////////each good info////////////////////

type Good struct {
	ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name  string        `json:"name"`
	Price int64         `json:"price"`
	Pic   string        `json:"pic"`
	Info  string        `json:"detail"`
}

///////////////////////////////////////////////////////////////////////////

//////////////////////////////get shops in a category//////////////////////
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

		log.Print("\n category query failed:\n")
		log.Print(err)
		return nil

	} else {
		return results

	}

}

/////////////////////////////////////////////////////////////////////////////

//////////////////////////get goods of each shop/////////////////////////////
func Get_goods(shopid string, cat string) (goods []Good) {

	var results []Good
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return nil
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("goods").C(cat)

	err = c.Find(bson.M{"shopid": shopid}).All(&results)

	if err != nil {

		log.Print("\n category query failed:\n")
		log.Print(err)
		return nil

	} else {
		return results

	}

}
