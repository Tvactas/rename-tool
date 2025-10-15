# Renamer (文件批量重命名工具)

由 Go 编写的 Windows 平台下文件批量重命名工具，支持富有的文件名处理功能和多语系界面（中文 / English / 日本語）。

## ✨ 功能特性

### 📁 基础重命名

* 添加前缀 / 后缀
* 数字编号（支持格式内单独编号）
* 保留原文件名
* 修改扩展名

### 💡 大小写转换

* 全部大写
* 全部小写
* 标题格式（Title Case）
* 驼锋格式（CamelCase）

### 🔤 字符操作

* 插入字符
* 删除字符
* 正则替换

### 🛠 其他功能

* 支持撤销操作
* 操作日志记录
* 实时预览
* 多语系界面（中文 / 英文 / 日文）
* 检测文件是否被占用

---

## 🧱 环境要求

* Go 1.22 或更高版本
* Windows 10 或以上

---

## 🚀 快速开始

### 1. 安装 Go

请前往 [Go 官网](https://go.dev/dl/) 下载并安装 Go 1.22 或更新版本。

验证安装：

```bash
go version
```

### 2. 安装 GCC

请前往 [GCC 官网](https://jmeubank.github.io/tdm-gcc/articles/2021-05/10.3.0-release/) 下载并安装 GCC对应系统版本。（需要勾选 ‘add to path'）



### 3. 获取项目依赖
如果在国内可能会出现超时现象,所以需要执行,将下载库源头改为阿里库

```
go env -w GOPROXY=https://goproxy.cn,direct
```

在项目根目录执行：

```bash
go mod tidy
```

### 4. 构建 EXE

使用以下命令生成无控制台窗口的 Windows 可执行文件（如需固定图标请加入 `.syso` 文件）：

```bash
go build -ldflags="-H windowsgui -s -w" -o renamer.exe
```

### 5. 使用 UPX 压缩（可选）

前往 [UPX Releases](https://github.com/upx/upx/releases/tag/v5.0.1) 下载适用于系统的压缩包，如：

* `upx-5.0.1-win64.zip`

解压后假设放在 `D:/upx-5.0.1-win64/`，执行以下命令：

```bash
D:/upx-5.0.1-win64/upx-5.0.1-win64/upx --best --lzma renamer.exe
```

---

## 📅 许可协议

MIT License

---

## 👤 作者

**Tvacats**
GitHub: [@Tvactas](https://github.com/Tvactas)
Gmail: tvacats@gmail.com
