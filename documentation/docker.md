### Docker
Plik 提供了一个简单的 Dockerfile，允许您在 docker 容器中运行它。

##### 从 docker 注册表获取镜像

```sh
$ docker pull lamfire8390/plik:latest
```

##### 构建 docker 镜像

首先，您需要构建 docker 镜像：   
```sh
$ make docker
```

##### 配置

然后您可以运行一个实例并将本地端口 8080 映射到 plik 端口：   
```sh
$ docker run -t -d -p 8080:8080 rootgg/plik
ab9b2c99da1f3e309cd3b12392b9084b5cafcca0325d7d47ff76f5b1e475d1b9
```

要使用不同的配置文件，您可以在运行时将单个文件映射到容器：   
这里，我们将本地文件 plikd.cfg 映射到 home/plik/server/plikd.cfg，这是容器中的默认配置文件位置：   
```sh
$ docker run -t -d -p 8080:8080 -v plikd.cfg:/home/plik/server/plikd.cfg rootgg/plik
ab9b2c99da1f3e309cd3b12392b9084b5cafcca0325d7d47ff76f5b1e475d1b9
```

您还可以使用卷将上传存储在容器外部：   
这里，我们将本地文件夹 /data 映射到容器的 /home/plik/server/files 文件夹，这是默认的上传目录：   
```sh
$ docker run -t -d -p 8080:8080 -v /data:/home/plik/server/files rootgg/plik
ab9b2c99da1f3e309cd3b12392b9084b5cafcca0325d7d47ff76f5b1e475d1b9
```


### 使用 docker-compose

使用此示例文件设置您的实例，包含所有持久数据/元数据。在此配置中，所有文件、账户和令牌都将持久保存。
根据您的需要调整目录。

```
$ cd ~
$ mkdir plik
$ curl https://raw.githubusercontent.com/lamfire/plik/master/server/plikd.cfg # 复制服务器配置
$ plik mkdir data # 创建目录以在 docker 镜像外保存文件和/或元数据
$ plik chown 1000:1000 data # 与 docker 匹配 UID
$ plik chown 1000:1000 plikd.cfg # 与 docker 匹配 UID
```

编辑 plikd.cfg 将元数据和/或数据指向您可以在 docker-compose 中匹配的挂载点（此示例中为 /data）
```
DataBackend = "file"
[DataBackendConfig]
    Directory = "/data/files" # <===

[MetadataBackendConfig]
    Driver = "sqlite3"
    ConnectionString = "/data/plik.db" # <===
```

创建包含以下内容的 docker-compose.yml 文件
```yaml
version: "2"
services:
  plik:
    image: lamfire8390/plik:latest
    container_name: plik
    volumes:
      - /home/{user}/plik/plikd.cfg:/home/plik/server/plikd.cfg
      - /home/{user}/plik/data:/data
    ports:
      - 8080:8080   
    restart: "unless-stopped"
```

```
$ docker-compose up
Starting plik ... done
Attaching to plik
plik    | [01/27/2022 10:48:26][INFO    ] Starting plikd server v...
plik    | [01/27/2022 10:48:26][INFO    ] Starting server at http://0.0.0.0:8080
```


