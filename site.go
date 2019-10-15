package anthropoi

// Site or domain.
type Site struct {
	ID     int64
	Name   string
	Groups map[string]Group
}

// Group for a site.
type Group struct {
	ID   int64
	Name string
}
