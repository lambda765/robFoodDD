package dd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type Msg struct {
	Filter struct {
		IsToAll bool   `json:"is_to_all"`
		TagId   string `json:"tag_id"`
	} `json:"filter"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	MsgType string `json:"msgtype"`
}

type Result struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgId   int32  `json:"msg_id"`
}

// cache token
var token *Token

// WeChatAccessToken .
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
func WeChatAccessToken(appId, appSecret string) (string, error) {
	now := time.Now().Unix()
	if token != nil && token.ExpiresIn > now {
		return token.AccessToken, nil
	}

	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appId + "&secret=" + appSecret
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	} else if resp.StatusCode != http.StatusOK {
		return "", errors.New("status code is not 200")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if bytes.Contains(body, []byte("access_token")) {
		token = &Token{}
		err = json.Unmarshal(body, &token)
		if err != nil {
			return "", err
		}

		token.ExpiresIn = now + token.ExpiresIn - 120

		return token.AccessToken, nil
	}

	return "", errors.New("get access token failed")
}

// WeChatPush .
// https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Batch_Sends_and_Originality_Checks.html
func WeChatPush(appId, appSecret, text string) error {
	accessToken, err := WeChatAccessToken(appId, appSecret)
	if err != nil {
		return err
	}

	msg := Msg{}
	msg.Filter.IsToAll = true
	msg.Text.Content = text
	msg.MsgType = "text"
	data, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	url := "https://api.weixin.qq.com/cgi-bin/message/mass/sendall?access_token=" + accessToken
	resp, err := http.Post(url, "application/json;charset=UTF-8", strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if bytes.Contains(body, []byte("errcode")) {
		result := Result{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			return err
		}

		if result.ErrCode != 0 {
			err = errors.New(result.ErrMsg)
		}
	}

	return err
}
