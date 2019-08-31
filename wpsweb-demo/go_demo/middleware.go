package main

import (
	//"fmt"
	//"net/http"
	//"net/url"

	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"time"
)

//检查token
func CheckToken(c *gin.Context) {
	token.RLock()
	defer token.RUnlock()
	access_token := c.Request.Header.Get("x-wps-weboffice-token")
	if access_token != token.key {
		fmt.Println("invalid access_token: ", access_token)
		fmt.Println("demo token: ", token.key)
		AbortWithErrorMessage(c, http.StatusBadRequest, "InvalidToken", "invalid token")
		return
	}

	if token.timeout-time.Now().Unix() < 0 {
		fmt.Println(" error: Token Time Out")
		AbortWithErrorMessage(c, http.StatusBadRequest, "TokenTimeOut", "token is timeout")
		return
	}
}

//检查签名
func CheckOpenSignature(c *gin.Context) {
	query := c.Request.URL.Query()
	query.Del("access_token")
	signature, err := url.PathUnescape(query.Get("_w_signature"))
	if err != nil {
		AbortWithErrorMessage(c, http.StatusUnauthorized, "InvalidRequest", "invaid signature")
		return
	}

	query.Del("_w_signature")
	authorization := Sign(query, App.Appsecret)
	if authorization != signature {
		fmt.Printf("error : authorization:%s and signature:%s \n", authorization, signature)
		AbortWithErrorMessage(c, http.StatusUnauthorized, "InvalidSignature", "signature mismatch")
		return
	}
}
