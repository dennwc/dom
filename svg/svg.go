package svg

import (
	"fmt"
	"strings"

	"github.com/dennwc/dom"
)

// NewElement creates a new SVG element.
func NewElement(tag string) *Element {
	return &Element{dom.Doc.CreateElementNS("http://www.w3.org/2000/svg", tag)}
}

// NewContainer creates an SVG element that provides container-like API (like "g").
func NewContainer(tag string) *Container {
	return &Container{*NewElement(tag)}
}

// New creates a new root SVG element with a given size.
func New(w, h dom.Unit) *SVG {
	e := NewContainer("svg")
	dom.Body.AppendChild(e.DOMElement())
	e.SetAttribute("width", w.String())
	e.SetAttribute("height", h.String())
	return &SVG{*e}
}

// NewFullscreen is like New, but the resulting element will try to fill the whole client area.
func NewFullscreen() *SVG {
	return New(dom.Perc(100), dom.Vh(98))
}

// Element is a common base for SVG elements.
type Element struct {
	e *dom.Element
}

// SetAttribute sets an attribute of SVG element.
func (e *Element) SetAttribute(k string, v interface{}) {
	e.e.SetAttribute(k, v)
}

// Style returns a style object for this element.
func (e *Element) Style() *dom.Style {
	return e.e.Style()
}

// Transform sets a list of transformations for SVG element.
// It will override an old value.
func (e *Element) Transform(arr ...Transform) {
	str := make([]string, 0, len(arr))
	for _, t := range arr {
		str = append(str, t.TransformString())
	}
	e.e.SetAttribute("transform", strings.Join(str, " "))
}

// Translate sets an SVG element transform to translation.
// It will override an old transform value.
func (e *Element) Translate(x, y float64) {
	e.Transform(Translate{X: x, Y: y})
}

// OnClick registers an onclick event listener.
func (e *Element) OnClick(h dom.MouseEventHandler) {
	e.e.OnClick(h)
}

// OnMouseDown registers an onmousedown event listener.
func (e *Element) OnMouseDown(h dom.MouseEventHandler) {
	e.e.OnMouseDown(h)
}

// OnMouseMove registers an onmousemove event listener.
func (e *Element) OnMouseMove(h dom.MouseEventHandler) {
	e.e.OnMouseMove(h)
}

// OnMouseUp registers an onmouseup event listener.
func (e *Element) OnMouseUp(h dom.MouseEventHandler) {
	e.e.OnMouseUp(h)
}

// NewG creates a detached SVG group element ("g").
func NewG() *G {
	return &G{*NewContainer("g")}
}

// NewCircle creates a detached SVG circle with a given radius.
func NewCircle(r int) *Circle {
	c := &Circle{*NewElement("circle")}
	c.SetR(r)
	return c
}

// NewRect creates a detached SVG rectangle with a given size.
func NewRect(w, h int) *Rect {
	r := &Rect{*NewElement("rect")}
	if w != 0 || h != 0 {
		r.SetSize(w, h)
	}
	return r
}

// NewLine creates a detached SVG line.
func NewLine() *Line {
	l := &Line{*NewElement("line")}
	l.SetStrokeWidth(1)
	l.SetAttribute("stroke", "#000")
	return l
}

// NewText creates a detached SVG text element.
func NewText(str string) *Text {
	t := &Text{*NewElement("text")}
	t.SetText(str)
	return t
}

// Container is a common base for SVG elements that can contain other elements.
type Container struct {
	Element
}

// NewG creates an SVG group element ("g") in this container.
func (c *Container) NewG() *G {
	g := NewG()
	c.e.AppendChild(g.DOMElement())
	return g
}

// NewCircle creates an SVG circle with a given radius in this container.
func (c *Container) NewCircle(r int) *Circle {
	ci := NewCircle(r)
	c.e.AppendChild(ci.DOMElement())
	return ci
}

// NewRect creates an SVG rectangle with a given size in this container.
func (c *Container) NewRect(w, h int) *Rect {
	r := NewRect(w, h)
	c.e.AppendChild(r.DOMElement())
	return r
}

// NewLine creates an SVG line in this container.
func (c *Container) NewLine() *Line {
	l := NewLine()
	c.e.AppendChild(l.DOMElement())
	return l
}

// NewText creates an SVG text element in this container.
func (c *Container) NewText(str string) *Text {
	t := NewText(str)
	c.e.AppendChild(t.DOMElement())
	return t
}

// SVG is a root SVG element.
type SVG struct {
	Container
}

// DOMElement returns a dom.Element associated with this SVG element.
func (e *Element) DOMElement() *dom.Element {
	return e.e
}

// G is an SVG group element.
type G struct {
	Container
}

// Circle is an SVG circle element.
type Circle struct {
	Element
}

func (c *Circle) SetR(r int) {
	c.SetAttribute("r", r)
}
func (c *Circle) SetPos(x, y int) {
	c.SetAttribute("cx", x)
	c.SetAttribute("cy", y)
}
func (c *Circle) Fill(cl dom.Color) {
	c.SetAttribute("fill", string(cl))
}
func (c *Circle) Stroke(cl dom.Color) {
	c.SetAttribute("stroke", string(cl))
}

// Rect is an SVG rectangle element.
type Rect struct {
	Element
}

func (c *Rect) SetPos(x, y int) {
	c.SetAttribute("x", x)
	c.SetAttribute("y", y)
}
func (c *Rect) SetSize(w, h int) {
	c.SetAttribute("width", w)
	c.SetAttribute("height", h)
}
func (c *Rect) SetRound(rx, ry int) {
	c.SetAttribute("rx", rx)
	c.SetAttribute("ry", ry)
}
func (c *Rect) Fill(cl dom.Color) {
	c.SetAttribute("fill", string(cl))
}
func (c *Rect) Stroke(cl dom.Color) {
	c.SetAttribute("stroke", string(cl))
}

// Line is an SVG line element.
type Line struct {
	Element
}

func (l *Line) SetStrokeWidth(w float64) {
	l.SetAttribute("stroke-width", w)
}
func (l *Line) SetPos1(p dom.Point) {
	l.SetAttribute("x1", p.X)
	l.SetAttribute("y1", p.Y)
}
func (l *Line) SetPos2(p dom.Point) {
	l.SetAttribute("x2", p.X)
	l.SetAttribute("y2", p.Y)
}
func (l *Line) SetPos(p1, p2 dom.Point) {
	l.SetPos1(p1)
	l.SetPos2(p2)
}

// Text is an SVG text element.
type Text struct {
	Element
}

func (t *Text) SetText(s string) {
	t.e.SetInnerHTML(s)
}
func (t *Text) SetPos(x, y int) {
	t.SetAttribute("x", x)
	t.SetAttribute("y", y)
}
func (t *Text) SetDPos(dx, dy dom.Unit) {
	if dx != nil {
		t.SetAttribute("dx", dx.String())
	}
	if dy != nil {
		t.SetAttribute("dy", dy.String())
	}
}
func (t *Text) Selectable(v bool) {
	if !v {
		t.Style().Set("user-select", "none")
	} else {
		t.Style().Set("user-select", "auto")
	}
}

// Transform is transformation that can be applied to SVG elements.
type Transform interface {
	TransformString() string
}

// Translate moves an element.
type Translate struct {
	X, Y float64
}

func (t Translate) TransformString() string {
	return fmt.Sprintf("translate(%v, %v)", t.X, t.Y)
}

// Scale scales an element.
type Scale struct {
	X, Y float64
}

func (t Scale) TransformString() string {
	return fmt.Sprintf("scale(%v, %v)", t.X, t.Y)
}

// Rotate rotates an element relative to the parent.
type Rotate struct {
	A float64
}

func (t Rotate) TransformString() string {
	return fmt.Sprintf("rotate(%v)", t.A)
}

// RotatePt rotates an element relative a point.
type RotatePt struct {
	A, X, Y float64
}

func (t RotatePt) TransformString() string {
	return fmt.Sprintf("rotate(%v, %v, %v)", t.A, t.X, t.Y)
}
