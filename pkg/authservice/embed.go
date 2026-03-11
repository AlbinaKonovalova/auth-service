package authservice

import _ "embed"

//go:embed v1/authservice.swagger.json
var SwaggerSpec []byte
