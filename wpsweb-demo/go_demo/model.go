package main

type FileModel struct {
	Id          string    `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name        string    `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Version     int32     `protobuf:"varint,3,opt,name=version" json:"version,omitempty"`
	Size        int64     `protobuf:"varint,4,opt,name=size" json:"size,omitempty"`
	Type        string    `protobuf:"bytes,5,opt,name=type" json:"type,omitempty"`
	DownloadUrl string    `protobuf:"bytes,6,opt,name=download_url,json=downloadUrl" json:"download_url,omitempty"`
	CreateTime  int64     `protobuf:"varint,7,opt,name=create_time,json=createTime" json:"create_time,omitempty"`
	ModifyTime  int64     `protobuf:"varint,8,opt,name=modify_time,json=modifyTime" json:"modify_time,omitempty"`
	Creator     string    `protobuf:"bytes,9,opt,name=creator" json:"creator,omitempty"`
	Modifier    string    `protobuf:"bytes,10,opt,name=modifier" json:"modifier,omitempty"`
	UserACL     UserACL   `protobuf:"bytes,11,opt,name=user_acl,json=userAcl" json:"user_acl,omitempty"`
	Watermark   Watermark `protobuf:"bytes,14,opt,name=watermark" json:"watermark,omitempty"`
}

type UserModel struct {
	Id         string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name       string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	AvatarUrl  string `protobuf:"bytes,3,opt,name=avatar_url,json=avatarUrl" json:"avatar_url,omitempty"`
	Permission string `protobuf:"bytes,4,opt,name=permission" json:"permission,omitempty"`
	Avatar     string `protobuf:"bytes,5,opt,name=avatar" json:"avatar,omitempty"`
}

type FileEditModel struct {
	File FileModel `protobuf:"bytes,1,opt,name=file" json:"file,omitempty"`
	User UserModel `protobuf:"bytes,2,opt,name=user" json:"user,omitempty"`
}

type GetUserInfoBatchOutput struct {
	Users []*UserModel `protobuf:"bytes,1,rep,name=users" json:"users,omitempty"`
}

type GetUserInfoBatchInput struct {
	Ids []string `protobuf:"bytes,1,rep,name=ids" json:"ids,omitempty"`
}

type PostFileOutput struct {
	File *FileModel `protobuf:"bytes,1,opt,name=file" json:"file,omitempty"`
}

type FileMetadata struct {
	Id          string     `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name        string     `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Version     int32      `protobuf:"varint,3,opt,name=version" json:"version,omitempty"`
	Size        int64      `protobuf:"varint,4,opt,name=size" json:"size,omitempty"`
	Type        string     `protobuf:"bytes,5,opt,name=type" json:"type,omitempty"`
	DownloadUrl string     `protobuf:"bytes,6,opt,name=download_url,json=downloadUrl" json:"download_url,omitempty"`
	CreateTime  int64      `protobuf:"varint,7,opt,name=create_time,json=createTime" json:"create_time,omitempty"`
	ModifyTime  int64      `protobuf:"varint,8,opt,name=modify_time,json=modifyTime" json:"modify_time,omitempty"`
	Creator     *UserModel `protobuf:"bytes,9,opt,name=creator" json:"creator,omitempty"`
	Modifier    *UserModel `protobuf:"bytes,10,opt,name=modifier" json:"modifier,omitempty"`
	UniqueId    string     `protobuf:"bytes,11,opt,name=unique_id,json=uniqueId" json:"unique_id,omitempty"`
	UserAcl     *UserACL   `protobuf:"bytes,12,opt,name=user_acl,json=userAcl" json:"user_acl,omitempty"`
	VerType     string     `protobuf:"bytes,13,opt,name=ver_type,json=verType" json:"ver_type,omitempty"`
	Watermark   *Watermark `protobuf:"bytes,14,opt,name=watermark" json:"watermark,omitempty"`
}

type UserACL struct {
	Read     int32 `protobuf:"varint,1,opt,name=read" json:"read,omitempty"`
	Update   int32 `protobuf:"varint,2,opt,name=update" json:"update,omitempty"`
	Download int32 `protobuf:"varint,3,opt,name=download" json:"download,omitempty"`
	Share    int32 `protobuf:"varint,4,opt,name=share" json:"share,omitempty"`
	Rename   int32 `protobuf:"varint,5,opt,name=rename" json:"rename,omitempty"`
	History  int32 `protobuf:"varint,6,opt,name=history" json:"history,omitempty"`
}

type Watermark struct {
	Type       int32   `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	Value      string  `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	Fillstyle  string  `protobuf:"bytes,3,opt,name=fillstyle" json:"fillstyle,omitempty"`
	Font       string  `protobuf:"bytes,4,opt,name=font" json:"font,omitempty"`
	Rotate     float32 `protobuf:"fixed32,5,opt,name=rotate" json:"rotate,omitempty"`
	Horizontal int32   `protobuf:"varint,6,opt,name=horizontal" json:"horizontal,omitempty"`
	Vertical   int32   `protobuf:"varint,7,opt,name=vertical" json:"vertical,omitempty"`
}

type GetFileHistoryVersionsResponse struct {
	Histories []*FileMetadata `protobuf:"bytes,1,rep,name=histories" json:"histories,omitempty"`
}

type GetFileHistoryVersionsRequest struct {
	Id     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Offset int32  `protobuf:"varint,2,opt,name=offset" json:"offset,omitempty"`
	Count  int32  `protobuf:"varint,3,opt,name=count" json:"count,omitempty"`
}


type GetFileVersionOutput struct {
	File *FileModel `protobuf:"bytes,1,opt,name=file" json:"file,omitempty"`
}

type PutFileInput struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

type GetTemplateInfo struct {
	RedirectUrl string `json:"redirect_url"`
	UserId      string `json:"user_id"`
}
