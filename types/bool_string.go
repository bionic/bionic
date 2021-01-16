package types

type BoolString bool

func (b *BoolString) UnmarshalCSV(csv string) (err error) {
	if csv == "true" {
		*b = true
	} else {
		*b = false
	}

	return nil
}
