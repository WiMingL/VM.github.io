package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

//token存活时间
const TokenExpiresTime = 600

func getListFileHandler(c *gin.Context) {
	s, err := GetAllFile(App.LocalDir, nil)
	if err != nil {
		fmt.Println("get list file err :", err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"name_files": s,
	})
	return
}

//返回路径下,以(docx|pptx|xlsx)结尾的文件名
func GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		var ref = regexp.MustCompile(".(docx|pptx|xlsx|pdf)$")
		if ref.MatchString(fi.Name()) {
			s = append(s, fi.Name())
		}
	}
	return s, nil
}

func getUrlAndTokenHandler(c *gin.Context) {
	_w_filepath := c.Query("_w_filepath")
	// 例: _w_filepath = testpath/1.docx
	// or _w_filepath = 1.docx
	if _w_filepath == "" {
		fmt.Println("err : missing parameter _w_filepath")
		c.Status(http.StatusBadRequest)
		return
	}

	webofficeUrl := getWpsUrl(_w_filepath)

	token.Lock()
	defer token.Unlock()
	timeout := time.Now().Add(TokenExpiresTime * time.Second).Unix()

	if token.key != "" && token.timeout-time.Now().Unix() >= 0 {
		token.timeout = timeout
		fmt.Println("get demo token :", token.key)
		c.JSON(http.StatusOK, gin.H{
			"token":      token.key,
			"expires_in": TokenExpiresTime,
			"wpsUrl":     webofficeUrl,
		})
		return
	}

	uuid, _ := uuid.NewV4()
	newtoken := NoHyphenString(uuid)

	token.key = newtoken
	token.timeout = timeout
	fmt.Println("get demo token :", token.key)
	c.JSON(http.StatusOK, gin.H{
		"token":      newtoken,
		"expires_in": TokenExpiresTime,
		"wpsUrl":     webofficeUrl,
	})
	return
}

func getWpsUrl(_w_filepath string) string {
	path_arr := strings.Split(_w_filepath, "/")
	var fname string
	if len(path_arr) > 1 {
		fname = path_arr[len(path_arr)-1]
	} else {
		fname = _w_filepath
	}

	//获取file_id
	file_id := fmt.Sprintf("%x", md5.Sum([]byte(fname)))
	if len(file_id) > 20 {
		file_id = file_id[0:20]
	}
	//t:文件类型
	arr := strings.Split(fname, ".")
	t := editorExtMap[arr[len(arr)-1]]

	//默认参数,需根据需求修改
	var values = make(url.Values)
	values.Add("_w_appid", App.Appid)
	values.Add("_w_tokentype", "1")
	values.Add("_w_filepath", _w_filepath)

	signature := Sign(values, App.Appsecret)
	webofficeUrl := fmt.Sprintf("%s/office/%s/%s?%s&_w_signature=%s", App.Domain, t, file_id, values.Encode(), url.QueryEscape(signature))

	return webofficeUrl
}

//下载文件
func getFileHanlder(c *gin.Context) {
	_w_filepath := c.Query("_w_filepath")
	if _w_filepath == "" {
		fmt.Println("error : _w_filepath is nil")
		c.Status(http.StatusBadRequest)
		return
	}

	path_arr := strings.Split(_w_filepath, "/")
	var fname string
	if len(path_arr) > 1 {
		fname = path_arr[len(path_arr)-1]
	} else {
		fname = _w_filepath
	}

	version := c.Query("version")
	var fpath string
	if version != "" {
		fpath = filepath.Join(App.LocalDir, version, fname)
	} else {
		fpath = filepath.Join(App.LocalDir, fname)
	}

	if !pathExist(fpath) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    40004,
			"message": "fileNotExists",
			"detail":  fmt.Sprintf("%s is not exists", fname),
		})
		return
	}
	c.File(fpath)
}

//上传文件
func fileUploadHandler(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println("get form file faild: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	_w_filepath := c.Query("_w_filepath")
	if _w_filepath == "" {
		fmt.Println("error : _w_filepath is nil")
		c.Status(http.StatusBadRequest)
		return
	}

	fpath := filepath.Join(App.LocalDir, _w_filepath)
	if pathExist(fpath) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "file already exists",
		})
		return
	}
	out, err := os.OpenFile(fpath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("create file faild: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	defer out.Close()
	// Copy数据
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println("copy file faild: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	getUrlHandler(c)
}

func getUrlHandler(c *gin.Context) {
	_w_filepath := c.Query("_w_filepath")
	if _w_filepath == "" {
		fmt.Println("error : _w_filepath is nil")
		c.Status(http.StatusBadRequest)
		return
	}

	path_arr := strings.Split(_w_filepath, "/")
	var fname string
	if len(path_arr) > 1 {
		fname = path_arr[len(path_arr)-1]
	} else {
		fname = _w_filepath
	}

	file_id := fmt.Sprintf("%x", md5.Sum([]byte(fname)))
	if len(file_id) > 20 {
		file_id = file_id[0:20]
	}
	arr := strings.Split(fname, ".")
	t := editorExtMap[arr[len(arr)-1]]
	query := c.Request.URL.Query()
	query.Add("_w_tokentype", "1")
	//剔除非wps的参数
	for k, _ := range query {
		if !re.MatchString(k) {
			query.Del(k)
		}
	}
	signature := Sign(query, App.Appsecret)
	webofficeUrl := fmt.Sprintf("%s/office/%s/%s?%s&_w_signature=%s", App.Domain, t, file_id, c.Request.URL.RawQuery, url.QueryEscape(signature))
	c.String(http.StatusOK, webofficeUrl)
}
