package gojson

import (
	"fmt"
	"strconv"
	"bytes"
)

type Lexer struct {
	/* a very good string */
	content string
	index   int
	char    uint8
}

func (this *Lexer) Init(content string) {
	this.content = content
	this.index = -1
	this.char = 99
	this.nextChar()
}

func (this *Lexer) match(x uint8) bool {
	return this.char == x
}

func (this *Lexer) accept(x uint8) {
	if this.match(x) {
		this.nextChar()
	} else {
		panic(fmt.Sprintf("expecting %v, got %v[%v]", string(x), string(this.char), this.index))
	}
}

func (this *Lexer) nextChar() {
	this.index += 1
	if this.index < len(this.content) {
		this.char = this.content[this.index]
	} else {
		this.char = 0
	}
}

func (this *Lexer) readString() string {
	var ret bytes.Buffer
	this.accept('"')
	for !this.match('"') {
		ret.WriteByte(this.char)
		this.nextChar()
	}

	this.accept('"')
	return ret.String()
}

func (this *Lexer) readInt() int {
	var ret bytes.Buffer
	for '0' <= this.char && this.char <= '9' {
		ret.WriteByte(this.char)
		this.nextChar()
	}

	i, error := strconv.Atoi(ret.String())
	if error != nil {
		fmt.Printf("the error is: %v", error)
	}
	return i
}

func (this *Lexer) readBoolean() bool {
	var ret bytes.Buffer
	if this.char == 't' || this.char == 'f' {
		for this.char != ',' && this.char != ' ' && this.char != '}' {
			ret.WriteByte(this.char)
			this.nextChar()
		}
	}

	if ret.String() == "true" {
		return true
	} else {
		return false
	}
}

func (this *Lexer) readObject() JSONObject {
	this.accept('{')
	var ret JSONObject;
	ret.pairs = map[string]interface{}{}
	name, value := this.readPair()
	ret.pairs[name] = value

	for this.match(',') {
		this.accept(',')
		name, value = this.readPair()
		ret.pairs[name] = value
	}
	this.skipBlank()
	this.accept('}')

	return ret
}

func (this *Lexer) readArray() JSONArray {
	this.accept('[')
	var ret JSONArray
	var values []interface{}

	value := this.readValue()
	values = append(values, value)
	for this.match(',') {
		this.accept(',')
		value = this.readValue()
		values = append(values, value)
	}
	this.accept(']')
	ret.values = values
	return ret
}

//
func (this *Lexer) readNil() {
	this.nextChar()
	this.nextChar()
	this.nextChar()
}

//
func (this *Lexer) readValue() (value interface{}) {
	if this.match('"') {
		value = this.readString()
	} else if '0' <= this.char && this.char <= '9' {
		value = this.readInt()
	} else if this.match('t') || this.match('f') {
		value = this.readBoolean()
	} else if this.match('{') {
		value = this.readObject()
	} else if this.match('[') {
		value = this.readArray()
	} else {
		value = nil
		this.readNil()
	}

	return value
}

func (this *Lexer) readPair() (name string, value interface{}) {
	name = this.readString()
	this.skipBlank()
	this.accept(':')
	this.skipBlank()
	value = this.readValue()
	return
}

func (this *Lexer) skipBlank() {
	for ; isBlank(this.char) ; {
		this.nextChar()
	}
}

func isBlank (x uint8) bool {
	return x == ' ' || x == '\t' || x == '\r' || x == '\n'
}

func Parse(str string) interface{} {
	lexer := new(Lexer)
	lexer.Init(str)
	if lexer.match('{') {
		return lexer.readObject()
	} else if lexer.match('[') {
		return lexer.readArray()
	}
	return nil
}

type JSONObject struct {
	pairs map[string]interface{}
}

type JSONArray struct {
	values []interface{}
}
