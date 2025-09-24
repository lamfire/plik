[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![Build](https://github.com/root-gg/plik/actions/workflows/master.yaml/badge.svg)](https://github.com/root-gg/plik/actions/workflows/master.yaml)
[![Go Report](https://img.shields.io/badge/Go_report-A+-brightgreen.svg)](http://goreportcard.com/report/root-gg/plik)
[![Docker Pulls](https://img.shields.io/docker/pulls/rootgg/plik.svg)](https://hub.docker.com/r/rootgg/plik)
[![GoDoc](https://godoc.org/github.com/root-gg/plik?status.svg)](https://godoc.org/github.com/root-gg/plik)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](http://opensource.org/licenses/MIT)

想和我们聊天吗？Telegram 频道：https://t.me/plik_rootgg

# Plik

Plik 是一个用 Go 语言编写的可扩展且友好的临时文件上传系统（类似 Wetransfer）。

### 主要功能
   - 功能强大的命令行客户端
   - 易于使用的 Web 界面
   - 多种数据后端：文件系统、OpenStack Swift、S3、Google Cloud Storage
   - 多种元数据后端：Sqlite3、PostgreSQL、MySQL
   - OneShot：文件在首次下载后销毁
   - Stream：文件从上传者流式传输到下载者（服务器端不存储任何内容）
   - Removable：允许上传者随时删除文件
   - TTL：自定义过期日期
   - Password：使用登录名/密码保护上传（基本认证）
   - Comments：添加自定义消息（支持 Markdown 格式）
   - 用户认证：本地 / Google / OVH
   - 上传限制：源 IP / Token
   - 管理员 CLI 和 Web 界面
   - 服务器端加密（使用 S3 数据后端）
   - 多架构构建和 Docker 镜像
   - [ShareX](https://getsharex.com/) 上传器：直接集成到 ShareX
   - [plikSharp](https://github.com/iss0/plikSharp)：Plik 的 .NET API 客户端
   - [Filelink for Plik](https://gitlab.com/joendres/filelink-plik)：Thunderbird 附件上传到 Plik 的插件

### 目录
1. [安装](#installation) 
2. [配置](#configuration)
3. [数据后端](#data-backends)
4. [元数据后端](#metadata-backends)
5. [Web 界面](#web-ui)
6. [客户端 CLI](#cli-client)
7. [Go 客户端](#go-client)
8. [HTTP API](#api)
9. [管理员 CLI](#admin-cli)
10. [认证](#authentication)
11. [安全](#security)
12. [交叉编译](#cross-compilation)
13. [常见问题](#faq)
14. [如何贡献](#how-to-contribute)

### 安装 <a name="installation"></a>

##### 从发布版本
运行 plik 非常简单：
```sh
$ wget https://github.com/root-gg/plik/releases/download/1.3.8/plik-1.3.8-linux-amd64.tar.gz
$ tar xzvf plik-1.3.8-linux-amd64.tar.gz
$ cd plik-1.3.8-linux-amd64/server
$ ./plikd
```
完成！您现在在 http://127.0.0.1:8080 上运行着一个功能完整的 Plik 实例。
您可以编辑 server/plikd.cfg 来根据您的需求调整配置（端口、SSL、TTL、后端参数等）

##### 从源码
要从源码编译 plik，您需要在系统上安装 golang 和 npm

Git clone 或 go get 项目，然后简单地运行 make：
```sh
$ make
$ cd server && ./plikd
```

##### Docker <a name="docker"></a>
Plik 提供了为 linux amd64/i386/arm/arm64 构建的多架构 Docker 镜像：
 - rootgg/plik:latest（最新发布版本）
 - rootgg/plik:{version}（发布版本）
 - rootgg/plik:dev（master 分支的最新提交）

请参阅 [Plik Docker 参考文档](documentation/docker.md)

Plik 还提供了一些有用的脚本来在独立的 Docker 实例中测试后端：

请参阅 [Plik Docker 后端测试](testing)

### 配置 <a name="configuration"></a>

配置使用 TOML 文件 [plikd.cfg](server/plikd.cfg) 进行管理

###### 使用环境变量定义配置参数

可以使用环境变量指定配置参数，配置参数使用大写蛇形命名法
```
    PLIKD_DEBUG_REQUESTS=true ./plikd
```

对于数组和配置映射，必须以 JSON 格式提供。
数组会被覆盖，但映射会被合并

```
    PLIKD_DATA_BACKEND_CONFIG='{"Directory":"/var/files"}' ./plikd
```

### 数据后端 <a name="data-backends"></a>

Plik 提供了多种用于上传文件的数据后端和用于上传元数据的元数据后端。

 - 文件数据后端：

将上传的文件存储在本地或挂载的文件系统目录中。

 - Openstack Swift 数据后端：http://docs.openstack.org/developer/swift/

Openstack Swift 是一个高可用、分布式、最终一致的对象/Blob 存储，支持服务器端加密

 - Amazon S3

 - Google Cloud Storage

### 元数据后端 <a name="metadata-backends"></a>

 - Sqlite3

适用于独立部署。

 - PostgreSQL / Mysql

适用于分布式 / 高可用部署。

### Web 界面 <a name="web-ui"></a>

默认情况下，Plikd 在与 API 相同的端口上提供 Angularjs Web 界面。
可以通过在 plikd.cfg 中设置 "NoWebInterface" 来禁用此行为。

可以通过在 plikd.cfg 中设置 "WebappDirectory" 来更改 WebUI 路径。
它默认为发布 .tar.gz 中的相对路径 '../webapp/dist'。

界面可以通过以下几种方式进行自定义：
- 可以在 `js/custom.js` 中更改标题
- 可以在 `css/custom.css` 中覆盖 CSS 样式（使用 !important）
- 可以更改背景：`img/background.jpg`
- 可以更改网站图标：`favicon.ico`

如果您使用 Docker，文件位于容器中的 `/home/plik/webapp/dist`。例如：
```sh
$ docker run -t -d -p 8080:8080 -v my_background.jpg:/home/plik/webapp/dist/img/background.jpg rootgg/plik
```

### CLI 客户端 <a name="cli-client"></a>
Plik 提供了一个功能强大的 Go 多平台 CLI 客户端（可在 Web 界面中下载）：

```
用法：
  plik [选项] [文件] ...

选项：
  -h --help                 显示此帮助
  -d --debug                启用调试模式
  -q --quiet                启用静默模式
  -o, --oneshot             启用 OneShot（每个文件在首次下载时删除）
  -r, --removable           启用可删除上传（任何人都可以随时删除每个文件）
  -S, --stream              启用流式传输（将阻塞直到远程用户开始下载）
  -t, --ttl TTL             过期前的时间（上传将在 m|h|d 后删除）
  -n, --name NAME           从 STDIN 管道时设置文件名
  --server SERVER           覆盖 plik url
  --token TOKEN             指定上传令牌
  --comments COMMENT        设置上传的注释（兼容 MarkDown）
  -p                        使用登录名和密码保护上传
  --password PASSWD         使用登录名:密码保护上传（如果省略，默认登录名为 "plik"）
  -a                        使用默认归档参数归档上传（参见 ~/.plikrc）
  --archive MODE            使用指定的归档后端归档上传：tar|zip
  --compress MODE           [tar] 压缩编解码器：gzip|bzip2|xz|lzip|lzma|lzop|compress|no
  --archive-options OPTIONS [tar|zip] 额外的命令行选项
  -s                        使用默认加密参数加密上传（参见 ~/.plikrc）
  --not-secure              无论 ~/.plikrc 配置如何都不加密上传
  --secure MODE             使用指定的归档后端归档上传：openssl|pgp
  --cipher CIPHER           [openssl] 使用的 Openssl 密码（参见 openssl 帮助）
  --passphrase PASSPHRASE   [openssl] 密码或 '-' 提示输入密码
  --recipient RECIPIENT     [pgp] 为 pgp 后端设置接收者（例如：--recipient Bob）
  --secure-options OPTIONS  [openssl|pgp] 额外的命令行选项
  --update                  更新客户端
  -v --version              显示客户端版本
```

例如，创建目录 tar.gz 归档并使用 openssl 加密：
```bash
$ plik -a -s mydirectory/
密码：30ICoKdFeoKaKNdnFf36n0kMH
上传成功创建：
    https://127.0.0.1:8080/#/?id=0KfNj6eMb93ilCrl

mydirectory.tar.gz：15.70 MB 5.92 MB/s

命令：
curl -s 'https://127.0.0.1:8080/file/0KfNj6eMb93ilCrl/q73tEBEqM04b22GP/mydirectory.tar.gz' | openssl aes-256-cbc -d -pass pass:30ICoKdFeoKaKNdnFf36n0kMH | tar xvf - --gzip
```

客户端配置和首选项存储在 ~/.plikrc 或 /etc/plik/plikrc（可通过 PLIKRC 环境变量覆盖）

### 仅使用 curl 快速上传

```bash
curl --form 'file=@/path/to/file' http://127.0.0.1:8080
```
当使用认证且启用了 NoAnonymousUploads 时，您可以使用用户令牌快速上传：
```bash
curl --form 'file=@/path/to/file' --header 'X-PlikToken: xxxx-xxx-xxxx-xxxxx-xxxxxxxx' http://127.0.0.1:8080
```

必须设置 DownloadDomain 配置选项才能正常工作。

### Go 客户端 <a name="go-client"></a>

Plik 现在提供了一个 Go 库，CLI 客户端基于此库构建

请参阅 [Plik 库参考文档](plik/README.md)

### API <a name="api"></a>
Plik 服务器公开了一个 HTTP API 来管理上传和获取文件：

请参阅 [Plik API 参考文档](documentation/api.md)

### 管理员 CLI <a name="admin-cli"></a>

使用 ./plikd 服务器二进制文件可以：
  - 创建/列出/删除本地账户
  - 创建/列出/删除用户 CLI 令牌
  - 创建/列出/删除文件和上传
  - 导入 / 导出元数据

更多详细信息请参阅帮助

### 认证 <a name="authentication"></a>

Plik 可以使用本地账户或使用 Google 或 OVH API 来认证用户。

要启用认证，请在 plikd.cfg 中将 FeatureAuthentication 设置为 "enabled"
要只允许认证用户上传文件，请在 plikd.cfg 中将 FeatureAuthentication 设置为 "forced"

如果启用了源 IP 地址限制，用户账户只能从受信任的 IP 创建，然后认证用户可以在没有源 IP 限制的情况下上传文件。

管理员用户可以访问管理员仪表板并操作所有上传。

   - **本地**：
      - 您可以使用服务器命令行操作本地用户
      
      ```sh
      $ ./plikd --config ./plikd.cfg user create --login root --name Admin --admin    
      为用户 root 生成的密码是 08ybEyh2KkiMho8dzpdQaJZm78HmvWGC
      ```
      
   - **Google**：
      - 您需要在 [Google 开发者控制台](https://console.developers.google.com) 中创建一个新应用程序
      - 您将获得一个 Google API ClientID 和一个 Google API ClientSecret，需要将它们放在 plikd.cfg 文件中。
      - 不要忘记为您的域名白名单有效的来源和重定向 URL（https://yourdomain/auth/google/callback）。
      - 可以只白名单一个或多个电子邮件域名。
   
   - **OVH**：
      - 您需要在 OVH API 中创建一个新应用程序：https://eu.api.ovh.com/createApp/
      - 您将获得一个 OVH 应用程序密钥和一个 OVH 应用程序秘密密钥，需要将它们放在 plikd.cfg 文件中。

认证后，用户可以生成上传令牌，可以在 ~/.plikrc 文件中指定这些令牌来认证命令行客户端。

```
Token = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
```

### 安全 <a name="security"></a>
Plik 允许用户按原样上传和提供任何内容，但托管不受信任的 HTML 会带来一些众所周知的安全问题。

Plik 将尝试通过将 Content-Type 覆盖为 "text-plain" 而不是 "text/html" 来避免 HTML 渲染。

默认情况下，Plik 设置了一些安全 HTTP 头，如 **X-Content-Type-Options、X-XSS-Protection、X-Frame-Options、Content-Security-Policy**，以禁用大多数现代浏览器的敏感功能，如资源加载、xhr 请求、iframe 等。
但是，这会破坏音频/视频播放、PDF 渲染等功能，因此可以通过将 EnhancedWebSecurity 配置参数设置为 false 来禁用此行为。

除此之外，还强烈建议在单独的（子）域上提供上传的文件，以对抗网络钓鱼链接，并使用 DownloadDomain 配置参数保护 Plik 的会话 cookie。

### 交叉编译 <a name="cross-compilation"></a>

所有二进制文件现在都是静态链接的。客户端可以安全地交叉编译到所有操作系统/架构，因为它们不依赖 GCO（sqlite）
服务器依赖 CGO/sqlite 需要交叉编译就绪环境。

 `make release` 将为 `amd64,i386,arm,arm64` 构建发布归档

要仅构建客户端的特定架构的发布
```
    CLIENT_TARGETS="linux/amd64" releaser/release.sh
```

要仅构建特定架构的发布
```
    TARGETS="linux/amd64" releaser/release.sh
```

要使用特定的交叉编译器工具链构建
```
    TARGETS="linux/arm/v6" CC=arm-linux-gnueabihf-gcc releaser/release.sh
```

### 常见问题 <a name="faq"></a>

* 为什么在多实例部署中流模式会失效？

因为流模式不是无状态的。由于上传者请求将阻塞在一个 plik 实例上，下载者请求**必须**转到同一实例才能成功。
负载平衡策略**必须**意识到这一点，并通过哈希文件 ID 将流请求路由到同一实例。

以下是使用 nginx 和一小段 LUA 实现此目的的示例。
确保您的 nginx 服务器构建时支持 LUA 脚本。
您可能想要安装带有内置 LUA 支持的 "nginx-extras" Debian 包（>1.7.2）。
```
upstream plik {
    server 127.0.0.1:8080;
    server 127.0.0.1:8081;
}

upstream stream {
    server 127.0.0.1:8080;
    server 127.0.0.1:8081;
    hash $hash_key;
}

server {
    listen 9000;

    location / {
        set $upstream "";
        set $hash_key "";
        access_by_lua '
            _,_,file_id = string.find(ngx.var.request_uri, "^/stream/[a-zA-Z0-9]+/([a-zA-Z0-9]+)/.*$")
            if file_id == nil then
                ngx.var.upstream = "plik"
            else
                ngx.var.upstream = "stream"
                ngx.var.hash_key = file_id
            end
        ';
        proxy_pass http://$upstream;
    }
}
```

* 使用 DownloadDomain 强制和反向代理的重定向循环

```
Invalid download domain 127.0.0.1:8080, expected plik.root.gg
```

DownloadDomain 检查传入 HTTP 请求的 Host 头，默认情况下，像 Nginx 或 Apache mod_proxy 这样的反向代理不会转发此头。检查以下配置指令：

```
Apache mod_proxy：ProxyPreserveHost On
Nginx：proxy_set_header Host $host;
```

* 从客户端上传时出现错误："Unable to upload file：HTTP error 411 Length Required"

在 nginx < 1.3.9 下，您必须启用 HttpChunkin 模块以允许传输编码 "chunked"。
您可能想要安装带有内置 HttpChunkin 模块的 "nginx-extras" Debian 包。

并在您的服务器配置中添加：

```sh
chunkin on;
error_page 411 = @my_411_error;
location @my_411_error {
        chunkin_resume;
}
```

* 如何禁用 nginx 缓冲？

默认情况下，nginx 将大型 HTTP 请求和响应缓冲到临时文件。这种行为会导致不必要的磁盘负载和较慢的传输。对于 /file 和 /stream 路径，应该关闭此功能（>1.7.12）。您可能还想增加缓冲区大小。

详细文档：http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_buffering
```
proxy_buffering off;
proxy_request_buffering off;
proxy_http_version 1.1;
proxy_buffer_size 1M;
proxy_buffers 8 1M;
client_body_buffer_size 1M;
```

* 为什么在设置 EnhancedWebSecurity 时，HTTP 连接的认证不工作？

当设置 EnhancedWebSecurity 时，Plik 会话 cookie 设置了 "secure" 标志，因此它们只能通过安全的 HTTPS 连接传输。

* 是否有用于将 plik 作为服务运行的 OpenRC 模板？

是的（必须创建 plikuser）：

```
#!/sbin/openrc-run

name=$RC_SVCNAME
description="Plik File Sharing Service"
command="/path/to/plik/server/plikd"
command_user="plikuser:plikuser"
command_background=true
directory="/path/to/plik/server"
pidfile=/run/${RC_SVCNAME}.pid
start_stop_daemon_args="--stdout /var/log/$RC_SVCNAME/${RC_SVCNAME}.log --stderr /var/log/$RC_SVCNAME/${RC_SVCNAME}.log"

depend() {
    use logger dns
    need net
    after firewall
}

start_pre() {
    checkpath --directory --owner $command_user --mode 0775 /var/log/$RC_SVCNAME
}
```

* 构建失败 "/usr/bin/env: 'node': No such file or directory"

Debian 用户可能需要安装 nodejs-legacy 包。

```
此包包含传统 Node.js 代码所需的符号链接
二进制文件为 /usr/bin/node（而不是 Debian 提供的 /usr/bin/nodejs）。
```

* 如何像老板一样截图并上传？

```
alias pshot="scrot -s -e 'plik -q \$f | xclip ; xclip -o ; rm \$f'"
```

需要您在 $PATH 中安装 plik、scrot 和 xclip。
scrot -s 允许您"使用鼠标交互式选择窗口或矩形"，然后
Plik 将上传截图，URL 将直接复制到您的剪贴板并由 xclip 显示。
然后从您的主目录中删除截图以避免垃圾。

### 如何为项目做贡献？<a name="how-to-contribute"></a>

欢迎贡献，请随时打开问题并/或提交拉取请求。
请确保也运行/更新测试套件：

```
    make fmt
    make lint
    make test
    make test-backends
```