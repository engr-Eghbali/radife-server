package magic

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mostafah/go-jalali/jalali"

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

/////////////////////preview of orders/////////////////////////////////
type PreOrderView struct {
	Total    string `json:"total"`
	DateIn   string `json:"date-in"`
	TimeIn   string `json:"time-in"`
	Recieved int32  `json:"recieved"`
}

/////////////////////////////////////each good info////////////////////

type Good struct {
	ID      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name    string        `json:"name"`
	Price   string        `json:"price"`
	Pic     string        `json:"pic"`
	Detail  string        `json:"detail"`
	Keyword string        `json:"keyword"`
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

///////////Shop status structur ///////////////////////////////////////////
type ShopStatus struct {
	Time      string   `json:"time"`
	Hood      string   `json:"hood"`
	Detail    string   `json:"detail"`
	Subcats   []string `json:"categories"`
	Followers []string `json:"followers"`
}

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
func Get_goods(shopid string, cat string, subcat string) (goods []Good) {

	var results []Good
	var temp []Good
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

	err = c.Find(bson.M{"shopid": shopid}).All(&temp)

	if err != nil {

		log.Print("\n category query failed:\n")
		log.Print(err)
		return nil

	} else {

		for _, result := range temp {
			if strings.Contains(result.Keyword, subcat) {

				results = append(results, result)
			}
		}
		return results

	}

}

////////////////////////////////////////////////////////////////////

//////////////////////////set order/////////////////////////////////

func Send_cart(shopID string, customer string, x string, y string, add string, total string, cart string) (orderID string) {

	order_array := strings.Split(cart, "#")
	originx := "0"
	originy := "0"
	destinationx := x
	destinationy := y
	totalPrice := total
	var order Order

	session, err := mgo.Dial("127.0.0.1")
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return "0"
	} else {

		///**check if any inorder cart ,remove it befor start new one
		c := session.DB("orderinfo").C("order")
		err = c.Find(bson.M{"customer": customer}).One(&order)

		if err != nil {

			log.Print("\n cancel order query for checking duplicate: :\n")
			log.Print(err)

		} else {
			order.Recieved = -1
			order.Pay = -1
			order.Courier = "-1"
			order.TimeOut = "-1"
			order.DateOut = "-1"
			c = session.DB("orderinfo").C("canceled")
			err = c.Insert(&order)

			if err != nil {

				log.Print("\n failed to auto remove order for :" + customer + "!!!!!!!!\n")

			} else {

				c = session.DB("orderinfo").C("order")
				err = c.Remove(bson.M{"customer": customer})

				if err != nil {
					fmt.Printf("\nQuery failed to auto remove order for : %v\n", err)

				} else {

					log.Print("\norder auto canceled by user:" + customer + "!!!!!!!!\n")

				}

			}
		}

		c = session.DB("orderinfo").C("order")

		// "Printed on 1392/04/02"
		// Get a new instance of ptime.Time using time.Time
		var orderID bson.ObjectId
		orderID = bson.NewObjectId()
		stringOrderID := orderID.Hex()

		err = c.Insert(&Order{Customer: customer, Shop: shopID, Courier: "0", Cart: order_array, Total: totalPrice, DateIn: jalali.Strftime("%Y/%b/%d", time.Now()), TimeIn: time.Now().Format("3:04PM"), DateOut: "0", TimeOut: "0", OriginX: originx, OriginY: originy, DestinationX: destinationx, DestinationY: destinationy, Recieved: 0, Pay: 0})

		if err != nil {

			log.Print("\n !!!!!!!!! new order failed by:" + customer + "!!!!!!!!\n")
			return "-1"

		} else {

			log.Print("\nnew order submited by:" + customer + "\n")
			return stringOrderID
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
			if !strings.Contains(element, "@") {
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
		order.Recieved = -1
		order.Pay = -1
		order.Courier = "-1"
		order.TimeOut = "-1"
		order.DateOut = "-1"
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

func UpdateName(phone string, name string, x string, y string) (flg bool) {

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return false
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("userinfo").C("users")

	colQuerier := bson.M{"phone": phone}
	change := bson.M{"$set": bson.M{"name": name, "x": x, "y": y}}
	err = c.Update(colQuerier, change)
	if err != nil {
		log.Print("\nupdate name failed...\n")
		log.Print(err)
		return false

	} else {
		return true
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////when user request history log //////////////////////////////////////////////
func ShowHisrory(customer string) (list []PreOrderView, flg bool) {

	var orders []PreOrderView
	var temp []PreOrderView
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return orders, false
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	//fetch canceled orders
	c := session.DB("orderinfo").C("canceled")
	err = c.Find(bson.M{"customer": customer}).All(&temp)
	if err == nil {
		orders = append(orders, temp...)
	} else {
		log.Print(err)
	}
	///////fetch in progress orders
	c = session.DB("orderinfo").C("inProgress")
	err1 := c.Find(bson.M{"customer": customer}).All(&temp)
	if err1 == nil {
		orders = append(orders, temp...)
	} else {
		log.Print(err1)
	}
	///fetch recieved orders
	c = session.DB("orderinfo").C("recieved")
	err2 := c.Find(bson.M{"customer": customer}).All(&temp)
	if err2 == nil {
		orders = append(orders, temp...)
	} else {
		log.Print(err2)
	}

	if err == nil || err1 == nil || err2 == nil {
		log.Print("\n log history visited:")
		log.Print(customer)
		return orders, true
	} else {
		log.Print("\n failed to quety preorderview:")
		log.Print(customer)
		return orders, false
	}

}

/////////// get shop status and subcategories (for header) /////

func GetShopStats(customer string, shopID string, cat string) (shopStats ShopStatus, liked bool, flg bool) {

	var results ShopStatus
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return results, false
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("shopinfo").C(cat)

	err = c.Find(bson.M{"phone": shopID}).One(&results)

	if err != nil {

		log.Print("\n Shop status query failed:\n")
		log.Print(err)
		return results, false

	} else {

		for i := range results.Followers {
			if results.Followers[i] == customer {

				return results, true, true
			}
		}

		return results, false, true

	}

}

func AddFollower(customer string, shopID string, category string) (flg bool) {

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {

		log.Print("\n!!!!-- DB connection error:")
		log.Print(err)
		log.Print("\n")
		return false
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("shopinfo").C(category)
	idQuerier := bson.M{"phone": shopID}
	change := bson.M{"$push": bson.M{"follower": customer}}
	err = c.Update(idQuerier, change)
	if err != nil {
		return false
	} else {
		return true
	}

}
