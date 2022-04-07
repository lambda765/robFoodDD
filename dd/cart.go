package dd

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Product struct {
	Id               string                   `json:"id"`
	ProductName      string                   `json:"-"`
	Price            string                   `json:"price"`
	Count            int                      `json:"count"`
	Sizes            []map[string]interface{} `json:"sizes"`
	TotalPrice       string                   `json:"total_money"`
	OriginPrice      string                   `json:"origin_price"`
	TotalOriginPrice string                   `json:"total_origin_money"`
}

func parseProduct(productMap gjson.Result) (error, Product) {
	var sizes []map[string]interface{}
	for _, size := range productMap.Get("sizes").Array() {
		sizes = append(sizes, size.Value().(map[string]interface{}))
	}
	product := Product{
		Id:          productMap.Get("id").Str,
		ProductName: productMap.Get("product_name").Str,
		Price:       productMap.Get("price").Str,
		Count:       int(productMap.Get("count").Num),
		TotalPrice:  productMap.Get("total_price").Str,
		OriginPrice: productMap.Get("origin_price").Str,
		Sizes:       sizes,
	}
	return nil, product
}

type Cart struct {
	ProdList        []Product `json:"effective_products"`
	ParentOrderSign string    `json:"parent_order_sign"`
}

func (s *DingdongSession) GetEffProd(result gjson.Result) error {
	var effProducts []Product
	effective := result.Get("data.product.effective").Array()
	for _, effProductMap := range effective {
		for _, productMap := range effProductMap.Get("products").Array() {
			_, product := parseProduct(productMap)
			effProducts = append(effProducts, product)
		}
	}
	s.Cart = Cart{
		ProdList:        effProducts,
		ParentOrderSign: result.Get("data.parent_order_info.parent_order_sign").Str,
	}
	return nil
}

func (s *DingdongSession) GetCheckProd(result gjson.Result) error {
	var products []Product
	orderProductList := result.Get("data.new_order_product_list").Array()
	for _, productList := range orderProductList {
		for _, productMap := range productList.Get("products").Array() {
			_, product := parseProduct(productMap)
			products = append(products, product)
		}
	}
	s.Cart = Cart{
		ProdList:        products,
		ParentOrderSign: result.Get("data.parent_order_info.parent_order_sign").Str,
	}
	return nil
}

func (s *DingdongSession) CheckCart() error {
	Url, _ := url.Parse("https://maicai.api.ddxq.mobi/cart/index")
	params := url.Values{}
	params.Set("station_id", s.Address.StationId)
	params.Set("city_number", s.Address.CityNumber)
	params.Set("api_version", "9.49.0")
	params.Set("app_version", "2.81.0")
	params.Set("applet_source", "")
	params.Set("app_client_id", "3")
	params.Set("h5_source", "")
	params.Set("sharer_uid", "")
	params.Set("s_id", "")
	params.Set("openid", "")
	params.Set("is_load", "1")
	params.Set("ab_config", "{\"key_onion\":\"D\",\"key_cart_discount_price\":\"C\"}")

	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	req, _ := http.NewRequest("GET", urlPath, nil)
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
			switch s.CartMode {
			case 1:
				return s.GetEffProd(result)
			case 2:
				return s.GetCheckProd(result)
			default:
				return errors.New("incorrect cart mode")
			}
		case -3000:
			return BusyErr
		default:
			return errors.New(string(body))
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
