package pbc

/*
#include <pbc/pbc.h>
*/
import "C"
import "runtime"

type Element struct {
	pairing *Pairing // Prevents garbage collection
	cptr    *C.struct_element_s

	checked   bool
	fieldPtr  *C.struct_field_s
	isInteger bool
}

func clearElement(element *Element) {
	C.element_clear(element.cptr)
}

func makeUncheckedElement(pairing *Pairing, initialize bool, field Field) *Element {
	element := &Element{
		cptr:    &C.struct_element_s{},
		pairing: pairing,
	}
	if initialize {
		switch field {
		case G1:
			C.element_init_G1(element.cptr, pairing.cptr)
		case G2:
			C.element_init_G2(element.cptr, pairing.cptr)
		case GT:
			C.element_init_GT(element.cptr, pairing.cptr)
		case Zr:
			C.element_init_Zr(element.cptr, pairing.cptr)
		default:
			panic(ErrUnknownField)
		}
	}
	runtime.SetFinalizer(element, clearElement)
	return element
}

func makeCheckedElement(pairing *Pairing, field Field, fieldPtr *C.struct_field_s) *Element {
	element := makeUncheckedElement(pairing, true, field)
	element.checked = true
	element.fieldPtr = fieldPtr
	element.isInteger = (field == Zr)
	return element
}

func checkFieldsMatch(f1, f2 *C.struct_field_s) {
	if f1 != f2 {
		panic(ErrIncompatible)
	}
}

func (el *Element) ensureChecked() {
	if !el.checked {
		panic(ErrUncheckedOp)
	}
}

func (el *Element) checkCompatible(other *Element) {
	other.ensureChecked()
	checkFieldsMatch(el.fieldPtr, other.fieldPtr)
}

func (el *Element) checkAllCompatible(elements ...*Element) {
	for _, other := range elements {
		el.checkCompatible(other)
	}
}

func (el *Element) checkInteger() {
	el.ensureChecked()
	if !el.isInteger {
		panic(ErrIllegalOp)
	}
}
