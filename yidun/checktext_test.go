package yidun

import (
	"encoding/json"
	"testing"

	"github.com/gamedev-embers/antispam/yidun/models"
	"github.com/stretchr/testify/assert"
)

func TestCheckTextResponse(t *testing.T) {
	{
		assert := assert.New(t)
		data := `{"code":200,"msg":"ok","result":{"antispam":{"taskId":"af38362590624c37a3baba72455b536c","action":2,"censorType":0,"isRelatedHit":false,"lang":[],"labels":[{"label":600,"level":2,"details":{"hitInfos":[],"hints":[{"hint":"草泥马","positions":[{"positionType":0,"startPos":0,"endPos":3}]}]},"subLabels":[{"subLabel":"600018"}]}],"censorLabels":[]},"emotionAnalysis":{},"anticheat":{},"userRisk":{}}}`
		obj := CheckTextV4Response{}
		err := json.Unmarshal([]byte(data), &obj)
		assert.NoError(err)
		assert.Equal(200, obj.Code)
		assert.Equal("ok", obj.Msg)
		assert.NotNil(obj.Result.AntiSpam)
		as := obj.Result.AntiSpam
		assert.Equal("af38362590624c37a3baba72455b536c", as.TaskID)
		assert.Equal(2, as.Action)
		assert.Equal(1, len(as.Labels))
		assert.Equal(1, len(as.Labels[0].Details.Hints))
		hint := &models.Hint{Hint: "草泥马", Positions: []*models.HintPosition{{PositionType: 0, StartPos: 0, EndPos: 3}}}
		assert.Equal(hint, as.Labels[0].Details.Hints[0])
		assert.False(obj.IsOK())
	}

	{
		assert := assert.New(t)
		data := `{"code":200,"msg":"ok","result":{"antispam":{"taskId":"af38362590624c37a3baba72455b536c","action":2,"censorType":0,"isRelatedHit":false,"lang":[],"labels":[{"label":600,"level":2,"details":{"hitInfos":[],"hints":[{"hint":"草泥马","positions":[{"positionType":999,"startPos":998,"endPos":3}]}]},"subLabels":[{"subLabel":"600018"}]}],"censorLabels":[]},"emotionAnalysis":{},"anticheat":{},"userRisk":{}}}`
		obj := CheckTextV4Response{}
		err := json.Unmarshal([]byte(data), &obj)
		assert.NoError(err)
		assert.Equal(200, obj.Code)
		assert.Equal("ok", obj.Msg)
		assert.NotNil(obj.Result.AntiSpam)
		as := obj.Result.AntiSpam
		assert.Equal("af38362590624c37a3baba72455b536c", as.TaskID)
		assert.Equal(2, as.Action)
		assert.Equal(1, len(as.Labels))
		assert.Equal(1, len(as.Labels[0].Details.Hints))
		hint := &models.Hint{Hint: "草泥马", Positions: []*models.HintPosition{{PositionType: 999, StartPos: 998, EndPos: 3}}}
		assert.Equal(hint, as.Labels[0].Details.Hints[0])
		assert.False(obj.IsOK())
	}

	{
		assert := assert.New(t)
		data := `{"code":200,"msg":"ok","result":{"antispam":{"taskId":"7ac996e5af464862a63625441f29050f","action":0,"censorType":0,"isRelatedHit":false,"lang":[],"labels":[],"censorLabels":[]},"emotionAnalysis":{},"anticheat":{},"userRisk":{}}}`
		obj := CheckTextV4Response{}
		err := json.Unmarshal([]byte(data), &obj)
		assert.NoError(err)
		assert.Equal(200, obj.Code)
		assert.Equal("ok", obj.Msg)
		assert.NotNil(obj.Result.AntiSpam)
		as := obj.Result.AntiSpam
		assert.Equal("7ac996e5af464862a63625441f29050f", as.TaskID)
		assert.Equal(0, as.Action)
		assert.Equal(0, len(as.Labels))
		assert.True(obj.IsOK())
	}
}

func TestCheckTextResponseFiltered(t *testing.T) {
	assert := assert.New(t)
	data := `{"code":200,"msg":"ok","result":{"antispam":{"taskId":"af38362590624c37a3baba72455b536c","action":2,"censorType":0,"isRelatedHit":false,"lang":[],"labels":[{"label":600,"level":2,"details":{"hitInfos":[],"hints":[{"hint":"草泥马","positions":[{"positionType":0,"startPos":0,"endPos":3}]}]},"subLabels":[{"subLabel":"600018"}]}],"censorLabels":[]},"emotionAnalysis":{},"anticheat":{},"userRisk":{}}}`
	obj := CheckTextV4Response{}
	err := json.Unmarshal([]byte(data), &obj)
	assert.NoError(err)

	text, err := obj.GetContentFiltered("草泥马")
	assert.NoError(err)
	assert.Equal("***", text)

	text, err = obj.GetContentFiltered("草泥马勒戈壁")
	assert.NoError(err)
	assert.Equal("***勒戈壁", text)
}
