package gojson

import (
	"testing"
	"fmt"
)

func TestReadObject(t *testing.T) {
	lexer := new(Lexer)
	lexer.Init(`{"a":149,"b":false,"c":"hello","d":[1,2,"foo"],"e":{"hello":"world"}}`)
	ret := lexer.readObject()
	if &ret == nil {
		t.Fail()
	}
	for name, value := range ret.pairs {
		fmt.Println(name, ": ", value)
	}
}


func TestReadArray(t *testing.T) {
	lexer := new(Lexer)
	lexer.Init(`[1,2,"hello"]`)
	ret := lexer.readArray()
	if &ret == nil {
		t.Fail()
	}
}

func TestParse(t *testing.T) {
	ret := parse(`{"a":149,"b":false,"c":"hello","d":[1,2,"foo"],"e":{"hello":"world"}}`)
	ret = parse(`[1,false,"hello"]`)
	if &ret == nil {
		t.Fail()
	}
}
