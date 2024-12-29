package static

import "embed"

//go:embed static/*
var Static embed.FS
