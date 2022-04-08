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

type ReserveTime struct {
	StartTimestamp int    `json:"start_timestamp"`
	EndTimestamp   int    `json:"end_timestamp"`
	SelectMsg      string `json:"select_msg"`
}

func (s *DingdongSession) GetMultiReserveTime(ddmcUid string) (error, []ReserveTime) {
	urlPath := "https://maicai.api.ddxq.mobi/order/getMultiReserveTime"
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
		}
		products = append(products, prod)
	}
	productsList := [][]map[string]interface{}{
		products,
	}
	productsJson, _ := json.Marshal(productsList)
	productsStr := string(productsJson)
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
	data.Add("group_config_id", "")
	data.Add("products", productsStr)
	data.Add("isBridge", "false")

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
	req.Header.Set("ddmc-uid", ddmcUid)
	req.Header.Set("cookie", s.Cookie)
	resp, err := s.Client.Do(req)
	if err != nil {
		return err, nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var reserveTimeList []ReserveTime
		result := gjson.Parse(string(body))
		for _, reserveTimeInfo := range result.Get("data.0.time.0.times").Array() {
			if reserveTimeInfo.Get("disableType").Num == 0 {
				reserveTime := ReserveTime{
					StartTimestamp: int(reserveTimeInfo.Get("start_timestamp").Num),
					EndTimestamp:   int(reserveTimeInfo.Get("end_timestamp").Num),
					SelectMsg:      reserveTimeInfo.Get("select_msg").Str,
				}
				reserveTimeList = append(reserveTimeList, reserveTime)
			}
		}
		return nil, reserveTimeList
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body)), nil
	}

}
