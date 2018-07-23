package dom

type DOMHighResTimeStamp float64

type Performance interface {
	EventTarget
	Now() DOMHighResTimeStamp
	TimeOrigin() DOMHighResTimeStamp
	ToJSON() Object
}

type WindowOrWorkerGlobalScope interface {
	Performance() Performance
}
