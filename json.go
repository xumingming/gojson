package gojson

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Number struct {
	// the actual data which holds the number
	data string
	// whether the number is negative
	negative bool
	// whether the number is a floating point number
	isFloat bool
}

func (n Number) Float64() float64 {
	f, _ := strconv.ParseFloat(n.data, 0)
	if n.negative {
		f = 0 - f
	}

	return f
}

func (this *Number) FloatPrecision() int {
	return len(this.data) - strings.Index(this.data, ".") - 1
}

func (n Number) Int64() int64 {
	i, _ := strconv.ParseInt(n.data, 10, 0)
	if n.negative {
		i = 0 - i
	}

	return i
}

type JSONObject struct {
	pairs map[string]interface{}
}

func (this JSONObject) String() string {
	f := Formatter{}
	f.formatJSONObject(this)
	return f.String()
}

type JSONArray struct {
	values []interface{}
}

func (this JSONArray) String() string {
	f := Formatter{}
	f.formatJSONArray(this)
	return f.String()
}

type Lexer struct {
	content string
	index   int
	char    uint8
}

func NewLexer(content string) *Lexer {
	lexer := new(Lexer)
	lexer.content = content
	lexer.index = -1
	lexer.char = 0
	lexer.nextChar()

	return lexer
}

func (this *Lexer) match(x uint8) bool {
	return this.char == x
}

func (this *Lexer) accept(x uint8) {
	if this.match(x) {
		this.nextChar()
	} else {
		panic(fmt.Sprintf("expecting '%v', got '%v'[index: %v]", string(x), string(this.char), this.index))
	}
}

func (this *Lexer) nextChar() uint8 {
	this.index += 1
	if this.index < len(this.content) {
		this.char = this.content[this.index]
	} else {
		this.char = 0
	}

	return this.char
}

func (this *Lexer) prevChar() uint8 {
	this.index -= 1
	if this.index > -1 {
		this.char = this.content[this.index]
	} else {
		this.char = 0
	}

	return this.char
}

func (this *Lexer) readString() string {
	var ret bytes.Buffer
	this.accept('"')
	escaping := false
	for !this.match('"') || escaping {
		if this.match('\\') {
			escaping = true
			this.nextChar()
			continue
		}

		if escaping {
			switch this.char {
			case '\\':
				ret.WriteByte('\\')
			case 'b':
				ret.WriteByte('\b')
			case 'f':
				ret.WriteByte('\f')
			case 'n':
				ret.WriteByte('\n')
			case 'r':
				ret.WriteByte('\r')
			case 't':
				ret.WriteByte('\t')
			case '"':
				ret.WriteByte('"')
			case 'u':
				i32 := int32(0)
				for i := 0; i < 4; i++ {
					i32 *= 16
					this.nextChar()
					if this.char <= '9' {
						i32 += int32(this.char - '0')
					} else {
						i32 += int32(this.char - 'a' + 10)
					}
				}
				ret.WriteRune(i32)
			}

			escaping = false
		} else {
			ret.WriteByte(this.char)
		}

		this.nextChar()
	}

	this.accept('"')
	return ret.String()
}

func (this *Lexer) readNumber() Number {
	var ret bytes.Buffer
	negative := false
	isFloat := false
	if this.match('-') {
		negative = true
		this.accept('-')
	}

	for '0' <= this.char && this.char <= '9' {
		ret.WriteByte(this.char)
		this.nextChar()
	}

	if this.match('.') {
		this.nextChar()
		if isDigit(this.char) {
			isFloat = true
			ret.WriteByte('.')
			for '0' <= this.char && this.char <= '9' {
				ret.WriteByte(this.char)
				this.nextChar()
			}
		} else {
			panic("unexpected symbol: '.'")
		}
	}

	return Number{ret.String(), negative, isFloat}
}

func (this *Lexer) readBoolean() bool {
	var ret bytes.Buffer
	if this.match('t') || this.match('f') {
		//      for this.char != ',' && this.char != ' ' && this.char != '}' {
		for this.char != ',' && !isBlank(this.char) && this.char != '}' && this.char != 0 {
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
	this.skipBlank()
	this.accept('{')
	var ret JSONObject
	ret.pairs = map[string]interface{}{}
	name, value := this.readPair()
	ret.pairs[name] = value

	for this.match(',') {
		this.skipBlank()
		this.accept(',')
		this.skipBlank()
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

func (this *Lexer) readNull() {
	this.nextChar()
	this.nextChar()
	this.nextChar()
}

func (this *Lexer) readValue() (value interface{}) {
	if this.match('"') {
		value = this.readString()
	} else if '0' <= this.char && this.char <= '9' || this.match('-') {
		value = this.readNumber()
	} else if this.match('t') || this.match('f') {
		value = this.readBoolean()
	} else if this.match('{') {
		value = this.readObject()
	} else if this.match('[') {
		value = this.readArray()
	} else {
		value = nil
		this.readNull()
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
	for isBlank(this.char) {
		this.nextChar()
	}
}

func isBlank(x uint8) bool {
	return x == ' ' || x == '\t' || x == '\r' || x == '\n'
}

func isDigit(x uint8) bool {
	return x >= '0' && x <= '9'
}

func Parse(str string) interface{} {
	lexer := NewLexer(str)
	if lexer.match('{') {
		return lexer.readObject()
	} else if lexer.match('[') {
		return lexer.readArray()
	}
	return nil
}

type Formatter struct {
	tabCnt int
	buf    bytes.Buffer
}

func (this *Formatter) incrTabCnt() {
	this.tabCnt++
}

func (this *Formatter) decrTabCnt() {
	this.tabCnt--
}

func (this *Formatter) String() string {
	return this.buf.String()
}

func (this *Formatter) newline() {
	this.buf.WriteByte('\n')
	for i := 0; i < this.tabCnt; i++ {
		this.buf.WriteByte('\t')
	}
}

func (this *Formatter) formatJSONObject(obj JSONObject) {
	this.buf.WriteByte('{')
	this.incrTabCnt()
	i := 0
	for name, value := range obj.pairs {
		if i > 0 {
			this.buf.WriteByte(',')
		}
		this.newline()
		this.formatPair(name, value)
		i++
	}
	this.decrTabCnt()
	this.newline()
	this.buf.WriteByte('}')
}

func (this *Formatter) formatJSONArray(arr JSONArray) {
	this.buf.WriteByte('[')
	for i, value := range arr.values {
		this.formatValue(value)
		if i < len(arr.values)-1 {
			this.buf.WriteByte(',')
		}
	}

	this.buf.WriteByte(']')
}

func (this *Formatter) formatValue(value interface{}) {
	switch value.(type) {
	case string:
		this.buf.WriteByte('"')
		this.buf.WriteString(value.(string))
		this.buf.WriteByte('"')
	case bool:
		this.buf.WriteString(strconv.FormatBool(value.(bool)))
	case JSONObject:
		this.formatJSONObject(value.(JSONObject))
	case JSONArray:
		this.formatJSONArray(value.(JSONArray))
	case Number:
		num := value.(Number)
		if num.isFloat {
			this.buf.WriteString(strconv.FormatFloat(num.Float64(), 'f', num.FloatPrecision(), 64))
		} else {
			this.buf.WriteString(strconv.FormatInt(num.Int64(), 10))
		}
	}
}

func (this *Formatter) formatPair(name string, value interface{}) {
	this.buf.WriteByte('"')
	this.buf.WriteString(name)
	this.buf.WriteString("\": ")
	this.formatValue(value)
}
