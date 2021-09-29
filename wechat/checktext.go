package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const urlOfMsgSecCheck = "https://api.weixin.qq.com/wxa/msg_sec_check"

// CheckText ...
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/sec-check/security.msgSecCheck.html
type CheckText struct {
	Version   int    `json:"version"`
	OpenID    string `json:"openid"`
	Scene     int    `json:"scene"`
	Content   string `json:"content"`
	NickName  string `json:"nickname,omitempty"`
	Title     string `json:"title,omitempty"`
	Signature string `json:"signature,omitempty"`
}

// CheckTextResponse ...
type CheckTextResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Result  struct {
		Suggest string `json:"suggest"`
		Label   string `json:"label"`
	} `json:"result"`
	Details []*checkTextDetail `json:"detail"`
}

type checkTextDetail struct {
	ErrCode  int    `json:"errcode"`
	Strategy string `json:"strategy"`
	Suggest  string `json:"suggest"`
	Label    string `json:"label"`
	Prob     int    `json:"prob"`
}

// Do ...
func (this *CheckText) Do(accessToken AccessToken) (*CheckTextResponse, error) {
	token := accessToken.GetToken()
	url := fmt.Sprintf("%s?accessToken=%s", urlOfMsgSecCheck, token)
	params, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(params)
	_resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer _resp.Body.Close()
	data, err := ioutil.ReadAll(_resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("url: " + url)
	log.Println("req: " + string(params))
	log.Println("resp:" + string(data))
	if len(data) <= 0 {
		return nil, fmt.Errorf("empty response")
	}
	resp := CheckTextResponse{ErrCode: -1}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("%w resp:%s", err, string(data))
	}
	log.Printf("resp(obj): %+v", resp)
	return &resp, nil
}
