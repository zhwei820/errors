package errors

import (
	"fmt"
	"strings"
)

type LangType string

const (
	ZhCn LangType = "zh-CN"
	EnUs LangType = "en-US"
	RuRu LangType = "ru-RU"
	ZhTW LangType = "zh-TW"
	JaJP LangType = "ja-JP"
	KoKR LangType = "ko-KR"
	ESES LangType = "es-ES"
	DEDE LangType = "de-DE"
)

func (l LangType) String() string {
	return string(l)
}

// TransInfo used to pack translate messages
type TransInfo struct {
	Lang LangType //
	Key  string   //
	Msg  string   //
}

var messageMap map[string]map[LangType]TransInfo

func init() {
	messageMap = map[string]map[LangType]TransInfo{}
}
func RegisterI18n(messages []TransInfo) {
	for _, msg := range messages {
		lang := ConvertLang(msg.Lang.String())
		if msgMap, ok := messageMap[msg.Key]; ok {
			msgMap[lang] = msg
		} else {
			messageMap[msg.Key] = map[LangType]TransInfo{lang: msg}
		}
	}
	fmt.Println(fmt.Sprintf("messageMap: %v", messageMap))
}

func Translate(langSpec LangType, key string) string {
	res, ok := messageMap[key]
	if !ok {
		fmt.Println(fmt.Sprintf("messageMap: %v", messageMap))
		fmt.Println(fmt.Sprintf("messageMap[key]: %v", messageMap[key]))
		return "msg not found: " + key
	}
	ret, ok := res[langSpec]
	if !ok {
		return "msg not found: " + key + " ," + langSpec.String()
	}
	return ret.Msg
}

func TranslateWithConvertLan(langRaw, key string) string {
	langSpec := ConvertLang(langRaw)
	return Translate(langSpec, key)
}

func ConvertLang(lang string) (langSpec LangType) {
	lang = strings.ReplaceAll(lang, "_", "-")
	switch lang {
	case "zh", "ZH", "cn", "CN", "zh_CN", ZhCn.String():
		langSpec = ZhCn
	case "en", "EN", "us", "US", "en_US", EnUs.String():
		langSpec = EnUs
	case "ru", "RU", "ru_RU", RuRu.String():
		langSpec = RuRu
	case JaJP.String(), KoKR.String(), ESES.String(), DEDE.String(), ZhTW.String():
		langSpec = LangType(lang)
	default:
		langSpec = EnUs
	}
	return langSpec
}
