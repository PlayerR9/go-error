package internal

type Key interface {
}

func NewKey() Key {
	type key struct{}

	return key{}
}

type Info struct {
	info map[Key]any
}

func (info *Info) Add() {

}

func (info Info) Get() {

}

func (info *Info) Set() {

}

func (info *Info) Edit() {

}

func (info *Info) Delete() {

}
