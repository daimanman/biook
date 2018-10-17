package requester

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"text/template"
)

type ParamInfo struct {
	Anum int
	Bnum int
}
type ZmmService struct {
}

func GetInt(key string) int {
	if key == "" {
		return 0
	}
	num, err := strconv.Atoi(strings.TrimSpace(key))
	if err != nil {
		log.Println(err)
		return 0
	}
	return num
}

func (zmmService *ZmmService) CreateSh(aStart, aEnd, bStart, bEnd, tpl string) *bytes.Buffer {
	sb := bytes.NewBufferString("")
	tmpl, _ := template.New("DxmTest").Parse(tpl + "\n")
	aNumStart := GetInt(aStart)
	aNumEnd := GetInt(aEnd)
	bNumStart := GetInt(bStart)
	bNumEnd := GetInt(bEnd)
	paramInfo := &ParamInfo{
		Anum: 0,
		Bnum: 0,
	}
	for i := aNumStart; i <= aNumEnd; i++ {
		for j := bNumStart; j <= bNumEnd; j++ {
			paramInfo.Anum = i
			paramInfo.Bnum = j
			tmpl.Execute(sb, paramInfo)
		}
	}
	return sb
}
