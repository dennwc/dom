package dom

import (
	"strconv"
)

var (
	_ Unit = Px(0)
	_ Unit = Em(0)
	_ Unit = Rem(0)
	_ Unit = Vw(0)
	_ Unit = Vh(0)
	_ Unit = Vmin(0)
	_ Unit = Vmax(0)
	_ Unit = Perc(0)
)

type Unit interface {
	String() string
}

type Auto struct{}

func (Auto) String() string {
	return "auto"
}

type Px int

func (v Px) String() string {
	return strconv.Itoa(int(v)) + "px"
}

type Em float64

func (v Em) String() string {
	return strconv.FormatFloat(float64(v), 'g', -1, 64) + "em"
}

type Rem int

func (v Rem) String() string {
	return strconv.Itoa(int(v)) + "rem"
}

type Vw int

func (v Vw) String() string {
	return strconv.Itoa(int(v)) + "vw"
}

type Vh int

func (v Vh) String() string {
	return strconv.Itoa(int(v)) + "vh"
}

type Vmin int

func (v Vmin) String() string {
	return strconv.Itoa(int(v)) + "vmin"
}

type Vmax int

func (v Vmax) String() string {
	return strconv.Itoa(int(v)) + "vmax"
}

type Perc int

func (v Perc) String() string {
	return strconv.Itoa(int(v)) + "%"
}
