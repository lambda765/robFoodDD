package dd

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type DingdongSession struct {
	Address         Address      `json:"address"`
	BarkId          string       `json:"bark_id"`
	WeChatAppId     string       `json:"wechat_app_id"`
	WeChatAppSecret string       `json:"wechat_app_secret"`
	Client          *http.Client `json:"client"`
	Cookie          string       `json:"cookie"`
	Cart            Cart         `json:"cart"`
	Order           Order        `json:"order"`
	PackageOrder    PackageOrder `json:"package_order"`
	PayType         int          `json:"pay_type"`
	CartMode        int          `json:"cart_mode"`
	NoticeMode      int          `json:"notice_mode"`
}

func (s *DingdongSession) InitSession() error {
	fmt.Println("########## 初始化 ##########")
	s.Client = &http.Client{}
	// https://activity.m.ddxq.mobi/#/coupon?code=VERxg&h5_source=caocao&btnType=jumpApp&path=https://u.100.me/m/maicai&random=4
	// 电脑浏览器从此地址登录后即可从 Cookie 取得 DDXQSESSID
	s.Cookie = os.Getenv("DDXQSESSID")
	if s.Cookie == "" {
		return errors.New("没有配置 DDXQSESSID 环境变量")
	}
	s.Cookie = "DDXQSESSID=" + s.Cookie

	s.BarkId = os.Getenv("BARKID")
	s.WeChatAppId = os.Getenv("WECHATAPPID")
	s.WeChatAppSecret = os.Getenv("WECHATAPPSECRET")

	// 根据不同环境变量设置通知方式，不设置不通知
	if s.BarkId != "" {
		s.NoticeMode = 1
	} else if s.WeChatAppId != "" && s.WeChatAppSecret != "" {
		s.NoticeMode = 2
	}

	err, addrList := s.GetAddress()
	if err != nil {
		return err
	}
	if len(addrList) == 0 {
		return errors.New("未查询到有效收货地址，请前往app添加或检查cookie是否正确！")
	}
	fmt.Println("########## 选择收货地址 ##########")
	for i, addr := range addrList {
		fmt.Printf("[%v] %s %s %s %s \n", i, addr.Name, addr.AddrDetail, addr.UserName, addr.Mobile)
	}
	var index int
	for true {
		fmt.Println("请输入地址序号（0, 1, 2...)：")
		stdin := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanln(stdin, &index)
		if err != nil {
			fmt.Printf("输入有误：%s!\n", err)
		} else if index >= len(addrList) {
			fmt.Println("输入有误：超过最大序号！")
		} else {
			break
		}
	}
	s.Address = addrList[index]
	fmt.Println("########## 选择支付方式 ##########")
	for true {
		fmt.Println("请输入支付方式序号（1：支付宝 2：微信)：")
		stdin := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanln(stdin, &index)
		if err != nil {
			fmt.Printf("输入有误：%s!\n", err)
		} else if index == 1 {
			s.PayType = 2
			break
		} else if index == 2 {
			s.PayType = 4
			break
		} else {
			fmt.Println("输入有误：序号无效！")
		}
	}
	fmt.Println("########## 选择购物车商品结算模式 ##########")
	for true {
		fmt.Println("请输入结算模式序号（1：结算所有有效商品（不包括换购） 2：结算所有勾选商品（包括换购)：")
		stdin := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanln(stdin, &index)
		if err != nil {
			fmt.Printf("输入有误：%s!\n", err)
		} else if index == 1 {
			s.CartMode = 1
			break
		} else if index == 2 {
			s.CartMode = 2
			break
		} else {
			fmt.Println("输入有误：序号无效！")
		}
	}
	return nil
}
