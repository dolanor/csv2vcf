package main

import (
	"bytes"
	_ "embed"
	"log"
	"testing"
)

//go:embed testdata/contacts.csv
var contacts string

//go:embed testdata/contacts.vcf
var want string

func TestParseCSVAndConvertToVCF(t *testing.T) {
	csvContactsReader := bytes.NewBufferString(contacts)
	var vcfWriter bytes.Buffer

	got, err := parseCSVAndConvertToVCF(&vcfWriter, csvContactsReader)
	if err != nil {
		t.Fatal("should work")
	}
	if got != 2 {
		t.Fatal("should have parsed 2 contacts")
	}

	if vcfWriter.String() != want {
		log.Fatalf("should have 2 contacts in correct VCF format.\ngot:\n%v\nwant:\n%v", vcfWriter.String(), want)
	}
}
