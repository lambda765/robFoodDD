package dd

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Address struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	StationId  string  `json:"station_id"`
	CityNumber string  `json:"city_number"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	UserName   string  `json:"user_name"`
	Mobile     string  `json:"mobile"`
	Address    string  `json:"address"`
	AddrDetail string  `json:"addr_detail"`
}

func parseAddress(addressMap gjson.Result) (error, Address) {
	address := Address{}
	address.Id = addressMap.Get("id").Str
	address.Name = addressMap.Get("location.name").Str
	address.StationId = addressMap.Get("station_id").Str
	address.CityNumber = addressMap.Get("city_number").Str
	address.Longitude = addressMap.Get("location.location.0").Num
	address.Latitude = addressMap.Get("location.location.1").Num
	address.UserName = addressMap.Get("user_name").Str
	address.Mobile = addressMap.Get("mobile").Str
	address.Address = addressMap.Get("location.address").Str
	address.AddrDetail = addressMap.Get("addr_detail").Str
	return nil, address
}

func (s *DingdongSession) GetAddress() (error, []Address) {
	Url, _ := url.Parse("https://sunquan.api.ddxq.mobi/api/v1/user/address/")
	params := url.Values{}
	params.Set("api_version", "9.49.0")
	params.Set("app_version", "2.81.0")
	params.Set("applet_source", "")
	params.Set("app_client_id", "3")
	params.Set("h5_source", "")
	params.Set("sharer_uid", "")
	params.Set("s_id", "")
	params.Set("openid", "")
	params.Set("source_type", "5")

	Url.RawQuery = params.Encode()
	urlPath := Url.String()

	req, _ := http.NewRequest("GET", urlPath, nil)
	req.Header.Set("Host", "sunquan.api.ddxq.mobi")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 9; LIO-AN00 Build/LIO-AN00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/92.0.4515.131 Mobile Safari/537.36 xzone/9.47.0 station_id/null")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
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
		return err, nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	resp.Body.Close()
	if resp.StatusCode == 200 {
		var addressList []Address
		result := gjson.Parse(string(body))
		validAddress := result.Get("data.valid_address").Array()
		for _, addressMap := range validAddress {
			_, address := parseAddress(addressMap)
			addressList = append(addressList, address)
		}
		return nil, addressList
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body)), nil
	}
}
