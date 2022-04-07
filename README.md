# robFoodDD
上海疫情买菜难，该脚本可帮助自动化抢购，接口调用存在封号风险，若非实在缺菜不建议使用！

该项目为本人第一个Golang项目，主要目的为学习语法与规范，如有代码相关问题欢迎讨论指导！

另外由于时间仓促，一些容错逻辑较为粗糙，结账商品数据未完整分析，可能存在部分类型商品无法结账的情况(如有发现可提issue)，后续有时间的话可能会继续完善。

## 使用方式
在main.go的main函数中修改该行代码
```
err := session.InitSession("DDXQSESSID=xxxxxxxxxxx", "xxxxxxxxxxxxx")
```
其中第一个参数为叮咚登录cookie，需要抓包获取，形式为```"DDXQSESSID=xxxxxxxxxxx""```

第二个参数为通知用的bark id，下载bark后从app界面获取

<img src="./assets/bark.jpg" width="300">

开始运行后按命令行提示操作即可。

## 声明
本项目仅供学习交流，严禁用作商业行为，特别禁止黄牛加价代抢等！

因违法违规等不当使用导致的后果与本人无关，如有任何问题可联系本人删除！