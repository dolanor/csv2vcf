// See RFC 2425 Mime Content-Type for Directory Information
package vcard

import (
	"io"
	"text/scanner"
)

type DirectoryInfoReader struct {
	scan *scanner.Scanner
}

func NewDirectoryInfoReader(reader io.Reader) *DirectoryInfoReader {
	var s scanner.Scanner
	s.Init(reader)
	return &DirectoryInfoReader{&s}
}

func (di *DirectoryInfoReader) ReadContentLine() *ContentLine {
	if di.scan.Peek() == scanner.EOF {
		return nil
	}
	group, name := di.readGroupName()
	params := make(map[string]Value)
	if di.scan.Peek() == ';' {
		params = di.readParameters()
	}
	di.scan.Next()
	value := di.readValues()
	return &ContentLine{group, name, params, value}
}

func (di *DirectoryInfoReader) readGroupName() (group, name string) {
	c := di.scan.Peek()
	var buf []rune
	for c != scanner.EOF {
		if c == '.' {
			group = string(buf)
			buf = []rune{}
		} else if c == ';' || c == ':' {
			name = string(buf)
			return
		} else {
			buf = append(buf, c)
		}
		di.scan.Next()
		c = di.scan.Peek()
	}
	return
}

func (di *DirectoryInfoReader) readParameters() (params map[string]Value) {
	lastChar := di.scan.Peek()
	c := lastChar
	var buf []rune
	var name string
	var value string
	params = make(map[string]Value)
	var values Value
	for c != scanner.EOF {
		if c == ',' {
			values = append(values, string(buf))
			buf = []rune{}
		} else if c == ';' || c == ':' {
			if name == "" {
				name = string(buf)
			} else {
				value = string(buf)
			}
			if name != "" {
				values = append(values, value)
				if _, ok := params[name]; ok {
					params[name] = append(params[name], values...)
				} else {
					params[name] = values
				}
			}
			if c == ':' {
				return
			}
			buf = []rune{}
			values = Value{}
			name = ""
			value = ""
		} else if c == '=' {
			name = string(buf)
			buf = []rune{}
		} else {
			buf = append(buf, c)
		}
		di.scan.Next()
		c = di.scan.Peek()
	}
	return
}

func (di *DirectoryInfoReader) readValues() (value StructuredValue) {
	lastChar := di.scan.Next()
	c := lastChar
	var buf []rune
	escape := false
	var val Value
	for c != scanner.EOF {
		if c == '\n' {
			la := di.scan.Peek()
			if la != 32 && la != 9 {
				// return
				if len(buf) > 0 {
					val = append(val, string(buf))
				}
				value = append(value, val)
				return
			} else {
				// unfold
				lastChar = la
				c = di.scan.Next()
				for c == 32 || c == 9 {
					c = di.scan.Next()
				}
			}
		}
		if c == '\\' {
			escape = true
		} else if escape {
			if c == 'n' || c == 'N' {
				c = '\n'
			}
			buf = append(buf, c)
			escape = false
		} else if c == ',' {
			if len(buf) > 0 {
				val = append(val, string(buf))
				buf = []rune{}
			}
		} else if c == ';' {
			if len(buf) > 0 {
				val = append(val, string(buf))
				buf = []rune{}
			}
			value = append(value, val)
			val = Value{}
		} else if c != '\n' && c != '\r' {
			buf = append(buf, c)
		}
		lastChar = c
		c = di.scan.Next()
	}
	return
}
