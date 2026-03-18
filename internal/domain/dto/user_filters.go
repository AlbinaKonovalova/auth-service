package dto

// UserListFilters содержит параметры фильтрации и пагинации списка пользователей.
type UserListFilters struct {
	// IsActive — фильтр по статусу активности; nil означает «все пользователи».
	IsActive *bool

	// Role — фильтр по коду роли; пустая строка означает «без фильтра по роли».
	Role string

	// Page — номер страницы (начиная с 1).
	Page int

	// PerPage — количество записей на страницу.
	PerPage int
}
