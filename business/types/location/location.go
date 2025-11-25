package location

type Location struct {
	venueName  string  // name of the venue
	address    string  // address of the venue (Calle de la Palma 123, etc.)
	city       string  // city of the venue (Barcelona, Madrid, etc.)
	state      string  // state/region of the venue (Catalonia, Aragon, etc.)
	postalCode string  // postal code of the venue (08003, 28001, etc.)
	country    string  // country of the venue (Spain, France, etc.)
	latitude   float64 // latitude of the venue (40.4167, -3.7028, etc.)
	longitude  float64 // longitude of the venue (2.1734, 103.8328, etc.)
}

func (l Location) String() string {
	return l.venueName
}
