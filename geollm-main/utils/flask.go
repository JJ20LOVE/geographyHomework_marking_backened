package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type EResponse struct {
	Result Result `json:"result"`
	Error  string `json:"error"`
	Detail string `json:"detail"`
}

type Result struct {
	Score     int    `json:"score"`
	Comment   string `json:"comment"`
	Structure int    `json:"structure"`
}

func Eva(question, correct_answer, student_answer string, full_score int) Result {
	url := "https://zer0.top/flask/llm_comment"

	// 定义要发送的 JSON 数据
	requestBody := map[string]interface{}{
		"question":       question,
		"full_score":     full_score,
		"correct_answer": correct_answer,
		"student_answer": student_answer,
		//"criteria":       "",
	}

	// 将请求数据转换为 JSON 字符串
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("无法将请求数据转换为 JSON: %v", err)
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("创建请求失败: %v", err)
	}

	// 设置请求头，声明我们发送的是 JSON 数据
	req.Header.Set("Content-Type", "application/json")

	// 使用 http.DefaultClient 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应失败: %v", err)
	}

	// 打印响应体（如果是 JSON 响应）
	//fmt.Println("响应:", string(body))

	// 将响应体解析为 JSON 对象

	var responseBody EResponse
	if err := json.Unmarshal(body, &responseBody); err != nil {
		log.Printf("解析响应 JSON 失败: %v", err)
		return Result{}
	}
	//fmt.Println(responseBody.Error, responseBody.Result, responseBody.Detail)
	if responseBody.Error != "" {
		log.Printf("服务器返回错误: %v", responseBody.Error)
		return Result{}
	}
	return responseBody.Result
}

func Problem(question string, comments []string, avg_score float64, full_score int) string {
	url := "https://zer0.top/flask/problem"

	// 定义要发送的 JSON 数据
	requestBody := map[string]interface{}{
		"question":   question,
		"comments":   comments,
		"avg_score":  avg_score,
		"full_score": full_score,
	}

	// 将请求数据转换为 JSON 字符串
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("无法将请求数据转换为 JSON: %v", err)
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("创建请求失败: %v", err)
	}

	// 设置请求头，声明我们发送的是 JSON 数据
	req.Header.Set("Content-Type", "application/json")

	// 使用 http.DefaultClient 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应失败: %v", err)
	}

	// 打印响应体（如果是 JSON 响应）
	//fmt.Println("响应:", string(body))

	// 将响应体解析为 JSON 对象

	var responseBody struct {
		Result string `json:"result"`
		Error  string `json:"error"`
		Detail string `json:"detail"`
	}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		log.Printf("解析响应 JSON 失败: %v", err)
		return ""
	}
	//fmt.Println(responseBody.Error, responseBody.Result, responseBody.Detail)
	if responseBody.Error != "" {
		log.Printf("服务器返回错误: %v", responseBody.Error)
		return ""
	}
	return responseBody.Result
}

func Knowledge(question, answer string) []string {
	url := "https://zer0.top/flask/knowledge"

	requestBody := map[string]interface{}{
		"question": question,
		"answer":   answer,
	}
	// 将请求数据转换为 JSON 字符串
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("无法将请求数据转换为 JSON: %v", err)
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("创建请求失败: %v", err)
	}

	// 设置请求头，声明我们发送的是 JSON 数据
	req.Header.Set("Content-Type", "application/json")

	// 使用 http.DefaultClient 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应失败: %v", err)
	}

	// 将响应体解析为 JSON 对象
	var responseBody struct {
		Result []string `json:"result"`
		Error  string   `json:"error"`
		Detail string   `json:"detail"`
	}
	if err := json.Unmarshal(body, &responseBody); err != nil {
		log.Printf("解析响应 JSON 失败: %v", err)
		return nil
	}
	if responseBody.Error != "" {
		log.Printf("服务器返回错误: %v", responseBody.Error)
		return nil
	}
	return responseBody.Result
}
