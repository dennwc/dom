package dom

type Event interface {
	Type() string
	Target() EventTarget
	SrcElement() EventTarget
	CurrentTarget() EventTarget
	ComposedPath() []EventTarget
	EventPhase() uint16
	StopPropagation()
	CancelBubble() bool
	SetCancelBubble(v bool)
	StopImmediatePropagation()
	Bubbles() bool
	Cancelable() bool
	ReturnValue() bool
	SetReturnValue(v bool)
	PreventDefault()
	DefaultPrevented() bool
	Composed() bool
	IsTrusted() bool
	TimeStamp() DOMHighResTimeStamp
	InitEvent(typ string, bubbles bool, cancelable bool)
}

const (
	Event_NONE            = 0
	Event_CAPTURING_PHASE = 1
	Event_AT_TARGET       = 2
	Event_BUBBLING_PHASE  = 3
)

type EventInit struct {
	Bubbles    bool
	Cancelable bool
	Composed   bool
}

type Window interface {
	Event() interface{}
}

type CustomEvent interface {
	Event
	Detail() interface{}
	InitCustomEvent(typ string, bubbles bool, cancelable bool, detail interface{})
}

type CustomEventInit struct {
	EventInit
	Detail interface{}
}

type EventTarget interface {
	AddEventListener(typ string, callback EventListener, options interface{ isAddEventListenerOptions() })
	RemoveEventListener(typ string, callback EventListener, options interface{ isRemoveEventListenerOptions() })
	DispatchEvent(event Event) bool
}

type AddEventListenerOptions1 struct {
	Value AddEventListenerOptions
}

func (AddEventListenerOptions1) isAddEventListenerOptions() {}

type AddEventListenerOptions2 struct {
	Value bool
}

func (AddEventListenerOptions2) isAddEventListenerOptions() {}

type RemoveEventListenerOptions1 struct {
	Value EventListenerOptions
}

func (RemoveEventListenerOptions1) isRemoveEventListenerOptions() {}

type RemoveEventListenerOptions2 struct {
	Value bool
}

func (RemoveEventListenerOptions2) isRemoveEventListenerOptions() {}

type EventListener interface {
	HandleEvent(event Event)
}

type EventListenerOptions struct {
	Capture bool
}

type AddEventListenerOptions struct {
	EventListenerOptions
	Passive bool
	Once    bool
}

type AbortController interface {
	Signal() AbortSignal
	Abort()
}

type AbortSignal interface {
	EventTarget
	Aborted() bool
	Onabort() EventHandler
	SetOnabort(v EventHandler)
}

type NodeList interface {
	Item(index uint32) Node
	Length() uint32
}

type HTMLCollection interface {
	Length() uint32
	Item(index uint32) Element
	NamedItem(name string) Element
}

type MutationObserver interface {
	Observe(target Node, options MutationObserverInit)
	Disconnect()
	TakeRecords() []MutationRecord
}

type MutationObserverInit struct {
	ChildList             bool
	Attributes            bool
	CharacterData         bool
	Subtree               bool
	AttributeOldValue     bool
	CharacterDataOldValue bool
	AttributeFilter       []string
}

type MutationRecord interface {
	Type() string
	Target() Node
	AddedNodes() NodeList
	RemovedNodes() NodeList
	PreviousSibling() Node
	NextSibling() Node
	AttributeName() string
	AttributeNamespace() string
	OldValue() string
}

type Node interface {
	EventTarget
	NodeType() uint16
	NodeName() string
	BaseURI() string
	IsConnected() bool
	OwnerDocument() Document
	GetRootNode(options GetRootNodeOptions) Node
	ParentNode() Node
	ParentElement() Element
	HasChildNodes() bool
	ChildNodes() NodeList
	FirstChild() Node
	LastChild() Node
	PreviousSibling() Node
	NextSibling() Node
	NodeValue() string
	SetNodeValue(v string)
	TextContent() string
	SetTextContent(v string)
	Normalize()
	CloneNode(deep bool) Node
	IsEqualNode(otherNode Node) bool
	IsSameNode(otherNode Node) bool
	CompareDocumentPosition(other Node) uint16
	Contains(other Node) bool
	LookupPrefix(namespace string) string
	LookupNamespaceURI(prefix string) string
	IsDefaultNamespace(namespace string) bool
	InsertBefore(node Node, child Node) Node
	AppendChild(node Node) Node
	ReplaceChild(node Node, child Node) Node
	RemoveChild(child Node) Node
}

const (
	Node_ELEMENT_NODE                              = 1
	Node_ATTRIBUTE_NODE                            = 2
	Node_TEXT_NODE                                 = 3
	Node_CDATA_SECTION_NODE                        = 4
	Node_ENTITY_REFERENCE_NODE                     = 5
	Node_ENTITY_NODE                               = 6
	Node_PROCESSING_INSTRUCTION_NODE               = 7
	Node_COMMENT_NODE                              = 8
	Node_DOCUMENT_NODE                             = 9
	Node_DOCUMENT_TYPE_NODE                        = 10
	Node_DOCUMENT_FRAGMENT_NODE                    = 11
	Node_NOTATION_NODE                             = 12
	Node_DOCUMENT_POSITION_DISCONNECTED            = 0x01
	Node_DOCUMENT_POSITION_PRECEDING               = 0x02
	Node_DOCUMENT_POSITION_FOLLOWING               = 0x04
	Node_DOCUMENT_POSITION_CONTAINS                = 0x08
	Node_DOCUMENT_POSITION_CONTAINED_BY            = 0x10
	Node_DOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC = 0x20
)

type GetRootNodeOptions struct {
	Composed bool
}

type Document interface {
	Node
	Implementation() DOMImplementation
	URL() string
	DocumentURI() string
	Origin() string
	CompatMode() string
	CharacterSet() string
	Charset() string
	InputEncoding() string
	ContentType() string
	Doctype() DocumentType
	DocumentElement() Element
	GetElementsByTagName(qualifiedName string) HTMLCollection
	GetElementsByTagNameNS(namespace string, localName string) HTMLCollection
	GetElementsByClassName(classNames string) HTMLCollection
	CreateElement(localName string, options interface{ isCreateElementOptions() }) Element
	CreateElementNS(namespace string, qualifiedName string, options interface{ isCreateElementNSOptions() }) Element
	CreateDocumentFragment() DocumentFragment
	CreateTextNode(data string) Text
	CreateCDATASection(data string) CDATASection
	CreateComment(data string) Comment
	CreateProcessingInstruction(target string, data string) ProcessingInstruction
	ImportNode(node Node, deep bool) Node
	AdoptNode(node Node) Node
	CreateAttribute(localName string) Attr
	CreateAttributeNS(namespace string, qualifiedName string) Attr
	CreateEvent(iface string) Event
	CreateRange() Range
	CreateNodeIterator(root Node, whatToShow uint32, filter NodeFilter) NodeIterator
	CreateTreeWalker(root Node, whatToShow uint32, filter NodeFilter) TreeWalker
}

type CreateElementOptions1 struct {
	Value string
}

func (CreateElementOptions1) isCreateElementOptions() {}

type CreateElementOptions2 struct {
	Value ElementCreationOptions
}

func (CreateElementOptions2) isCreateElementOptions() {}

type CreateElementNSOptions1 struct {
	Value string
}

func (CreateElementNSOptions1) isCreateElementNSOptions() {}

type CreateElementNSOptions2 struct {
	Value ElementCreationOptions
}

func (CreateElementNSOptions2) isCreateElementNSOptions() {}

type XMLDocument interface {
	Document
}

type ElementCreationOptions struct {
	Is string
}

type DOMImplementation interface {
	CreateDocumentType(qualifiedName string, publicId string, systemId string) DocumentType
	CreateDocument(namespace string, qualifiedName string, doctype DocumentType) XMLDocument
	CreateHTMLDocument(title string) Document
	HasFeature() bool
}

type DocumentType interface {
	Node
	Name() string
	PublicId() string
	SystemId() string
}

type DocumentFragment interface {
	Node
}

type ShadowRoot interface {
	DocumentFragment
	Mode() ShadowRootMode
	Host() Element
}

type ShadowRootMode string

const (
	ShadowRootMode_1 = ShadowRootMode("open")
	ShadowRootMode_2 = ShadowRootMode("closed")
)

type Element interface {
	Node
	NamespaceURI() string
	Prefix() string
	LocalName() string
	TagName() string
	Id() string
	SetId(v string)
	ClassName() string
	SetClassName(v string)
	ClassList() DOMTokenList
	Slot() string
	SetSlot(v string)
	HasAttributes() bool
	Attributes() NamedNodeMap
	GetAttributeNames() []string
	GetAttribute(qualifiedName string) string
	GetAttributeNS(namespace string, localName string) string
	SetAttribute(qualifiedName string, value string)
	SetAttributeNS(namespace string, qualifiedName string, value string)
	RemoveAttribute(qualifiedName string)
	RemoveAttributeNS(namespace string, localName string)
	ToggleAttribute(qualifiedName string, force bool) bool
	HasAttribute(qualifiedName string) bool
	HasAttributeNS(namespace string, localName string) bool
	GetAttributeNode(qualifiedName string) Attr
	GetAttributeNodeNS(namespace string, localName string) Attr
	SetAttributeNode(attr Attr) Attr
	SetAttributeNodeNS(attr Attr) Attr
	RemoveAttributeNode(attr Attr) Attr
	AttachShadow(init ShadowRootInit) ShadowRoot
	ShadowRoot() ShadowRoot
	Closest(selectors string) Element
	Matches(selectors string) bool
	WebkitMatchesSelector(selectors string) bool
	GetElementsByTagName(qualifiedName string) HTMLCollection
	GetElementsByTagNameNS(namespace string, localName string) HTMLCollection
	GetElementsByClassName(classNames string) HTMLCollection
	InsertAdjacentElement(where string, element Element) Element
	InsertAdjacentText(where string, data string)
}

type ShadowRootInit struct {
	Mode ShadowRootMode
}

type NamedNodeMap interface {
	Length() uint32
	Item(index uint32) Attr
	GetNamedItem(qualifiedName string) Attr
	GetNamedItemNS(namespace string, localName string) Attr
	SetNamedItem(attr Attr) Attr
	SetNamedItemNS(attr Attr) Attr
	RemoveNamedItem(qualifiedName string) Attr
	RemoveNamedItemNS(namespace string, localName string) Attr
}

type Attr interface {
	Node
	NamespaceURI() string
	Prefix() string
	LocalName() string
	Name() string
	Value() string
	SetValue(v string)
	OwnerElement() Element
	Specified() bool
}

type CharacterData interface {
	Node
	Data() string
	SetData(v string)
	Length() uint32
	SubstringData(offset uint32, count uint32) string
	AppendData(data string)
	InsertData(offset uint32, data string)
	DeleteData(offset uint32, count uint32)
	ReplaceData(offset uint32, count uint32, data string)
}

type Text interface {
	CharacterData
	SplitText(offset uint32) Text
	WholeText() string
}

type CDATASection interface {
	Text
}

type ProcessingInstruction interface {
	CharacterData
	Target() string
}

type Comment interface {
	CharacterData
}

type AbstractRange interface {
	StartContainer() Node
	StartOffset() uint32
	EndContainer() Node
	EndOffset() uint32
	Collapsed() bool
}

type StaticRange interface {
	AbstractRange
}

type Range interface {
	AbstractRange
	CommonAncestorContainer() Node
	SetStart(node Node, offset uint32)
	SetEnd(node Node, offset uint32)
	SetStartBefore(node Node)
	SetStartAfter(node Node)
	SetEndBefore(node Node)
	SetEndAfter(node Node)
	Collapse(toStart bool)
	SelectNode(node Node)
	SelectNodeContents(node Node)
	CompareBoundaryPoints(how uint16, sourceRange Range) int16
	DeleteContents()
	ExtractContents() DocumentFragment
	CloneContents() DocumentFragment
	InsertNode(node Node)
	SurroundContents(newParent Node)
	CloneRange() Range
	Detach()
	IsPointInRange(node Node, offset uint32) bool
	ComparePoint(node Node, offset uint32) int16
	IntersectsNode(node Node) bool
}

const (
	Range_START_TO_START = 0
	Range_START_TO_END   = 1
	Range_END_TO_END     = 2
	Range_END_TO_START   = 3
)

type NodeIterator interface {
	Root() Node
	ReferenceNode() Node
	PointerBeforeReferenceNode() bool
	WhatToShow() uint32
	Filter() NodeFilter
	NextNode() Node
	PreviousNode() Node
	Detach()
}

type TreeWalker interface {
	Root() Node
	WhatToShow() uint32
	Filter() NodeFilter
	CurrentNode() Node
	SetCurrentNode(v Node)
	ParentNode() Node
	FirstChild() Node
	LastChild() Node
	PreviousSibling() Node
	NextSibling() Node
	PreviousNode() Node
	NextNode() Node
}

type NodeFilter interface {
	AcceptNode(node Node) uint16
}

const (
	NodeFilter_FILTER_ACCEPT               = 1
	NodeFilter_FILTER_REJECT               = 2
	NodeFilter_FILTER_SKIP                 = 3
	NodeFilter_SHOW_ALL                    = 0xFFFFFFFF
	NodeFilter_SHOW_ELEMENT                = 0x1
	NodeFilter_SHOW_ATTRIBUTE              = 0x2
	NodeFilter_SHOW_TEXT                   = 0x4
	NodeFilter_SHOW_CDATA_SECTION          = 0x8
	NodeFilter_SHOW_ENTITY_REFERENCE       = 0x10
	NodeFilter_SHOW_ENTITY                 = 0x20
	NodeFilter_SHOW_PROCESSING_INSTRUCTION = 0x40
	NodeFilter_SHOW_COMMENT                = 0x80
	NodeFilter_SHOW_DOCUMENT               = 0x100
	NodeFilter_SHOW_DOCUMENT_TYPE          = 0x200
	NodeFilter_SHOW_DOCUMENT_FRAGMENT      = 0x400
	NodeFilter_SHOW_NOTATION               = 0x800
)

type DOMTokenList interface {
	Length() uint32
	Item(index uint32) string
	Contains(token string) bool
	Add(tokens ...string)
	Remove(tokens ...string)
	Toggle(token string, force bool) bool
	Replace(token string, newToken string) bool
	Supports(token string) bool
	Value() string
	SetValue(v string)
}
