package vcard

import (
	"io"
)

// Permit to serialize Directory Information data as defined by RFC 2425
type DirectoryInfoWriter struct {
	writer io.Writer
}

//create a new DirectoryInfoWriter 
func NewDirectoryInfoWriter(writer io.Writer) *DirectoryInfoWriter {
	return &DirectoryInfoWriter{writer}
}

func (di *DirectoryInfoWriter) WriteContentLine(contentLine *ContentLine) {
	if contentLine.Group != "" {
		io.WriteString(di.writer, contentLine.Group)
		io.WriteString(di.writer, ".")
	}
	io.WriteString(di.writer, contentLine.Name)
	if contentLine.Params != nil {
		for key, values := range contentLine.Params {
			io.WriteString(di.writer, ";")
			io.WriteString(di.writer, key)
			if len(values) > 0 {
				io.WriteString(di.writer, "=")
				for vi := 0; vi < len(values); vi++ {
					io.WriteString(di.writer, values[vi])
					if vi+1 < len(values) {
						io.WriteString(di.writer, ",")
					}
				}
			}
		}
	}
	io.WriteString(di.writer, ":")
	for si := 0; si < len(contentLine.Value); si++ {
		for vi := 0; vi < len(contentLine.Value[si]); vi++ {
			di.WriteValue(contentLine.Value[si][vi])
			if vi+1 < len(contentLine.Value[si]) {
				io.WriteString(di.writer, ",")
			}
		}
		if si+1 < len(contentLine.Value) {
			io.WriteString(di.writer, ";")
		}
	}
	io.WriteString(di.writer, "\r\n")
}

// this function escape '\n' '\r' ';' ',' character with the '\\' character 
func (di *DirectoryInfoWriter) WriteValue(value string) {
	i := 0
	for _, c := range value {
		if i == 76 {
			// if line to long fold value on multiple line 
			io.WriteString(di.writer, "\n  ")
			i = 0
		}
		var e string
		switch c {
		case '\r':
			e = `\r`
		case '\n':
			e = `\n`
		case ';':
			e = `\;`
		case ':':
			e = `\:`
		case ',':
			e = `\,`
		default:
			// c is an int representing a Unicode code point. 
			// convert it to string (UTF-8 encoded character) 
			e = string(c)
		}
		io.WriteString(di.writer, e)
		i++
	}
}
