# Excel to Markdown Table Converter

A command-line tool to convert Excel table data (TSV/CSV) to Markdown table format. Perfect for quickly converting spreadsheet data for documentation, README files, or any Markdown-based content.

[中文文档](README.zh.md) | [English](README.md)

## ✨ Features

- ✅ **Auto-detect format**: Supports both **TSV** (tab-separated) and **CSV** (comma-separated) formats
- ✅ **Clipboard integration**: Read from and write to clipboard (cross-platform)
- ✅ **Column alignment**: Support left, center, and right alignment via `^l`, `^c`, `^r` markers
- ✅ **Auto column width**: Automatically calculates optimal column widths
- ✅ **CSV handling**: Properly handles quotes, escaping, and fields containing commas
- ✅ **Bilingual support**: Help messages and error messages in both English and Chinese
- ✅ **Cross-platform**: Works on macOS, Linux, and Windows

## 🚀 Installation

### Option 1: Build from source

```bash
git clone https://github.com/lyuangg/excel-to-markdown.git
cd excel-to-markdown
go build -o excel-to-markdown
```

### Option 2: Download pre-built binaries

Download the latest release from the [Releases](https://github.com/lyuangg/excel-to-markdown/releases) page.

## 📖 Usage

### Basic Usage

```bash
# Read from stdin and output to terminal
echo -e "Name\tTitle\nJane\tCEO" | ./excel-to-markdown

# Read from clipboard and convert (auto-write back to clipboard)
./excel-to-markdown -clipboard

# Read from stdin and write to clipboard
cat data.csv | ./excel-to-markdown -copy
```

### Command-line Options

- `-clipboard`: Read data from clipboard (cross-platform support). When used, automatically writes result back to clipboard.
- `-copy`: Write result to clipboard (cross-platform support). Useful when reading from stdin.

### Examples

#### Example 1: Quick Excel Conversion (Most Common)

1. Copy table from Excel (Cmd+C / Ctrl+C)
2. Run `./excel-to-markdown -clipboard`
3. Result is automatically copied to clipboard, paste (Cmd+V / Ctrl+V) into your Markdown editor

#### Example 2: Process Files

```bash
# Process TSV file
cat data.tsv | ./excel-to-markdown > output.md

# Process CSV file and copy to clipboard
cat data.csv | ./excel-to-markdown -copy
```

#### Example 3: Column Alignment

In Excel header row, use alignment markers:

| animal | ^rweight | ^ccolor  |
|--------|----------|----------|
| dog    | 30lb     | tan      |
| cat    | 18lb     | calico   |

After conversion:

```markdown
| animal | weight | color  |
|--------|-------:|:------:|
| dog    | 30lb   | tan    |
| cat    | 18lb   | calico |
```

## 🔧 Supported Formats

### TSV (Tab-Separated Values)

Default format when copying from Excel.

```bash
printf "Name\tTitle\tEmail\nJane\tCEO\tjane@acme.com\n" | ./excel-to-markdown
```

### CSV (Comma-Separated Values)

Supports quoted fields and escaping.

```bash
# Simple CSV
printf "Name,Title,Email\nJane,CEO,jane@acme.com\n" | ./excel-to-markdown

# CSV with quotes (fields containing commas)
printf '"Name","Description","Price"\n"Apple","Red, sweet fruit","$1.50"\n' | ./excel-to-markdown

# Escaped quotes ("" represents a single quote in CSV)
printf '"Name","Quote"\n"John","He said ""Hello"""\n' | ./excel-to-markdown
```

The tool automatically detects the format by checking for quotes, commas, or tabs.

## 🎯 Alignment Markers

- `^l` - Left align (default)
- `^c` - Center align
- `^r` - Right align

Place these markers at the beginning of header cells in your Excel table.

## 🌍 Cross-Platform Clipboard Support

### macOS

- ✅ Built-in support using `pbpaste` and `pbcopy`
- ✅ No additional tools required

### Linux

Requires one of the following clipboard tools:

```bash
# Ubuntu/Debian
sudo apt-get install xclip
# or
sudo apt-get install xsel

# Fedora/CentOS
sudo yum install xclip
# or
sudo yum install xsel
```

The tool automatically detects and uses the available tool (prefers `xclip`).

### Windows

- ✅ Windows 10+ includes `clip.exe` (no installation needed)
- ✅ Older versions use PowerShell (usually pre-installed)

## 📋 Use Cases

1. **Documentation**: Quickly convert spreadsheet data to Markdown tables for README files
2. **Blog Posts**: Convert data tables for blog posts or articles
3. **GitHub Issues**: Format data tables in GitHub issues and pull requests
4. **Scripts**: Integrate into automation workflows and scripts

## 🛠️ Development

### Build

```bash
go build -o excel-to-markdown
```

### Test

```bash
go test -v
```

### Run

```bash
go run . -h  # Show help
```

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## ⭐ Star History

If you find this tool useful, please consider giving it a star ⭐!

## 🔗 Links

- [GitHub Repository](https://github.com/lyuangg/excel-to-markdown)
- [Issues](https://github.com/lyuangg/excel-to-markdown/issues)

---

Made with ❤️ by [lyuangg](https://github.com/lyuangg)
