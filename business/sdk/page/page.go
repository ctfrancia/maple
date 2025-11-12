package page

import (
	"fmt"
	"strconv"
)

// Page represents a page of data.
type Page struct {
	number int
	rows   int
}

func Parse(page, rowsPerPage string) (Page, error) {
	number := 1

	if page != "" {
		var err error
		number, err = strconv.Atoi(page)
		if err != nil {
			return Page{}, fmt.Errorf("page conversion: %w", err)
		}
	}

	rows := 10
	if rowsPerPage != "" {
		var err error
		rows, err = strconv.Atoi(rowsPerPage)
		if err != nil {
			return Page{}, fmt.Errorf("rows conversion: %w", err)
		}
	}

	if number <= 0 {
		return Page{}, fmt.Errorf("page value too small, must be larger than 0")
	}

	if rows <= 0 {
		return Page{}, fmt.Errorf("rows value too small, must be larger than 0")
	}

	if rows > 100 {
		return Page{}, fmt.Errorf("rows value too large, must be less than 100")
	}

	p := Page{
		number: number,
		rows:   rows,
	}

	return p, nil
}

// MustParse is the same as Parse but panics if an error occurs.
func MustParse(page, rowsPerPage string) Page {
	pg, err := Parse(page, rowsPerPage)
	if err != nil {
		panic(err)
	}

	return pg
}

// String implements the stringer interface
func (p Page) String() string {
	return fmt.Sprintf("page: %d rows: %d", p.number, p.rows)
}

// Number returns the rows per page.
func (p Page) Number() int {
	return p.number
}

// RowsPerPage returns the number of rows per page.
func (p Page) RowsPerPage() int {
	return p.rows
}
