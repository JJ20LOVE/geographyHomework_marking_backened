package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Result []string `json:"result"`
	Error  string   `json:"error"`
	Detail string   `json:"detail"`
}

func StartOcr(fileHeaders []*multipart.FileHeader) []string {
	if len(fileHeaders) == 0 {
		fmt.Println("OCR错误: 没有上传文件")
		return nil
	}

	// 根据配置选择OCR服务商
	switch OcrProvider {
	case "aliyun":
		fmt.Println("使用阿里云OCR服务")
		return StartAliyunOcr(fileHeaders)
	case "baidu":
		fmt.Println("使用百度OCR服务")
		// 可以在这里添加百度OCR调用
		return nil
	case "tencent":
		fmt.Println("使用腾讯云OCR服务")
		// 可以在这里添加腾讯云OCR调用
		return nil
	default:
		fmt.Printf("不支持的OCR服务商: %s\n", OcrProvider)
		return nil
	}
}

// 保留原有的通用OCR函数作为备用
func StartGenericOcr(fileHeaders []*multipart.FileHeader) []string {
	// 创建一个buffer来存储multipart form的数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	t := strconv.Itoa(len(fileHeaders) - 1)
	// 遍历所有传入的文件头
	for _, fileHeader := range fileHeaders {
		// 打开文件
		file, err := fileHeader.Open()
		if err != nil {
			return nil
		}
		defer file.Close()

		// 创建一个文件字段
		part, err := writer.CreateFormFile("files", fileHeader.Filename)
		if err != nil {
			return nil
		}

		// 复制文件内容到part中
		_, err = io.Copy(part, file)
		if err != nil {
			return nil
		}
	}

	// 关闭multipart writer，确保写入结束
	err := writer.Close()
	if err != nil {
		return nil
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", OcrApi+t, body)
	if err != nil {
		return nil
	}

	// 设置Content-Type为multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求（添加超时）
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil
	}
	if response.Error != "" {
		return nil
	}
	// 返回结果
	return response.Result
}
