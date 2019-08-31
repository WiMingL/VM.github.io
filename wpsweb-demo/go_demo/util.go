package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type Any interface{}

// JSONError golint
type JSONError struct {
	StatusCode int         `json:"-"`
	Result     string      `json:"result"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
}

// Error golint
func (err *JSONError) Error() string {
	return err.Result
}

// AbortWithErrorMessage golint
func AbortWithErrorMessage(c *gin.Context, status int, result string, message string) {
	err := &JSONError{
		StatusCode: status,
		Result:     result,
		Message:    message,
	}
	AbortWithError(c, err)
}

// AbortWithError golint
func AbortWithError(c *gin.Context, err error) {
	jsonErr, ok := err.(*JSONError)
	if !ok {
		jsonErr = &JSONError{
			StatusCode: http.StatusInternalServerError,
			Result:     "Unknown",
			Message:    err.Error(),
		}
	}
	c.JSON(jsonErr.StatusCode, jsonErr)
	c.Abort()
}

func pathExist(_path string) bool {
	indexMatches, err := filepath.Glob(_path)
	return err == nil && len(indexMatches) > 0
}

func Sign(values url.Values, secretKey string) string {
	contents := StringToSign(values, secretKey)
	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(contents))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	fmt.Println("signature:",signature)
	return signature
}

func StringToSign(values url.Values, secretKey string) []byte {
	keys := []string{}
	for k, _ := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := &bytes.Buffer{}
	for _, k := range keys {
		fmt.Fprintf(buf, "%s=%s", k, values.Get(k))
	}
	fmt.Fprintf(buf, "_w_secretkey=%s", secretKey)
	fmt.Println("StringToSign:",buf.String())
	return buf.Bytes()
}


func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func GetLatestVersion(ctx context.Context, fname string) int32 {
	var max int
	files, _ := ioutil.ReadDir(App.LocalDir)
	for _, file := range files {
		if file.IsDir() {
			version, _ := strconv.Atoi(file.Name())
			max = Max(version, max)
		}
	}
	for i := 1; i <= max; i++ {
		fpath := filepath.Join(App.LocalDir, strconv.Itoa(i), fname)
		if !pathExist(fpath) {
			return int32(i)
		}
	}
	max += 1
	return int32(max)
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func NoHyphenString(u uuid.UUID) string {
	return fmt.Sprintf("%x%x%x%x%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}