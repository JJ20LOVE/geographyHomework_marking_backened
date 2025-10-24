package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

// 新的OCR响应结构 - 根据实际响应调整
type AliyunOCRResponse struct {
	Content        string     `json:"content,omitempty"`
	Words          string     `json:"words,omitempty"`
	Text           string     `json:"text,omitempty"`
	PrismWnum      int        `json:"prism_wnum,omitempty"`
	PrismWordsInfo []WordInfo `json:"prism_wordsInfo,omitempty"`
	RequestID      string     `json:"requestId,omitempty"`
	Success        bool       `json:"success,omitempty"`
	Code           int        `json:"code,omitempty"`
	Message        string     `json:"message,omitempty"`
	Height         int        `json:"height,omitempty"`
	Width          int        `json:"width,omitempty"`
}

type WordInfo struct {
	Word      string `json:"word,omitempty"`
	Direction int    `json:"direction,omitempty"`
}

// 配置信息 - 更新为新的OCR服务
const (
	AppCode = "yourcode"
	APIHost = "gjbsb.market.alicloudapi.com"
	APIPath = "/ocrservice/advanced"

	// 图片尺寸限制
	MinImageSize = 50              // 最小边长
	MaxImageSize = 4096            // 最大边长（调整为4096更安全）
	MaxFileSize  = 4 * 1024 * 1024 // 最大文件大小 4MB
)

// StartAliyunOcr 使用新的阿里云OCR服务
func StartAliyunOcr(fileHeaders []*multipart.FileHeader) []string {
	var results []string

	fmt.Printf("使用阿里云OCR服务（通用文字识别）\n")
	fmt.Printf("开始处理 %d 个文件的OCR识别\n", len(fileHeaders))

	for i, fileHeader := range fileHeaders {
		fmt.Printf("处理第 %d 个文件: %s\n", i+1, fileHeader.Filename)

		// 打开文件
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Printf("打开文件失败: %v\n", err)
			results = append(results, "")
			continue
		}

		// 读取文件内容
		fileData, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			fmt.Printf("读取文件失败: %v\n", err)
			results = append(results, "")
			continue
		}

		fmt.Printf("文件大小: %d bytes\n", len(fileData))

		// 检查文件大小
		if len(fileData) > MaxFileSize {
			fmt.Printf("文件过大，进行压缩处理\n")
			fileData, err = compressImage(fileData, fileHeader.Filename)
			if err != nil {
				fmt.Printf("图片压缩失败: %v\n", err)
				results = append(results, "")
				continue
			}
			fmt.Printf("压缩后文件大小: %d bytes\n", len(fileData))
		}

		// 检查图片尺寸并调整
		processedData, err := checkAndResizeImage(fileData, fileHeader.Filename)
		if err != nil {
			fmt.Printf("图片尺寸调整失败: %v\n", err)
			results = append(results, "")
			continue
		}

		// 编码为base64
		imageBase64 := base64.StdEncoding.EncodeToString(processedData)
		fmt.Printf("Base64长度: %d\n", len(imageBase64))

		// 构建请求体 - 根据新的API格式
		requestBody := map[string]interface{}{
			"img":      imageBase64,
			"url":      "", // 不使用URL方式
			"prob":     false,
			"charInfo": false,
			"rotate":   false,
			"table":    false,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			fmt.Printf("构建JSON请求体失败: %v\n", err)
			results = append(results, "")
			continue
		}

		// 创建HTTP请求
		apiURL := "https://" + APIHost + APIPath
		fmt.Printf("请求URL: %s\n", apiURL)

		req, err := http.NewRequest("POST", apiURL, bytes.NewReader(jsonBody))
		if err != nil {
			fmt.Printf("创建请求失败: %v\n", err)
			results = append(results, "")
			continue
		}

		// 设置请求头 - 根据新的API要求
		req.Header.Set("Authorization", "APPCODE "+AppCode)
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		fmt.Printf("请求头: Authorization=APPCODE %s..., Content-Type=%s\n",
			AppCode[:8], req.Header.Get("Content-Type"))

		// 发送请求
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("OCR请求失败: %v\n", err)
			results = append(results, "")
			continue
		}
		defer resp.Body.Close()

		// 读取响应
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("读取响应失败: %v\n", err)
			results = append(results, "")
			continue
		}

		fmt.Printf("OCR响应状态: %d\n", resp.StatusCode)
		fmt.Printf("OCR响应体: %s\n", string(respBody))

		// 处理响应
		if resp.StatusCode != 200 {
			fmt.Printf("OCR服务返回错误状态码: %d, 错误信息: %s\n", resp.StatusCode, string(respBody))
			results = append(results, "")
			continue
		}

		// 解析响应
		var ocrResp AliyunOCRResponse
		if err := json.Unmarshal(respBody, &ocrResp); err != nil {
			fmt.Printf("解析OCR响应失败: %v\n", err)
			// 尝试通用方式解析
			content := extractOCRContentGeneric(respBody)
			results = append(results, content)
			continue
		}

		fmt.Printf("解析后的OCR响应: %+v\n", ocrResp)

		// 提取识别结果
		content := extractOCRContentNew(ocrResp)
		results = append(results, content)
		fmt.Printf("OCR识别结果: %s\n", content)
	}

	fmt.Printf("OCR处理完成，返回 %d 个结果\n", len(results))
	return results
}

// 从新的OCR响应结构中提取内容
func extractOCRContentNew(ocrResp AliyunOCRResponse) string {
	// 优先从content字段提取
	if ocrResp.Content != "" {
		fmt.Printf("从content字段提取到内容，长度: %d\n", len(ocrResp.Content))
		return ocrResp.Content
	}

	// 其次从words字段提取
	if ocrResp.Words != "" {
		fmt.Printf("从words字段提取到内容，长度: %d\n", len(ocrResp.Words))
		return ocrResp.Words
	}

	// 然后从text字段提取
	if ocrResp.Text != "" {
		fmt.Printf("从text字段提取到内容，长度: %d\n", len(ocrResp.Text))
		return ocrResp.Text
	}

	// 如果有prism_wordsInfo，提取所有单词
	if len(ocrResp.PrismWordsInfo) > 0 {
		var text string
		for _, wordInfo := range ocrResp.PrismWordsInfo {
			if wordInfo.Word != "" {
				text += wordInfo.Word + " "
			}
		}
		if text != "" {
			fmt.Printf("从prism_wordsInfo提取到内容，长度: %d\n", len(text))
			return strings.TrimSpace(text)
		}
	}

	// 如果都没有有效内容，返回提示信息
	fmt.Printf("无法提取OCR内容，响应结构: %+v\n", ocrResp)
	return "OCR识别完成但无法提取文本内容"
}

// 通用方式提取OCR内容（兼容旧格式）
func extractOCRContentGeneric(respBody []byte) string {
	var genericResp map[string]interface{}
	if err := json.Unmarshal(respBody, &genericResp); err != nil {
		fmt.Printf("通用解析失败: %v\n", err)
		return "解析OCR响应失败"
	}

	fmt.Printf("通用解析结果 - 可用字段: ")
	for key := range genericResp {
		fmt.Printf("%s ", key)
	}
	fmt.Printf("\n")

	// 尝试不同的字段提取内容
	if content, ok := genericResp["content"].(string); ok && content != "" {
		fmt.Printf("从content字段提取到内容，长度: %d\n", len(content))
		return content
	}
	if words, ok := genericResp["words"].(string); ok && words != "" {
		fmt.Printf("从words字段提取到内容，长度: %d\n", len(words))
		return words
	}
	if text, ok := genericResp["text"].(string); ok && text != "" {
		fmt.Printf("从text字段提取到内容，长度: %d\n", len(text))
		return text
	}

	// 尝试从prism_wordsInfo中提取
	if wordsInfo, ok := genericResp["prism_wordsInfo"].([]interface{}); ok {
		var text string
		for i, word := range wordsInfo {
			if wordMap, ok := word.(map[string]interface{}); ok {
				if wordStr, ok := wordMap["word"].(string); ok && wordStr != "" {
					text += wordStr + " "
					if i < 5 { // 只打印前5个单词用于调试
						fmt.Printf("单词 %d: %s\n", i+1, wordStr)
					}
				}
			}
		}
		if text != "" {
			fmt.Printf("从prism_wordsInfo提取到内容，总长度: %d\n", len(text))
			return strings.TrimSpace(text)
		}
	}

	fmt.Printf("无法提取OCR内容，响应结构: %+v\n", genericResp)
	return "OCR识别完成但无法提取文本内容"
}

// 检查并调整图片尺寸
func checkAndResizeImage(fileData []byte, filename string) ([]byte, error) {
	// 解码图片
	img, format, err := image.Decode(bytes.NewReader(fileData))
	if err != nil {
		return nil, fmt.Errorf("图片解码失败: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	fmt.Printf("图片原始尺寸: %dx%d, 格式: %s\n", width, height, format)

	// 检查尺寸是否符合要求
	minSide := min(width, height)
	maxSide := max(width, height)

	if minSide < MinImageSize || maxSide > MaxImageSize {
		fmt.Printf("图片尺寸不符合要求，进行缩放. 最小边: %d, 最大边: %d\n", minSide, maxSide)

		// 计算缩放比例
		scale := 1.0
		if minSide < MinImageSize {
			scale = float64(MinImageSize) / float64(minSide)
		}
		if maxSide > MaxImageSize {
			scale = minFloat(scale, float64(MaxImageSize)/float64(maxSide))
		}

		newWidth := uint(float64(width) * scale)
		newHeight := uint(float64(height) * scale)

		fmt.Printf("缩放比例: %.2f, 新尺寸: %dx%d\n", scale, newWidth, newHeight)

		// 使用Lanczos3算法进行高质量缩放
		resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

		// 编码回字节
		var buf bytes.Buffer
		switch strings.ToLower(format) {
		case "jpeg", "jpg":
			err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 85})
		case "png":
			err = png.Encode(&buf, resizedImg)
		default:
			// 默认使用JPEG
			err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 85})
		}

		if err != nil {
			return nil, fmt.Errorf("图片编码失败: %v", err)
		}

		return buf.Bytes(), nil
	}

	fmt.Printf("图片尺寸符合要求，无需调整\n")
	return fileData, nil
}

// 压缩图片
func compressImage(fileData []byte, filename string) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(fileData))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 75})
	case "png":
		// PNG压缩比较困难，可以转换为JPEG
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	default:
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 75})
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
