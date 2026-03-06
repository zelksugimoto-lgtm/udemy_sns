package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate 構造体のバリデーションを実行し、日本語エラーメッセージを返す
func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// 最初のエラーを日本語メッセージに変換
			for _, e := range validationErrors {
				return fmt.Errorf("%s", translateError(e))
			}
		}
		return err
	}
	return nil
}

// translateError バリデーションエラーを日本語メッセージに変換
func translateError(err validator.FieldError) string {
	field := err.Field()

	// フィールド名を日本語に変換
	fieldName := translateFieldName(field)

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%sは必須です", fieldName)
	case "email":
		return fmt.Sprintf("%sの形式が正しくありません", fieldName)
	case "min":
		if err.Type().Kind() == 24 { // string
			return fmt.Sprintf("%sは%s文字以上である必要があります", fieldName, err.Param())
		}
		return fmt.Sprintf("%sは%s以上である必要があります", fieldName, err.Param())
	case "max":
		if err.Type().Kind() == 24 { // string
			return fmt.Sprintf("%sは%s文字以内である必要があります", fieldName, err.Param())
		}
		return fmt.Sprintf("%sは%s以内である必要があります", fieldName, err.Param())
	case "alphanum":
		return fmt.Sprintf("%sは英数字のみ使用できます", fieldName)
	case "oneof":
		return fmt.Sprintf("%sの値が不正です", fieldName)
	default:
		return fmt.Sprintf("%sが不正です", fieldName)
	}
}

// translateFieldName フィールド名を日本語に変換
func translateFieldName(field string) string {
	switch strings.ToLower(field) {
	case "email":
		return "メールアドレス"
	case "password":
		return "パスワード"
	case "username":
		return "ユーザー名"
	case "displayname":
		return "表示名"
	case "content":
		return "内容"
	case "visibility":
		return "公開範囲"
	default:
		return field
	}
}
