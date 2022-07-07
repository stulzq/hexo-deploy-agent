# Hexo Deploy Agent

[English](README.md)|中文

Hexo 部署 Agent，基于 Github Actions 可实现完全自动化部署 Hexo 博客，每次提交都会自动打包、部署、更新和刷新 CDN 缓存。

特性：

- 支持 Github Action 或者 Jenkins 等自动化工具
- 通过 CURL 上传部署包
- 解压、动态更新网站文件
- 目录级别刷新 CDN（目前仅支持腾讯云）
- 支持部署消息推送（目前仅支持钉钉）

Demo: https://xcmaster.com/

>刷新 CDN 的目的：因为 hexo 是以生成静态文件部署的，CDN 默认是全部缓存了的，如果有变更需要主动刷新，一般采用目录刷新的方式。

## 快速开始

### 部署 Agent

部署 Agent 需要虚拟机或者轻量应用服务器，支持二进制和 Docker 方式运行

#### 二进制

````shell

export agent_version=v0.2.0

wget https://github.com/stulzq/hexo-deploy-agent/releases/download/$hexo_version/hexo_deploy_agent_$(agent_version)_linux_amd64.tar.gz

tar -xzvf hexo_deploy_agent_$(agent_version)_linux_amd64.tar.gz

cd hexo_deploy_agent_$(agent_version)_linux_amd64

# 修改配置 conf/config.yml

chmod +x hexo_deploy_agent
./hexo_deploy_agent

````

#### Docker

````shell
mkdir -p /data/hexo-deploy-agent/conf

curl https://raw.githubusercontent.com/stulzq/hexo-deploy-agent/main/conf/config.yml -o /data/hexo-deploy-agent/config.yml

# 修改配置 /data/hexo-deploy-agent/config.yml

docker run --name hexo-deploy-agent \
  -v /data/hexo-deploy-agent/conf:/app/conf \
  -v /data/hexo-deploy-agent/logs:/app/logs \
  -d stulzq/hexo-deploy-agent:v0.2.0

````

### 修改配置

````yaml
log:
  level: Debug # 日志级别

deploy:
  blog_dir: /wwwroot/blog # 网站根目录
  cdn:
    enable: false # 是否启用腾讯云 cdn 目录刷新 https://console.cloud.tencent.com/cdn/refresh
    accessKey: # 腾讯云 ak & sk https://console.cloud.tencent.com/cam/capi
    secretKey:
    flushType: flush
    refreshPaths:
      - https://xcmaster.com/ # 刷新路径
  dingtalk:
    enable: false # 是否发送钉钉机器人消息
    url: # 钉钉机器人 url
````

## Github Actions 配置

在你的博客根目录下新建文件夹

````shell
mkdir -p .github/workflows
cd .github/workflows
````

新建配置文件

````shell
touch deploy.yml
````

添加配置

````yaml
name: Deploy

on: 
  push:
    branches:
      - master

jobs:
  build:
    name: build and package
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-node@v3
      with:
        node-version: 16
        registry-url: https://registry.npmjs.org/
        cache: 'npm'
    - name: Install dependencies
      run: npm ci
    - name: Deploy
      run: npm run deploy
    - name: Package
      run: |
          mkdir /home/runner/work/release
          cd public
          zip -r /home/runner/work/release/site.zip ./*
          cd ../
    - name: Upload artifacts
      uses: actions/upload-artifact@v2
      with:
        name: site
        path: /home/runner/work/release
    - name: Clean
      run: |
          rm -rf public
          rm -rf /home/runner/work/release

  publish:
    name: publish blog
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Download build artifacts
        uses: actions/download-artifact@v1
        with:
          name: site
      - name: upload
        env: 
          UPLOAD_URL: ${{ secrets.UPLOAD_URL }}
        run: curl -X POST -F "f=@site/site.zip" $UPLOAD_URL

````

该配置依赖 Github Action Secret

![img.png](img/img.png)

进入项目仓库的 `Settings -> Secrets -> Actions` 新建一个 Secret:

- 名称： `UPLOAD_URL`
- 值：`http://<agentIp>:<agentPort>/deploy/upload` 示例：`http://127.0.0.1:9190/deploy/upload`

> 可以直接使用国内云服务器，POST 部署包速度也是很快的

配置完成！



