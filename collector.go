package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

const CARREFOUR_ES_API_KEY = "3_Ns3U5-wXeiSQL-vZtu1Fd2DpWBsEdB78mYs2dn0_kyFFwwSJAZZd1EHUm9kodfND"
const CARREFOUR_ES_AUTHORITY = "ss-mya-p3.carrefour.es"
const CARREFOU_ES_PURCHASES_AUTHORITY = "apig-pro.secure.hd.carrefour.es"
const CARREFOUR_ES_ORIGIN = "https://www.carrefour.es"
const CARREFOUR_ES_REFERER = "https://www.carrefour.es/"
const SEARCH_DATE_FROM = "2000-08-20T14:37:15.672Z"
const SEARCH_DATE_TO = "2023-08-21T14:37:15.672Z"
const CARREFOUR_ES_LOGIN_URL = "https://ss-mya-p3.carrefour.es/accounts.login"

var client = resty.New()

type PurchaseTypeDetail struct {
	code        string
	description string
}

type SessionInfo struct {
	cookieName  string
	cookieValue string
}

type LoginResponse struct {
	statusCode  int
	sessionInfo SessionInfo
}

type OrderSourceDetail struct {
	code        string
	description string
}

type Purchase struct {
	amount            float64 `db:"amount"`
	customerId        string  `db:"customer_id"`
	mallAddress       string  `db:"mall_address"`
	mallId            string  `db:"mall_id"`
	mallName          string  `db:"mall_name"`
	orderSource       string  `db:"order_source"`
	orderSourceDetail struct {
		code        string
		description string
	}
	purchaseDate       string `db:"purchase_date"`
	purchaseId         string `db:"purchase_id"`
	purchaseType       string `db:"purchase_type"`
	purchaseTypeDetail struct {
		code        string
		description string
	}
}

type PurchaseResponse struct {
	purchases []Purchase
}

func CollectFromCarrefour(username string, password string) ([]Purchase, error) {
	cookie := Login(username, password)
	token := GetJWTToken(cookie)
	purchases, err := ListPurchases(token)
	fmt.Printf("Found %d purchases: %v\n", len(purchases), purchases)
	if err != nil {
		fmt.Printf("Error Listing Purchases: %v\n", err)
	}
	return purchases, nil
}

func MapToPurchase(input map[string]interface{}) Purchase {
	var purchase Purchase
	purchase.amount = input["amount"].(float64)
	purchase.customerId = input["customerId"].(string)
	purchase.mallAddress = input["mallAddress"].(string)
	purchase.mallId = input["mallId"].(string)
	purchase.mallName = input["mallName"].(string)
	purchase.orderSource = input["orderSource"].(string)
	orderDetail := input["orderSourceDetail"].(map[string]interface{})
	purchase.orderSourceDetail = OrderSourceDetail{
		code:        orderDetail["code"].(string),
		description: orderDetail["description"].(string),
	}
	purchase.purchaseDate = input["purchaseDate"].(string)
	purchase.purchaseId = input["purchaseId"].(string)
	purchase.purchaseType = input["purchaseType"].(string)
	purchaseTypeDetail := input["purchaseTypeDetail"].(map[string]interface{})
	purchase.purchaseTypeDetail = PurchaseTypeDetail{
		code:        purchaseTypeDetail["code"].(string),
		description: purchaseTypeDetail["description"].(string),
	}
	return purchase
}

func ListPurchases(token string) ([]Purchase, error) {
	var offset = 0
	var pageSize = 10
	fmt.Printf("Listing purchases with pageSize %d\n", pageSize)
	var allPurchases []Purchase
	for {
		fmt.Printf("Querying offset %d\n", offset)
		resp, err := client.R().
			SetHeader("authority", CARREFOU_ES_PURCHASES_AUTHORITY).
			SetHeader("accept", "application/json, text/plain, */*").
			SetHeader("accept-language", "en-US,en;q=0.9").
			SetHeader("authorization", "bearer "+token).
			SetHeader("content-type", "application/json").
			SetHeader("origin", "https://www.carrefour.es").
			SetHeader("pragma", "no-cache").
			SetHeader("referer", "https://www.carrefour.es/").
			SetHeader("requestorigin", "MYA").
			SetQueryParams(map[string]string{
				"from":               SEARCH_DATE_FROM,
				"to":                 SEARCH_DATE_TO,
				"atgfOffset":         "0",
				"atgnfOffset":        "0",
				"currentAtgfOrders":  "0",
				"currentAtgnfOrders": "0",
				"currentTickets":     "0",
				"ticketOffset":       fmt.Sprintf("%d", offset),
				"count":              fmt.Sprintf("%d", pageSize),
			}).
			Get("https://apig-pro.secure.hd.carrefour.es/md-purchasesAccount-v1/purchases")

		var ret []Purchase = extractListOfPurchases(resp.Body())
		allPurchases = append(allPurchases, ret...)
		if err != nil {
			fmt.Printf("Error Listing Purchases: %v\n", err)
		}
		if len(ret) < pageSize {
			break
		}
		offset = offset + pageSize
	}

	return allPurchases, nil
}

func extractListOfPurchases(body []byte) []Purchase {
	var ret []Purchase = []Purchase{}
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	purchases := result["purchases"]
	for _, key := range purchases.([]interface{}) {
		element := key.(map[string]interface{})
		purchase := MapToPurchase(element)
		ret = append(ret, purchase)
	}

	return ret
}

func Login(username string, password string) string {
	resp, err := client.R().
		SetHeader("authority", CARREFOUR_ES_AUTHORITY).
		SetHeader("accept", "*/*").
		SetHeader("accept-language", "en-US,en;q=0.9").
		SetHeader("origin", CARREFOUR_ES_ORIGIN).
		SetHeader("referer", CARREFOUR_ES_REFERER).
		SetHeader("content-type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"APIKey":            CARREFOUR_ES_API_KEY,
			"loginID":           username,
			"password":          password,
			"sessionExpiration": "-1",
			"include":           "profile,data,emails,subscriptions,preferences,",
			"includeUserInfo":   "true",
			"format":            "json",
			"loginMode":         "standard",
			"source":            "showScreenSet",
			"sdkBuild":          "15170",
		}).
		Post(CARREFOUR_ES_LOGIN_URL)

	if err != nil {
		fmt.Printf("Login Error: %v\n", err)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body(), &result)
	sessionInfo := result["sessionInfo"].(map[string]interface{})
	return sessionInfo["cookieValue"].(string)
}

func GetJWTToken(cookie string) string {
	resp, err := client.R().
		SetHeader("authority", "ss-mya-p3.carrefour.es").
		SetHeader("accept", "*/*").
		SetHeader("accept-language", "en-US,en;q=0.9").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("origin", "https://www.carrefour.es").
		SetHeader("referer", "https://www.carrefour.es/").
		SetFormData(map[string]string{
			"APIKey":      CARREFOUR_ES_API_KEY,
			"login_token": cookie,
			"fields":      "data.GR,profile.email,data.DQ,data.acceptedCustomerPolicies,data.ID_ATG",
			"expiration":  "1800",
			"authMode":    "cookie",
			"pageURL":     "https://www.carrefour.es/access/#/area-privada/dashboard?source=login_nonFood&redirect=https://www.carrefour.es/&back=https://www.carrefour.es/",
			"sdkBuild":    "15170",
			"format":      "json",
		}).
		Post("https://ss-mya-p3.carrefour.es/accounts.getJWT")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	var result map[string]interface{}
	json.Unmarshal(resp.Body(), &result)

	return result["id_token"].(string)
}
