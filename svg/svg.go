package svg

import (
	"fmt"
	"github.com/dennwc/dom"
	"strings"
)

func NewElement(tag string) *Element {
	return &Element{dom.Doc.CreateElementNS("http://www.w3.org/2000/svg", tag)}
}

func NewContainer(tag string) *Container {
	return &Container{*NewElement(tag)}
}

func New(w, h dom.Unit) *SVG {
	e := NewContainer("svg")
	dom.Body.AppendChild(e.DOMElement())
	e.SetAttribute("width", w.String())
	e.SetAttribute("height", h.String())
	return &SVG{*e}
}

type Element struct {
	e *dom.Element
}

func (e *Element) SetAttribute(k string, v interface{}) {
	e.e.SetAttribute(k, v)
}
func (e *Element) Style() *dom.Style {
	return e.e.Style()
}
func (e *Element) Transform(arr ...Transform) {
	str := make([]string, 0, len(arr))
	for _, t := range arr {
		str = append(str, t.TransformString())
	}
	e.e.SetAttribute("transform", strings.Join(str, " "))
}
func (e *Element) Translate(x, y float64) {
	e.Transform(Translate{X: x, Y: y})
}

type Container struct {
	Element
}

func (c *Container) NewCircle(r int) *Circle {
	ci := &Circle{*NewElement("circle")}
	ci.SetR(r)
	c.e.AppendChild(ci.DOMElement())
	return ci
}
func (c *Container) NewLine() *Line {
	l := &Line{*NewElement("line")}
	l.SetStrokeWidth(1)
	l.SetAttribute("stroke", "#000")
	c.e.AppendChild(l.DOMElement())
	return l
}

func (c *Container) NewG() *G {
	g := &G{*NewContainer("g")}
	c.e.AppendChild(g.DOMElement())
	return g
}

func (c *Container) NewText(str string) *Text {
	t := &Text{*NewElement("text")}
	t.SetText(str)
	c.e.AppendChild(t.DOMElement())
	return t
}

type SVG struct {
	Container
}

func (e *Element) DOMElement() *dom.Element {
	return e.e
}

type G struct {
	Container
}

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

type Line struct {
	Element
}

func (l *Line) SetStrokeWidth(w float64) {
	l.SetAttribute("stroke-width", w)
}
func (l *Line) SetPos1(x, y int) {
	l.SetAttribute("x1", x)
	l.SetAttribute("y1", y)
}
func (l *Line) SetPos2(x, y int) {
	l.SetAttribute("x2", x)
	l.SetAttribute("y2", y)
}
func (l *Line) SetPos(x1, y1, x2, y2 int) {
	l.SetPos1(x1, y1)
	l.SetPos2(x2, y2)
}

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

type Transform interface {
	TransformString() string
}

type Translate struct {
	X, Y float64
}

func (t Translate) TransformString() string {
	return fmt.Sprintf("translate(%v, %v)", t.X, t.Y)
}

type Scale struct {
	X, Y float64
}

func (t Scale) TransformString() string {
	return fmt.Sprintf("scale(%v, %v)", t.X, t.Y)
}

type Rotate struct {
	A float64
}

func (t Rotate) TransformString() string {
	return fmt.Sprintf("rotate(%v)", t.A)
}

type RotatePt struct {
	A, X, Y float64
}

func (t RotatePt) TransformString() string {
	return fmt.Sprintf("rotate(%v, %v, %v)", t.A, t.X, t.Y)
}
