package magic

import (
	"fmt"
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

	Stars string `json:"star"`

	Avatar string `json:"avatar"`

	Off string `json:"off"`

	Delivery int64 `json:"delivery"`
}

////////////////////////////////////////////////////////////////////////

/////////////////////////////////////each good info////////////////////

type Good struct {
	ID     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name   string        `json:"name"`
	Price  string        `json:"price"`
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
	Total        string   `json:"total"`
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

type User struct {
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

	Delivery string `json:"delivery"`

	X string `json:"x"`

	Y string `json:"y"`
}

///////////////////////////////////////////////////////////////////////////

////////////////calculate delivery cost////////////////////////////////////
func calcDelivery(userInfoX string, userInfoY string, shopInfoX string, shopInfoy string) (deliveryCost float64) {

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
	totalPrice := total
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
func Get_factor(customer string, cat string) (items []Recep, promo string, delivery string, off string, total string) {

	var order Order
	var itemsTemp []Recep
	var goodTemp Good
	var basketTotal float64
	var userInfo User
	var deliveryCost float64
	var shopInfo Shop
	var offSale float64
	var promoTemp float64 = 0
	var singleItem Recep

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return nil, "", "", "", ""
	}

	/////////get order info
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("orderinfo").C("order")
	err = c.Find(bson.M{"customer": customer}).One(&order)

	if err != nil {

		log.Print("\n order query failed:\n")
		log.Print(err)
		return nil, "", "", "", ""

	} else {

		basketTotal, _ = strconv.ParseFloat(order.Total, 64)

		c = session.DB("goods").C(cat)

		///////get goods info for each element in cart
		for _, element := range order.Cart {

			log.Print(element)
			if element == "undefined" {
				break
			}
			NoGood := strings.Split(element, "@")
			singleItem.No = NoGood[0]
			singleItem.ID = NoGood[1]
			err2 := c.FindId(bson.ObjectIdHex(singleItem.ID)).One(&goodTemp)

			if err2 != nil {
				log.Print("\n get  element  query failed:\n")
				log.Print(err2)
				return nil, "", "", "", ""

			} else {

				singleItem.Name = goodTemp.Name
				singleItem.Price = goodTemp.Price
			}

			itemsTemp = append(itemsTemp, singleItem)

		}

		////////get user info
		c = session.DB("userinfo").C("users")
		err3 := c.Find(bson.M{"phone": customer}).One(&userInfo)
		if err3 != nil {

			log.Print("\n get  user  query failed:\n")
			log.Print(err3)
			return nil, "", "", "", ""

		} else {

			promoTemp, _ = strconv.ParseFloat(userInfo.Promo, 64)
		}
		///////get shop info
		c = session.DB("shopinfo").C(cat)
		log.Print(order.Shop)
		err4 := c.Find(bson.M{"phone": order.Shop}).One(&shopInfo)

		if err4 != nil {

			log.Print("\n get  shop  query failed:\n")
			log.Print(err4)
			return nil, "", "", "", ""

		} else {

			offSale, _ = strconv.ParseFloat(shopInfo.Off, 64)
			deliveryCost = calcDelivery(userInfo.X, userInfo.Y, shopInfo.X, shopInfo.Y)
		}

		if deliveryCost > promoTemp {

			total = strconv.FormatFloat(deliveryCost-promoTemp+(basketTotal*((100-offSale)/100)), 'f', 0, 64)
			log.Print("\ncond1:\n")
			log.Print(total)
			log.Print(basketTotal)
			log.Print(offSale)
			log.Print(deliveryCost - promoTemp + (basketTotal * ((100 - offSale) / 100)))

		} else {

			total = strconv.FormatFloat(basketTotal*((100-offSale)/100), 'f', 0, 64)

		}

		return itemsTemp, userInfo.Promo, strconv.FormatFloat(deliveryCost, 'f', 0, 64), shopInfo.Off, total

	}

}

///////////////////////////////////////////////////////////////////
///////////////////cancel an order///////////////////////////////
func CancelOrder(customer string) (flg bool) {

	var order Order

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return true
	}

	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("orderinfo").C("order")
	err = c.Find(bson.M{"customer": customer}).One(&order)

	if err != nil {

		log.Print("\n cancel order query failed:\n")
		log.Print(err)
		return false

	} else {

		c = session.DB("orderinfo").C("canceled")
		err = c.Insert(&order)

		if err != nil {

			log.Print("\n updated canceled order DB Failed" + customer + "!!!!!!!!\n")
			return false

		} else {

			c = session.DB("orderinfo").C("order")
			err = c.Remove(bson.M{"customer": customer})

			if err != nil {
				fmt.Printf("\norder remove fail %v\n", err)
				return false
			} else {

				log.Print("\norder canceled by user:" + customer + "!!!!!!!!\n")
				return true
			}

		}
	}
	// Error check on every access

}

func Profile(customerId string) (person User, flg bool) {

	var userInfo User
	log.Print(customerId)
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return userInfo, false
	}
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("userinfo").C("users")
	err3 := c.Find(bson.M{"phone": customerId}).One(&userInfo)
	if err3 != nil {

		log.Print("\nprofile 404:")
		log.Print(customerId)
		return userInfo, false
	} else {
		log.Print("\nprofile visited:")
		log.Print(customerId)

		return userInfo, true
	}

}
