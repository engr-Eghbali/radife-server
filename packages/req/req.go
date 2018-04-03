package magic

import (
	"log"
	"strconv"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	ID     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name   string        `json:"name"`
	Price  int64         `json:"price"`
	Pic    string        `json:"pic"`
	Detail string        `json:"detail"`
}

///////////////////////////////////////////////////////////////////////////

////////////////////////////////////each cart info////////////////////////
type Order struct {
	Customer     string   `json:"customer"`
	Shop         string   `json:"shop"`
	Courier      string   `json:"courier"`
	Cart         []string `json:"cart"`
	Total        int64    `json:"total"`
	DateIn       string   `json:"date-in"`
	TimeIn       string   `json:"time-in"`
	DateOut      string   `json:"date-out"`
	TimeOut      string   `json:"time-out"`
	OriginX      string   `json:"originX"`
	OriginY      string   `json:"originY"`
	DestinationX string   `json:"destinationX"`
	DestinationY string   `json:"destinationY"`
	Recieved     int32    `json:"recieved"`
	Pay          int32    `json:"pay"`
}

///////////////////////////////////////////////////////////////////////////

////////////////////////////factor(recep) each item row info////////////////////////////
type Recep struct {
	ID    string
	Name  string
	Price string
	No    string
}

///////////////////////////////////////////////////////////////////////////

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

	Favorit []string `json:"favorit"`

	Wallet string `json:"wallet"`

	Promo string `json:"promo"`
}

/////////////////////////////////////////////////////////////////////////////

/////////////////////////shopInfo struct////////////////////////////////////
type Shop struct {
	Name string `json:"name"`

	Add string `json:"add"`

	Phone string `json:"phone"`

	Avatar string `json:"avatar"`

	Off string `json:"off"`

	Delivery int64 `json:"delivery"`

	X int64 `json:"x"`

	Y int64 `json:"y"`
}

///////////////////////////////////////////////////////////////////////////

////////////////calculate delivery cost////////////////////////////////////
func calcDelivery(userInfoX string, userInfoY string, shopInfoX string, shopInfoy string) {

	log.Print(userInfoX + "&" + userInfoY + "&" + shopInfoX + "&" + shopInfoy + "\n")
	return 5000
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

////////////////////////////////////////////////////////////////////

//////////////////////////set order/////////////////////////////////

func Send_cart(shopID string, customer string, x string, y string, add string, total string, cart string) (flg bool) {

	order_array := strings.Split(cart, "#")
	originx := "0"
	originy := "0"
	destinationx := x
	destinationy := y
	totalPrice, _ := strconv.ParseInt(total, 10, 64)
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

		err = c.Insert(&Order{Customer: customer, Shop: shopID, Courier: "0", Cart: order_array, Total: totalPrice, DateIn: time.Now().Local().Format("2006-01-02"), TimeIn: time.Now().Format("3:04PM"), DateOut: "0", TimeOut: "0", OriginX: originx, OriginY: originy, DestinationX: destinationx, DestinationY: destinationy, Recieved: 0, Pay: 0})

		if err != nil {

			log.Print("\n !!!!!!!!! new order failed by:" + customer + "!!!!!!!!\n")
			return false

		} else {

			log.Print("\nnew order submited by:" + customer + "\n")
			return true
		}

	}
}

////////////////////////get from DB and process factor////////////////////////////
func Get_factor(customer, cat) (items []recep, promo string, delivery string, off string, total string) {

	var order Order
	var itemsTemp []Recep
	var goodTemp Good
	var basketTotal int64
	var shopID string
	var userInfo user
	var deliveryCost string
	var shopInfo Shop
	var offSale string

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return nil, "", "", "", nil
	}

	/////////get order info
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("orderinfo").C("order")
	err = c.Find(bson.M{"customer": customer}).All(&order)

	if err != nil {

		log.Print("\n order query failed:\n")
		log.Print(err)
		return nil, "", "", "", nil

	} else {

		basketTotal = order.Total
		shopID = order.Shop
		c = session.DB("goods").C(cat)

		///////get goods info for each element in cart
		for i, element := range order.Cart {

			NoGood := strings.Split(element, "@")
			items[i].No = NoGood[0]
			items[i].ID = NoGood[1]
			err2 := c.FindId(bson.ObjectIdHex(items[i].ID)).One(&goodTemp)

			if err2 != nil {
				log.Print("\n get  element  query failed:\n")
				log.Print(err2)
				return nil, "", "", "", nil

			} else {

				items[i].Name = goodTemp.Name
				items[i].Price = goodTemp.Price
			}

		}

		////////get user info
		c = session.DB("userinfo").C("users")
		err3 := c.Find(bson.M{"phone": customer}).One(&userInfo)
		if err3 != nil {

			log.Print("\n get  user  query failed:\n")
			log.Print(err3)
			return nil, "", "", "", nil

		} else {

			promo = userInfo.Promo
		}
		///////get shop info
		c = session.DB("shopinfo").C(cat)
		err4 := c.Find(bson.M{"phone": order.Shop}).One(&shopInfo)

		if err4 != nil {

			log.Print("\n get  shop  query failed:\n")
			log.Print(err4)
			return nil, "", "", "", nil

		} else {

			deliveryCost = calcDelivery(userInfo.X, userInfo.Y, shopInfo.X, shopInfo.y)
		}

	}

}
