package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// looksLikeTable 检查数据是否看起来像表格
// 简单检查：如果包含制表符或逗号，可能是表格数据
func looksLikeTable(data string) bool {
	return strings.Contains(data, "\t") || strings.Contains(data, ",")
}

// getLanguage detects the preferred language (zh or en)
func getLanguage() string {
	// Check environment variable first
	lang := os.Getenv("LANG")
	if lang != "" {
		lang = strings.ToLower(lang)
		if strings.Contains(lang, "zh") || strings.Contains(lang, "cn") {
			return "zh"
		}
	}

	// Check LC_ALL
	lcAll := os.Getenv("LC_ALL")
	if lcAll != "" {
		lcAll = strings.ToLower(lcAll)
		if strings.Contains(lcAll, "zh") || strings.Contains(lcAll, "cn") {
			return "zh"
		}
	}

	// Default to English
	return "en"
}

// errorMsg returns error message based on language
func errorMsg(lang, zhMsg, enMsg string) string {
	if lang == "zh" {
		return zhMsg
	}
	return enMsg
}

// printError prints error message and exits
func printError(lang, zhMsg, enMsg string) {
	fmt.Fprintf(os.Stderr, "%s\n", errorMsg(lang, zhMsg, enMsg))
	os.Exit(1)
}

// printErrorf prints formatted error message and exits
func printErrorf(lang, zhMsg, enMsg string, args ...interface{}) {
	msg := errorMsg(lang, zhMsg, enMsg)
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// readInput reads input from clipboard or stdin
func readInput(fromClipboard bool, lang string) string {
	if fromClipboard {
		input, err := readFromClipboard()
		if err != nil {
			printError(lang,
				"无法从剪贴板读取: "+err.Error()+"\n请使用标准输入或安装剪贴板工具",
				"Failed to read from clipboard: "+err.Error()+"\nPlease use stdin or install clipboard tools")
		}
		return input
	}

	// Read from stdin
	// Go's slice will automatically grow as needed, starting with a small capacity
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		printErrorf(lang, "读取输入时出错: %v", "Error reading input: %v", err)
	}
	return strings.Join(lines, "\n")
}

// validateInput validates that input is not empty
func validateInput(input string, lang string) {
	if input == "" {
		printError(lang,
			"错误: 没有输入数据\n用法: echo '表格数据' | "+os.Args[0]+"\n或者: "+os.Args[0]+" -clipboard",
			"Error: No input data\nUsage: echo 'table data' | "+os.Args[0]+"\nOr: "+os.Args[0]+" -clipboard")
	}
}

// convertTable converts input table data to markdown
func convertTable(input string, lang string) string {
	if !looksLikeTable(input) {
		printError(lang, "输入数据不是表格格式", "Input data is not in table format")
	}

	converter := NewConverter()
	rows, err := converter.ParseTable(input)
	if err != nil {
		printErrorf(lang, "错误: 解析表格数据失败: %v", "Error: Failed to parse table data: %v", err)
	}
	if len(rows) == 0 {
		printError(lang, "错误: 无法解析表格数据", "Error: Unable to parse table data")
	}

	return converter.ConvertToMarkdown(rows)
}

// outputResult outputs markdown to clipboard or stdout
func outputResult(markdown string, shouldCopy, fromClipboard bool, lang string) {
	if !shouldCopy {
		fmt.Println(markdown)
		return
	}

	err := writeToClipboard(markdown)
	if err != nil {
		// Even if write fails, output to stdout
		printErrorf(lang, "无法写入剪贴板: %v", "Failed to write to clipboard: %v", err)
		fmt.Println(markdown)
		os.Exit(1)
	}

	// Success message
	successMsg := errorMsg(lang, "✓ Markdown 表格已复制到剪贴板", "✓ Markdown table copied to clipboard")
	fmt.Fprintf(os.Stderr, "%s\n", successMsg)

	// If read from clipboard, also output to stdout for viewing
	if fromClipboard {
		fmt.Println(markdown)
	}
}

// setupUsage sets up custom usage information with bilingual support
func setupUsage() {
	lang := getLanguage()

	flag.Usage = func() {
		if lang == "zh" {
			// Chinese help
			fmt.Fprintf(os.Stderr, "用法: %s [选项]\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "将 CSV 或 TSV 格式的表格数据转换为 Markdown 表格格式。\n\n")
			fmt.Fprintf(os.Stderr, "选项:\n")
			flag.PrintDefaults()
			fmt.Fprintf(os.Stderr, "\n示例:\n")
			fmt.Fprintf(os.Stderr, "  # 从标准输入读取并输出到终端\n")
			fmt.Fprintf(os.Stderr, "  echo -e \"Name\\tTitle\\nJane\\tCEO\" | %s\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "  # 从剪贴板读取并转换（自动写回剪贴板）\n")
			fmt.Fprintf(os.Stderr, "  %s -clipboard\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "  # 从标准输入读取并写入剪贴板\n")
			fmt.Fprintf(os.Stderr, "  cat data.csv | %s -copy\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "支持的格式:\n")
			fmt.Fprintf(os.Stderr, "  - TSV (制表符分隔): Excel 复制时的默认格式\n")
			fmt.Fprintf(os.Stderr, "  - CSV (逗号分隔): 自动检测格式\n\n")
			fmt.Fprintf(os.Stderr, "列对齐标记:\n")
			fmt.Fprintf(os.Stderr, "  在表头使用 ^l (左对齐), ^c (居中), ^r (右对齐)\n")
			fmt.Fprintf(os.Stderr, "  例如: \"^r价格\" 表示右对齐的价格列\n\n")
			fmt.Fprintf(os.Stderr, "更多信息请查看: https://github.com/lyuangg/excel-to-markdown\n")
		} else {
			// English help
			fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "Convert CSV or TSV table data to Markdown table format.\n\n")
			fmt.Fprintf(os.Stderr, "Options:\n")
			flag.PrintDefaults()
			fmt.Fprintf(os.Stderr, "\nExamples:\n")
			fmt.Fprintf(os.Stderr, "  # Read from stdin and output to terminal\n")
			fmt.Fprintf(os.Stderr, "  echo -e \"Name\\tTitle\\nJane\\tCEO\" | %s\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "  # Read from clipboard and convert (automatically write back to clipboard)\n")
			fmt.Fprintf(os.Stderr, "  %s -clipboard\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "  # Read from stdin and write to clipboard\n")
			fmt.Fprintf(os.Stderr, "  cat data.csv | %s -copy\n\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "Supported formats:\n")
			fmt.Fprintf(os.Stderr, "  - TSV (tab-separated): Default format when copying from Excel\n")
			fmt.Fprintf(os.Stderr, "  - CSV (comma-separated): Auto-detected format\n\n")
			fmt.Fprintf(os.Stderr, "Column alignment markers:\n")
			fmt.Fprintf(os.Stderr, "  Use ^l (left), ^c (center), ^r (right) in header row\n")
			fmt.Fprintf(os.Stderr, "  Example: \"^rPrice\" for right-aligned price column\n\n")
			fmt.Fprintf(os.Stderr, "For more information, visit: https://github.com/lyuangg/excel-to-markdown\n")
		}
	}
}

func main() {
	lang := getLanguage()

	// Set flag descriptions based on language
	clipboardDesc := errorMsg(lang, "从剪贴板读取数据（跨平台支持）", "Read data from clipboard (cross-platform support)")
	copyDesc := errorMsg(lang, "将结果写入剪贴板（跨平台支持）", "Write result to clipboard (cross-platform support)")

	fromClipboard := flag.Bool("clipboard", false, clipboardDesc)
	toClipboard := flag.Bool("copy", false, copyDesc)
	setupUsage()
	flag.Parse()

	// Read and validate input
	input := readInput(*fromClipboard, lang)
	validateInput(input, lang)

	// Convert table to markdown
	markdown := convertTable(input, lang)

	// Output result
	shouldCopy := *toClipboard || *fromClipboard
	outputResult(markdown, shouldCopy, *fromClipboard, lang)
}

// getLinuxClipboardCommand returns the appropriate Linux clipboard command
func getLinuxClipboardCommand(read bool) (*exec.Cmd, error) {
	if _, err := exec.LookPath("xclip"); err == nil {
		if read {
			return exec.Command("xclip", "-selection", "clipboard", "-out"), nil
		}
		return exec.Command("xclip", "-selection", "clipboard"), nil
	}
	if _, err := exec.LookPath("xsel"); err == nil {
		if read {
			return exec.Command("xsel", "--clipboard", "--output"), nil
		}
		return exec.Command("xsel", "--clipboard", "--input"), nil
	}
	lang := getLanguage()
	if lang == "zh" {
		return nil, fmt.Errorf("需要安装 xclip 或 xsel: sudo apt-get install xclip 或 sudo apt-get install xsel")
	}
	return nil, fmt.Errorf("xclip or xsel required: sudo apt-get install xclip or sudo apt-get install xsel")
}

// getUnsupportedOSError returns error message for unsupported OS
func getUnsupportedOSError() error {
	lang := getLanguage()
	if lang == "zh" {
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
	return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
}

// readFromClipboard 从剪贴板读取数据（跨平台实现）
func readFromClipboard() (string, error) {
	var cmd *exec.Cmd
	var err error

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbpaste")
	case "linux":
		cmd, err = getLinuxClipboardCommand(true)
		if err != nil {
			return "", err
		}
	case "windows":
		cmd = exec.Command("powershell", "-Command", "Get-Clipboard")
	default:
		return "", getUnsupportedOSError()
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// writeToClipboard 将数据写入剪贴板（跨平台实现）
func writeToClipboard(data string) error {
	var cmd *exec.Cmd
	var err error

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(data)
	case "linux":
		cmd, err = getLinuxClipboardCommand(false)
		if err != nil {
			return err
		}
		cmd.Stdin = strings.NewReader(data)
	case "windows":
		if _, err := exec.LookPath("clip"); err == nil {
			cmd = exec.Command("clip")
			cmd.Stdin = strings.NewReader(data)
		} else {
			cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Set-Clipboard -Value %q", data))
		}
	default:
		return getUnsupportedOSError()
	}

	return cmd.Run()
}
