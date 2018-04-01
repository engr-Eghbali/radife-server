package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	auth "./packages/auth"
	req "./packages/req"
	"github.com/unrolled/secure"
)

// imports end

////////////////submit control////////////////////////////////////////////
func submit_ctrl(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {

		var phone string
		r.ParseForm()
		phone = r.Form["phone"][0]

		if len(phone) == 11 {

			code := auth.Verify_phone(phone)

			if code == "0" {
				fmt.Fprintf(w, "0")
				return
			}
			if code == "-1" {
				fmt.Fprintf(w, "-1")
				return
			} else {
				fmt.Fprintf(w, code)
				return
			}

		} else {

			fmt.Fprintf(w, "شماره تلفن شما صحیح نمی باشد.")
			return

		}

	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////update control///////////////////////////////////////////////////////////
func setAddress_ctrl(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {

		var phone, add, x, y string
		r.ParseForm()
		phone = r.Form["phone"][0]
		add = r.Form["add"][0]
		x = r.Form["x"][0]
		y = r.Form["y"][0]

		if len(phone) == 11 {

			flg := auth.Update_add(phone, add, x, y)
			if flg {
				fmt.Fprintf(w, "1")
			} else {
				fmt.Fprintf(w, "0")
			}

		} else {

			fmt.Fprintf(w, "invalid input data")
		}

	}
}

//////////////////////////////////////////////////////////////////////////
///////////////////////////category request control/////////////////////////

func category_ctrl(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {
		var response string
		var shops []req.PreShop
		var ranking string
		var i int64 = 0
		r.ParseForm()
		cat := r.Form["category"][0]
		shops = req.Get_category(cat)

		for _, shop := range shops {

			offer := shop.Off + "%%تخفیف"
			delivery := "ارسال رایگان"
			if shop.Delivery != 0 {
				delivery = strconv.FormatInt(shop.Delivery, 10) + "هزینه ارسال"
			}
			for i < 6 {

				if shop.Stars > i {
					ranking = ranking + "<span class=\"glyphicon glyphicon-star-empty\"></span>"
				} else {
					ranking = ranking + "<span class=\"glyphicon glyphicon-star\"></span>"
				}

				i++
			}

			response = response + "<div class=\"market\" onclick=\"getgoods('" + shop.Phone + "')\"><img class=\"marketicon\" src=\"" + shop.Avatar + "\"><p class=\"restname\">" + shop.Name + "</p><p class=\"restlocation\"><span class=\"glyphicon glyphicon-pushpin\"></span>" + shop.Add + "</p><div id=\"reststatus\"><p class=\"stars\">" + ranking + "</p><p class=\"restoff\">" + offer + "</p><p class=\"bike\"><i class=\"fas fa-motorcycle\"></i>" + delivery + "</p></div></div>"

		}
		fmt.Fprintf(w, response)
	}
}

///////////////////////////////////////////////////////////////////////////////////

///////////////////////////get goods of each shop//////////////////////////////

func goods_ctrl(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {
		var response string
		var goods []req.Good

		r.ParseForm()

		shopID := r.Form["shopID"][0]
		category := r.Form["cat"][0]
		goods = req.Get_goods(shopID, category)

		for _, good := range goods {

			response = response + "<div class=\"good\"><img id=\"goodicon" + good.ID.Hex() + "\"class=\"goodicon\" src=\"" + good.Pic + "\"><div class=\"goodinfo\"></br></br><p1 class=\"goodname\" id=\"goodname" + good.ID.Hex() + "\">" + good.Name + "</p1></br><p2 class=\"gooddetail\" id=\"gooddetail" + good.ID.Hex() + "\">" + good.Detail + "</p2></br><p3 class=\"goodprice\" id=\"goodprice" + good.ID.Hex() + "\">" + strconv.FormatInt(good.Price, 10) + "تومان</p3></div><div id=\"addremove\"><button class=\"adder\" onclick=\"addgood('" + good.ID.Hex() + "')\"></button><button id=\"remover" + good.ID.Hex() + "\" class=\"remover\" onclick=\"removegood('" + good.ID.Hex() + "')\"></button></div></div>"

		}
		fmt.Fprintf(w, response)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////

func main() {

	//initial route security
	submit_rout := secure.New(secure.Options{
		AllowedHosts:            []string{"ssl.example.com"},                     // AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
		HostsProxyHeaders:       []string{"X-Forwarded-Hosts"},                   // HostsProxyHeaders is a set of header keys that may hold a proxied hostname value for the request.
		SSLRedirect:             true,                                            // If SSLRedirect is set to true, then only allow HTTPS requests. Default is false.
		SSLTemporaryRedirect:    false,                                           // If SSLTemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
		SSLHost:                 "ssl.example.com",                               // SSLHost is the host name that is used to redirect HTTP requests to HTTPS. Default is "", which indicates to use the same host.
		SSLProxyHeaders:         map[string]string{"X-Forwarded-Proto": "https"}, // SSLProxyHeaders is set of header keys with associated values that would indicate a valid HTTPS request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
		STSSeconds:              315360000,                                       // STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
		STSIncludeSubdomains:    true,                                            // If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
		STSPreload:              true,                                            // If STSPreload is set to true, the `preload` flag will be appended to the Strict-Transport-Security header. Default is false.
		ForceSTSHeader:          false,                                           // STS header is only included when the connection is HTTPS. If you want to force it to always be added, set to true. `IsDevelopment` still overrides this. Default is false.
		FrameDeny:               true,                                            // If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
		CustomFrameOptionsValue: "SAMEORIGIN",                                    // CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option. Default is "".
		ContentTypeNosniff:      true,                                            // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
		BrowserXssFilter:        true,                                            // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
		CustomBrowserXssValue:   "1; report=https://example.com/xss-report",      // CustomBrowserXssValue allows the X-XSS-Protection header value to be set with a custom value. This overrides the BrowserXssFilter option. Default is "".
		ContentSecurityPolicy:   "default-src 'self'",                            // ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "".
		// PublicKey: `pin-sha256="base64+primary=="; pin-sha256="base64+backup=="; max-age=5184000; includeSubdomains; report-uri="https://www.example.com/hpkp-report"`, // PublicKey implements HPKP to prevent MITM attacks with forged certificates. Default is "".
		//ReferrerPolicy: "same-origin" // ReferrerPolicy allows the Referrer-Policy header with the value to be set with a custom value. Default is "".
		//***triger***
		IsDevelopment: true, // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
	})

	setAddress_rout := secure.New(secure.Options{
		AllowedHosts:            []string{"ssl.example.com"},                     // AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
		HostsProxyHeaders:       []string{"X-Forwarded-Hosts"},                   // HostsProxyHeaders is a set of header keys that may hold a proxied hostname value for the request.
		SSLRedirect:             true,                                            // If SSLRedirect is set to true, then only allow HTTPS requests. Default is false.
		SSLTemporaryRedirect:    false,                                           // If SSLTemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
		SSLHost:                 "ssl.example.com",                               // SSLHost is the host name that is used to redirect HTTP requests to HTTPS. Default is "", which indicates to use the same host.
		SSLProxyHeaders:         map[string]string{"X-Forwarded-Proto": "https"}, // SSLProxyHeaders is set of header keys with associated values that would indicate a valid HTTPS request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
		STSSeconds:              315360000,                                       // STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
		STSIncludeSubdomains:    true,                                            // If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
		STSPreload:              true,                                            // If STSPreload is set to true, the `preload` flag will be appended to the Strict-Transport-Security header. Default is false.
		ForceSTSHeader:          false,                                           // STS header is only included when the connection is HTTPS. If you want to force it to always be added, set to true. `IsDevelopment` still overrides this. Default is false.
		FrameDeny:               true,                                            // If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
		CustomFrameOptionsValue: "SAMEORIGIN",                                    // CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option. Default is "".
		ContentTypeNosniff:      true,                                            // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
		BrowserXssFilter:        true,                                            // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
		CustomBrowserXssValue:   "1; report=https://example.com/xss-report",      // CustomBrowserXssValue allows the X-XSS-Protection header value to be set with a custom value. This overrides the BrowserXssFilter option. Default is "".
		ContentSecurityPolicy:   "default-src 'self'",                            // ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "".
		// PublicKey: `pin-sha256="base64+primary=="; pin-sha256="base64+backup=="; max-age=5184000; includeSubdomains; report-uri="https://www.example.com/hpkp-report"`, // PublicKey implements HPKP to prevent MITM attacks with forged certificates. Default is "".
		//ReferrerPolicy: "same-origin" // ReferrerPolicy allows the Referrer-Policy header with the value to be set with a custom value. Default is "".
		//***triger***
		IsDevelopment: true, // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
	})

	category_rout := secure.New(secure.Options{
		AllowedHosts:            []string{"ssl.example.com"},                     // AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
		HostsProxyHeaders:       []string{"X-Forwarded-Hosts"},                   // HostsProxyHeaders is a set of header keys that may hold a proxied hostname value for the request.
		SSLRedirect:             true,                                            // If SSLRedirect is set to true, then only allow HTTPS requests. Default is false.
		SSLTemporaryRedirect:    false,                                           // If SSLTemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
		SSLHost:                 "ssl.example.com",                               // SSLHost is the host name that is used to redirect HTTP requests to HTTPS. Default is "", which indicates to use the same host.
		SSLProxyHeaders:         map[string]string{"X-Forwarded-Proto": "https"}, // SSLProxyHeaders is set of header keys with associated values that would indicate a valid HTTPS request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
		STSSeconds:              315360000,                                       // STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
		STSIncludeSubdomains:    true,                                            // If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
		STSPreload:              true,                                            // If STSPreload is set to true, the `preload` flag will be appended to the Strict-Transport-Security header. Default is false.
		ForceSTSHeader:          false,                                           // STS header is only included when the connection is HTTPS. If you want to force it to always be added, set to true. `IsDevelopment` still overrides this. Default is false.
		FrameDeny:               true,                                            // If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
		CustomFrameOptionsValue: "SAMEORIGIN",                                    // CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option. Default is "".
		ContentTypeNosniff:      true,                                            // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
		BrowserXssFilter:        true,                                            // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
		CustomBrowserXssValue:   "1; report=https://example.com/xss-report",      // CustomBrowserXssValue allows the X-XSS-Protection header value to be set with a custom value. This overrides the BrowserXssFilter option. Default is "".
		ContentSecurityPolicy:   "default-src 'self'",                            // ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "".
		// PublicKey: `pin-sha256="base64+primary=="; pin-sha256="base64+backup=="; max-age=5184000; includeSubdomains; report-uri="https://www.example.com/hpkp-report"`, // PublicKey implements HPKP to prevent MITM attacks with forged certificates. Default is "".
		//ReferrerPolicy: "same-origin" // ReferrerPolicy allows the Referrer-Policy header with the value to be set with a custom value. Default is "".
		//***triger***
		IsDevelopment: true, // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
	})

	goods_rout := secure.New(secure.Options{
		AllowedHosts:            []string{"ssl.example.com"},                     // AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
		HostsProxyHeaders:       []string{"X-Forwarded-Hosts"},                   // HostsProxyHeaders is a set of header keys that may hold a proxied hostname value for the request.
		SSLRedirect:             true,                                            // If SSLRedirect is set to true, then only allow HTTPS requests. Default is false.
		SSLTemporaryRedirect:    false,                                           // If SSLTemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
		SSLHost:                 "ssl.example.com",                               // SSLHost is the host name that is used to redirect HTTP requests to HTTPS. Default is "", which indicates to use the same host.
		SSLProxyHeaders:         map[string]string{"X-Forwarded-Proto": "https"}, // SSLProxyHeaders is set of header keys with associated values that would indicate a valid HTTPS request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
		STSSeconds:              315360000,                                       // STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
		STSIncludeSubdomains:    true,                                            // If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
		STSPreload:              true,                                            // If STSPreload is set to true, the `preload` flag will be appended to the Strict-Transport-Security header. Default is false.
		ForceSTSHeader:          false,                                           // STS header is only included when the connection is HTTPS. If you want to force it to always be added, set to true. `IsDevelopment` still overrides this. Default is false.
		FrameDeny:               true,                                            // If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
		CustomFrameOptionsValue: "SAMEORIGIN",                                    // CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option. Default is "".
		ContentTypeNosniff:      true,                                            // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
		BrowserXssFilter:        true,                                            // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
		CustomBrowserXssValue:   "1; report=https://example.com/xss-report",      // CustomBrowserXssValue allows the X-XSS-Protection header value to be set with a custom value. This overrides the BrowserXssFilter option. Default is "".
		ContentSecurityPolicy:   "default-src 'self'",                            // ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "".
		// PublicKey: `pin-sha256="base64+primary=="; pin-sha256="base64+backup=="; max-age=5184000; includeSubdomains; report-uri="https://www.example.com/hpkp-report"`, // PublicKey implements HPKP to prevent MITM attacks with forged certificates. Default is "".
		//ReferrerPolicy: "same-origin" // ReferrerPolicy allows the Referrer-Policy header with the value to be set with a custom value. Default is "".
		//***triger***
		IsDevelopment: true, // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
	})

	app1 := submit_rout.Handler(http.HandlerFunc(submit_ctrl))
	app2 := setAddress_rout.Handler(http.HandlerFunc(setAddress_ctrl))
	app3 := category_rout.Handler(http.HandlerFunc(category_ctrl))
	app4 := goods_rout.Handler(http.HandlerFunc(goods_ctrl))
	http.Handle("/submit", app1)
	http.Handle("/setAddress", app2)
	http.Handle("/category", app3)
	http.Handle("/goods", app4)
	log.Fatal(http.ListenAndServe(":80", nil))

}
