package main

import (
	"encoding/csv"
	"regexp"
	"strings"
	"unicode"
)

var (
	// alignmentRegex matches alignment markers in header: ^l, ^c, ^r
	alignmentRegex = regexp.MustCompile(`(?i)^(\^[lcr])`)
)

// Converter 表格转换器
type Converter struct{}

// NewConverter 创建新的转换器实例
func NewConverter() *Converter {
	return &Converter{}
}

// DisplayWidth 计算字符串的显示宽度（考虑中文字符）
// 中文字符通常占用 2 个显示宽度，英文字符占用 1 个
func (c *Converter) DisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		// 判断是否为全角字符（包括中文、日文、韩文等）
		// 全角字符通常占用 2 个显示宽度
		if unicode.Is(unicode.Han, r) || // 中文字符
			unicode.Is(unicode.Hiragana, r) || // 日文平假名
			unicode.Is(unicode.Katakana, r) || // 日文片假名
			unicode.Is(unicode.Hangul, r) || // 韩文字符
			(r >= 0xFF00 && r <= 0xFFEF) { // 全角符号
			width += 2
		} else {
			// 半角字符占用 1 个显示宽度
			width += 1
		}
	}
	return width
}

// ColumnWidth 计算指定列的最大显示宽度
func (c *Converter) ColumnWidth(rows [][]string, columnIndex int) int {
	maxWidth := 0
	for _, row := range rows {
		if columnIndex < len(row) {
			cellWidth := c.DisplayWidth(row[columnIndex])
			if cellWidth > maxWidth {
				maxWidth = cellWidth
			}
		}
	}
	return maxWidth
}

// DetectFormat 检测输入格式是 CSV、TSV 还是 Column（空格对齐）
func (c *Converter) DetectFormat(data string) string {
	// 检查是否包含引号（CSV 的特征）
	hasQuotes := strings.Contains(data, `"`)
	// 检查是否包含逗号
	hasCommas := strings.Contains(data, ",")
	// 检查是否包含制表符
	hasTabs := strings.Contains(data, "\t")
	// 检查是否包含多个连续空格（column 命令对齐格式）
	hasMultipleSpaces := strings.Contains(data, "  ")

	// 如果包含引号，很可能是 CSV
	if hasQuotes {
		return "csv"
	}
	// 如果包含制表符，优先认为是 TSV
	if hasTabs {
		return "tsv"
	}
	// 如果包含逗号但不包含制表符，可能是 CSV
	if hasCommas && !hasTabs {
		return "csv"
	}
	// 如果包含多个连续空格且没有制表符和逗号，可能是 column 格式
	if hasMultipleSpaces && !hasTabs && !hasCommas {
		// 进一步验证：检查是否有多行且每行都有多个空格对齐
		if c.looksLikeColumnFormat(data) {
			return "column"
		}
	}
	// 如果只有逗号，默认认为是 CSV
	if hasCommas {
		return "csv"
	}
	// 默认返回 TSV（保持向后兼容）
	return "tsv"
}

// looksLikeColumnFormat 检查数据是否像 column 命令对齐的格式
// column 格式特征：多行数据，每行都有多个连续空格，且列数大致一致
func (c *Converter) looksLikeColumnFormat(data string) bool {
	data = normalizeLineEndings(data)
	lines := strings.Split(data, "\n")

	// 至少需要2行才能形成表格
	if len(lines) < 2 {
		return false
	}

	// 统计每行的列数（通过多个连续空格分割）
	columnCounts := make([]int, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 使用正则表达式分割多个连续空格
		fields := regexp.MustCompile(`\s{2,}`).Split(line, -1)
		// 过滤空字段
		nonEmptyFields := 0
		for _, field := range fields {
			if strings.TrimSpace(field) != "" {
				nonEmptyFields++
			}
		}
		if nonEmptyFields > 1 {
			columnCounts = append(columnCounts, nonEmptyFields)
		}
	}

	// 如果有多行且列数一致（或大致一致），可能是 column 格式
	if len(columnCounts) >= 2 {
		// 检查列数是否大致一致（允许1列的差异）
		firstCount := columnCounts[0]
		consistent := true
		for _, count := range columnCounts[1:] {
			if count != firstCount && count != firstCount-1 && count != firstCount+1 {
				consistent = false
				break
			}
		}
		return consistent
	}

	return false
}

// ParseCSV 解析 CSV 格式的表格数据
func (c *Converter) ParseCSV(data string) ([][]string, error) {
	// 处理各种换行符
	data = normalizeLineEndings(data)

	reader := csv.NewReader(strings.NewReader(data))
	// 允许字段数量不一致
	reader.FieldsPerRecord = -1
	// 使用逗号作为分隔符
	reader.Comma = ','

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// 过滤空行
	return c.filterEmptyRows(rows), nil
}

// ParseTSV 解析 TSV 格式的表格数据
func (c *Converter) ParseTSV(data string) [][]string {
	data = strings.TrimSpace(data)
	// 处理各种换行符
	data = normalizeLineEndings(data)

	// 按换行符分割成行
	lines := strings.Split(data, "\n")

	var rows [][]string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue // 跳过空行
		}
		// 按制表符分割
		cells := strings.Split(line, "\t")
		rows = append(rows, cells)
	}
	return rows
}

// ParseColumn 解析 column 命令对齐格式的表格数据（使用多个空格对齐）
func (c *Converter) ParseColumn(data string) [][]string {
	data = strings.TrimSpace(data)
	// 处理各种换行符
	data = normalizeLineEndings(data)

	// 按换行符分割成行
	lines := strings.Split(data, "\n")

	// 首先分析所有行，找出列的分割位置
	// 通过分析每行中多个连续空格的位置来确定列边界
	var rows [][]string
	for _, line := range lines {
		line = strings.TrimRight(line, " \t") // 保留行尾空格用于对齐分析
		if strings.TrimSpace(line) == "" {
			continue // 跳过空行
		}

		// 使用正则表达式分割多个连续空格（2个或更多）
		// 但需要保留单个空格在字段内容中
		// 策略：找到所有连续2个或更多空格的位置，作为列分隔符
		fields := regexp.MustCompile(`\s{2,}`).Split(line, -1)

		// 清理每个字段的前后空格
		var cleanedFields []string
		for _, field := range fields {
			cleaned := strings.TrimSpace(field)
			if cleaned != "" || len(cleanedFields) == 0 {
				// 保留第一个字段即使为空（可能是对齐导致的）
				cleanedFields = append(cleanedFields, cleaned)
			}
		}

		// 如果分割后只有一个字段，尝试用单个空格分割（可能是单空格分隔的简单表格）
		if len(cleanedFields) == 1 && strings.Contains(line, " ") {
			// 回退到单空格分割
			fields = strings.Fields(line)
			cleanedFields = fields
		}

		if len(cleanedFields) > 0 {
			rows = append(rows, cleanedFields)
		}
	}

	return rows
}

// ParseTable 解析表格数据，自动检测格式（CSV、TSV 或 Column）
func (c *Converter) ParseTable(data string) ([][]string, error) {
	format := c.DetectFormat(data)

	switch format {
	case "csv":
		return c.ParseCSV(data)
	case "column":
		rows := c.ParseColumn(data)
		return rows, nil
	default:
		// TSV 格式（默认）
		rows := c.ParseTSV(data)
		return rows, nil
	}
}

// ConvertToMarkdown 将表格数据转换为 Markdown 格式
func (c *Converter) ConvertToMarkdown(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	colAlignments, columnWidths := c.processHeader(rows)

	// 生成 Markdown 行
	var markdownRows []string

	// 生成表头行
	markdownRows = append(markdownRows, c.generateHeaderRow(rows[0], columnWidths))

	// 生成分隔行
	markdownRows = append(markdownRows, c.generateSeparatorRow(columnWidths, colAlignments))

	// 生成数据行
	for i := 1; i < len(rows); i++ {
		markdownRows = append(markdownRows, c.generateDataRow(rows[i], columnWidths))
	}

	return strings.Join(markdownRows, "\n")
}

// processHeader 处理表头，提取对齐信息和计算列宽
func (c *Converter) processHeader(rows [][]string) ([]string, []int) {
	colAlignments := make([]string, len(rows[0]))
	columnWidths := make([]int, len(rows[0]))

	// 处理表头，提取对齐信息
	for i := range rows[0] {
		alignment := "l" // 默认左对齐
		column := rows[0][i]

		matches := alignmentRegex.FindStringSubmatch(column)
		if len(matches) > 0 {
			align := strings.ToLower(string(matches[1][1]))
			switch align {
			case "c":
				alignment = "c"
			case "r":
				alignment = "r"
			}
			// 移除对齐标记
			rows[0][i] = alignmentRegex.ReplaceAllString(column, "")
		}
		colAlignments[i] = alignment
		columnWidths[i] = c.ColumnWidth(rows, i)
	}

	return colAlignments, columnWidths
}

// generateHeaderRow 生成表头行
func (c *Converter) generateHeaderRow(header []string, columnWidths []int) string {
	var cells []string
	for i, cell := range header {
		cellDisplayWidth := c.DisplayWidth(cell)
		padding := strings.Repeat(" ", columnWidths[i]-cellDisplayWidth+1)
		cells = append(cells, cell+padding)
	}
	return "| " + strings.Join(cells, " | ") + " |"
}

// generateSeparatorRow 生成分隔行
func (c *Converter) generateSeparatorRow(columnWidths []int, colAlignments []string) string {
	var cells []string
	for i, width := range columnWidths {
		prefix := ""
		postfix := ""
		adjust := 0
		alignment := colAlignments[i]

		switch alignment {
		case "r":
			postfix = ":"
			adjust = 1
		case "c":
			prefix = ":"
			postfix = ":"
			adjust = 2
		}
		dashes := strings.Repeat("-", width+3-adjust)
		cells = append(cells, prefix+dashes+postfix)
	}
	return "|" + strings.Join(cells, "|") + "|"
}

// generateDataRow 生成数据行
func (c *Converter) generateDataRow(row []string, columnWidths []int) string {
	var cells []string
	for j, cell := range row {
		if j < len(columnWidths) {
			cellDisplayWidth := c.DisplayWidth(cell)
			padding := strings.Repeat(" ", columnWidths[j]-cellDisplayWidth+1)
			cells = append(cells, cell+padding)
		} else {
			cells = append(cells, cell)
		}
	}
	return "| " + strings.Join(cells, " | ") + " |"
}

// filterEmptyRows 过滤空行
func (c *Converter) filterEmptyRows(rows [][]string) [][]string {
	var filteredRows [][]string
	for _, row := range rows {
		// 检查是否所有字段都为空
		allEmpty := true
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				allEmpty = false
				break
			}
		}
		if !allEmpty {
			filteredRows = append(filteredRows, row)
		}
	}
	return filteredRows
}

// normalizeLineEndings 标准化换行符
func normalizeLineEndings(data string) string {
	data = strings.ReplaceAll(data, "\r\n", "\n")   // Windows 换行符
	data = strings.ReplaceAll(data, "\r", "\n")     // Mac 旧式换行符
	data = strings.ReplaceAll(data, "\u0085", "\n") // NEXT LINE
	data = strings.ReplaceAll(data, "\u2028", "\n") // LINE SEPARATOR
	data = strings.ReplaceAll(data, "\u2029", "\n") // PARAGRAPH SEPARATOR
	return data
}
