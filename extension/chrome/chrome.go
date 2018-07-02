package chrome

import "github.com/dennwc/dom/js"

var chrome = js.Get("chrome")

type WindowID int

const CurrentWindow = WindowID(0)
