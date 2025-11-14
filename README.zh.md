# Excel 转 Markdown 表格工具

一个命令行工具，用于将 Excel 表格数据（TSV/CSV）转换为 Markdown 表格格式。非常适合快速将电子表格数据转换为文档、README 文件或任何基于 Markdown 的内容。

[English](README.md) | [中文文档](README.zh.md)

## ✨ 功能特性

- ✅ **自动检测格式**：支持 **TSV**（制表符分隔）和 **CSV**（逗号分隔）格式
- ✅ **剪贴板集成**：支持从剪贴板读取和写入（跨平台）
- ✅ **列对齐**：通过 `^l`、`^c`、`^r` 标记支持左对齐、居中、右对齐
- ✅ **自动列宽**：自动计算最佳列宽
- ✅ **CSV 处理**：正确处理引号、转义和包含逗号的字段
- ✅ **双语支持**：帮助信息和错误信息支持中英文
- ✅ **跨平台**：支持 macOS、Linux 和 Windows

## 🚀 安装

### 方式一：从源码编译

```bash
git clone https://github.com/lyuangg/excel-to-markdown.git
cd excel-to-markdown
go build -o excel-to-markdown
```

### 方式二：下载预编译二进制文件

从 [Releases](https://github.com/lyuangg/excel-to-markdown/releases) 页面下载最新版本。

#### macOS 安全提示

如果在 macOS 上运行下载的二进制文件时看到安全警告：

1. **快速解决**：移除隔离属性：
   ```bash
   xattr -d com.apple.quarantine excel-to-markdown
   ```

2. **替代方法**：右键点击文件 → 打开 → 在安全对话框中点击"打开"（仅首次需要）

这是 macOS 的安全机制（Gatekeeper），适用于所有未签名的二进制文件。该二进制文件是安全的。

## 📖 使用方法

### 基本用法

```bash
# 从标准输入读取并输出到终端
echo -e "Name\tTitle\nJane\tCEO" | ./excel-to-markdown

# 从剪贴板读取并转换（自动写回剪贴板）
./excel-to-markdown -clipboard

# 从标准输入读取并写入剪贴板
cat data.csv | ./excel-to-markdown -copy
```

### 命令行选项

- `-clipboard`: 从剪贴板读取数据（跨平台支持）。使用此选项时会自动将结果写回剪贴板。
- `-copy`: 将结果写入剪贴板（跨平台支持）。适用于从标准输入读取时写入剪贴板。

### 示例

#### 示例 1：快速 Excel 转换（最常用）

1. 在 Excel 中复制表格（Cmd+C / Ctrl+C）
2. 运行 `./excel-to-markdown -clipboard`
3. 结果已自动复制到剪贴板，直接粘贴（Cmd+V / Ctrl+V）到 Markdown 编辑器即可

#### 示例 2：处理文件

```bash
# 处理 TSV 文件
cat data.tsv | ./excel-to-markdown > output.md

# 处理 CSV 文件并复制到剪贴板
cat data.csv | ./excel-to-markdown -copy
```

#### 示例 3：列对齐

在 Excel 表头行中使用对齐标记：

| animal | ^rweight | ^ccolor  |
|--------|----------|----------|
| dog    | 30lb     | tan      |
| cat    | 18lb     | calico   |

转换后：

```markdown
| animal | weight | color  |
|--------|-------:|:------:|
| dog    | 30lb   | tan    |
| cat    | 18lb   | calico |
```

## 🔧 支持的格式

### TSV（制表符分隔值）

从 Excel 复制时的默认格式。

```bash
printf "Name\tTitle\tEmail\nJane\tCEO\tjane@acme.com\n" | ./excel-to-markdown
```

### CSV（逗号分隔值）

支持引号包围的字段和转义。

```bash
# 简单 CSV
printf "Name,Title,Email\nJane,CEO,jane@acme.com\n" | ./excel-to-markdown

# 带引号的 CSV（包含逗号的字段）
printf '"Name","Description","Price"\n"Apple","Red, sweet fruit","$1.50"\n' | ./excel-to-markdown

# 转义引号（CSV 中使用 "" 表示一个引号）
printf '"Name","Quote"\n"John","He said ""Hello"""\n' | ./excel-to-markdown
```

工具会自动检测输入格式（通过检查是否包含引号、逗号或制表符）。

## 🎯 对齐标记说明

- `^l` - 左对齐（默认）
- `^c` - 居中对齐
- `^r` - 右对齐

在 Excel 表格的表头单元格开头使用这些标记。

## 🌍 跨平台剪贴板支持

### macOS

- ✅ 使用系统自带的 `pbpaste` 和 `pbcopy`
- ✅ 无需安装额外工具

### Linux

需要安装以下剪贴板工具之一：

```bash
# Ubuntu/Debian
sudo apt-get install xclip
# 或
sudo apt-get install xsel

# Fedora/CentOS
sudo yum install xclip
# 或
sudo yum install xsel
```

工具会自动检测并使用可用的工具（优先使用 `xclip`）。

### Windows

- ✅ Windows 10+ 自带 `clip.exe`（无需安装）
- ✅ 较旧版本使用 PowerShell（通常已预装）

## 📋 使用场景

1. **文档编写**：快速将电子表格数据转换为 Markdown 表格，用于 README 文件
2. **博客文章**：为博客文章或文章转换数据表格
3. **GitHub Issues**：在 GitHub Issues 和 Pull Requests 中格式化数据表格
4. **脚本集成**：集成到自动化工作流和脚本中

## 🛠️ 开发

### 编译

```bash
go build -o excel-to-markdown
```

### 测试

```bash
go test -v
```

### 运行

```bash
go run . -h  # 显示帮助信息
```

## 📝 许可证

MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎贡献！请随时提交 Pull Request。

## ⭐ Star 历史

如果你觉得这个工具有用，请考虑给它一个星标 ⭐！

## 🔗 链接

- [GitHub 仓库](https://github.com/lyuangg/excel-to-markdown)
- [问题反馈](https://github.com/lyuangg/excel-to-markdown/issues)

---

由 [lyuangg](https://github.com/lyuangg) 用 ❤️ 制作

