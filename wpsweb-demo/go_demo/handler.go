package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//虚拟参数
const (
	filesize   = 1024 * 1024
	creator    = "0"
	createTime = 1136185445
)

var re = regexp.MustCompile("^_w_")

//支持的文件格式
var etExts = []string{"et", "xls", "xlt", "xlsx", "xlsm", "xltx", "xltm", "csv"}
var wpsExts = []string{"doc", "docx", "txt", "dot", "wps", "wpt", "dotx", "docm", "dotm"}
var wppExts = []string{"ppt", "pptx", "pptm", "pptm", "ppsm", "pps", "potx", "potm", "dpt", "dps"}
var pdfExts = []string{"pdf"}
var editorExtMap = map[string]string{}

func init() {
	for _, ext := range etExts {
		editorExtMap[ext] = "s"
	}
	for _, ext := range wpsExts {
		editorExtMap[ext] = "w"
	}
	for _, ext := range wppExts {
		editorExtMap[ext] = "p"
	}
	for _, ext := range pdfExts {
		editorExtMap[ext] = "f"
	}
}

//文件重命名
func putFileName(c *gin.Context) {
	var args PutFileInput
	c.BindJSON(&args)
	fname := args.Name
	fmt.Println(fname)
	c.Status(200)
}

//获取文件信息
func fileHandler(c *gin.Context) {

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

	//第三方对接模块按需求,自行控制
	permission := "write"
	uid := "user1"

	fpath := filepath.Join(App.LocalDir, _w_filepath)
	if !pathExist(fpath) {
		fmt.Println("file not exist: %v", fpath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    40004,
			"message": "fileNotExists",
			"detail":  fmt.Sprintf("%s:fileNotExists", _w_filepath),
		})
		return
	}
	file_id := fmt.Sprintf("%x", md5.Sum([]byte(fname)))
	if len(file_id) > 20 {
		file_id = file_id[0:20]
	}
	version := GetLatestVersion(c.Request.Context(), fname)
	userACL := &UserACL{
		Rename:  1,
		History: 0,
	}

	watermark := &Watermark{}

	file := &FileModel{
		Id:          file_id,
		Name:        fname,
		Size:        filesize,
		DownloadUrl: fmt.Sprintf("%s/v1/file?_w_fname=%s", App.DownloadHost, fname),
		Creator:     creator,
		CreateTime:  createTime,
		Modifier:    uid,
		ModifyTime:  time.Now().Unix(),
		Version:     version,
		UserACL:     *userACL,
		Watermark:   *watermark,
	}

	if version == 0 {
		file.Version = 1
	}

	user := &UserModel{
		Id:         uid,
		Permission: permission,
		AvatarUrl:  "https://picsum.photos/100/100/?image=" + uid,
		Name:       "wps_user-" + uid,
	}
	fem := &FileEditModel{
		File: *file,
		User: *user,
	}

	c.JSON(http.StatusOK, fem)
}

//获取特定版本的文件信息
func fileVersion(c *gin.Context) {
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

	uid := c.Query("_w_userid")
	version, _ := strconv.Atoi(c.Param("version"))
	file_id := fmt.Sprintf("%x", md5.Sum([]byte(fname)))
	if len(file_id) > 20 {
		file_id = file_id[0:20]
	}
	file := &FileModel{
		Id:          file_id,
		Name:        fname,
		Size:        filesize,
		DownloadUrl: fmt.Sprintf("%s/v1/file?_w_filepath=%s&version=%d", App.DownloadHost, _w_filepath, version),
		Creator:     "0",
		CreateTime:  createTime,
		Modifier:    uid,
		ModifyTime:  time.Now().Unix(),
		Version:     int32(version),
	}
	out := &GetFileVersionOutput{
		File: file,
	}
	c.JSON(http.StatusOK, out)
}

//上传最新版文件到demo，并且储存上一个版本的文件作为历史版本
//每一个版本会生成一个文件夹,从1开始递增,开发者可以根据自己需求去修改功能控制
func postFile(c *gin.Context) {
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

	uid := "user1"

	file_id := fmt.Sprintf("%x", md5.Sum([]byte(fname)))
	if len(file_id) > 20 {
		file_id = file_id[0:20]
	}
	version := GetLatestVersion(c.Request.Context(), fname)

	//迁移当前文档到历史版本
	pathTmp := filepath.Join(App.LocalDir, strconv.Itoa(int(version)))
	if !pathExist(pathTmp) {
		err := os.Mkdir(pathTmp, 0777)
		if err != nil {
			fmt.Println("get file faild: ", err.Error())
			c.Status(http.StatusBadRequest)
			return
		}
	}
	historypath := filepath.Join(pathTmp, fname)
	fpath := filepath.Join(App.LocalDir, fname)
	CopyFile(historypath, fpath)

	//更新文档
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println("get form file faild: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	if !pathExist(fpath) {
		fmt.Println("file not exist: %v", fpath)
		c.Status(http.StatusBadRequest)
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
	version += 1
	fm := &FileModel{
		Id:          file_id,
		Name:        fname,
		Size:        filesize,
		DownloadUrl: App.DownloadHost + "/v1/file?_w_fname=" + fname,
		Creator:     creator,
		CreateTime:  createTime,
		Modifier:    uid,
		ModifyTime:  time.Now().Unix(),
		Version:     version,
	}
	c.JSON(http.StatusOK, gin.H{
		"file": fm,
	})
}

//批量获取用户信息
func getUserBatch(c *gin.Context) {
	var users []*UserModel
	var in GetUserInfoBatchInput
	err := c.BindJSON(&in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, id := range in.Ids {
		var User UserModel
		User.Id = id
		User.Permission = "write"
		User.AvatarUrl = "https://picsum.photos/100/100/?image=" + id
		uid, _ := strconv.Atoi(id)
		User.Name = fmt.Sprintf("wps_user-%d", uid)
		users = append(users, &User)
	}
	out := GetUserInfoBatchOutput{
		Users: users,
	}
	c.JSON(http.StatusOK, out)
}

//上传当前正在编辑的用户信息
func PostFileOnline(c *gin.Context) {
	var in GetUserInfoBatchInput
	var users []*UserModel
	err := c.BindJSON(&in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, id := range in.Ids {
		user := &UserModel{
			Id:         id,
			Permission: "write",
			AvatarUrl:  "https://picsum.photos/100/100/?image=" + id,
			Name:       "wps_user-" + id,
		}
		users = append(users, user)
	}
	c.String(http.StatusOK, "success")
}

//查看所有历史版本
func GetFileHistoryVersions(c *gin.Context) {
	var in GetFileHistoryVersionsRequest
	err := c.BindJSON(&in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

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

	uid := "user1"

	file_id := fmt.Sprintf("%x", md5.Sum([]byte(fname)))
	if len(file_id) > 10 {
		file_id = file_id[0:10]
	}
	version := GetLatestVersion(c.Request.Context(), fname)
	if err != nil {
		fmt.Println("get latest version faild: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}

	history := []*FileMetadata{}
	total := in.Count
	if total > version {
		total = version
	}
	for i := in.Offset; i < total; i++ {
		version_ := int32(total - i)
		md := &FileMetadata{
			Id:          file_id,
			Name:        fname,
			Size:        filesize,
			DownloadUrl: fmt.Sprintf("%s/v1/file?_w_fname=%s&version=%d", App.DownloadHost, fname, version_),
			Version:     version_,
			Type:        "file",
			CreateTime:  createTime,
			ModifyTime:  time.Now().Unix(),
		}
		md.Creator = &UserModel{
			Id:        creator,
			Name:      "wps_user-" + creator,
			AvatarUrl: "https://picsum.photos/100/100/?image=" + creator,
		}
		md.Modifier = &UserModel{
			Id:        uid,
			Name:      "wps_user-" + uid,
			AvatarUrl: "https://picsum.photos/100/100/?image=" + uid,
		}
		history = append(history, md)
	}
	out := GetFileHistoryVersionsResponse{
		Histories: history,
	}
	c.JSON(http.StatusOK, out)
}

func postNewFile(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println("get form file faild: ", err.Error())
		c.Status(http.StatusBadRequest)
		return
	}
	fname := c.Request.Form.Get("name")
	fpath := filepath.Join(App.LocalDir, fname)
	uid := c.Query("_w_userid")
	appid := c.Query("_w_appid")
	out, err := os.OpenFile(fpath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("create file faild: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println("copy file faild: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	arr := strings.Split(fname, ".")
	t := editorExtMap[arr[len(arr)-1]]
	file_id := fmt.Sprintf("%x", md5.Sum([]byte(fname)))
	query := fmt.Sprintf("_w_fname=%s&_w_userid=%s&_w_appid=%s&_w_permission=write", fname, uid, appid)
	urlquery, _ := url.ParseQuery(query)
	signature := Sign(urlquery, App.Appsecret)
	redirectUrl := fmt.Sprintf("%s/office/%s/%s?%s&_w_signature=%s", App.Domain, t, file_id, query, url.QueryEscape(signature))
	getTemplateInfo := GetTemplateInfo{
		RedirectUrl: redirectUrl,
		UserId:      uid,
	}
	c.JSON(http.StatusOK, getTemplateInfo)
}
