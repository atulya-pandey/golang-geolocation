package geolocation

type Importer interface {
	Parse(string) (string, error)
	Import(string) (string, error)
	Statistics()
}
