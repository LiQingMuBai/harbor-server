package cdn

import (
	"cointrade/internal/bootstrap/shared"
	"cointrade/utils"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Options struct {
	Port   int
	Domain string
}

func OptionsFromEnv() (Options, error) {
	options := Options{
		Port:   shared.GetenvInt("CDN_PORT", 9999),
		Domain: strings.TrimSpace(shared.Getenv("CDN_DOMAIN", "")),
	}
	if options.Domain == "" {
		return Options{}, errors.New("missing CDN_DOMAIN")
	}
	return options, nil
}

func Run(options Options) error {
	router := gin.Default()
	router.Use(crossDomain)
	router.POST("/", checkSid, upload(options.Domain))
	router.OPTIONS("/", func(r *gin.Context) {
		r.JSON(http.StatusOK, nil)
	})
	router.Static("/static", "./static")
	router.Static("/pdf", "./pdf")
	router.Static("/whitepaper", "./whitepaper")
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

func upload(domain string) gin.HandlerFunc {
	return func(r *gin.Context) {
		file, err := r.FormFile("file")
		if err != nil {
			r.JSON(http.StatusForbidden, "upload error")
			return
		}
		filename := file.Filename
		parts := strings.Split(filename, ".")
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
		saveDir := "static/" + time.Now().Format("20060102") + "/" + strconv.Itoa(time.Now().Hour()) + "/"
		if !utils.FileExists(saveDir) {
			if os.MkdirAll(saveDir, os.ModePerm) != nil {
				r.JSON(http.StatusForbidden, "dir make error")
				return
			}
		}
		savePath := saveDir + utils.RandName() + "." + fileType
		if r.SaveUploadedFile(file, savePath) != nil {
			r.JSON(http.StatusForbidden, "savefile error")
			return
		}
		r.JSON(http.StatusOK, map[string]string{"path": domain + savePath})
	}
}
