package dom

func (d *Document) NewInput(typ string) *Input {
	e := d.CreateElement("input")
	inp := &Input{*e}
	inp.SetType(typ)
	return inp
}

type Input struct {
	Element
}

func (inp *Input) Value() string {
	return inp.v.Get("value").String()
}
func (inp *Input) SetType(typ string) {
	inp.SetAttribute("type", typ)
}
func (inp *Input) SetName(name string) {
	inp.SetAttribute("name", name)
}
func (inp *Input) SetValue(val interface{}) {
	inp.v.Set("value", val)
}
func (inp *Input) OnChange(h EventHandler) {
	inp.AddEventListener("change", h)
}
func (inp *Input) OnInput(h EventHandler) {
	inp.AddEventListener("input", h)
}

func (d *Document) NewButton(s string) *Button {
	e := d.CreateElement("button")
	b := &Button{*e}
	b.SetInnerHTML(s)
	return b
}

type Button struct {
	Element
}

func (b *Button) OnClick(h EventHandler) {
	b.AddEventListener("click", h)
}
