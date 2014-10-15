package DhtCrawler

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	b64 "encoding/base64"
	"fmt"
	curl "github.com/andelf/go-curl"
	"time"
)

type MSQ struct {
	accessKeyId     string
	accessKeySecret string
	accessOwnerId   string
}

func NewMSQ(accessKeyId string, accessKeySecret string, accessOwnerId string) *MSQ {
	mq := new(MSQ)
	mq.accessKeyId = accessKeyId
	mq.accessKeySecret = accessKeySecret
	mq.accessOwnerId = accessOwnerId
	return mq
}

func (mq *MSQ) sign(verb, content_md5, content_type, date_gmt, mqs_headers, request_resource string) string {
	str2sign := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		verb,
		content_md5,
		content_type,
		date_gmt,
		mqs_headers,
		request_resource)
	sign := hash_hmac([]byte(str2sign), []byte(mq.accessKeySecret))
	return "MQS " + mq.accessKeyId + ":" + base64_encode(sign)
}

func (mq *MSQ) addMessage(hash string, queueIdx int) {
	layout := `<?xml version="1.0" encoding="UTF-8"?>
<Message xmlns="http://mqs.aliyuncs.com/doc/v1/">
  <MessageBody>%s</MessageBody>
  <DelaySeconds>0</DelaySeconds>
</Message>
`

	queueName := fmt.Sprintf("dht-infohash-%d", queueIdx)
	verb := "POST"
	content := fmt.Sprintf(layout, hash)

	content_md5 := base64_encode(md5_encode(content))
	content_type := "text/xml;utf-8"
	date_gmt := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 MST")
	mqs_headers := "x-mqs-version:2014-07-08"
	request_resource := "/" + queueName + "/messages"
	sign := mq.sign(verb, content_md5, content_type, date_gmt, mqs_headers, request_resource)

	easy := curl.EasyInit()
	defer easy.Cleanup()

	// set URL to get
	easy.Setopt(curl.OPT_URL, fmt.Sprintf("http://%s.mqs-cn-hangzhou.aliyuncs.com/%s/messages", mq.accessOwnerId, queueName))

	easy.Setopt(curl.OPT_HTTPHEADER, []string{
		fmt.Sprintf("Host:%s.mqs-cn-hangzhou.aliyuncs.com", mq.accessOwnerId),
		"Date:" + date_gmt,
		"x-mqs-version:2014-07-08",
		"Content-Type:" + content_type,
		"Content-MD5:" + content_md5,
		"Authorization:" + sign})

	easy.Setopt(curl.OPT_POSTFIELDS, content)

	//easy.Setopt(curl.OPT_RTSP_TRANSPORT, false)
	easy.Setopt(curl.OPT_CUSTOMREQUEST, "POST")
	fooTest := func(buf []byte, userdata interface{}) bool {
		fmt.Println(string(buf))
		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

	if err := easy.Perform(); err != nil {
		println("MSQ add msg ERROR: ", err.Error())
	}
}

func hash_hmac(msg, key []byte) string {
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(msg))
	expectedMAC := mac.Sum(nil)
	return string(expectedMAC)
}

func md5_encode(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	x := string(h.Sum(nil))
	return fmt.Sprintf("%x", x)
}

func base64_encode(str string) string {
	return b64.StdEncoding.EncodeToString([]byte(str))
}
