# csv2vcf

Because entering contacts 1 by 1 in vi in csv is faster than doing it in ownCloud and go is fun.

## Use

```console
csv2vcf <file.csv> <file.vcf>
```

## CSV format

```console
FILE.CSV should not have a column header and should have column ordered as
- 0: FirstName FamilyName # firstname and family name separated by space.
                          # Will take the whole column as a formatted name.
- 1: Street info ZipCode City # mostly works for french address. Very limited.
                              # Contributions welcomed :)
- 2: email@address.com
- 3: TelephoneNumber # using international number would be more useful when
                     # travelling abroad and calling friends.
- 4: OrganizationName
```
