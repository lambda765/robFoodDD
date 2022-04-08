package dd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type DingdongSession struct {
	Address      Address      `json:"address"`
	BarkId       string       `json:"bark_id"`
	Client       *http.Client `json:"client"`
	Cookie       string       `json:"cookie"`
	Cart         Cart         `json:"cart"`
	Order        Order        `json:"order"`
	PackageOrder PackageOrder `json:"package_order"`
	PayType      int          `json:"pay_type"`
	CartMode     int          `json:"cart_mode"`
}

func (s *DingdongSession) InitSession(cookie, barkId string, addressNum, payMethod, settlementMode int) error {
	fmt.Println("########## 初始化 ##########")
	s.Client = &http.Client{}
	s.Cookie = cookie
	s.BarkId = barkId
	err, addrList := s.GetAddress()
	if err != nil {
		return err
	}

	addrLenNum := len(addrList) // app中获取的地址列表长度
	if addrLenNum == 0 {
		return errors.New("未查询到有效收货地址，请前往app添加或检查cookie是否正确！")
	}
	//fmt.Println("########## 选择收货地址 ##########")
	for i, addr := range addrList {
		fmt.Printf("[%v] %s %s %s %s \n", i, addr.Name, addr.AddrDetail, addr.UserName, addr.Mobile)
	}
	//var index int
	/*
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
		}*/
	if addressNum < 0 || addressNum >= addrLenNum {
		log.Fatalln("请根据地址列表，选择合适的地址索引。")
	}
	s.Address = addrList[addressNum]
	fmt.Printf("已选配送地址：%v\n", s.Address)
	//fmt.Println("########## 选择支付方式 ##########")

	switch payMethod {
	case 1:
		s.PayType = 2
		fmt.Println("支付方式:支付宝")
		break
	case 2:
		s.PayType = 4
		fmt.Println("支付方式:微信支付")
		break
	default:
		s.PayType = 4
		fmt.Println("支付方式:微信支付")
	}

	/*for true {
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
	}*/
	//fmt.Println("########## 选择购物车商品结算模式 ##########")

	switch settlementMode {
	case 1:
		s.CartMode = 1
		fmt.Println("结算所有有效商品（不包括换购）")
		break
	case 2:
		s.CartMode = 2
		fmt.Println("结算所有勾选商品（包括换购)")
		break
	default:
		s.CartMode = 1
		fmt.Println("结算所有有效商品（不包括换购）")
		break
	}

	/*for true {
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
	}*/
	return nil
}
