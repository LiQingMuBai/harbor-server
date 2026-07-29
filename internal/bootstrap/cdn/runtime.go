package cdn

import (
	"cointrade/internal/bootstrap/shared"
	"cointrade/utils"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Options struct {
	Port          int
	Domain        string
	StaticDir     string
	PDFDir        string
	WhitepaperDir string
}

func OptionsFromEnv() (Options, error) {
	options := Options{
		Port:          shared.GetenvInt("CDN_PORT", 9999),
		Domain:        strings.TrimSpace(shared.Getenv("CDN_DOMAIN", "")),
		StaticDir:     strings.TrimSpace(shared.Getenv("CDN_STATIC_DIR", "./static")),
		PDFDir:        strings.TrimSpace(shared.Getenv("CDN_PDF_DIR", "./pdf")),
		WhitepaperDir: strings.TrimSpace(shared.Getenv("CDN_WHITEPAPER_DIR", "./whitepaper")),
	}
	if options.Domain == "" {
		return Options{}, errors.New("missing CDN_DOMAIN")
	}
	return options, nil
}

func Run(options Options) error {
	router := gin.Default()
	router.Use(crossDomain)
	router.POST("/", checkSid, upload(options.Domain, options.StaticDir))
	router.OPTIONS("/", func(r *gin.Context) {
		r.JSON(http.StatusOK, nil)
	})
	router.Static("/static", options.StaticDir)
	router.Static("/pdf", options.PDFDir)
	router.Static("/whitepaper", options.WhitepaperDir)
	return router.Run(fmt.Sprintf(":%d", options.Port))
}

func crossDomain(r *gin.Context) {
	r.Header("Access-Control-Allow-Origin", "*")
	r.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	r.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	r.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	r.Header("Access-Control-Allow-Credentials", "true")
	r.Next()
}

func checkSid(r *gin.Context) {
	sid := r.DefaultQuery("sid", "")
	if sid == "" {
		r.Abort()
		return
	}
	r.Next()
}

func upload(domain string, staticDir string) gin.HandlerFunc {
	return func(r *gin.Context) {
		file, err := r.FormFile("file")
		if err != nil {
			r.JSON(http.StatusForbidden, "upload error")
			return
		}
		originalFilename := file.Filename
		parts := strings.Split(originalFilename, ".")
		fileType := strings.ToLower(parts[len(parts)-1])
		allowed := false
		for _, value := range []string{"jpg", "jpeg", "png", "gif", "pdf"} {
			if value == fileType {
				allowed = true
				break
			}
		}
		if !allowed {
			r.JSON(http.StatusForbidden, "file type is not allowed")
			return
		}

		date := time.Now().Format("20060102")
		hour := strconv.Itoa(time.Now().Hour())
		saveDir := filepath.Join(staticDir, date, hour)
		if !utils.FileExists(saveDir) {
			if os.MkdirAll(saveDir, os.ModePerm) != nil {
				r.JSON(http.StatusForbidden, "dir make error")
				return
			}
		}

		filename := utils.RandName() + "." + fileType
		savePath := filepath.Join(saveDir, filename)
		if r.SaveUploadedFile(file, savePath) != nil {
			r.JSON(http.StatusForbidden, "savefile error")
			return
		}

		domainPrefix := strings.TrimRight(domain, "/")
		publicPath := "static/" + date + "/" + hour + "/" + filename
		r.JSON(http.StatusOK, map[string]string{"path": domainPrefix + "/" + publicPath})
	}
}
