package sl

import (
	"log/slog"

	
)

func Err(err error) slog.Attr { // принимает на вход ошибку и возвращает атрибут пакета слог
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()), // текстовое содержимое ошибки
	}
}
