package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const PdfPrefix = `JVBER`

func main() {
	r := gin.New()
	r.POST("/decode", decodeFileHandler)
	r.Run()
}

type Document struct {
	base64data string `json:"base64data"`
	name       string `json:"name"`
}

func decodeFileHandler(ctx *gin.Context) {
	var document Document
	var b64 string

	if err := ctx.ShouldBindJSON(&document); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(document.base64data) != 0 {
		b64 = document.base64data
	} else {
		b64 = readBase64FromTextFile()
	}

	fileName := decodeAndWriteFile(b64, document.name)
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Type", "application/text/plain")
	ctx.Header("Accept-Length", fmt.Sprintf("%d", len(fileBytes)))
	ctx.Writer.Write([]byte(fileBytes))
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Download file successfully",
	})
}

func decodeAndWriteFile(b64 string, name string) string {
	docType := `png`

	if strings.HasPrefix(b64, PdfPrefix) {
		docType = `pdf`
	}
	fileName := "output." + docType

	dec, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}

	if len(name) != 0 {
		fileName = name
	}
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	return fileName
}

func readBase64FromTextFile() string {
	content, err := os.ReadFile("input_image.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(string(content)))
	return string(content)
}
