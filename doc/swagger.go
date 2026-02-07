package doc

import (
	"embed"
	"net/http"
)

//go:embed swagger.yaml
var swagger embed.FS
var Swagger = http.FileServer(http.FS(swagger))
