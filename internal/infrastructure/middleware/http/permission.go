package middleware

import (
	"net/http"
	"strings"

	"github.com/AlbinaKonovalova/auth-service/internal/domain/value"
)

func AdminPermissionGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		required, ok := requiredAdminPermission(r.Method, r.URL.Path)
		if !ok {
			writeJSON(w, http.StatusForbidden, apiError(7, "permission denied"))
			return
		}

		claims, ok := value.ClaimsFromContext(r.Context())
		if !ok {
			writeJSON(w, http.StatusUnauthorized, apiError(16, "missing auth claims"))
			return
		}

		if !hasPermission(claims.Permissions, required) {
			writeJSON(w, http.StatusForbidden, apiError(7, "permission denied"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func requiredAdminPermission(method, path string) (string, bool) {
	switch {
	case method == http.MethodGet && path == "/api/v1/admin/users":
		return "users.read", true

	case method == http.MethodPost && path == "/api/v1/admin/users":
		return "users.write", true

	// GET /api/v1/admin/users/{id}
	case method == http.MethodGet && isAdminUserByID(path):
		return "users.read", true

	// POST /api/v1/admin/users/{id}/activate
	case method == http.MethodPost && isAdminUserAction(path, "activate"):
		return "users.write", true

	// POST /api/v1/admin/users/{id}/deactivate
	case method == http.MethodPost && isAdminUserAction(path, "deactivate"):
		return "users.write", true

	// GET /api/v1/admin/users/{id}/roles
	case method == http.MethodGet && isAdminUserAction(path, "roles"):
		return "users.read", true

	// POST /api/v1/admin/users/{id}/roles
	case method == http.MethodPost && isAdminUserAction(path, "roles"):
		return "users.write", true

	// DELETE /api/v1/admin/users/{id}/roles/{role_code}
	case method == http.MethodDelete && isAdminUserRoleAction(path):
		return "users.write", true

	// GET /api/v1/admin/roles
	case method == http.MethodGet && path == "/api/v1/admin/roles":
		return "users.read", true

	// POST /api/v1/admin/roles
	case method == http.MethodPost && path == "/api/v1/admin/roles":
		return "users.write", true

	// GET /api/v1/admin/roles/{role_code}/permissions
	case method == http.MethodGet && isAdminRolePermissionsPath(path):
		return "users.read", true

	// POST /api/v1/admin/roles/{role_code}/permissions
	case method == http.MethodPost && isAdminRolePermissionsPath(path):
		return "users.write", true

	// DELETE /api/v1/admin/roles/{role_code}/permissions/{permission_code}
	case method == http.MethodDelete && isAdminRolePermissionRevokePath(path):
		return "users.write", true

	// GET /api/v1/admin/permissions
	case method == http.MethodGet && path == "/api/v1/admin/permissions":
		return "users.read", true

	default:
		// Unknown admin routes are denied by default.
		// Add an explicit case above when a new admin route is introduced.
		return "", false
	}
}

// isAdminUserByID возвращает true если path точно соответствует /api/v1/admin/users/{id}:
// непустой id, без вложенных сегментов.
func isAdminUserByID(path string) bool {
	const prefix = "/api/v1/admin/users/"
	if !strings.HasPrefix(path, prefix) {
		return false
	}
	rest := path[len(prefix):]
	return rest != "" && !strings.Contains(rest, "/")
}

// isAdminUserAction возвращает true если path точно соответствует
// /api/v1/admin/users/{id}/{action},
// где {id} непустой, {action} совпадает с переданным action и хвоста после него нет.
func isAdminUserAction(path, action string) bool {
	const prefix = "/api/v1/admin/users/"
	if !strings.HasPrefix(path, prefix) {
		return false
	}
	rest := path[len(prefix):]

	slash := strings.Index(rest, "/")
	if slash < 1 {
		// slash < 1: либо нет слэша вообще, либо id пустой (slash == 0)
		return false
	}

	id, suffix := rest[:slash], rest[slash+1:]
	// id не пустой, action точно совпадает, хвоста после action нет
	return id != "" && suffix == action
}

func hasPermission(perms []string, required string) bool {
	for _, p := range perms {
		if p == required {
			return true
		}
	}

	return false
}

// isAdminUserRoleAction возвращает true если path точно соответствует
// /api/v1/admin/users/{id}/roles/{role_code}:
// непустой id, сегмент "roles", непустой role_code, хвоста нет.
func isAdminUserRoleAction(path string) bool {
	const prefix = "/api/v1/admin/users/"
	if !strings.HasPrefix(path, prefix) {
		return false
	}
	rest := path[len(prefix):]

	// rest: {id}/roles/{role_code}
	slash1 := strings.Index(rest, "/")
	if slash1 < 1 {
		return false
	}
	id, rest := rest[:slash1], rest[slash1+1:]

	if id == "" {
		return false
	}

	// rest: roles/{role_code}
	if !strings.HasPrefix(rest, "roles/") {
		return false
	}
	roleCode := rest[len("roles/"):]

	return roleCode != "" && !strings.Contains(roleCode, "/")
}

// isAdminRolePermissionsPath возвращает true если path точно соответствует
// /api/v1/admin/roles/{role_code}/permissions:
// непустой role_code, сегмент "permissions", хвоста нет.
func isAdminRolePermissionsPath(path string) bool {
	const prefix = "/api/v1/admin/roles/"
	if !strings.HasPrefix(path, prefix) {
		return false
	}
	rest := path[len(prefix):]

	// rest: {role_code}/permissions
	slash := strings.Index(rest, "/")
	if slash < 1 {
		return false
	}
	roleCode, suffix := rest[:slash], rest[slash+1:]

	return roleCode != "" && suffix == "permissions"
}

// isAdminRolePermissionRevokePath возвращает true если path точно соответствует
// /api/v1/admin/roles/{role_code}/permissions/{permission_code}:
// непустой role_code, сегмент "permissions", непустой permission_code, хвоста нет.
func isAdminRolePermissionRevokePath(path string) bool {
	const prefix = "/api/v1/admin/roles/"
	if !strings.HasPrefix(path, prefix) {
		return false
	}
	rest := path[len(prefix):]

	// rest: {role_code}/permissions/{permission_code}
	slash := strings.Index(rest, "/")
	if slash < 1 {
		return false
	}
	roleCode, rest := rest[:slash], rest[slash+1:]

	if roleCode == "" {
		return false
	}

	// rest: permissions/{permission_code}
	if !strings.HasPrefix(rest, "permissions/") {
		return false
	}
	permCode := rest[len("permissions/"):]

	return permCode != "" && !strings.Contains(permCode, "/")
}
