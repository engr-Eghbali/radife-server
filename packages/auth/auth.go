package magic

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

////////////////user info data type//////////////////////////////////////////

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

	Login int32 `json:"login"`
}

/////////////////////////////verify phone number/////////////////////////////////

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

		if err.Error() == "not found" {

			inlog := []string{"null", "null"}
			infav := []string{"null", "null"}
			log.Print(err)
			//build crypted verification code and return and send *SMS*
			err = c.Insert(&user{Phone: phone, Name: "نام", Add: "آدرس", X: "0", Y: "0", Rank: "b", Level: "1", Pending: "null", Avatar: "avatar.jpg", Log: inlog, Favorit: infav, Wallet: "0", Promo: "0", Login: 1, Key: "12345"})
			log.Print("\nnew user submited:" + phone + "\n")

			return "12345"
		} else {
			log.Print("++submit auth@user database error")
			log.Print(err)
			return "-1"
		}

	} else {

		if result.Login == 1 {
			log.Print("\n**++duplicate user try to submit...++**\n")
			return "0"
		} else {
			log.Print("\nlog outed user come back...\n")
			//build crypted verification code and return and send *SMS*
			// Update
			colQuerier := bson.M{"phone": phone}
			change := bson.M{"$set": bson.M{"key:", "12345"}}
			err = c.Update(colQuerier, change)
			if err != nil {
				log.Print("query update key failed on:")
				log.Print(phone)
				return -1
			} else {
				return "12345"
			}
			return "12345"
		}
	}

}

///////////////////////////////////////////////////////////////////////////////////////////////////////////

///////////////////////////update address//////////////////////////////////

func Update_add(phone string, add string, x string, y string) (flg bool) {

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
	change := bson.M{"$set": bson.M{"add": add, "x": x, "y": y}}
	err = c.Update(colQuerier, change)
	if err != nil {
		log.Print("\nupdate address failed...\n")
		log.Print(err)
		return false

	} else {
		return true
	}

}
////////////////////////////////////////////////////////////////////////////////

/////////////////////////log out func//////////////////////////////////////////

func Logout(phone string, key string) (flg bool) {

	result user
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

	
	err = c.Find(bson.M{"phone": phone}).Select(bson.M{"key": key}).One(&result)
	if err != nil {
	
	 log.Printf("!! logout failed,DB phone+key query failed")
	 return false
	
	 }else{
		
	colQuerier := bson.M{"phone": phone}
	change := bson.M{"$set": bson.M{"login": 0,"key":"nil"}}
	err = c.Update(colQuerier, change)
	
	if err != nil {
		log.Print("\nupdate key failed...\n")
		log.Print(err)
		return false

	} else {
		return true
	}


	}


}