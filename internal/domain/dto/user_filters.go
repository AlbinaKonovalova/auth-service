package dto

type UserListFilters struct {
	IsActive *bool

	Role string

	Page int

	PerPage int
}
