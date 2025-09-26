### API
Plik 服务器提供 REST-full API 来管理上传和获取文件：

获取和创建上传：
 
   - **POST**        /upload
     - 参数（请求体中的 json 对象）：
      - oneshot (bool)
      - stream (bool)
      - removable (bool)
      - ttl (int)
      - login (string)
      - password (string)
      - files (见下文)
     - 返回：
         JSON 格式的上传对象。
         重要字段：
           - id (上传文件时必需)
           - uploadToken (上传/删除文件时必需)
           - files (见下文)

   对于流模式，您需要在上传开始前知道文件 ID，因为它会阻塞。
   文件大小和/或文件类型也需要在上传开始前知道，因为它们必须在 
   HTTP 响应头中打印。
   要获取文件 ID，请传递一个包含您即将上传的每个文件的 "files" json 对象。
   用任意字符串填充 reference 字段，以避免使用 fileName 字段匹配文件 ID。
   这也用于在文件上传尚未完成或失败时通知缺失文件。
  ```
  "files" : [
    {
      "fileName": "file.txt",
      "fileSize": 12345,
      "fileType": "text/plain",
      "reference": "0"
    },...
  ]
  ```
  
   - **GET** /upload/:uploadid:
     - 获取上传元数据（文件列表、上传日期、ttl 等）

上传文件：

   - **POST** /$mode/:uploadid:/:fileid:/:filename:
     - 请求体必须是包含名为 "file" 部分的 multipart 请求，该部分包含文件数据。

   - **POST** /file/:uploadid:
     - 与上面相同，但不传递文件 ID，不适用于流模式。
     
   - **POST** /:
     - 快速模式，自动创建具有默认参数的上传并将文件添加到其中。

获取文件：

  - **HEAD** /$mode/:uploadid:/:fileid:/:filename:
    - 仅返回 HTTP 头。用于在不下载文件的情况下了解 Content-Type 和 Content-Length。特别是当上传启用了 OneShot 选项时。

  - **GET**  /$mode/:uploadid:/:fileid:/:filename:
    - 下载文件。文件名**必须**匹配。浏览器可能会尝试显示文件，例如如果是 jpeg。您可以尝试在 URL 中使用 ?dl=1 强制下载。

  - **GET**  /archive/:uploadid:/:filename:
    - 以 zip 压缩包形式下载上传的文件。:filename: 必须以 .zip 结尾

删除文件：

   - **DELETE** /$mode/:uploadid:/:fileid:/:filename:
     - 删除文件。上传**必须**启用 "removable" 选项。

显示服务器详情：

   - **GET** /version
     - 显示 plik 服务器版本和一些构建信息（构建主机、日期、git 版本等）

   - **GET** /config
     - 显示 plik 服务器配置（ttl 值、最大文件大小等）

   - **GET** /stats
     - 获取服务器统计信息（上传/文件数量、用户数量、使用的总大小）
     - 仅管理员

用户认证：

   - 
   Plik 可以使用 Google 和/或 OVH 第三方 API 对用户进行身份验证。   
   /auth API 是为 Plik Web 应用程序设计的，但如果您想要自动化它，请确保提供有效的
   Referrer HTTP 头并转发所有会话 cookie。   
   Plik 会话 cookie 设置了 "secure" 标志，因此只能通过安全的 HTTPS 连接传输。   
   为避免 CSRF 攻击，plik-xsrf cookie 的值必须复制到每个
   认证请求的 X-XSRFToken HTTP 头中。   
   一旦认证，用户可以生成上传令牌。这些令牌可以在 X-PlikToken HTTP 头中使用，用于将
   上传链接到用户账户。它可以放在 Plik 命令行客户端的 ~/.plikrc 文件中。   
   
   - **本地** :
      - 您需要使用服务器命令行创建用户
   
   - **Google** :
      - 您需要在 [Google 开发者控制台](https://console.developers.google.com) 中创建新应用程序
      - 您将获得 Google API ClientID 和 Google API ClientSecret，需要将它们放在 plikd.cfg 文件中
      - 不要忘记为您的域名将有效的来源和重定向 URL（https://yourdomain/auth/google/callback）加入白名单
   
   - **OVH** :
      - 您需要在 OVH API 中创建新应用程序：https://eu.api.ovh.com/createApp/
      - 您将获得 OVH 应用程序密钥和 OVH 应用程序秘密密钥，需要将它们放在 plikd.cfg 文件中

   - **GET** /auth/google/login
      - 获取 Google 用户同意 URL。用户必须访问此 URL 进行身份验证

   - **GET** /auth/google/callback
     - 用户同意对话框的回调
     - 在此调用结束时，用户将被重定向回带有 Plik 会话 cookie 的 Web 应用程序

   - **GET** /auth/ovh/login
     - 获取 OVH 用户同意 URL。用户必须访问此 URL 进行身份验证
     - 响应将包含一个临时会话 cookie，用于将 API 端点和 OVH 消费者密钥转发到回调

   - **GET** /auth/ovh/callback
     - 用户同意对话框的回调。
     - 在此调用结束时，用户将被重定向回带有 Plik 会话 cookie 的 Web 应用程序

   - **POST** /auth/local/login
     - 参数：
       - login : 用户登录名
       - password : 用户密码

   - **GET** /auth/logout
     - 使 Plik 会话 cookie 失效

   - **GET** /me
     - 返回基本用户信息（ID、姓名、邮箱）和令牌

   - **DELETE** /me
     - 删除用户账户。

   - **GET** /me/token
     - 列出用户令牌
      - 此调用使用分页

   - **POST** /me/token
     - 创建新的上传令牌
     - 可以在 json 体中传递备注

   - **DELETE** /me/token/{token}
     - 撤销上传令牌

   - **GET** /me/uploads
     - 列出用户上传
     - 参数：
        - token : 按令牌过滤
      - 此调用使用分页

   - **DELETE** /me/uploads
     - 删除链接到用户账户的所有上传
     - 参数：
        - token : 按令牌过滤

   - **GET** /me/stats
     - 获取用户统计信息（上传/文件数量、使用的总大小）

   - **GET** /users
     - 列出所有用户
     - 此调用使用分页
     - 仅管理员

二维码：

   - **GET** /qrcode
     - 从 URL 生成二维码图像
     - 参数：
        - url  : 您想要存储在二维码中的 URL
        - size : 生成图像的像素大小（默认：250，最大：1000）


$mode 可以是 "file" 或 "stream"，取决于是否启用流模式。更多详情请参见 FAQ。

示例：
```sh
创建上传（在 json 响应中，您将获得上传 ID 和上传令牌）
$ curl -X POST http://127.0.0.1:8080/upload

创建 OneShot 上传
$ curl -X POST -d '{ "OneShot" : true }' http://127.0.0.1:8080/upload

上传文件到上传
$ curl -X POST --header "X-UploadToken: M9PJftiApG1Kqr81gN3Fq1HJItPENMhl" -F "file=@test.txt" http://127.0.0.1:8080/file/IsrIPIsDskFpN12E

获取头信息
$ curl -I http://127.0.0.1:8080/file/IsrIPIsDskFpN12E/sFjIeokH23M35tN4/test.txt
HTTP/1.1 200 OK
Content-Disposition: filename=test.txt
Content-Length: 3486
Content-Type: text/plain; charset=utf-8
Date: Fri, 15 May 2015 09:16:20 GMT

```
