package dd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Order struct {
	Products []Product `json:"products"`
	Price    string    `json:"price"`
}

type Package struct {
	FirstSelectedBigTime string                   `json:"first_selected_big_time"`
	Products             []map[string]interface{} `json:"products"`
	EtaTraceId           string                   `json:"eta_trace_id"`
	PackageId            int                      `json:"package_id"`
	ReservedTimeStart    int                      `json:"reserved_time_start"`
	ReservedTimeEnd      int                      `json:"reserved_time_end"`
	SoonArrival          int                      `json:"soon_arrival"`
	PackageType          int                      `json:"package_type"`
}

type PaymentOrder struct {
	ReservedTimeStart    int    `json:"reserved_time_start"`
	ReservedTimeEnd      int    `json:"reserved_time_end"`
	FreightDiscountMoney string `json:"freight_discount_money"`
	FreightMoney         string `json:"freight_money"`
	OrderFreight         string `json:"order_freight"`
	AddressId            string `json:"address_id"`
	UsedPointNum         int    `json:"used_point_num"`
	ParentOrderSign      string `json:"parent_order_sign"`
	PayType              int    `json:"pay_type"`
	OrderType            int    `json:"order_type"`
	IsUseBalance         int    `json:"is_use_balance"`
	ReceiptWithoutSku    string `json:"receipt_without_sku"`
	Price                string `json:"price"`
}

type PackageOrder struct {
	Packages     []*Package   `json:"packages"`
	PaymentOrder PaymentOrder `json:"payment_order"`
}

type AddNewOrderReturnData struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    struct {
		PackageOrder     PackageOrder `json:"package_order"`
		StockoutProducts []Product    `json:"stockout_products"`
	} `json:"data"`
}

func (s *DingdongSession) CheckOrder() error {
	urlPath := "https://maicai.api.ddxq.mobi/order/checkOrder"

	var products []map[string]interface{}
	for _, product := range s.Order.Products {
		prod := map[string]interface{}{
			"id":                   product.Id,
			"total_money":          product.TotalPrice,
			"total_origin_money":   product.OriginPrice,
			"count":                product.Count,
			"price":                product.Price,
			"instant_rebate_money": "0.00",
			"origin_price":         product.OriginPrice,
			"sizes":                product.Sizes,
		}
		products = append(products, prod)
	}
	packagesInfo := []map[string]interface{}{
		{
			"package_type": 1,
			"package_id":   1,
			"products":     products,
		},
	}
	packagesJson, _ := json.Marshal(packagesInfo)
	packagesStr := string(packagesJson)

	data := url.Values{}
	data.Add("station_id", s.Address.StationId)
	data.Add("city_number", s.Address.CityNumber)
	data.Add("api_version", "9.49.0")
	data.Add("app_version", "2.81.0")
	data.Add("applet_source", "")
	data.Add("app_client_id", "3")
	data.Add("h5_source", "")
	data.Add("sharer_uid", "")
	data.Add("s_id", "")
	data.Add("openid", "")
	data.Add("user_ticket_id", "default")
	data.Add("freight_ticket_id", "default")
	data.Add("is_use_point", "0")
	data.Add("is_use_balance", "0")
	data.Add("is_buy_vip", "0")
	data.Add("coupons_id", "")
	data.Add("is_buy_coupons", "0")
	data.Add("packages", packagesStr)
	data.Add("check_order_type", "0")
	data.Add("is_support_merge_payment", "1")
	data.Add("showData", "true")
	data.Add("showMsg", "false")

	req, _ := http.NewRequest("POST", urlPath, strings.NewReader(data.Encode()))
	req.Header.Set("Host", "maicai.api.ddxq.mobi")
	req.Header.Set("ddmc-city-number", s.Address.CityNumber)
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 9; LIO-AN00 Build/LIO-AN00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/92.0.4515.131 Mobile Safari/537.36 xzone/9.47.0 station_id/null")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("ddmc-os-version", "undefined")
	req.Header.Set("ddmc-channel", "undefined")
	req.Header.Set("ddmc-build-version", "2.81.0")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("ddmc-app-client-id", "3")
	req.Header.Set("ddmc-api-version", "9.49.0")
	req.Header.Set("ddmc-station-id", s.Address.StationId)
	req.Header.Set("origin", "https://wx.m.ddxq.mobi")
	req.Header.Set("x-requested-with", "com.yaya.zone")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", "https://wx.m.ddxq.mobi/")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("cookie", s.Cookie)
	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode == 200 {
		result := gjson.Parse(string(body))
		switch result.Get("code").Num {
		case 0:
			s.Order.Price = result.Get("data.order.total_money").Str
			return nil
		case -3000:
			return BusyErr
		case -3100:
			return DataLoadErr
		default:
			return errors.New(string(body))
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}

func (s *DingdongSession) GeneratePackageOrder() {
	var products []map[string]interface{}
	for _, product := range s.Order.Products {
		prod := map[string]interface{}{
			"id":                   product.Id,
			"total_money":          product.TotalPrice,
			"total_origin_money":   product.OriginPrice,
			"count":                product.Count,
			"price":                product.Price,
			"instant_rebate_money": "0.00",
			"origin_price":         product.OriginPrice,
			"sizes":                product.Sizes,
		}
		products = append(products, prod)
	}

	p := Package{
		FirstSelectedBigTime: "0",
		Products:             products,
		EtaTraceId:           "",
		PackageId:            1,
		PackageType:          1,
	}
	paymentOrder := PaymentOrder{
		FreightDiscountMoney: "5.00",
		FreightMoney:         "5.00",
		OrderFreight:         "0.00",
		AddressId:            s.Address.Id,
		UsedPointNum:         0,
		ParentOrderSign:      s.Cart.ParentOrderSign,
		PayType:              s.PayType,
		OrderType:            1,
		IsUseBalance:         0,
		ReceiptWithoutSku:    "1",
		Price:                s.Order.Price,
	}
	packageOrder := PackageOrder{
		Packages: []*Package{
			&p,
		},
		PaymentOrder: paymentOrder,
	}
	s.PackageOrder = packageOrder
}

func (s *DingdongSession) UpdatePackageOrder(reserveTime ReserveTime) {
	s.PackageOrder.PaymentOrder.ReservedTimeStart = reserveTime.StartTimestamp
	s.PackageOrder.PaymentOrder.ReservedTimeEnd = reserveTime.EndTimestamp
	for _, p := range s.PackageOrder.Packages {
		p.ReservedTimeStart = reserveTime.StartTimestamp
		p.ReservedTimeEnd = reserveTime.EndTimestamp
	}
}

func (s *DingdongSession) AddNewOrder() error {
	urlPath := "https://maicai.api.ddxq.mobi/order/addNewOrder"

	packageOrderJson, _ := json.Marshal(s.PackageOrder)
	packageOrderStr := string(packageOrderJson)

	data := url.Values{}
	data.Add("uid", "")
	data.Add("station_id", s.Address.StationId)
	data.Add("city_number", s.Address.CityNumber)
	data.Add("api_version", "9.49.0")
	data.Add("app_version", "2.81.0")
	data.Add("applet_source", "")
	data.Add("app_client_id", "3")
	data.Add("h5_source", "")
	data.Add("sharer_uid", "")
	data.Add("s_id", "")
	data.Add("openid", "")
	data.Add("package_order", packageOrderStr)
	data.Add("showData", "true")
	data.Add("showMsg", "false")
	data.Add("ab_config", "{\"key_onion\":\"C\"}")
	req, _ := http.NewRequest("POST", urlPath, strings.NewReader(data.Encode()))
	req.Header.Set("Host", "maicai.api.ddxq.mobi")
	req.Header.Set("ddmc-city-number", s.Address.CityNumber)
	req.Header.Set("user-agent", fmt.Sprintf("Mozilla/5.0 (Linux; Android 9; LIO-AN00 Build/LIO-AN00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/92.0.4515.131 Mobile Safari/537.36 xzone/9.47.0 station_id/%s", s.Address.StationId))
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("ddmc-os-version", "undefined")
	req.Header.Set("ddmc-channel", "undefined")
	req.Header.Set("ddmc-build-version", "2.81.0")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("ddmc-app-client-id", "3")
	req.Header.Set("ddmc-api-version", "9.49.0")
	req.Header.Set("ddmc-station-id", s.Address.StationId)
	req.Header.Set("origin", "https://wx.m.ddxq.mobi")
	req.Header.Set("x-requested-with", "com.yaya.zone")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", "https://wx.m.ddxq.mobi/")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("cookie", s.Cookie)
	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode == 200 {
		result := AddNewOrderReturnData{}
		err := json.Unmarshal(body, &result)
		if err != nil {
			return err
		}
		switch result.Code {
		case 0:
			return nil
		case 5001:
			s.PackageOrder = result.Data.PackageOrder
			return OOSErr
		case 5003:
			return ProdInfoErr
		case 5004:
			fmt.Println(result.Msg)
			return TimeExpireErr
		default:
			return errors.New(string(body))
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
