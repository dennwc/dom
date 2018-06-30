package dom

var (
	Doc  = GetDocument()
	Body = Doc.GetElementsByTagName("body")[0]
)

func init() {
	Body.Style().SetMarginsRaw("0")
}
