package tvmaze

type network struct {
	ID      int
	Name    string
	Country country
}

type country struct {
	Name     string
	Code     string
	Timezone string
}
