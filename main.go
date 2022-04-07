package main

import (
	"fmt"
	"math/rand"
	"robFoodDD/dd"
	"time"
)

func main() {
	session := dd.DingdongSession{}
	err := session.InitSession("DDXQSESSID=xxxxxxxxxxx", "xxxxxxxxxxxxx")
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
			fmt.Println("购物车中无有效商品，请先前往app添加或勾选！")
			return
		}
		for index, prod := range session.Cart.ProdList {
			fmt.Printf("[%v] %s 数量：%v 总价：%s\n", index, prod.ProductName, prod.Count, prod.TotalPrice)
		}
		session.Order.Products = session.Cart.ProdList
		for i := 0; i < 10; i++ {
			fmt.Printf("########## 生成订单信息【%s】 ###########\n", time.Now().Format("15:04:05"))
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
			fmt.Printf("########## 获取可预约时间【%s】 ###########\n", time.Now().Format("15:04:05"))
			err, multiReserveTime := session.GetMultiReserveTime()
			if err != nil {
				fmt.Println(err)
				continue
			}
			if len(multiReserveTime) == 0 {
				sleepInterval := 3 + rand.Intn(6)
				fmt.Printf("暂无可预约时间，%v秒后重试！\n", sleepInterval)
				time.Sleep(time.Duration(sleepInterval) * time.Second)
				continue
			} else {
				fmt.Println("发现可用的配送时段!")
			}
			for _, reserveTime := range multiReserveTime {
				session.UpdatePackageOrder(reserveTime)
			OrderLoop:
				for i := 0; i < 15; i++ {
					fmt.Printf("########## 提交订单中【%s】 ###########\n", time.Now().Format("15:04:05"))
					err = session.AddNewOrder()
					switch err {
					case nil:
						fmt.Println("抢购成功，请前往app付款！")
						for true {
							err = session.PushSuccess()
							if err == nil {
								break
							} else {
								fmt.Println(err)
							}
							time.Sleep(1 * time.Second)
						}
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
