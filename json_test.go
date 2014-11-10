package gojson

import (
    "testing"
)

func TestReadInt(t *testing.T) {
    var lexer *Lexer
    var i int
    lexer = NewLexer("1234567")
    i = lexer.readInt()
    if i != 1234567 {
        t.Fail()
    }

    // test a big int(larger than 2 ^ 32)
    lexer = NewLexer("8589934592")
    i = lexer.readInt()
    if i != 8589934592 {
        t.Fail()
    }

	// test negative number
    lexer = NewLexer("-100")
    i = lexer.readInt()
    if i != -100 {
        t.Fail()
    }
}

func TestReadString(t *testing.T) {
    var lexer *Lexer
    var str string

    // test normal case
    lexer = NewLexer(`"hello"`)
    str = lexer.readString()
    if str != "hello" {
        t.Fail()
    }

    // test \n
    lexer = NewLexer(`"hello world\n"`)
    str = lexer.readString()
    if str != "hello world\n" {
        t.Fail()
    }

    // test \"
    lexer = NewLexer(`"hello\""`)
    str = lexer.readString()
    if str != "hello\"" {
        t.Fail()
    }

    // test unicode
    lexer = NewLexer(`"你好\u554a"`)
    str = lexer.readString()
    if str != "你好啊" {
        t.Fail()
    }
}

func TestReadBoolean(t *testing.T) {
    var lexer *Lexer
    var ret bool
    
    lexer = NewLexer(`true`)
    ret = lexer.readBoolean()
    if ret != true {
        t.Fail()
    }

    lexer = NewLexer(`false`)
    ret = lexer.readBoolean()
    if ret != false {
        t.Fail()
    }
    
    lexer = NewLexer(`true1`)
    ret = lexer.readBoolean()
    if ret == true {
        t.Fail()
    }
}

func TestReadNil(t *testing.T) {

}

func TestReadValue(t *testing.T) {

}

func TestReadPair(t *testing.T) {

}

func TestReadObject(t *testing.T) {
    var lexer *Lexer
    var ret JSONObject
    
    lexer = NewLexer(`{"a": 149,"b":false,"c":"hello" }`)
    ret = lexer.readObject()
    if len(ret.pairs) != 3 {
        t.Fail()
    }
    
    if ret.pairs["a"] != 149 {
        t.Fail()
    }

    if ret.pairs["b"] != false {
        t.Fail()
    }

    if ret.pairs["c"] != "hello" {
        t.Fail()
    }

    // test leading and trailing blanks
    lexer = NewLexer(` 
                      {"a": 149} `)
    ret = lexer.readObject()
    if len(ret.pairs) != 1 {
        t.Fail()
    }
    // for name, value := range ret.pairs {
    //  fmt.Println(name, ": ", value)
    // }
}

func TestReadNestedObject(t *testing.T) {
    lexer := NewLexer(`{"a":149,"b":false,"c":"hello","d":[1,2,"foo"],"e":{"hello":"world"}}`)
    ret := lexer.readObject()
    if &ret == nil {
        t.Fail()
    }
    // for name, value := range ret.pairs {
    //  fmt.Println(name, ": ", value)
    // }
}

func TestReadArray(t *testing.T) {
    lexer := NewLexer(`[1,2,"hello"]`)
    ret := lexer.readArray()
    if &ret == nil {
        t.Fail()
    }
}

func TestParse(t *testing.T) {
    ret := Parse(`{"a":149,"b":false,"c":"hello","d":[1,2,"foo"],"e":{"hello":"world"}}`)
    ret = Parse(`[1,false,"hello"]`)
    if &ret == nil {
        t.Fail()
    }
}
