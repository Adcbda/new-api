package router

import (
	"embed"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// ThemeAssets holds the embedded frontend assets for both themes.
type ThemeAssets struct {
	DefaultBuildFS   embed.FS
	DefaultIndexPage []byte
	ClassicBuildFS   embed.FS
	ClassicIndexPage []byte
}

func SetWebRouter(router *gin.Engine, assets ThemeAssets) {
	defaultFS := common.EmbedFolder(assets.DefaultBuildFS, "web/default/dist")
	classicFS := common.EmbedFolder(assets.ClassicBuildFS, "web/classic/dist")
	themeFS := common.NewThemeAwareFS(defaultFS, classicFS)

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.GlobalWebRateLimit())
	router.Use(middleware.Cache())
	router.Use(static.Serve("/", themeFS))
	router.NoRoute(func(c *gin.Context) {
		c.Set(middleware.RouteTagKey, "web")
		if strings.HasPrefix(c.Request.RequestURI, "/v1") || strings.HasPrefix(c.Request.RequestURI, "/api") || strings.HasPrefix(c.Request.RequestURI, "/assets") {
			controller.RelayNotFound(c)
			return
		}
		c.Header("Cache-Control", "no-cache")
		if common.GetTheme() == "classic" {
			c.Data(http.StatusOK, "text/html; charset=utf-8", prepareIndexPage(assets.ClassicIndexPage))
		} else {
			c.Data(http.StatusOK, "text/html; charset=utf-8", prepareIndexPage(assets.DefaultIndexPage))
		}
	})
}

func prepareIndexPage(indexPage []byte) []byte {
	if common.AppBasePath == "" {
		return indexPage
	}

	html := prefixRootHTMLReferences(string(indexPage), common.AppBasePath)
	basePathJSON, err := common.Marshal(common.AppBasePath)
	if err != nil {
		basePathJSON = []byte(`""`)
	}
	injection := `<base href="` + common.AppBasePath + `/"><script>window.__APP_BASE_PATH__=` + string(basePathJSON) + `;</script>`
	if strings.Contains(html, "</head>") {
		html = strings.Replace(html, "</head>", injection+"</head>", 1)
	} else {
		html = injection + html
	}
	return []byte(html)
}

func prefixRootHTMLReferences(html string, basePath string) string {
	for _, attr := range []string{`href="`, `src="`} {
		html = prefixRootHTMLAttribute(html, attr, basePath)
	}
	return html
}

func prefixRootHTMLAttribute(html string, attr string, basePath string) string {
	search := attr + "/"
	var builder strings.Builder
	offset := 0
	for {
		idx := strings.Index(html[offset:], search)
		if idx < 0 {
			builder.WriteString(html[offset:])
			break
		}
		idx += offset
		valueStart := idx + len(attr)
		builder.WriteString(html[offset:valueStart])
		if strings.HasPrefix(html[valueStart:], basePath+"/") || strings.HasPrefix(html[valueStart:], basePath+`"`) {
			builder.WriteString("/")
			offset = valueStart + 1
			continue
		}
		builder.WriteString(basePath)
		builder.WriteString("/")
		offset = valueStart + 1
	}
	return builder.String()
}
