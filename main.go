package main

import (
	"fmt"
	"math/rand"
	"robFoodDD/dd"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	//init config
	dd.InitConfigX()
	fmt.Printf("users : %+v \n", dd.GetConfigX().Users)

	for _, user := range dd.GetConfigX().Users {
		wg.Add(1)
		time.Sleep(5) //间隔账号之间的登录时间，同时运行会报鉴权异常
		go robFood(user, &wg)
	}
	wg.Wait()

}

//robDood: cookie 用户信息 action: 操作动作
func robFood(user dd.UserModel, wg *sync.WaitGroup) {
	session := dd.DingdongSession{}
	err := session.InitSession(user.Cookie, user.BarkId, user.AddressNum, user.PayMethodNum, user.SettlementMode)
	if err != nil {
		fmt.Println(err)
		return
	}
cartLoop:
	for true {
		fmt.Printf("########## 获取购物车中有效商品【%s】 ###########\n", time.Now().Format("15:04:05"))
		err = session.CheckCart()
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		if len(session.Cart.ProdList) == 0 {
			fmt.Printf("%s  ==> 购物车中无有效商品，请先前往app添加或勾选！\n", user.UserName)
			return
		}
		for index, prod := range session.Cart.ProdList {
			fmt.Printf("%s  ==> [%v] %s 数量：%v 总价：%s\n", user.UserName, index, prod.ProductName, prod.Count, prod.TotalPrice)
		}
		session.Order.Products = session.Cart.ProdList
		for i := 0; i < 10; i++ {
			fmt.Printf("%s  ==> ########## 生成订单信息【%s】 ###########\n", user.UserName, time.Now().Format("15:04:05"))
			err = session.CheckOrder()
			if err != nil {
				fmt.Println(err)
				time.Sleep(2 * time.Second)
				continue
			} else {
				break
			}
		}
		if err != nil {
			continue
		}
		fmt.Printf("订单总金额：%v\n", session.Order.Price)
		session.GeneratePackageOrder()
		for i := 0; i < 60; i++ {
			fmt.Printf("%s  ==> ########## 获取可预约时间【%s】 ###########\n", user.UserName, time.Now().Format("15:04:05"))
			err, multiReserveTime := session.GetMultiReserveTime(user.DdmcUid)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if len(multiReserveTime) == 0 {
				sleepInterval := 3 + rand.Intn(6)
				fmt.Printf("%s  ==> 暂无可预约时间，%v秒后重试！\n", user.UserName, sleepInterval)
				time.Sleep(time.Duration(sleepInterval) * time.Second)
				continue
			} else {
				fmt.Printf("%s  ==> 发现可用的配送时段!\n", user.UserName)
			}
			for _, reserveTime := range multiReserveTime {
				session.UpdatePackageOrder(reserveTime)
			OrderLoop:
				for i := 0; i < 15; i++ {
					fmt.Printf("%s  ==> ########## 提交订单中【%s】 ###########\n", user.UserName, time.Now().Format("15:04:05"))
					err = session.AddNewOrder()
					switch err {
					case nil:
						fmt.Printf("%s  ==> 抢购成功，请前往app付款！", user.UserName)
						for true {
							err = session.PushSuccess()
							if err == nil {
								break
							} else {
								fmt.Println(err)
							}
							time.Sleep(1 * time.Second)
						}
						wg.Done()
						return
					case dd.TimeExpireErr:
						fmt.Printf("[%s] %s\n", reserveTime.SelectMsg, err)
						break OrderLoop
					case dd.ProdInfoErr:
						fmt.Println(err)
						continue cartLoop
					default:
						fmt.Println(err)
					}
					time.Sleep(1 * time.Second)
				}
			}
		}
	}
}
