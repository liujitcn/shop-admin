package assets

import "embed"

//go:embed openapi.yaml
var OpenApiData []byte

//go:embed web/*
var WebAssets embed.FS
