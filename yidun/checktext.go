package yidun

// https://support.dun.163.com/documents/2018041901?docId=424375611814748160
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gamedev-embers/wordfilter/yidun/models"
	"github.com/google/uuid"
)

const (
	urlCheckText = "http://as.dun.163.com/v4/text/check"
)

type Options struct {
	SecretID   string
	SecretKey  string
	BusinessID string
}

func (o *Options) Check() error {
	if o.SecretID == "" {
		return fmt.Errorf("empty secretId")
	}
	if o.SecretKey == "" {
		return fmt.Errorf("empty secretKey")
	}
	if o.BusinessID == "" {
		return fmt.Errorf("empty businessId")
	}
	return nil
}

// https://support.dun.163.com/documents/2018041901?docId=424375611814748160
type CheckTextV4 struct {
	Content string

	Version     string
	DataID      string
	Title       string
	DataType    string
	PublishTime int64
	Callback    string
	CallbackUrl string
	CheckLabels string
	Category    string
	SubProdduct string
}

func (this *CheckTextV4) newParams() url.Values {
	params := url.Values{}
	addif := func(k, v string) {
		if v != "" {
			params.Set(k, v)
		}
	}
	if this.DataID == "" {
		this.DataID = uuid.NewString()
	}
	if this.Version == "" {
		this.Version = "v4.2"
	}
	if this.Content == "" {
		panic(fmt.Errorf("empty content"))
	}
	params.Set("dataId", this.DataID)
	params.Set("version", this.Version)
	params.Set("content", this.Content)
	addif("title", this.Title)
	addif("dataType", this.DataType)
	addif("publishTime", strconv.Itoa(int(this.PublishTime)))
	addif("callback", this.Callback)
	addif("callbackUrl", this.CallbackUrl)
	addif("category", this.Category)
	addif("subProduct", this.SubProdduct)
	return params
}

func (this *CheckTextV4) Do(c *Client) (*CheckTextV4Response, error) {
	params := this.newParams()
	c.BuildParams(params)
	url := urlCheckText
	_resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	defer _resp.Body.Close()

	data, err := ioutil.ReadAll(_resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println("url: " + url)
	log.Println("req: " + params.Encode())
	log.Println("resp:" + string(data))
	if len(data) <= 0 {
		return nil, fmt.Errorf("empty response")
	}
	resp := CheckTextV4Response{}
	resp.Result.AntiSpam = &models.AntiSpam{Action: -1}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("%w resp:%s", err, string(data))
	}
	return &resp, nil
}

type CheckTextV4Response struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		AntiSpam *models.AntiSpam `json:"antispam"`
	} `json:"result"`
}

func (ctr *CheckTextV4Response) IsOK() bool {
	if ctr == nil {
		return false
	}
	if ctr.Code != 200 {
		return false
	}
	if ctr.Result.AntiSpam.Action != 0 {
		return false
	}
	return true
}

func (ctr *CheckTextV4Response) GetContentFiltered(content string) (string, error) {
	if ctr == nil {
		return "", fmt.Errorf("nil response")
	}

	if len(ctr.Result.AntiSpam.Labels) <= 0 {
		return "", fmt.Errorf("no labels")
	}
	if len(ctr.Result.AntiSpam.Labels[0].Details.Hints) <= 0 {
		return "", fmt.Errorf("no labels.details.hints")
	}

	chars := []rune(content)
	fill := []rune("*")[0]
	for _, label := range ctr.Result.AntiSpam.Labels {
		for _, hint := range label.Details.Hints {
			for _, pst := range hint.Positions {
				if pst == nil {
					return "", fmt.Errorf("nil position")
				}
				if pst.StartPos >= len(content) || pst.EndPos > len(chars) {
					return "", fmt.Errorf("invalid position %+v", pst)
				}
				for i := pst.StartPos; i < pst.EndPos; i++ {
					chars[i] = fill
				}
				log.Printf("%s  pst=%d,%d", string(chars), pst.StartPos, pst.EndPos)
			}
		}
	}
	return string(chars), nil
}
