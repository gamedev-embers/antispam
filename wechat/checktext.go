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
		Label   int    `json:"label"`
	} `json:"result"`
	Details []*checkTextDetail `json:"detail"`
}

func (res *CheckTextResponse) IsPass() bool {
	if res.ErrCode != 0 {
		return false
	}
	return res.Result.Suggest == "pass"
}

type checkTextDetail struct {
	ErrCode  int    `json:"errcode"`
	Strategy string `json:"strategy"`
	Suggest  string `json:"suggest"`
	Label    int    `json:"label"`
	Prob     int    `json:"prob"`
}

// Do ...
func (this *CheckText) Do(accessToken AccessToken) (*CheckTextResponse, error) {

	do := func() (*CheckTextResponse, error) {
		token := accessToken.GetToken()
		url := fmt.Sprintf("%s?access_token=%s", urlOfMsgSecCheck, token)
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
		if len(data) <= 0 {
			return nil, fmt.Errorf("empty response")
		}
		resp := CheckTextResponse{ErrCode: -1}
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, fmt.Errorf("%w resp:%s", err, string(data))
		}
		log.Println("\turl: " + url)
		log.Println("\treq: " + string(params))
		log.Println("\tresp:" + string(data))
		log.Printf("\tresp(obj): %+v", &resp)
		return &resp, nil
	}
	resp, err := do()
	if err != nil {
		return resp, err
	}
	switch resp.ErrCode {
	case 0:
		return resp, nil
	case 87014:
		return resp, fmt.Errorf("errcode: %d", resp.ErrCode)
	case 40001:
		// retry
		accessToken.RefreshIf(true)
		resp, err = do()
	default:
		return resp, fmt.Errorf("errcode: %d resp:%v", resp.ErrCode, resp)
	}
	return resp, nil
}
