package service

import (
	"bytes"
	"dbdemo/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AIService struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// AI请求结构
type AIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AI响应结构
type AIResponse struct {
	Choices []Choice `json:"choices"`
	Error   *AIError `json:"error"`
}

type Choice struct {
	Message Message `json:"message"`
}

type AIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func NewAIService(apiKey, baseURL string) *AIService {
	return &AIService{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetSimilarQuestions 调用AI API获取同类题目推荐
func (s *AIService) GetSimilarQuestions(questionText, knowledgePoint string, limit int) ([]model.SimilarQuestion, error) {
	if s.APIKey == "" {
		// 如果没有配置API Key，返回模拟数据
		return s.getMockSimilarQuestions(questionText, knowledgePoint, limit)
	}

	prompt := fmt.Sprintf(`请根据以下地理题目推荐%d道同类题目：

原题目：%s
知识点：%s

要求：
1. 题目类型和难度与原题相似
2. 围绕相同的知识点
3. 返回JSON格式，包含question_id, question_text, knowledge_point, difficulty字段
4. question_id从1001开始递增
5. difficulty可以是"简单"、"中等"、"困难"

请直接返回JSON数组，不要其他解释：`, limit, questionText, knowledgePoint)

	aiReq := AIRequest{
		Model: "deepseek-chat", // 或其他模型
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   2000,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(aiReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回错误: %s", string(body))
	}

	var aiResp AIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if aiResp.Error != nil {
		return nil, fmt.Errorf("AI服务错误: %s", aiResp.Error.Message)
	}

	if len(aiResp.Choices) == 0 {
		return nil, fmt.Errorf("AI返回空结果")
	}

	// 解析AI返回的JSON
	var similarQuestions []model.SimilarQuestion
	content := aiResp.Choices[0].Message.Content

	// 清理可能的markdown代码块
	cleanContent := s.cleanJSONResponse(content)

	if err := json.Unmarshal([]byte(cleanContent), &similarQuestions); err != nil {
		// 如果解析失败，返回模拟数据
		fmt.Printf("AI响应解析失败，使用模拟数据: %v\n", err)
		return s.getMockSimilarQuestions(questionText, knowledgePoint, limit)
	}

	return similarQuestions, nil
}

// cleanJSONResponse 清理AI返回的JSON响应
func (s *AIService) cleanJSONResponse(content string) string {
	// 移除可能的markdown代码块标记
	cleaned := content
	if len(cleaned) >= 7 && cleaned[:7] == "```json" {
		cleaned = cleaned[7:]
	}
	if len(cleaned) >= 3 && cleaned[len(cleaned)-3:] == "```" {
		cleaned = cleaned[:len(cleaned)-3]
	}
	return cleaned
}

// getMockSimilarQuestions 返回模拟的相似题目（用于测试或API不可用时）
func (s *AIService) getMockSimilarQuestions(questionText, knowledgePoint string, limit int) ([]model.SimilarQuestion, error) {
	mockQuestions := []model.SimilarQuestion{
		{
			QuestionID:     1001,
			QuestionText:   "分析导致珠江三角洲人口密集的自然因素。",
			KnowledgePoint: knowledgePoint,
			Difficulty:     "中等",
		},
		{
			QuestionID:     1002,
			QuestionText:   "说明四川盆地人口稠密的自然条件。",
			KnowledgePoint: knowledgePoint,
			Difficulty:     "中等",
		},
		{
			QuestionID:     1003,
			QuestionText:   "比较长江三角洲和珠江三角洲的人口分布特征。",
			KnowledgePoint: knowledgePoint,
			Difficulty:     "困难",
		},
		{
			QuestionID:     1004,
			QuestionText:   "分析黄河流域人口分布的主要影响因素。",
			KnowledgePoint: knowledgePoint,
			Difficulty:     "中等",
		},
		{
			QuestionID:     1005,
			QuestionText:   "说明地形对华北平原人口分布的影响。",
			KnowledgePoint: knowledgePoint,
			Difficulty:     "简单",
		},
	}

	if limit > 0 && limit < len(mockQuestions) {
		return mockQuestions[:limit], nil
	}
	return mockQuestions, nil
}
