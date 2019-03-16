package dom

import "github.com/dennwc/dom/js"

// https://developer.mozilla.org/en-US/docs/Web/API/HTMLElement

func (e *Element) AsHTMLElement() *HTMLElement {
	if e == nil {
		return nil
	}
	// TODO: check if the type matches
	return &HTMLElement{Element: *e}
}

type HTMLElement struct {
	Element
}

// Properties

// AccessKey is a DOMString representing the access key assigned to the element.
func (e *HTMLElement) AccessKey() string {
	return e.v.Get("accessKey").String()
}

// SetAccessKey is a DOMString representing the access key assigned to the element.
func (e *HTMLElement) SetAccessKey(v string) {
	e.v.Set("accessKey", v)
}

// AccessKeyLabel returns a DOMString containing the element's assigned access key.
func (e *HTMLElement) AccessKeyLabel() string {
	return e.v.Get("accessKeyLabel").String()
}

// ContentEditable is a DOMString, where a value of "true" means the element is editable and a value of "false" means it isn't.
func (e *HTMLElement) ContentEditable() string {
	return e.v.Get("contentEditable").String()
}

// SetContentEditable is a DOMString, where a value of "true" means the element is editable and a value of "false" means it isn't.
func (e *HTMLElement) SetContentEditable(v string) {
	e.v.Set("contentEditable", v)
}

// IsContentEditable returns a Boolean that indicates whether or not the content of the element can be edited.
func (e *HTMLElement) IsContentEditable() bool {
	return e.v.Get("isContentEditable").Bool()
}

// Dataset returns a DOMStringMap with which script can read and write the element's custom data attributes (data-*) .
func (e *HTMLElement) Dataset() js.Value {
	return e.v.Get("dataset")
}

// Dir is a DOMString, reflecting the dir global attribute, representing the directionality of the element. Possible values are "ltr", "rtl", and "auto".
func (e *HTMLElement) Dir() string {
	return e.v.Get("dir").String()
}

// SetDir is a DOMString, reflecting the dir global attribute, representing the directionality of the element. Possible values are "ltr", "rtl", and "auto".
func (e *HTMLElement) SetDir(v string) {
	e.v.Set("dir", v)
}

// Draggable is a Boolean indicating if the element can be dragged.
func (e *HTMLElement) Draggable() bool {
	return e.v.Get("draggable").Bool()
}

// SetDraggable is a Boolean indicating if the element can be dragged.
func (e *HTMLElement) SetDraggable(v bool) {
	e.v.Set("draggable", v)
}

// Dropzone returns a DOMSettableTokenList reflecting the dropzone global attribute and describing the behavior of the element regarding a drop operation.
func (e *HTMLElement) Dropzone() *TokenList {
	return AsTokenList(e.v.Get("dropzone"))
}

// Hidden is a Boolean indicating if the element is hidden or not.
func (e *HTMLElement) Hidden() bool {
	return e.v.Get("hidden").Bool()
}

// SetHidden is a Boolean indicating if the element is hidden or not.
func (e *HTMLElement) SetHidden(v bool) {
	e.v.Set("hidden", v)
}

// Inert is a Boolean indicating whether the user agent must act as though the given node is absent for the purposes of user interaction events, in-page text searches ("find in page"), and text selection.
func (e *HTMLElement) Inert() bool {
	return e.v.Get("inert").Bool()
}

// SetInert is a Boolean indicating whether the user agent must act as though the given node is absent for the purposes of user interaction events, in-page text searches ("find in page"), and text selection.
func (e *HTMLElement) SetInert(v bool) {
	e.v.Set("inert", v)
}

// InnerText represents the "rendered" text content of a node and its descendants. As a getter, it approximates the text the user would get if they highlighted the contents of the element with the cursor and then copied it to the clipboard.
func (e *HTMLElement) InnerText() string {
	return e.v.Get("innerText").String()
}

// SetInnerText represents the "rendered" text content of a node and its descendants. As a getter, it approximates the text the user would get if they highlighted the contents of the element with the cursor and then copied it to the clipboard.
func (e *HTMLElement) SetInnerText(v string) {
	e.v.Set("innerText", v)
}

// ItemScope  is a Boolean representing the item scope.
func (e *HTMLElement) ItemScope() bool {
	return e.v.Get("itemScope").Bool()
}

// SetItemScope  is a Boolean representing the item scope.
func (e *HTMLElement) SetItemScope(v bool) {
	e.v.Set("itemScope", v)
}

// ItemType returns a DOMSettableTokenList…
func (e *HTMLElement) ItemType() *TokenList {
	return AsTokenList(e.v.Get("itemType"))
}

// ItemId  is a DOMString representing the item ID.
func (e *HTMLElement) ItemId() string {
	return e.v.Get("itemId").String()
}

// SetItemId  is a DOMString representing the item ID.
func (e *HTMLElement) SetItemId(v string) {
	e.v.Set("itemId", v)
}

// ItemRef returns a DOMSettableTokenList…
func (e *HTMLElement) ItemRef() *TokenList {
	return AsTokenList(e.v.Get("itemRef"))
}

// ItemProp returns a DOMSettableTokenList…
func (e *HTMLElement) ItemProp() *TokenList {
	return AsTokenList(e.v.Get("itemProp"))
}

// ItemValue  returns a Object representing the item value.
func (e *HTMLElement) ItemValue() js.Value {
	return e.v.Get("itemValue")
}

// SetItemValue  returns a Object representing the item value.
func (e *HTMLElement) SetItemValue(v js.Value) {
	e.v.Set("itemValue", v)
}

// Lang is a DOMString representing the language of an element's attributes, text, and element contents.
func (e *HTMLElement) Lang() string {
	return e.v.Get("lang").String()
}

// SetLang is a DOMString representing the language of an element's attributes, text, and element contents.
func (e *HTMLElement) SetLang(v string) {
	e.v.Set("lang", v)
}

// NoModule is a Boolean indicating whether an import script can be executed in user agents that support module scripts.
func (e *HTMLElement) NoModule() bool {
	return e.v.Get("noModule").Bool()
}

// SetNoModule is a Boolean indicating whether an import script can be executed in user agents that support module scripts.
func (e *HTMLElement) SetNoModule(v bool) {
	e.v.Set("noModule", v)
}

// Nonce returns the cryptographic number used once that is used by Content Security Policy to determine whether a given fetch will be allowed to proceed.
func (e *HTMLElement) Nonce() js.Value {
	return e.v.Get("nonce")
}

// SetNonce returns the cryptographic number used once that is used by Content Security Policy to determine whether a given fetch will be allowed to proceed.
func (e *HTMLElement) SetNonce(v js.Value) {
	e.v.Set("nonce", v)
}

// OffsetHeight returns a double containing the height of an element, relative to the layout.
func (e *HTMLElement) OffsetHeight() float64 {
	return e.v.Get("offsetHeight").Float()
}

// OffsetLeft returns a double, the distance from this element's left border to its offsetParent's left border.
func (e *HTMLElement) OffsetLeft() float64 {
	return e.v.Get("offsetLeft").Float()
}

// OffsetParent returns a Element that is the element from which all offset calculations are currently computed.
func (e *HTMLElement) OffsetParent() *Element {
	return AsElement(e.v.Get("offsetParent"))
}

// OffsetTop returns a double, the distance from this element's top border to its offsetParent's top border.
func (e *HTMLElement) OffsetTop() float64 {
	return e.v.Get("offsetTop").Float()
}

// OffsetWidth returns a double containing the width of an element, relative to the layout.
func (e *HTMLElement) OffsetWidth() float64 {
	return e.v.Get("offsetWidth").Float()
}

// Properties returns a HTMLPropertiesCollection…
// func (e *HTMLElement) Properties() HTMLPropertiesCollection {
// 	return e.v.Get("properties")
// }

// Spellcheck is a Boolean that controls spell-checking. It is present on all HTML elements, though it doesn't have an effect on all of them.
func (e *HTMLElement) Spellcheck() bool {
	return e.v.Get("spellcheck").Bool()
}

// SetSpellcheck is a Boolean that controls spell-checking. It is present on all HTML elements, though it doesn't have an effect on all of them.
func (e *HTMLElement) SetSpellcheck(v bool) {
	e.v.Set("spellcheck", v)
}

// Style is a CSSStyleDeclaration, an object representing the declarations of an element's style attributes.
func (e *HTMLElement) Style() *Style {
	return AsStyle(e.v.Get("style"))
}

// SetStyle is a CSSStyleDeclaration, an object representing the declarations of an element's style attributes.
func (e *HTMLElement) SetStyle(v *Style) {
	e.v.Set("style", v.v)
}

// TabIndex is a long representing the position of the element in the tabbing order.
func (e *HTMLElement) TabIndex() int {
	return e.v.Get("tabIndex").Int()
}

// SetTabIndex is a long representing the position of the element in the tabbing order.
func (e *HTMLElement) SetTabIndex(v int) {
	e.v.Set("tabIndex", v)
}

// Title is a DOMString containing the text that appears in a popup box when mouse is over the element.
func (e *HTMLElement) Title() string {
	return e.v.Get("title").String()
}

// SetTitle is a DOMString containing the text that appears in a popup box when mouse is over the element.
func (e *HTMLElement) SetTitle(v string) {
	e.v.Set("title", v)
}

// Translate  is a Boolean representing the translation.
func (e *HTMLElement) Translate() bool {
	return e.v.Get("translate").Bool()
}

// SetTranslate  is a Boolean representing the translation.
func (e *HTMLElement) SetTranslate(v bool) {
	e.v.Set("translate", v)
}

// Methods

// TODO
