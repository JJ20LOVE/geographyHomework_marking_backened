package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

type Response struct {
	Result []string `json:"result"`
	Error  string   `json:"error"`
	Detail string   `json:"detail"`
}

func StartOcr(fileHeaders []*multipart.FileHeader) []string {
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

	// 发送请求
	client := &http.Client{}
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
