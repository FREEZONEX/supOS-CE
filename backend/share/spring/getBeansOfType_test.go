package spring

import (
	"fmt"
	"reflect"
	"testing"
)

type Animal interface {
	Move()
}
type Bird interface {
	Move()
	Fly()
}
type Eagle struct {
	id int64
}

func (e *Eagle) Move() {
	fmt.Printf("The Eagle[%d] is moving on the glass land..\n", e.id)
}
func (e *Eagle) Fly() {
	fmt.Printf("The Eagle[%d] is Flying in the sky.\n", e.id)
}

type Nightingale struct {
	id int64
}

func (e *Nightingale) Move() {
	fmt.Printf("The Nightingale[%d] is moving on the land..\n", e.id)
}
func (e *Nightingale) Fly() {
	fmt.Printf("The Nightingale[%d] is Flying in the forest.\n", e.id)
}

type Swallow struct {
	id int64
}

func (e *Swallow) Move() {
	fmt.Printf("The Swallow[%d] is moving on the land..\n", e.id)
}
func (e *Swallow) Fly() {
	fmt.Printf("The Swallow[%d] is Flying in the forest.\n", e.id)
}
func TestTypeCast(t *testing.T) {
	egleA := &Eagle{id: 125}
	var a any = egleA
	if bd, ok := a.(Bird); ok {
		t.Logf("egleA[%v] bird", bd)
	}
	rs, ok := a.(Bird)
	var BD Bird
	Interface := reflect.TypeOf(&BD).Elem()
	t.Log("egle is bird ? ", rs, ok)
	t.Log("Implements ? ", reflect.TypeOf(a).Implements(Interface))
	t.Log("convertTo ? ", reflect.TypeOf(a).ConvertibleTo(Interface))
}

func TestRegisterBean(t *testing.T) {
	egleA := &Eagle{id: 125}
	ng := &Nightingale{id: 126}
	RegisterBean(egleA)
	RegisterLazy(func() *Nightingale {
		return ng
	})
	RegisterLazy[*Swallow](func() *Swallow {
		return nil
	})
	//rsBird := GetBean[Bird]()
	//rsBird.Fly()
	list := GetBeansOfType[Bird]()
	for i, a := range list {
		//if a != nil {
		//	a.Fly()
		//}
		t.Logf("Bean[%d]: %+v\n", i, a)
	}
}
