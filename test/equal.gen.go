package test

import (
	"bytes"
)

var _ = bytes.MinRead

func derivEqualPtrToA(this, that *A) bool {
	return (this == nil && that == nil) ||
		(this != nil && that != nil) &&
			this.B == that.B &&
			((this.C == nil && that.C == nil) || (this.C != nil && that.C != nil && *this.C == *that.C)) &&
			this.D == that.D &&
			((this.E == nil && that.E == nil) || (this.E != nil && that.E != nil && *this.E == *that.E)) &&
			derivEqualSliceOfbool(this.I, that.I) &&
			derivEqualSliceOfPtrToB(this.J, that.J) &&
			derivEqualSliceOfB(this.K, that.K) &&
			this.L.Equal(that.L) &&
			this.M.Equal(&that.M) &&
			derivEqualMapOfintToB(this.N, that.N) &&
			derivEqualMapOfstringToPtrToB(this.O, that.O) &&
			derivEqualMapOfint64Tostring(this.P, that.P)
}
func derivEqualPtrToB(this, that *B) bool {
	return (this == nil && that == nil) ||
		(this != nil && that != nil) &&
			bytes.Equal(this.Bytes, that.Bytes) &&
			derivEqualMapOfintToB(this.N, that.N)
}
func derivEqualMapOfint64Tostring(this, that map[int64]string) bool {
	if this == nil {
		if that == nil {
			return true
		} else {
			return false
		}
	} else if that == nil {
		return false
	}
	if len(this) != len(that) {
		return false
	}
	for k, v := range this {
		thatv, ok := that[k]
		if !ok {
			return false
		}
		if !(v == thatv) {
			return false
		}
	}
	return true

}
func derivEqualSliceOfbool(this, that []bool) bool {
	if this == nil {
		if that == nil {
			return true
		} else {
			return false
		}
	} else if that == nil {
		return false
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(this[i] == that[i]) {
			return false
		}
	}
	return true

}
func derivEqualSliceOfPtrToB(this, that []*B) bool {
	if this == nil {
		if that == nil {
			return true
		} else {
			return false
		}
	} else if that == nil {
		return false
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(this[i].Equal(that[i])) {
			return false
		}
	}
	return true

}
func derivEqualSliceOfB(this, that []B) bool {
	if this == nil {
		if that == nil {
			return true
		} else {
			return false
		}
	} else if that == nil {
		return false
	}
	if len(this) != len(that) {
		return false
	}
	for i := 0; i < len(this); i++ {
		if !(this[i].Equal(&that[i])) {
			return false
		}
	}
	return true

}
func derivEqualMapOfintToB(this, that map[int]B) bool {
	if this == nil {
		if that == nil {
			return true
		} else {
			return false
		}
	} else if that == nil {
		return false
	}
	if len(this) != len(that) {
		return false
	}
	for k, v := range this {
		thatv, ok := that[k]
		if !ok {
			return false
		}
		if !(v.Equal(&thatv)) {
			return false
		}
	}
	return true

}
func derivEqualMapOfstringToPtrToB(this, that map[string]*B) bool {
	if this == nil {
		if that == nil {
			return true
		} else {
			return false
		}
	} else if that == nil {
		return false
	}
	if len(this) != len(that) {
		return false
	}
	for k, v := range this {
		thatv, ok := that[k]
		if !ok {
			return false
		}
		if !(v.Equal(thatv)) {
			return false
		}
	}
	return true

}
