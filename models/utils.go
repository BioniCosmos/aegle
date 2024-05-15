package models

type Query struct {
	Skip  int64 `query:"skip"`
	Limit int64 `query:"limit"`
}
