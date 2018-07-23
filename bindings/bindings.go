package dom

//go:generate go run ../cmd/webidl-gen/main.go DOM.widl High_Resolution_Time.widl

type EventHandler func()

type Object interface{}
