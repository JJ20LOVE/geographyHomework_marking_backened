package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

// AliyunOCRResponse 阿里云OCR响应结构
type AliyunOCRResponse struct {
	Data struct {
		Content string `json:"content"`
	} `json:"data"`
	RequestID string `json:"requestId"`
	Success   bool   `json:"success"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

// StartAliyunOcr 使用阿里云OCR服务
func StartAliyunOcr(fileHeaders []*multipart.FileHeader) []string {
	var results []string

	// 检查配置
	if AliyunAccessKeyId == "" || AliyunAccessKeySecret == "" {
		fmt.Println("阿里云OCR配置缺失，请设置 ALIYUN_ACCESS_KEY_ID 和 ALIYUN_ACCESS_KEY_SECRET")
		return results
	}

	// 创建阿里云客户端
	client, err := sdk.NewClientWithAccessKey(
		AliyunRegion,
		AliyunAccessKeyId,
		AliyunAccessKeySecret,
	)
	if err != nil {
		fmt.Printf("创建阿里云客户端失败: %v\n", err)
		return results
	}

	for _, fileHeader := range fileHeaders {
		// 打开文件
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Printf("打开文件失败: %v\n", err)
			results = append(results, "")
			continue
		}
		defer file.Close()

		// 读取文件内容
		fileData, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("读取文件失败: %v\n", err)
			results = append(results, "")
			continue
		}

		// 编码为base64
		imageBase64 := base64.StdEncoding.EncodeToString(fileData)

		// 创建OCR请求
		request := requests.NewCommonRequest()
		request.Domain = "ocr.cn-shanghai.aliyuncs.com"
		request.Version = "2019-12-30"
		request.ApiName = "RecognizeCharacter"
		request.QueryParams["Image"] = imageBase64
		request.QueryParams["MinHeight"] = "10"
		request.QueryParams["OutputProbability"] = "false"

		// 发送请求
		response, err := client.ProcessCommonRequest(request)
		if err != nil {
			fmt.Printf("阿里云OCR请求失败: %v\n", err)
			results = append(results, "")
			continue
		}

		// 解析响应
		var ocrResp AliyunOCRResponse
		if err := json.Unmarshal(response.GetHttpContentBytes(), &ocrResp); err != nil {
			fmt.Printf("解析OCR响应失败: %v\n", err)
			results = append(results, "")
			continue
		}

		if !ocrResp.Success {
			fmt.Printf("阿里云OCR错误: %s (代码: %s)\n", ocrResp.Message, ocrResp.Code)
			results = append(results, "")
			continue
		}

		results = append(results, ocrResp.Data.Content)
	}

	return results
}
