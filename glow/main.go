package main

import "fmt"

func main() {
	fmt.Println("hello")
}

func NewMsg(args ...interface{}) Msg {
	return Msg(args)
}

type Msg []interface{}

func (m Msg) At(index int) Msg {
	if index < len(m) {
		if m2, ok := m[index].(Msg); ok {
			return m2
		}
		return NewMsg(m[index])
	}
	return Msg{}
}

func (m Msg) IsBang() bool {
	return len(m) == 0
}

func (m Msg) IsInt() (ok bool) {
	if len(m) == 1 {
		_, ok = m[0].(int)
	}
	return
}

func (m Msg) IsString() (ok bool) {
	if len(m) == 1 {
		_, ok = m[0].(string)
	}
	return
}

func (m Msg) AsInt() int {
	if m.IsInt() {
		return m[0].(int)
	}
	//fmt.Println("not an int:", m)
	return 0
}

func (m Msg) AsString() string {
	if m.IsString() {
		return m[0].(string)
	}
	//fmt.Println("not a string:", m)
	return ""
}
