package js

// Type is an analog for JS "typeof" operator.
func (v Value) Type() Type {
	return v.Ref.Type()
}
