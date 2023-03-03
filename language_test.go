package errors

import (
	"fmt"
	"testing"
)

func TestRegisterI18n(t *testing.T) {
	messages := []TransInfo{
		{
			Lang: EnUs,
			Key:  "100000101",
			Msg:  "demo error message",
		},
		{
			Lang: ZhCn,
			Key:  "100000101",
			Msg:  "测试消息",
		},
		{
			Lang: ZhTW,
			Key:  "100000101",
			Msg:  "测试消息",
		},
		{
			Lang: RuRu,
			Key:  "100000101",
			Msg:  "demo error message",
		},
	}
	RegisterI18n(messages)

	fmt.Println(TranslateWithConvertLan("en", "100000101"))
	fmt.Println(TranslateWithConvertLan("zh", "100000101"))
	fmt.Println(TranslateWithConvertLan("ru", "100000101"))
	fmt.Println(TranslateWithConvertLan("zh-TW", "100000101"))
}
