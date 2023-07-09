package pagination

type Pagination struct {
	Page    int
	PerPage int
	Sort    string
}

func (p *Pagination) Validate() error {
	return nil
}

func (p *Pagination) Limit() int {
	return p.PerPage
}

func (p *Pagination) OffSet() int {
	return (p.Page - 1) * p.PerPage
}
