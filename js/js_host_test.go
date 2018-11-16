//+build !js

package js

import (
	"testing"

	"github.com/dennwc/dom/js/jstest"
)

func TestJS(t *testing.T) {
	jstest.RunTestNodeJS(t)
}
