package qyxw

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/yuanpengchao/go-curl"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// 单个图文消息结构
type MsgNews struct {
	Title       string `json:"title"`       // "title" : "中秋节礼品领取",
	Description string `json:"description"` //"description" : "今年中秋节公司有豪礼相送",
	URL         string `json:"url"`         // "url" : "URL",
	PicURL      string `json:"picurl"`      //"picurl" : "http://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png"
}

type QywxRobot struct {
}

func NewQywxRobot() {

}

// 发送文本消息
func (wx *QywxRobot) SendText(ctx context.Context, msg string, mentionedMobileList []string, robotHookURL string) {
	type MsgRobotsMsgTextData struct {
		Content             string   `json:"content"`
		MentionedMobileList []string `json:"mentioned_mobile_list"`
	}
	var data = MsgRobotsMsgTextData{Content: msg, MentionedMobileList: mentionedMobileList}
	postData := map[string]interface{}{
		"msgtype": "text",
		"text":    data,
	}
	wx.sendMsg(ctx, postData, robotHookURL)
}

// 发送图片
func (wx *QywxRobot) SendImg(ctx context.Context, imgPath string, robotHookURL string) {
	//读原图片
	ff, err := os.Open(imgPath)
	defer ff.Close()
	if err != nil {
		log.Println("读取文件内容出错", "err", err)
		return
	}
	var sourceBuffer []byte
	sourceBuffer, err = ioutil.ReadAll(ff)
	if err != nil {
		log.Println("读取文件内容出错", "err", err)
		return
	}
	//base64压缩
	base64Str := base64.StdEncoding.EncodeToString(sourceBuffer)
	//写入临时文件
	type imageData struct {
		Base64 string `json:"base64"`
		Md5    string `json:"md5"`
	}
	w := md5.New()
	_, err = io.WriteString(w, string(sourceBuffer))
	if err != nil {
		log.Println("转换内容到md5失败了", "err", err)
		return
	}
	//将str写入到w中
	imgMd5 := fmt.Sprintf("%x", w.Sum(nil))
	// data:image/png;base64,
	data := imageData{Base64: base64Str, Md5: imgMd5}
	postData := map[string]interface{}{
		"msgtype": "image",
		"image":   data,
	}
	wx.sendMsg(ctx, postData, robotHookURL)
}

// 发送图文消息
func (wx *QywxRobot) SendNews(ctx context.Context, msgNews []MsgNews, robotHookURL string) {
	type news struct {
		Articles []MsgNews `json:"articles"`
	}
	data := news{Articles: msgNews}
	postData := map[string]interface{}{
		"msgtype": "news",
		"news":    data,
	}
	wx.sendMsg(ctx, postData, robotHookURL)
}

// 执行发送请求
func (wx *QywxRobot) sendMsg(ctx context.Context, postData map[string]interface{}, webHook string) {
	log.Println("发送企业微信机器人消息：", "data", postData)
	request := curl.NewRequest()
	if len(webHook) == 0 {
		log.Println("机器人WebHook必传！：", "robotHookURL", webHook)
		return
	}
	_, err := request.
		SetPostData(postData).
		SetUrl(webHook).
		SetInsecureSkipVerify(true).
		Post()
	if err != nil {
		log.Println("发送微信机器人消息错误:", "err", err)
		return
	}
}
