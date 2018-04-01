package magic

import (
	"log"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type coord struct {
	x string `json:"x"`
	y string `json:"y"`
}

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
	ID     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name   string        `json:"name"`
	Price  int64         `json:"price"`
	Pic    string        `json:"pic"`
	Detail string        `json:"detail"`
}

///////////////////////////////////////////////////////////////////////////

////////////////////////////////////each cart info////////////////////////
type Order struct {
	Customer    string   `json:"customer"`
	Courier     string   `json:"courier"`
	Cart        []string `json:"cart"`
	Total       int64    `json:"total"`
	DateIn      string   `json:"date-in"`
	TimeIn      string   `json:"time-in"`
	DateOut     string   `json:"date-out"`
	TimeOut     string   `json:"time-out"`
	Origin      coord    `json:"origin"`
	Destination coord    `json:"destination"`
	Recieved    int32    `json:"recieved"`
	Pay         int32    `json:"pay"`
}

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

////////////////////////////////////////////////////////////////////

//////////////////////////set order/////////////////////////////////

func Send_cart(shopID string, customer string, x string, y string, add string, total string, cart string) (flg bool) {

	var result Order
	var origin coord
	var destination coord

	order_array := strings.Split(cart, "#")
	origin.x = "0"
	origin.y = "0"
	destination.x = x
	destination.y = y

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return false
	} else {

		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB("orderinfo").C("order")

		err = c.Insert(&Order{Customer: customer, Courier: "0", Cart: order_array, Total: total, DateIn: time.Now().Local().Format("2006-01-02"), TimeIn: time.Now().Format("3:04PM"), DateOut: "0", TimeOut: "0", Origin: origin, Destination: destination, Recieved: 0, Pay: 0})

		if err != nil {

			log.Print("\n !!!!!!!!! new order failed by:" + customer + "!!!!!!!!\n")
			return false

		} else {

			log.Print("\nnew order submited by:" + customer + "\n")
			return true
		}

	}
}
