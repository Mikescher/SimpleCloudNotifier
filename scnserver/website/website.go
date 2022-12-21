package website

import "embed"

//go:embed *
//go:embed css/*
//go:embed js/*
var Assets embed.FS
