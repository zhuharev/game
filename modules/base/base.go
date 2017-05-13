package base

import "fmt"

var (
	DefaultItemsPerPageCount = 10
)

// PageToOffset convert page and items_per_page to limit offset for db query
func PageToOffset(page, itemsPerPage int) (offset int, err error) {
	if page == 0 {
		page = 1
	}
	if itemsPerPage < 1 {
		err = fmt.Errorf("Items per page cannot be 0")
		return
	}
	offset = (page - 1) * itemsPerPage
	return
}
