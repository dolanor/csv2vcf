package main

import (
	"bitbucket.org/llg/vcard"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	csvFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	rawCsvData, err := reader.ReadAll()

	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	vcf, err := os.Create(os.Args[2])
	if err != nil {
		log.Print(err)
		os.Exit(3)
	}
	defer vcf.Close()

	vciw := vcard.NewDirectoryInfoWriter(vcf)

	recNb := 0
	for _, record := range rawCsvData {
		var vc vcard.VCard
		if strings.Trim(record[0], " \t") != "" {
			names := strings.Split(record[0], " ")
			vc.FormattedName = record[0]
			vc.GivenNames = append(vc.GivenNames, names[0])
			if len(names) > 1 {
				vc.FamilyNames = append(vc.FamilyNames, names[1])
			}
		}
		if strings.Trim(record[1], " \t") != "" {
			addRec := strings.Fields(record[1])
			var streetEnd int
			for idx, val := range addRec {
				if match, err := regexp.Match("[0-9]{5}", []byte(val)); err == nil {
					if match == true {
						streetEnd = idx
					}
				}
			}
			street := strings.Join(addRec[:streetEnd], " ")
			address := vcard.Address{Street: street, PostalCode: addRec[streetEnd], Locality: addRec[streetEnd+1]}
			vc.Addresses = append(vc.Addresses, address)
		}
		if strings.Trim(record[2], " \t") != "" {
			vc.Emails = append(vc.Emails, vcard.Email{Address: record[2]})
		}
		if strings.Trim(record[3], " \t") != "" {
			vc.Telephones = append(vc.Telephones, vcard.Telephone{Number: record[3]})
		}
		if strings.Trim(record[4], " \t") != "" {
			vc.Org = append(vc.Org, record[4])
		}
		vc.WriteTo(vciw)
		recNb++
	}

	fmt.Println(strconv.Itoa(recNb) + " contact(s) written")
}
