package errors

import (
	"fmt"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	// ZhCn zh-CN, language type Chinese
	ZhCn = "zh-CN"
	// EnUs en-US, language type English
	EnUs = "en-US"
	// RuRu ru-RU  language type Russian
	RuRu = "ru-RU"
	// ZhTW zh-TW language type traditional chinese
	ZhTW = "zh-TW"
	// JaJP language type JaJP
	JaJP = "ja-JP"
	// KoKR language type Korean
	KoKR = "ko-KR"
	// ESES  language type Spain
	ESES = "es-ES"
	// DEDE language type Germany
	DEDE = "de-DE"
)

var _PrinterStore map[string]*message.Printer

func getPrinter(lanSpec string) *message.Printer {
	printer, ok := _PrinterStore[lanSpec]
	if ok {
		return printer
	}
	tag, _, _ := matcher.Match(language.MustParse(lanSpec))
	return message.NewPrinter(tag)
}

func init() {
	matcher = language.NewMatcher(message.DefaultCatalog.Languages())
	_PrinterStore = make(map[string]*message.Printer)
}

var matcher language.Matcher

// TransInfo used to pack translate messages
type TransInfo struct {
	Tag string      // 语言
	Key string      // 翻译key 可以是biz code
	Msg interface{} // 翻译后的内容
}

// RegisterTranslate 各业务注册 自定义码以及翻译内容
func RegisterTranslate(messages []TransInfo) {
	for _, msg := range messages {
		tag := language.MustParse(ConvertLanHeader(msg.Tag))
		switch m := msg.Msg.(type) {
		case string:
			if err := message.SetString(tag, msg.Key, m); err != nil {
				fmt.Printf("error occurred, set translate failed, Tag is %v, Key is %s \n", tag, msg.Key)
			}
		}
	}
	matcher = language.NewMatcher(message.DefaultCatalog.Languages())
	refreshPrinterStore()
}

// RegisterOne 各业务注册 自定义码以及翻译内容
func RegisterOne(msg TransInfo) {
	tag := language.MustParse(ConvertLanHeader(msg.Tag))
	switch m := msg.Msg.(type) {
	case string:
		if err := message.SetString(tag, msg.Key, m); err != nil {
			fmt.Printf("error occurred, set translate failed, Tag is %v, Key is %s \n", tag, msg.Key)
		}
	}
	matcher = language.NewMatcher(message.DefaultCatalog.Languages())
	refreshPrinterStore()
}

// translateWithLanHeader 根据请求头部的语言信息翻译文本，需预先在messages中设定翻译内容
func translateWithLanHeader(lanHeader, key string, args ...interface{}) string {
	return getPrinter(ConvertLanHeader(lanHeader)).Sprintf(key, args...)
}

// Translate 按语言翻译，需预先在messages中设定翻译内容
func Translate(lanSpec, key string, args ...interface{}) string {
	return getPrinter(lanSpec).Sprintf(key, args...)
}

// TranslateWithConvertLan 按语言翻译，需预先在messages中设定翻译内容
func TranslateWithConvertLan(lanSpec, key string, args ...interface{}) string {
	return getPrinter(ConvertLanHeader(lanSpec)).Sprintf(key, args...)
}

// ConvertLanHeader 根据请求头中的LANGUAGE-TYPE获取一个标准的语言标识
func ConvertLanHeader(lan string) (lanSpec string) {
	lan = strings.Replace(lan, "_", "-", -1)
	switch lan {
	case "0", "zh-CN", "zh", "ZH", "cn", "CN", "zh_CN":
		lanSpec = ZhCn
	case "1", "en-US", "en", "EN", "us", "US", "en_US":
		lanSpec = EnUs
	case "2", "ru-RU", "ru", "RU", "ru_RU":
		lanSpec = RuRu
	case "5", ZhTW:
		lanSpec = ZhTW
	case JaJP, KoKR, ESES, DEDE:
		lanSpec = lan
	default:
		lanSpec = EnUs
	}
	return lanSpec
}

func refreshPrinterStore() {
	for _, lan := range []string{ZhCn, EnUs, RuRu, ZhTW, JaJP, KoKR, ESES, DEDE} {
		tag, _, _ := matcher.Match(language.MustParse(lan))
		_PrinterStore[lan] = message.NewPrinter(tag)
	}
}
