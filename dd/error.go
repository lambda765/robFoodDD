package dd

import "errors"

var TimeExpireErr = errors.New("送达时间已失效")
var OOSErr = errors.New("部分商品已缺货")
var BusyErr = errors.New("当前人多拥挤")
var DataLoadErr = errors.New("部分数据加载失败")
var ProdInfoErr = errors.New("商品信息有变化")
