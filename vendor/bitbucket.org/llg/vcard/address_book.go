package vcard

import (
	"log"
)

type AddressBook struct {
	Contacts []VCard
}

func (ab *AddressBook) LastContact() *VCard {
	if len(ab.Contacts) > 0 {
		return &ab.Contacts[len(ab.Contacts)-1]
	}
	return nil
}

func (ab *AddressBook) ReadFrom(di *DirectoryInfoReader) {
	contentLine := di.ReadContentLine()
	for contentLine != nil {
		switch contentLine.Name {
		case "BEGIN":
			if contentLine.Value.GetText() == "VCARD" {
				var vcard VCard
				vcard.ReadFrom(di)
				ab.Contacts = append(ab.Contacts, vcard)
			}
		default:
			log.Printf("Not read %s, %s: %s\n", contentLine.Group, contentLine.Name, contentLine.Value)
		}
		contentLine = di.ReadContentLine()
	}
}

func (ab *AddressBook) WriteTo(di *DirectoryInfoWriter) {
	for _, vcard := range ab.Contacts {
		vcard.WriteTo(di)
	}
}
