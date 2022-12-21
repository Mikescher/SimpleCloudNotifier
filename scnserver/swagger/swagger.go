package swagger

import (
	"blackforestbytes.com/simplecloudnotifier/api/ginresp"
	"embed"
	_ "embed"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//go:embed *.html
//go:embed *.json
//go:embed *.yaml
//go:embed *.js
//go:embed *.css
//go:embed themes/*
var assets embed.FS

func getAsset(fn string) ([]byte, string, bool) {
	data, err := assets.ReadFile(fn)
	if err != nil {
		return nil, "", false
	}

	mime := "text/plain"

	lowerFN := strings.ToLower(fn)
	if strings.HasSuffix(lowerFN, ".html") || strings.HasSuffix(lowerFN, ".htm") {
		mime = "text/html"
	} else if strings.HasSuffix(lowerFN, ".css") {
		mime = "text/css"
	} else if strings.HasSuffix(lowerFN, ".js") {
		mime = "text/javascript"
	} else if strings.HasSuffix(lowerFN, ".json") {
		mime = "application/json"
	} else if strings.HasSuffix(lowerFN, ".jpeg") || strings.HasSuffix(lowerFN, ".jpg") {
		mime = "image/jpeg"
	} else if strings.HasSuffix(lowerFN, ".png") {
		mime = "image/png"
	} else if strings.HasSuffix(lowerFN, ".svg") {
		mime = "image/svg+xml"
	}

	return data, mime, true
}

func Handle(g *gin.Context) ginresp.HTTPResponse {
	type uri struct {
		Filename string `uri:"sub"`
	}

	var u uri
	if err := g.ShouldBindUri(&u); err != nil {
		return ginresp.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	u.Filename = strings.TrimLeft(u.Filename, "/")

	if u.Filename == "" {
		index, _, _ := getAsset("index.html")
		return ginresp.Data(http.StatusOK, "text/html", index)
	}

	if data, mime, ok := getAsset(u.Filename); ok {
		return ginresp.Data(http.StatusOK, mime, data)
	}

	return ginresp.JSON(http.StatusNotFound, gin.H{"error": "AssetNotFound", "filename": u.Filename})
}
