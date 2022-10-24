package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"bitbucket.org/llg/vcard"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func csv2vcfCommand(ctx context.Context, args []string) error {
	csvFilePath := args[0]
	vcfFilePath := args[1]

	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	vcf, err := os.Create(vcfFilePath)
	if err != nil {
		return err
	}
	defer vcf.Close()

	recNb, err := parseCSVAndConvertToVCF(vcf, csvFile)
	if err != nil {
		return err
	}

	log.Printf("%d contact(s) written", recNb)
	return nil
}

func parseCSVAndConvertToVCF(vcfw io.Writer, csvr io.Reader) (int, error) {
	reader := csv.NewReader(csvr)
	reader.FieldsPerRecord = -1

	// use .Read() instead if memory issue with giant CSV files. Don't believe there are people
	// with so many contacts, though.
	rawCsvData, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("reader read all: %w", err)
	}

	vciw := vcard.NewDirectoryInfoWriter(vcfw)

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
				match, err := regexp.Match("[0-9]{5}", []byte(val))
				if err != nil {
					return 0, fmt.Errorf("regexp match: %w", err)
				}
				if err == nil {
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

	return recNb, nil
}

func main() {
	rootCmd := &ffcli.Command{
		Exec:       csv2vcfCommand,
		ShortUsage: "csv2vcf <file.csv> <file.vcf>",
		LongHelp: `
FILE.CSV should not have a column header and should have column ordered as
- 0: FirstName FamilyName # firstname and family name separated by space.
                          # Will take the whole column as a formatted name.
- 1: Street info ZipCode City # mostly works for french address. Very limited.
                              # Contributions welcomed :)
- 2: email@address.com
- 3: TelephoneNumber # using international number would be more useful when
                     # travelling abroad and calling friends.
- 4: OrganizationName`,
	}

	if len(os.Args) < 3 {
		log.Println("not enough arguments to the command line")
		log.Fatal(rootCmd.ShortUsage)

	}

	err := rootCmd.ParseAndRun(context.Background(), os.Args[1:])
	if err != nil {
		log.Println(err)
		log.Fatal(rootCmd.LongHelp)
	}
}
