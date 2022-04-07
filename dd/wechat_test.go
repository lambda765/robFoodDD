package dd

import "testing"

// 通过微信公众平台(测试号)推送
// 访问 https://mp.weixin.qq.com/debug/cgi-bin/sandbox?t=sandbox/login 微信扫一扫开通并关注测试号
func TestWeChatPush(t *testing.T) {
	err := WeChatPush("wx*****", "****", "微信公众平台(测试号)推送测试")
	t.Log(err)
}
