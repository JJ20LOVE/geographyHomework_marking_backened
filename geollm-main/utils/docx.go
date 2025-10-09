package utils

import (
	"baliance.com/gooxml/document"
	"fmt"
	"regexp"
	"strings"
)

// 定义小题结构体
type Question struct {
	Number  string `json:"Number"`  // 小题号，只包含数字
	Content string `json:"Content"` // 小题内容
}

// 定义大题结构体
type Section struct {
	Title     string     `json:"Title"`     // 大题题干
	Questions []Question `json:"Questions"` // 该大题下的小题
}

func Extractor(filename string) ([]Section, error) {
	doc, err := document.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening document: %w", err)
	}

	// 使用原始字符串字面量定义正则表达式
	sectionRegex := regexp.MustCompile(`^第?[一二三四五六七八九十百零〇]+[、.]?`) // 匹配 "一、"、"二、" 等大题号
	questionRegex := regexp.MustCompile(`^\d+[.、．\s]?`)            // 匹配 "1."、"1、" 等
	choiceRegex := regexp.MustCompile(`^[A-D][.．\s]`)              // 匹配 "A."、"A．"

	// 用于提取纯数字的小题号
	numberOnlyRegex := regexp.MustCompile(`\d+`)

	// 存储大题和小题的结构
	var sections []Section
	var currentSection *Section
	var currentQuestion *Question

	for _, para := range doc.Paragraphs() {
		var paragraphText strings.Builder

		// 拼接段落中的文字
		for _, run := range para.Runs() {
			text := strings.TrimSpace(run.Text())
			if text != "" {
				paragraphText.WriteString(text)
			}
		}

		text := paragraphText.String()
		if text == "" {
			continue // 跳过空行
		}

		// 处理大题号
		if sectionRegex.MatchString(text) {
			// 保存上一个大题
			if currentSection != nil {
				// 保存上一个小题
				if currentQuestion != nil {
					currentSection.Questions = append(currentSection.Questions, *currentQuestion)
					currentQuestion = nil
				}
				sections = append(sections, *currentSection)
			}
			// 创建新的大题
			currentSection = &Section{
				Title:     text,
				Questions: []Question{},
			}
			continue
		}

		// 处理小题号
		if currentSection != nil && questionRegex.MatchString(text) {
			// 保存上一个小题
			if currentQuestion != nil {
				currentSection.Questions = append(currentSection.Questions, *currentQuestion)
			}
			// 开始新的小题
			questionNumber := questionRegex.FindString(text)
			// 提取纯数字的小题号
			numberOnly := numberOnlyRegex.FindString(questionNumber)
			questionContent := strings.TrimSpace(text[len(questionNumber):])
			currentQuestion = &Question{
				Number:  numberOnly,
				Content: questionContent,
			}
		} else if currentSection != nil && choiceRegex.MatchString(text) {
			// 拼接选择题选项
			if currentQuestion != nil {
				currentQuestion.Content += "\n" + text
			} else {
				// 如果没有当前小题，可能是题干的一部分
				currentSection.Title += "\n" + text
			}
		} else if currentSection != nil {
			// 拼接到当前小题或大题引言
			if currentQuestion != nil {
				currentQuestion.Content += "\n" + text
			} else {
				currentSection.Title += "\n" + text
			}
		}
	}

	// 保存最后一个小题和大题
	if currentSection != nil {
		if currentQuestion != nil {
			currentSection.Questions = append(currentSection.Questions, *currentQuestion)
		}
		sections = append(sections, *currentSection)
	}

	return sections, nil
}
