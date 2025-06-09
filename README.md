# 文件重命名工具 (File Rename Tool)
可直接下载renamer.exe进行使用，也可自行打包
一个功能强大的文件重命名工具，支持多种重命名方式和多语言界面。

Power by Tvacats

## 功能特点

- 批量重命名
  - 支持添加前缀和后缀
  - 支持数字编号
  - 支持保留原文件名
  - 支持按格式单独编号

- 扩展名修改
  - 支持修改文件扩展名
  - 支持批量处理

- 大小写转换
  - 大写转换
  - 小写转换
  - 标题格式转换
  - 驼峰格式转换

- 字符操作
  - 插入字符
  - 删除字符
  - 正则替换

- 其他功能
  - 撤销操作
  - 操作日志
  - 多语言支持（中文、英文、日文）
  - 文件占用检测
  - 实时预览

## 系统要求

- Windows 10 或更高版本
- 至少 100MB 可用磁盘空间
- 最小分辨率 800x600
- Go 1.20 或更高版本（开发环境）
- Fyne 2.4.0 或更高版本
- 安装 handle.exe（用于文件占用检测）

## 开发环境配置

1. 安装 Go 环境
   ```bash
   # 下载并安装 Go
   https://golang.org/dl/
   ```

2. 安装 Fyne 工具
   ```bash
   go install fyne.io/fyne/v2/cmd/fyne@latest
   ```

3. 安装 handle.exe
   - 下载 [Handle](https://learn.microsoft.com/en-us/sysinternals/downloads/handle)
   - 将 handle.exe 放入系统 PATH 目录或程序目录

4. 安装依赖
   ```bash
   go mod tidy
   ```

## 打包说明

1. 准备资源文件
   - 将字体文件放入 `src/font/` 目录
   - 将图片文件放入 `src/img/` 目录

2. 编译打包
   ```bash
   # Windows
   fyne package -os windows -icon icon.png -name "Rename Tool" -appID com.yourdomain.renametool

   # 生成安装包
   fyne package -os windows -icon icon.png -name "Rename Tool" -appID com.yourdomain.renametool -release
   ```

3. 打包后的文件结构
   ```
   rename-tool/
   ├── rename-tool.exe
   ├── handle.exe
   ├── src/
   │   ├── font/
   │   │   ├── JP.TTF
   │   │   ├── TIMES.TTF
   │   │   └── STXINGKA.TTF
   │   └── img/
   │       └── cat.png
   └── README.md
   ```

## 安装说明

1. 下载最新版本的发布包
2. 解压到任意目录
3. 确保 handle.exe 在程序目录或系统 PATH 中
4. 运行 `rename-tool.exe`

## 使用说明

### 批量重命名
1. 点击"批量重命名"按钮
2. 选择目标目录
3. 扫描文件格式
4. 设置前缀、后缀和编号
5. 预览效果
6. 执行重命名

### 扩展名修改
1. 点击"修改扩展名"按钮
2. 选择目标目录
3. 选择要修改的文件格式
4. 输入新的扩展名
5. 预览效果
6. 执行修改

### 大小写转换
1. 选择转换类型（大写/小写/标题/驼峰）
2. 选择目标目录
3. 预览效果
4. 执行转换

### 字符操作
1. 选择操作类型（插入/删除/替换）
2. 选择目标目录
3. 设置操作参数
4. 预览效果
5. 执行操作

## 注意事项

- 重命名操作前请先预览
- 建议先备份重要文件
- 如遇到文件占用，可以使用内置的进程管理功能
- 所有操作都可以通过"撤销"功能恢复
- 确保 handle.exe 可用，否则文件占用检测功能将无法使用

## 常见问题

Q: 为什么某些文件无法重命名？
A: 可能是文件被其他程序占用，请关闭相关程序后重试。如果问题持续，请确保 handle.exe 已正确安装。

Q: 如何恢复误操作？
A: 使用"撤销"功能可以恢复最近的操作。

Q: 如何查看操作历史？
A: 使用"日志"功能可以查看和导出操作历史。

Q: 程序无法启动怎么办？
A: 请检查是否已安装所有必要的依赖，特别是 handle.exe 是否在正确的位置。

## 技术支持

如有问题或建议，请提交 Issue 或 Pull Request。

## 许可证

MIT License

## 作者

Tvacats 