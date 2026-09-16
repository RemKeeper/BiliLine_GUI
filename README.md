# 哔哩哔哩排队姬 GUI版本

可运行版本，请前往哔哩哔哩饭饭下载

https://play-live.bilibili.com/details/1662382413323



### 使用教程

感谢@幻芯al 制作教程
https://space.bilibili.com/431919235


如有使用问题，请进群 659559309

### 构建与发布

正式发布产物以 **Windows GUI `.exe`** 为主（`go build -ldflags "-w -s -H=windowsgui"`）。

开放平台密钥（`AccessKey` / `AccessSecret` / `AppID`）在源码中保持空占位，**不要写入仓库**。GitHub Actions 在构建时用仓库 Secrets 注入：

- `ACCESSKEY`
- `ACCESSSECRET`
- `APPID`

推送 `v*.*.*` 标签，或在 Actions 里手动 `workflow_dispatch`，会构建 Windows 可执行文件并用 `GITHUB_TOKEN` 创建 GitHub Release（无需再配置个人 `ACCESS_TOKEN`）。Linux amd64 为可选附加产物，失败不会挡住 Windows 发布。

本地自行编译（需 Go 1.22+，Fyne 需要 CGO）：

```bash
go build -ldflags "-w -s"
# Windows GUI 发布构建：
# go build -ldflags "-w -s -H=windowsgui"
```

队列 / 弹幕 HTTP 展示默认先监听 `100` 端口（兼容现有 Windows / OBS 浏览器源）。在 Linux 等非 root 环境若绑定失败，会回退到 `10100`。可用环境变量 `BILILINE_HTTP_PORT` 指定端口。主界面复制的排队 / 弹幕 URL 会使用实际监听端口。音乐插件仍使用 `99` 端口。




![fe32c688d8d4f1b434c61c9643ea57ab_0](https://github.com/RemKeeper/BiliLine_GUI/assets/78486275/d30e6b9e-93b6-460c-842e-9ca837f61d1a)

![6cd47a93e62e4cd1f5f8dcb067c35233_0](https://github.com/RemKeeper/BiliLine_GUI/assets/78486275/d9c4ee2a-6b13-42bf-9cdb-ec79f79634b2)

![c249310d3dc1433b0234cbdff36d4ff4_0](https://github.com/RemKeeper/BiliLine_GUI/assets/78486275/4ce009bc-90dc-45ea-90ed-9a02135a9210)

![257b09711aa677031abcc371e20c7f25_0](https://github.com/RemKeeper/BiliLine_GUI/assets/78486275/c603ee5e-835d-480b-ab3c-570d29dac164)

![5ddd4c17e3052ae3f2a84309144892e1_0](https://github.com/RemKeeper/BiliLine_GUI/assets/78486275/103deb70-3759-4aa3-88a6-99b9a4c6297b)
