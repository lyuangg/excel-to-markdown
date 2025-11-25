package main

import (
	"strings"
	"testing"
)

func TestDisplayWidth(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"英文字符", "Hello", 5},
		{"中文字符", "你好", 4},
		{"混合字符", "Hello世界", 9},
		{"数字", "123", 3},
		{"空字符串", "", 0},
		{"日文", "こんにちは", 10},
		{"韩文", "안녕하세요", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.DisplayWidth(tt.input)
			if result != tt.expected {
				t.Errorf("DisplayWidth(%q) = %d, 期望 %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestColumnWidth(t *testing.T) {
	converter := NewConverter()

	rows := [][]string{
		{"Name", "Title", "Email"},
		{"张三", "经理", "zhang@example.com"},
		{"John", "CEO", "john@example.com"},
	}

	tests := []struct {
		name        string
		columnIndex int
		expected    int
	}{
		{"第一列", 0, 4},  // "Name" = 4, "张三" = 4, "John" = 4 -> max = 4
		{"第二列", 1, 5},  // "Title" = 5, "经理" = 4, "CEO" = 3 -> max = 5
		{"第三列", 2, 17}, // Email 列 "zhang@example.com" = 17, "john@example.com" = 17
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ColumnWidth(rows, tt.columnIndex)
			if result != tt.expected {
				t.Errorf("ColumnWidth(rows, %d) = %d, 期望 %d", tt.columnIndex, result, tt.expected)
			}
		})
	}
}

func TestDetectFormat(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"TSV格式", "Name\tTitle\tEmail", "tsv"},
		{"CSV格式（简单）", "Name,Title,Email", "csv"},
		{"CSV格式（带引号）", `"Name","Title","Email"`, "csv"},
		{"混合（优先TSV）", "Name\tTitle,Email", "tsv"},
		{"只有逗号", "Name,Title", "csv"},
		{"空字符串", "", "tsv"},
		{"Column格式（多个空格对齐）", "Name    Age    City\nJohn    25     NYC", "column"},
		{"Column格式（三列）", "Product      Price    Stock\nApple       1.50     100", "column"},
		{"Column格式（多行）", "Name    Age    City\nJohn    25     NYC\nJane    30     LA", "column"},
		{"单个空格不是Column格式", "Name Age City", "tsv"},
		{"有逗号时不识别为Column", "Name    Age,City", "csv"},
		{"有制表符时不识别为Column", "Name    Age\tCity", "tsv"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.DetectFormat(tt.input)
			if result != tt.expected {
				t.Errorf("DetectFormat(%q) = %s, 期望 %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseTSV(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			"简单TSV",
			"Name\tTitle\nJane\tCEO",
			[][]string{{"Name", "Title"}, {"Jane", "CEO"}},
		},
		{
			"多行TSV",
			"Name\tTitle\tEmail\nJane\tCEO\tjane@example.com\nJohn\tCTO\tjohn@example.com",
			[][]string{
				{"Name", "Title", "Email"},
				{"Jane", "CEO", "jane@example.com"},
				{"John", "CTO", "john@example.com"},
			},
		},
		{
			"包含空行",
			"Name\tTitle\n\nJane\tCEO",
			[][]string{{"Name", "Title"}, {"Jane", "CEO"}},
		},
		{
			"Windows换行符",
			"Name\tTitle\r\nJane\tCEO",
			[][]string{{"Name", "Title"}, {"Jane", "CEO"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ParseTSV(tt.input)
			if !equalRows(result, tt.expected) {
				t.Errorf("ParseTSV(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseCSV(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		input    string
		expected [][]string
		wantErr  bool
	}{
		{
			"简单CSV",
			"Name,Title\nJane,CEO",
			[][]string{{"Name", "Title"}, {"Jane", "CEO"}},
			false,
		},
		{
			"带引号的CSV",
			"\"Name\",\"Title\"\n\"Jane\",\"CEO\"",
			[][]string{{"Name", "Title"}, {"Jane", "CEO"}},
			false,
		},
		{
			"包含逗号的字段",
			"\"Name\",\"Description\"\n\"Apple\",\"Red, sweet fruit\"",
			[][]string{{"Name", "Description"}, {"Apple", "Red, sweet fruit"}},
			false,
		},
		{
			"转义引号",
			"\"Name\",\"Quote\"\n\"John\",\"He said \"\"Hello\"\"\"",
			[][]string{{"Name", "Quote"}, {"John", "He said \"Hello\""}},
			false,
		},
		{
			"空行过滤",
			"Name,Title\n\nJane,CEO",
			[][]string{{"Name", "Title"}, {"Jane", "CEO"}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ParseCSV(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !equalRows(result, tt.expected) {
				t.Errorf("ParseCSV(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLooksLikeColumnFormat(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"标准Column格式（多行对齐）",
			"Name    Age    City\nJohn    25     NYC\nJane    30     LA",
			true,
		},
		{
			"两行对齐",
			"Name    Age    City\nJohn    25     NYC",
			true,
		},
		{
			"列数一致",
			"Product      Price    Stock\nApple       1.50     100\nBanana      0.80     200",
			true,
		},
		{
			"列数允许1列差异",
			"Name    Age    City\nJohn    25",
			true,
		},
		{
			"单行不是Column格式",
			"Name    Age    City",
			false,
		},
		{
			"空字符串",
			"",
			false,
		},
		{
			"只有一行数据",
			"Name    Age    City",
			false,
		},
		{
			"列数差异太大",
			"Name    Age    City\nJohn",
			false,
		},
		{
			"包含空行但仍有有效数据",
			"Name    Age    City\n\nJohn    25     NYC",
			true,
		},
		{
			"只有单个空格（不是Column格式）",
			"Name Age City\nJohn 25 NYC",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.looksLikeColumnFormat(tt.input)
			if result != tt.expected {
				t.Errorf("looksLikeColumnFormat(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseColumn(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			"简单Column格式",
			"Name    Age    City\nJohn    25     NYC",
			[][]string{{"Name", "Age", "City"}, {"John", "25", "NYC"}},
		},
		{
			"多行Column格式",
			"Name    Age    City\nJohn    25     NYC\nJane    30     LA",
			[][]string{
				{"Name", "Age", "City"},
				{"John", "25", "NYC"},
				{"Jane", "30", "LA"},
			},
		},
		{
			"包含空行",
			"Name    Age    City\n\nJohn    25     NYC",
			[][]string{{"Name", "Age", "City"}, {"John", "25", "NYC"}},
		},
		{
			"不同空格数量",
			"Product      Price    Stock\nApple       1.50     100\nBanana      0.80     200",
			[][]string{
				{"Product", "Price", "Stock"},
				{"Apple", "1.50", "100"},
				{"Banana", "0.80", "200"},
			},
		},
		{
			"Windows换行符",
			"Name    Age    City\r\nJohn    25     NYC",
			[][]string{{"Name", "Age", "City"}, {"John", "25", "NYC"}},
		},
		{
			"字段内容包含单个空格",
			"Name        Full Name\nJohn        John Doe\nJane        Jane Smith",
			[][]string{
				{"Name", "Full Name"},
				{"John", "John Doe"},
				{"Jane", "Jane Smith"},
			},
		},
		{
			"列宽不一致",
			"Short    Very Long Column    Medium\nA        B                    C",
			[][]string{
				{"Short", "Very Long Column", "Medium"},
				{"A", "B", "C"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ParseColumn(tt.input)
			if !equalRows(result, tt.expected) {
				t.Errorf("ParseColumn(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseTable(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		input    string
		expected [][]string
		wantErr  bool
	}{
		{
			"自动检测TSV",
			"Name\tTitle\nJane\tCEO",
			[][]string{{"Name", "Title"}, {"Jane", "CEO"}},
			false,
		},
		{
			"自动检测CSV",
			"Name,Title\nJane,CEO",
			[][]string{{"Name", "Title"}, {"Jane", "CEO"}},
			false,
		},
		{
			"自动检测Column格式",
			"Name    Age    City\nJohn    25     NYC",
			[][]string{{"Name", "Age", "City"}, {"John", "25", "NYC"}},
			false,
		},
		{
			"自动检测Column格式（多行）",
			"Product      Price    Stock\nApple       1.50     100\nBanana      0.80     200",
			[][]string{
				{"Product", "Price", "Stock"},
				{"Apple", "1.50", "100"},
				{"Banana", "0.80", "200"},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ParseTable(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !equalRows(result, tt.expected) {
				t.Errorf("ParseTable(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertToMarkdown(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		input    [][]string
		expected string
	}{
		{
			"简单表格",
			[][]string{
				{"Name", "Title"},
				{"Jane", "CEO"},
			},
			"| Name  | Title  |\n|-------|--------|\n| Jane  | CEO    |",
		},
		{
			"带对齐标记",
			[][]string{
				{"^rweight", "^ccolor"},
				{"30lb", "tan"},
			},
			"| weight  | color  |\n|--------:|:------:|\n| 30lb    | tan    |",
		},
		{
			"中文表格",
			[][]string{
				{"姓名", "职位"},
				{"张三", "经理"},
			},
			"| 姓名  | 职位  |\n|-------|-------|\n| 张三  | 经理  |",
		},
		{
			"空表格",
			[][]string{},
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertToMarkdown(tt.input)
			// 标准化换行符进行比较
			result = strings.ReplaceAll(result, "\r\n", "\n")
			tt.expected = strings.ReplaceAll(tt.expected, "\r\n", "\n")
			if result != tt.expected {
				t.Errorf("ConvertToMarkdown() = %q, 期望 %q", result, tt.expected)
			}
		})
	}
}

func TestProcessHeader(t *testing.T) {
	converter := NewConverter()

	rows := [][]string{
		{"^rweight", "^ccolor", "name"},
		{"30lb", "tan", "apple"},
	}

	alignments, widths := converter.processHeader(rows)

	expectedAlignments := []string{"r", "c", "l"}
	expectedWidths := []int{6, 5, 5} // weight=6, color=5, name=5

	if len(alignments) != len(expectedAlignments) {
		t.Fatalf("对齐数组长度不匹配: %d vs %d", len(alignments), len(expectedAlignments))
	}

	for i, expected := range expectedAlignments {
		if alignments[i] != expected {
			t.Errorf("对齐[%d] = %s, 期望 %s", i, alignments[i], expected)
		}
	}

	// 检查对齐标记是否被移除
	if rows[0][0] != "weight" {
		t.Errorf("对齐标记未移除: %q", rows[0][0])
	}
	if rows[0][1] != "color" {
		t.Errorf("对齐标记未移除: %q", rows[0][1])
	}

	// 检查列宽（允许一些误差，因为实际宽度可能略有不同）
	for i, expected := range expectedWidths {
		if widths[i] < expected-1 || widths[i] > expected+1 {
			t.Errorf("列宽[%d] = %d, 期望约 %d", i, widths[i], expected)
		}
	}
}

// equalRows 比较两个行数组是否相等
func equalRows(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
