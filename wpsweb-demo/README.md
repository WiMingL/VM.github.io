# WebOffice开放平台接口文档

## 一.简述
   
   WebOffice开放平台，是为了让第三方企业能够接入wps的在线编辑服务，让用户享受无需安装客户端，通过分享在线文档进行跨平台的多人协作的便捷服务。平台还支持同时接入多家第三方企业。

**企业在接入第三方平台开发时，申请和上线流程如下：**

***1.申请Appid和SecretKey***
> 需要第三方企业提供对接模块的外网访问域名。

***2.实现对接模***
>根据下文提供的一系列接口，以及企业自身的业务需求，实现一个对接模块，供wps侧weboffice开放平台服务端调用。

***3.申请全网发布并上线***
>对接模块开发完成之后，上线并发布到外网，可通过先前提供的域名进行访问。

***4.生成带签名的访问url***
>根据wps侧提供的appid和secretkey和对接模块需要透传的参数，按要求生成签名，然后生成一个可以在线编辑文档的url，访问以使用wps提供的在线编辑服务。

***5.文件格式生成的url***

表格文件url:
>https://wwo.wps.cn/office/s/:file_id?...

文字文件url:
>https://wwo.wps.cn/office/w/:file_id?...

演示文件url:
>https://wwo.wps.cn/office/p/:file_id?...

PDF文件url:
>https://wwo.wps.cn/office/f/:file_id?...



***

## 二.支持格式

| **文件预览支持格式** | |
| :---: | :---: |
| 表格文件	| xls, xlt, et, xlsx, xltx, csv, xlsm, xltm |
| 文字文件 | doc, dot, wps, wpt, docx, dotx, docm, dotm |  
| 演示文件 |	ppt,pptx,pptm,ppsx,ppsm,pps,potx,potm,dpt,dps |
| PDF文件 |	pdf |


| **文件编辑支持格式** |  |
| :---: | :---: |
| 表格文件 |	xls, xlt, et, xlsx, xltx, csv, xlsm, xltm |
| 文字文件 |	doc, dotx，txt, dot, wps, wpt, docx, docm, dotm | 
| 演示文件 |	ppt,pptx,pptm,ppsx,ppsm,pps,potx,potm,dpt,dps |
| PDF文件 |	pdf |
***
## 三.对接流程
### 3.1.接口调用流程图

![接口调用流程图](https://raw.githubusercontent.com/zhongxuan123/weboffice_openplatform/master/images/3_1_p1.png)


***


### 3.2.API列表

| **API地址** | **方法** | **描述** | **预览是否实现** | **编辑是否实现** |
| --- | --- |  --- |  --- |  --- |
| /v1/3rd/file/info | 	GET	 | 获取文件元数据 | 	是	 | 是 |
| /v1/3rd/user/info | 	POST	 | 获取用户信息 | 	否	 | 是 |
| /v1/3rd/file/save | 	POST	 | 上传文件新版本 | 	否	 | 是 |
| /v1/3rd/file/online | 	POST	 | 上传在线编辑用户信息 | 	是	 | 是 |
| /v1/3rd/file/version/:version | 	 GET |	获取特定版本的文件信息 | 	否	 | 是 |
| /v1/3rd/file/rename | 	PUT	 | 文件重命名 | 	否	 | 是 |
| /v1/3rd/file/history | 	POST	 | 获取所有历史版本文件信息 | 	否	 | 是 |
| /v1/3rd/file/new |	POST | 	新建文件 |  	否	 | 是 |

>请开发者关注：所有对接模块相关的参数都要以”\_w_”作为前缀，否则容易导致签名不能通过验证。 


***


#### 3.2.1.关于token
token是用来校验用户身份的凭证，它是用来限制文档只能由有权限的用户访问，以提高文档安全性。Web office开放平台会在请求所有接口时，请求头带上此token。

目前获取第三方token方式,是第三方直接提供获取token的接口。具体接入方式由第三方的授权接入文档为准。

提供参考整体交互流程图如下：

![token交互流程图](https://raw.githubusercontent.com/zhongxuan123/weboffice_openplatform/master/images/token%E4%BA%A4%E4%BA%92%E6%B5%81%E7%A8%8B%E5%9B%BE.png)



#### Demo中的交互:
>接口:/getUrlAndToken
>
>随机产生一个token保存在内存,10分钟内没有使用则会过期,也会返回生存时间给到前端APP.并且会根据传入的文件名_w_fname生成可以访问的url到前端APP.



Demo的token交互流程图:

![demo交互流程图](https://raw.githubusercontent.com/zhongxuan123/weboffice_openplatform/master/images/demo%E4%BA%A4%E4%BA%92%E6%B5%81%E7%A8%8B%E5%9B%BE.png)

***



#### 3.2.2.获取文件元数据
描述：
>获取文件元数据，包括当前文件信息和当前用户信息。其中，当user的permission参数返回“write”时，进入在线编辑模式，返回“read”时进入预览模式。

API地址：
>/v1/3rd/file/info

调用方法：
>GET

请求头（由weboffice开放平台写入）：

| **Header**	| **值**	| **描述** |
| :---: | :---: |  --- |
| User-Agent | wps-weboffice-openplatform | 	用户代理 |
| x-wps-weboffice-token | 	xxxxxxx | 	校验身份的token，值根据对接的企业而定 |
| x-weboffice-file-id | 	xxxxxxx	 | 文件id |

请求参数：

| **参数名**	| **参数类型**	| **数据类型**	| **描述**	| **可选** |
| :---: | :---: |  :---: | --- | :---: |
| _w_signature | 	query | 	string | 	请求签名 | 	否 |
| _w_appid | 	query | 	string	| 应用id  	 | 否 |
| *（任意参数） | 	query | 	*（任意类型） | 	按照需求传入参数  | 	是 |
示例：
>http://www.xxx.cn/v1/3rd/file/info?_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&_w_userid=1000&_w_fname=example.doc&_w_signature=tBnhFqvVrC1LcJKye1m7GjzvqoA%3D
>
>此处的_w_fname等等是对接模块处理的参数，可以对接模块自己定义，weboffice服务不会对此参数进行处理,http://www.xxx.cn 为企业对应的域名和nginx转发地址。
>
返回值：
>
>status code为200表示获取数据成功，返回下面的信息。其余值表示失败，需要在返回中指定code,msg等信息。
>
正确返回响应信息示例:
>
>
>(其中user_acl和watermark为非必须返回参数，user_acl用于控制用户权限，不返回则为系统默认值，watermark只用于预览时添加第三方水印）
````
{
    file: {
        id:  "132aa30a87064",	 			//文件id,字符串长度小于40
        name : "example.doc",				//文件名
        version: 1,				        //当前版本号，位数小于11
        size: 200,				        //文件大小，单位为kb	
        creator: "id0",			        	//创建者id，字符串长度小于40
        create_time: 1136185445,			//创建时间，时间戳，单位为秒
        modifier: "id1000",			    	//修改者id，字符串长度小于40
        modify_time: 1551409818,			//修改时间，时间戳，单位为秒
        download_url: "http://www.xxx.cn/v1/file?fid=f132aa30a87064",  //文档下载地址
        user_acl: {
            rename: 1,    //重命名权限，1为打开该权限，0为关闭该权限，默认为0
            history: 1    //历史版本权限，1为打开该权限，0为关闭该权限,默认为1 
        },
        watermark:  {
            type: 1,		 		//水印类型， 0为无水印； 1为文字水印
            value: "禁止传阅",  		          //文字水印的文字，当type为1时此字段必选
            fillstyle: "rgba( 192, 192, 192, 0.6 )",     //水印的透明度，非必选，有默认值
            font: "bold 20px Serif",		//水印的字体，非必选，有默认值
            rotate: -0.7853982,			//水印的旋转度，非必选，有默认值
            horizontal: 50,			 //水印水平间距，非必选，有默认值
            vertical: 100			 //水印垂直间距，非必选，有默认值
        }
    },	
    user: {                              
        id: "id1000",				//用户id，长度小于40
        name: "wps-1000",			//用户名称
        permission: "write",			//用户操作权限，write：可编辑，read：预览
        avatar_url: "http://xxx.cn/id=1000"	//用户头像地址
    }
}
````

***


#### 3.2.3.获取用户信息
描述：

>批量获取当前正在编辑和编辑过文档的用户信息，以数组的形式返回响应。

API地址：

>/v1/3rd/user/info

调用方法：

>POST

请求头（由weboffice开放平台写入）：


| **Header**	| **值**	| **描述** |
| :---: | :---: |  --- |
| User-Agent | 	wps-weboffice-openplatform | 	用户代理 |
| x-wps-weboffice-token | 	xxxxxxx	 | 校验身份的token，值根据对接的企业而定 |
| x-weboffice-file-id | 	xxxxxxx	 | 文件id |

请求参数示例：


| **参数名**	| **参数类型**	| **数据类型**	| **描述**	| **可选** |
| :---: | :---: | :---:| --- | :---: |
| _w_signature	| 	query	| 	string	| 	请求签名		| 否	| 
| _w_appid		| query		| string	| 	应用id		| 否	| 
| *（任意参数）	| 	query	| 	*（任意类型）	| 	按照需求传入参数	| 	是
| ids		| body	| 	string[]	| 	用户id数组	| 	否	| 

示例：

>http://www.xxx.cn/v1/3rd/user/info?_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&_w_userid=1000&_w_fname=example.doc&_w_signature=tBnhFqvVrC1LcJKye1m7GjzvqoA%3D
>
>此处的_w_fname等等是对接模块处理的参数，可以对接模块自己定义，weboffice服务不会对此参数进行处理,http://www.xxx.cn 为企业对应的域名和nginx转发地址。

POST请求示例：
````
{
	ids:["id1000", "id2000"]
	
}
````

返回值：

>status code为200表示获取数据成功，返回下面的信息。其余值表示失败，需要在返回中指定code,msg等信息。

正确返回响应信息示例：
````
{
   "users":[
                {   
                    id: "id1000", 			    //用户ID，字符串长度小于40
                    name: "wps-1000",		    	    //用户名	
                    avatar_url: "http://xxx.cn/?id=1000"    //用户头像
                },
                {
                    id: "id2000", 		 	    //用户ID，字符串长度小于40
                    name: "wps-2000",		  	    //用户名	
                    avatar_url: "http://xxx.cn/?id=2000"    //用户头像
                }
   ]
}
````

***


#### 3.2.4.上传在线编辑用户信息
描述：
>当有用户加入或者退出协作的时候 ，上传当前文档协作者的用户信息，可以用作上下线通知。此接口可以根据需要实现，若不实现直接返回http响应码200。

API地址: 
>/v1/3rd/file/online

调用方法：
>POST

请求头（由weboffice开放平台写入）：

| **Header**	| **值**	| **描述** |
| :---: | :---: | --- |
| User-Agent	| 	wps-weboffice-openplatform		| 用户代理 |
| x-wps-weboffice-token		| xxxxxxx	| 	校验身份的token，值根据对接的企业而定 |
| x-weboffice-file-id	| 	xxxxxxx		| 文件id |

请求参数示例：


| **参数名**	| **参数类型**	| **数据类型**	| **描述**	| **可选**	|
| :---: | :---: | :---: | --- | :---: |
| _w_signature	 | query | 	string | 	请求签名 | 	否 | 
| _w_appid	 | query | 	string | 	应用id | 	否 | 
| *（任意参数） | 	query	 | *（任意类型） | 	按照需求传入参数	 | 是 | 
| ids	 | body	 | string[]	 | 用户id数组 | 	否 | 

示例：

>http://www.xxx.cn/v1/3rd/file/online?_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&_w_userid=1000&_w_fname=example.doc&_w_signature=tBnhFqvVrC1LcJKye1m7GjzvqoA%3D
>
>此处的_w_fname等等是对接模块处理的参数，可以对接模块自己定义，weboffice服务不会对此参数进行处理,http://www.xxx.cn 为企业对应的域名和nginx转发地址。

POST请求示例：
````
{
	ids:["id1000", "id2000"]  		 //当前协作用户id	
}
````
                                                                                          
返回值：
>status code为200表示获取数据成功。其余值表示失败，需要在返回中指定code,msg等信息。


***


#### 3.2.5.上传文件新版本
描述：
>当文档在线编辑并保存之后，上传该文档最新版本到对接模块，同时版本号进行相应的变更，需要对接方内部实现文件多版本管理机制。

API地址: 
>/v1/3rd/file/save

调用方法：
>POST

请求头（由weboffice开放平台写入）：

| **Header**	| **值**	| **描述** |
| :---: | :---: | --- |
| User-Agent | 	wps-weboffice-openplatform | 	用户代理 | 
| x-wps-weboffice-token	 | xxxxxxx | 	校验身份的token，值根据对接的企业而定 | 
| x-weboffice-file-id	 | xxxxxxx | 	文件id | 

请求参数示例：

| **参数名**	| **参数类型**	| **数据类型**	| **描述**	| **可选**	|
| :---: | :---: | :---: | --- | :---: |
| _w_signature | 	query | 	string	 | 请求签名 | 	否 | 
| _w_appid	 | query	 | string	 | 应用id | 	否 | 
| *（任意参数） | 	query | 	*（任意类型） | 	按照需求传入参数 | 	是 | 
| file	 | body	 | file | 	新版本的文件 | 	否 | 

示例：

>http://www.xxx.cn/office/v1/3rd/file/save?_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&_w_userid=1000&_w_fname=example.doc&_w_signature=tBnhFqvVrC1LcJKye1m7GjzvqoA%3D
>
>此处的_w_fname等等是对接模块处理的参数，可以对接模块自己定义，weboffice服务不会对此参数进行处理,http://www.xxx.cn 为企业对应的域名和nginx转发地址。

返回值：

>status code为200表示获取数据成功，返回下面的信息。其余值表示失败，需要在返回中指定code,msg等信息。


响应信息示例：
````
{
    file: {
            id: "f132aa30a87064",              		//文件id，字符串长度小于40
            name: "example.doc",		        //文件名
            version: 2,					//当前版本号，位数小于11
            size: 200,					//文件大小，单位时kb
            download_url: "http://www.xxx.cn/v1/file?fid=f132aa30a87064" //文件下载地址
    }
}
````

***



#### 3.2.6.获取特定版本的文件信息
描述：
>在历史版本预览和回滚历史版本的时候，获取特定版本文档的文件信息。

API地址：
>/v1/3rd/file/version/:version

调用方法：
>GET

请求头（由weboffice开放平台写入）：

| **Header**	| **值**	| **描述** |
| :---: | :---: | --- |
 | User-Agent | 	wps-weboffice-openplatform | 	用户代理 | 
 | x-wps-weboffice-token	 | xxxxxxx	 | 校验身份的token，值根据对接的企业而定 | 
 | x-weboffice-file-id	 | xxxxxxx	 | 文件id | 


请求参数：

| **参数名**	| **参数类型**	| **数据类型**	| **描述**	| **可选**	|
| :---: | :---: | :---: | --- | :---: |
 | _w_signature	 | query	 | string	 | 请求签名	 | 否 | 
 | _w_appid	 | query	 | string	 | 应用id	 | 否 | 
 | *（任意参数）	 | query	 | *（任意类型） | 	按照需求传入参数	 | 是 | 

示例：
>http://www.xxx.cn/v1/3rd/file/version/6?_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&_w_userid=1000&_w_fname=example.doc&_w_signature=tBnhFqvVrC1LcJKye1m7GjzvqoA%3D
>
>此处的_w_fname等等是对接模块处理的参数，可以对接模块自己定义，weboffice服务不会对此参数进行处理,http://www.xxx.cn 为企业对应的域名和nginx转发地址。

返回值：
>status code为200表示获取数据成功，返回下面的信息。其余值表示失败，需要在返回中指定code,msg等信息。

正确返回响应信息示例：
````
{
    file: {
        id:”f132aa30a87064”,   		        //文件id,字符串长度小于40
        name: "example.doc",			//文件名
        version: 6,				//当前版本号，位数小于11
        size: 200,				//文件大小，单位为kb		
        create_time: 1136185445,		//创建时间，时间戳，单位为秒
        creator: "id0",			        //创建者id，字符串长度小于40	
        modify_time: 1551409818,		//修改时间，时间戳，单位为秒
        modifier: "id1000",			//修改者id，字符串长度小于40
        download_url: "http://www.xxx.cn/v1/file?fid=f132aa30a87064&version=6" //文档下载地址
    }
}
````


***


#### 3.2.7.文件重命名
描述：
>用户在h5页面修改了文件名后，把新的文件名上传到服务端保存。

API地址：
>/v1/3rd/file/rename

调用方法：
>PUT

请求头（由weboffice开放平台写入）：

| **Header**	| **值**	| **描述** |
| :---: | :---: | --- |
 | User-Agent	 | wps-weboffice-openplatform	 | 用户代理 | 
 | x-wps-weboffice-token	 | xxxxxxx | 	校验身份的token，值根据对接的企业而定 | 
 | x-weboffice-file-id | 	xxxxxxx | 	文件id | 


请求参数：

| **参数名**	| **参数类型**	| **数据类型**	| **描述**	| **可选**	|
| :---: | :---: | :---: | --- | :---: |
 | _w_signature	 | query	 | string	 | 请求签名 | 	否 | 
 | _w_appid	 | query	 | string | 	应用id | 	否 | 
 | *（任意参数）	 | query	 | *（任意类型） | 	按照需求传入参数 | 	是 | 
 | name	 | body	 | string	 | 文件新名称 | 	否 | 

PUT请求示例：
````
{ 
    "name": "rename.doc" 	 //文件新名称
} 	
````
	                                                                         
示例：
>http://www.xxx.cn/v1/3rd/file/rename?_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&&_w_userid=1000&_w_fname=example.doc&_w_signature=tBnhFqvVrC1LcJKye1m7GjzvqoA%3D
>
>此处的_w_fname等等是对接模块处理的参数，可以对接模块自己定义，weboffice服务不会对此参数进行处理,http://www.xxx.cn 为企业对应的域名和nginx转发地址。

返回值：
>status code为200表示获取数据成功。其余值表示失败，需要在返回中指定code,msg等信息。


***


#### 3.2.8.获取所有历史版本文件信息
描述：
>获取当前文档所有历史版本的文件信息，以数组的形式,按版本号从大到小的顺序返回响应。

API地址：
>/v1/3rd/file/history

调用方法：
>POST

请求头（由weboffice开放平台写入）：

| **Header**	| **值**	| **描述** |
| :---: | :---: | --- |
| User-Agent	 | wps-weboffice-openplatform | 	用户代理 | 
| x-wps-weboffice-token | xxxxxxx	 | 校验身份的token，值根据对接的企业而定 | 
| x-weboffice-file-id	 | xxxxxxx	 | 文件id | 

请求参数：

 | **参数名**	| **参数类型**	| **数据类型**	| **描述**	| **可选**	|
 | :---: | :---: | :---: | --- | :---: |
 | _w_signature	 | query	 | string	 | 请求签名 | 	否 | 
 | _w_appid	 | query	 | string	 | 应用id	 | 否
 | id	 | body	 | string	 | 文件id | 	否 | 
 | offset	 | body	 | int	 | 记录偏移量	 | 否 | 
 | count	 | body	 | int | 	记录总数	 | 否 | 
 | *（任意参数）	 | query	 | *（任意类型） | 	按照需求传入参数	 | 是 | 

示例：
>http://www.xxx.cn/v1/3rd/file/history?_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&_w_userid=1000&_w_fname=example.docx&_w_signature=tBnhFqvVrC1LcJKye1m7GjzvqoA%3D

POST请求示例：
````
{
    id: "f132aa30a87064”,
    offset: 0,
    count: 3
}
```` 

                    
返回值：
>status code为200表示获取数据成功，返回下面的信息。其余值表示失败，需要在返回中指定code,msg等信息。

正确返回响应信息示例：
````
{
    histories: [
        {
            id: "f132aa30a87064",		//文件id,字符串长度小于40
            name: "example.doc",		//文件名
            version: 3,		        	//版本号，位数小于11
            size: 200,			        //文件大小	
            download_url: "http://www.xxx.cn/v1/file?fid=f132aa30a87064&version=3",  //文档下载地址
            create_time: 1136185445,		//创建时间，以时间戳表示，单位为秒
            modify_time: 1539847453,		//修改时间，以时间戳表示，单位为秒
            creator: {
                id: "id0",			        //创建者id，字符串长度小于40	
                name: "wps-0",			        //创建者名字	
                avatar_url: "http://xxx.cn/?id=0"	//创建者头像地址
            },	
            modifier: {		
                id: "id1000",			//修改者id，字符串长度小于40	
                name: "wps-1000",		//修改者名字	
                avatar_url: "http://xxx.cn/?id=1000"	//修改者头像地址
            }
        },
        {
            id: "f132aa30a87064",		//文件id,字符串长度小于40
            name: "example.doc",		//文件名
            version: 2,			//版本号，位数小于11
            size: 200,			//文件大小	
            download_url: "http://www.xxx.cn/v1/file?fid=f132aa30a87064&version=2",
					//文档下载地址
            create_time: 1136185445,		//创建时间，以时间戳表示，单位为秒
            modify_time: 1539847453,		//修改时间，以时间戳表示，单位为秒
            creator: {
                id: "id0",			//创建者id，字符串长度小于40	
                name: "wps-0",			//创建者名字	
                avatar_url: "http://xxx.cn/?id=0"	//创建者头像地址
            },	
            modifier: {		
                id: "id1000",			//修改者id，字符串长度小于40	
                name: "wps-1000",		//修改者名字	
                avatar_url: "http://xxx.cn/?id=1000"	//修改者头像地址
            }
        },
        {
            id: "f132aa30a87064",		//文件id,字符串长度小于40
            name: "example.doc",		//文件名
            version: 1,			//版本号，位数小于11
            size: 200,			//文件大小	
            download_url: "http://www.xxx.cn/v1/file?fid=f132aa30a87064&version=1",
					//文档下载地址
            create_time: 1136185445,		//创建时间，以时间戳表示，单位为秒
            modify_time: 1539847453,		//修改时间，以时间戳表示，单位为秒
            creator: {
                id: "id0",			//创建者id，字符串长度小于40	
                name: "wps-0",			//创建者名字	
                avatar_url: "http://xxx.cn/?id=0"	//创建者头像地址
            },	
            modifier: {		
                id: "id1000",			//修改者id，字符串长度小于40	
                name: "wps-1000",		//修改者名字	
                avatar_url: "http://xxx.cn/?id=1000"	//修改者头像地址
            }
        }    
    ]
}
````


***


#### 3.2.9.新建文件
描述：
>在模板页选择对应的模板后，将新创建的文件上传到对接模块，返回访问此文件的跳转url。

文字模板列表地址：
>https://wwo.wps.cn/office/w/new/0?_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxappid&_w_signature=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&......(对接模块需要的自定义参数)


表格模板列表地址：
>https://wwo.wps.cn/office/s/new/0?_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxappid&_w_signature=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&......(对接模块需要的自定义参数)

API地址: 
>/v1/3rd/file/new

调用方法：
>POST

请求头（由weboffice开放平台写入）：

| **Header**	| **值**	| **描述** |
| :---: | :---: | --- |
 | User-Agent	 | wps-weboffice-openplatform	 | 用户代理 | 
 | x-wps-weboffice-token	 | xxxxxxx	 | 校验身份的token，值根据对接的企业而定 | 
 | x-weboffice-file-id	 | xxxxxxx	 | 文件id | 

请求参数示例：

 | **参数名**	| **参数类型**	| **数据类型**	| **描述**	| **可选**	|
 | :---: | :---: | :---: | --- | :---: |
 | _w_signature	 | query	 | string	 | 请求签名	 | 否 | 
 | _w_appid	 | query	 | string	 | 应用id | 	否 | 
 | *（任意参数）	 | query	 | *（任意类型） | 	按照需求传入参数 | 	是 | 
 | file	 | body	 | file	 | 新创建的文件 | 	否 | 
 | name	 | body	 | string	 | 新创建的文件名 | 	否 | 

示例：
>http://www.xxx.cn/v1/3rd/file/new?_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&_w_signature=tBnhFqvVrC1LcJKye1m7GjzvqoA%3D
>
>此处的_w_fname等等是对接模块处理的参数，可以对接模块自己定义，weboffice服务不会对此参数进行处理,http://www.xxx.cn 为企业对应的域名和nginx转发地址。

POST请求示例：
````
Content-Disposition: form-data; name="file"; filename="/home/wps/dir/1.docx”       //待上传文件对象
Content-Disposition: form-data; name="name" 1.docx   		//新文件名称
````
返回值：
>status code为200表示获取数据成功，返回下面的信息。其余值表示失败，需要在返回中指定code,msg等信息。


响应信息示例：
````
{
    redirect_url: "http://wwo.wps.cn/office/w/<:fileid>?_w_fname=example.doc&_w_userid=1000&_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&_w_signature=878e966c8e729a2a28e699a3455a57f2",
		  //根据此url，可以访问到对应创建的文档
    user_id: "id1000"  //创建此文档的用户id 
}
````

***


### 3.3.错误处理
返回的错误信息，http status 返回200表示处理正常，其余的都表示错误，返回下面的错误码，以定位问题。

| 参数名 | 数据类型 | 描述 |
| :---: | :---: | :--- |
| code | 	int	 | 错误码
| message | 	string | 	错误信息

示例：
````
{
    code: 40001,
    message: "User not logged in"
}
````


***


#### 3.3.1.错误码
预览服务返回的错误信息

| 错误码 | 错误信息 | 描述 |
| :---: | :---: | :--- |	
| 40001	 | 用户未登陆 | 	对接模块认证用户身份失败 |
| 40002	 | token过期 | 	对接模块校验用户token过期 |
| 40003	 | 用户无权限访问	 | 对接模块校验用户没有访问此文档的权限 |
| 40004	 | 资源不存在 | 	对接模块校验文档为无效文件或文件不存在 |

***


#### 3.3.2.错误信息
返回的错误字符串，是对错误信息的详细描述，方便定位问题。

***

## 四.签名生成说明
以如下请求为例，来说明signature的生成步骤：
>https://wwo.wps.cn/office/w/1?_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxappid&_w_userid=1000&_w_fname=example.doc&_w_signature=tBnhFqvVrC1LcJKye1m7GjzvqoA%3D

 ***1. 将以”_w_”作为前缀所有参数按key进行字典升序排列，将排序后的key,value字符串以%s=%s格式拼接起来。如下：***
> _w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxappid_w_fname=example.doc_w_userid=1000 
>
>请开发者关注：只有以”_w_”作为前缀才需要签名，否则容易导致签名不能通过验证。

***2. 将appsecret加到最后，得到最终加密的字符串。如下：***

>_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxappid_w_fname=example.doc_w_userid=1000_w_secretkey=xxxxxxxxxxxxxxxxxxxxxxxsecretkey
其中xxxxxxxxxxxxxxxxxxxxxxxxxxxappid为实际申请的appid，xxxxxxxxxxxxxxxxxxxxxxxsecretkey为实际申请的appsecret。

***3. 生成签名值***

   使用HMAC-SHA1哈希算法，使用注册的appsecret密钥对上一步骤的源字符串进行加密。（注：一般程序语言中会内置HMAC-SHA1加密算法的函数，例如PHP5.1.2之后的版本可直接调用hash_hmac函数。）

***4. 然后将加密后的字符串经过Base64编码。 （注：一般程序语言中会内置Base64编码函数，例如PHP中可直接调用 base64_encode() 函数。）***

***5. 得到的签名值结果如下：***

>tBnhFqvVrC1LcJKye1m7GjzvqoA=

***6. 将上面生成的字符串进行URL编码。***

请开发者关注：URL编码注意事项，否则容易导致后面签名不能通过验证。 

例如：
>tBnhFqvVrC1LcJKye1m7GjzvqoA%3D 



经过以上步骤后，即生成示例中的signature。JAVA示例代码：
````
//签名方法
import org.apache.commons.codec.digest.HmacUtils;
import java.util.*;
import static org.apache.commons.codec.binary.Base64.encodeBase64String;

public class Main {
       private static String getSignature(Map<String, String> params, String appId, String appSecret) {
        List<String> keys=new ArrayList();
        for (Map.Entry<String, String> entry : params.entrySet()) {
            keys.add(entry.getKey());
        }

        // 将所有参数按key的升序排序
        Collections.sort(keys, new Comparator<String>() {
            public int compare(String o1, String o2) {
                return o1.compareTo(o2);
            }
        });

        // 构造签名的源字符串
        StringBuilder contents=new StringBuilder("");
        for (String key : keys) {

            if (key=="_w_signature"){
                continue;
            }
            contents.append(key+"=").append(params.get(key));
            System.out.println("key:"+key+",value:"+params.get(key));
        }
        contents.append("_w_secretkey=").append(appSecret);


        System.out.println(appSecret);
        System.out.println(contents.toString());

        // 进行hmac sha1 签名
        byte[] bytes= HmacUtils.
                hmacSha1(appSecret.getBytes(),contents.toString().getBytes());

        //字符串经过Base64编码
        String sign= encodeBase64String(bytes);
        System.out.println(sign);
        return sign;
    }
    public static Map<String, String> paramToMap(String paramStr) {
        String[] params = paramStr.split("&");
        Map<String, String> resMap = new HashMap<String, String>();
        for (int i = 0; i < params.length; i++) {
            String[] param = params[i].split("=");
            if (param.length >= 2) {
                String key = param[0];
                String value = param[1];
                for (int j = 2; j < param.length; j++) {
                    value += "=" + param[j];
                }
                resMap.put(key, value);
            }
        }
        return resMap;
    }
}
````
***

## 五.访问地址
##### 1. 文档访问的url地址根据此格式生成：

>https://wwo.wps.cn/office/<:type>/<:fileid>?_w_appid=xxxxxxxxxxxxxxxxxxxxxxxxxxxappid&_w_signature=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx&......(对接模块需要的自定义参数)

##### 2. 文档访问是进入在线编辑模式还是预览模式，取决于对接模块返回的文件元数据中permission返回的值，为“write”进入编辑模式，为“read”进入预览模式，详见3.2.1。

文字文件示例：

>https://wwo.wps.cn/office/w/471eba50307c1f9dc540?_w_fname=%E4%BC%9A%E8%AE%AE%E7%BA%AA%E8%A6%81.docx&_w_userid=33&_w_appid=d8f99daa999e47965f4b9727e32ddaa8&_w_permission=read&_w_signature=I%2BnTsXn8F5wp%2FBqzOP4fX6E2s2M%3D

表格文件示例：

>https://wwo.wps.cn/office/s/56922f98b164afbd0daf?_w_fname=%E8%80%83%E5%8B%A4%E8%A1%A8.xlsx&_w_userid=33&_w_appid=d8f99daa999e47965f4b9727e32ddaa8&_w_permission=read&_w_signature=g4fJKz5GToEV2UpekVytveCI8oI%3D

演示文件示例：

>https://wwo.wps.cn/office/p/b08631b2e65b0c1b8fb6?_w_fname=%E5%B7%A5%E4%BD%9C%E6%80%BB%E7%BB%93.pptx&_w_userid=31&_w_appid=d8f99daa999e47965f4b9727e32ddaa8&_w_permission=read&_w_signature=JlUqYGu5eyPhIAaa4kRNp1Q8yEo%3D

PDF文件示例：

>https://wwo.wps.cn/office/f/548c695ebac24f746a91?_w_fname=%E7%BA%A2%E5%A4%B4%E6%96%87%E4%BB%B6.pdf&_w_userid=0&_w_appid=d8f99daa999e47965f4b9727e32ddaa8&_w_permission=write&_w_signature=9E8gQ4myphpfqB5tIGbrVFn6b1Q%3D

