package yidun

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/tjfoc/gmsm/sm3"
)

type Client struct {
	secretId   string
	secretKey  string
	businessId string
	rand       *rand.Rand
}

func NewClient(secretId, secretKey, businessId string) *Client {
	return &Client{
		secretId:   secretId,
		secretKey:  secretKey,
		businessId: businessId,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (this *Client) BuildParams(params url.Values) {
	params.Set("secretId", this.secretId)
	params.Set("businessId", this.businessId)
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixNano()/1000000, 10))
	params.Set("nonce", strconv.FormatInt(rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000), 10))
	params.Set("signatureMethod", "SM3")
	params.Set("signature", this.Sign(params, "SM3"))

}

func (this *Client) Sign(params url.Values, method string) string {
	var paramStr string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		paramStr += key + params[key][0]
	}
	paramStr += this.secretKey
	if method == "SM3" {
		sm3Reader := sm3.New()
		sm3Reader.Write([]byte(paramStr))
		return hex.EncodeToString(sm3Reader.Sum(nil))
	} else {
		md5Reader := md5.New()
		md5Reader.Write([]byte(paramStr))
		return hex.EncodeToString(md5Reader.Sum(nil))
	}
}
