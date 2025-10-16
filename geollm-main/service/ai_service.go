package service

import (
	"bytes"
	"dbdemo/model"
	"dbdemo/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AIService struct {
	APIKey      string
	BaseURL     string
	Model       string
	MaxTokens   int
	Temperature float64
	Client      *http.Client
}

// AI请求结构 - 适配DeepSeek API
type AIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	Stream      bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AI响应结构 - 适配DeepSeek API
type AIResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
	Error   *AIError `json:"error,omitempty"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type AIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code,omitempty"`
}

func NewAIService(apiKey, baseURL, model string, maxTokens int, temperature float64) *AIService {
	return &AIService{
		APIKey:      apiKey,
		BaseURL:     baseURL,
		Model:       model,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		Client: &http.Client{
			Timeout: 60 * time.Second, // 增加超时时间
		},
	}
}

// NewDefaultAIService 使用配置的默认值创建AI服务
func NewDefaultAIService() *AIService {
	return NewAIService(
		utils.DeepSeekApiKey,
		utils.DeepSeekBaseUrl,
		utils.AIModel,
		utils.AIMaxTokens,
		utils.AITemperature,
	)
}

// GetSimilarQuestions 调用DeepSeek API获取同类题目推荐
func (s *AIService) GetSimilarQuestions(questionText, knowledgePoint string, limit int) ([]model.SimilarQuestion, error) {
	if s.APIKey == "" {
		// 如果没有配置API Key，返回模拟数据
		fmt.Println("DeepSeek API Key未配置，使用模拟数据")
		return s.getMockSimilarQuestions(questionText, knowledgePoint, limit)
	}

	prompt := fmt.Sprintf(`你是一个地理学科教育专家，请根据以下地理题目推荐%d道同类题目：

原题目：%s
知识点：%s

要求：
1. 题目类型和难度与原题相似，都是地理主观题
2. 围绕相同的知识点或相关地理概念
3. 返回格式必须是纯JSON数组，不要有任何其他文字
4. JSON数组包含%d个对象，每个对象包含以下字段：
   - question_id: 从1001开始递增的数字
   - question_text: 题目内容
   - knowledge_point: 知识点
   - difficulty: 难度级别，只能是"简单"、"中等"、"困难"之一

请确保返回的是纯JSON格式，不要有任何markdown标记或额外解释：`, limit, questionText, knowledgePoint, limit)

	aiReq := AIRequest{
		Model: s.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   s.MaxTokens,
		Temperature: s.Temperature,
		Stream:      false,
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
	req.Header.Set("Accept", "application/json")

	fmt.Printf("调用DeepSeek API，模型: %s\n", s.Model)

	resp, err := s.Client.Do(req)
	if err != nil {
		fmt.Printf("API请求失败: %v\n", err)
		return nil, fmt.Errorf("API请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API返回错误状态码: %d, 响应: %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("API返回错误: %s", string(body))
	}

	var aiResp AIResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		fmt.Printf("解析响应失败: %v, 原始响应: %s\n", err, string(body))
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if aiResp.Error != nil {
		fmt.Printf("AI服务错误: %s (类型: %s)\n", aiResp.Error.Message, aiResp.Error.Type)
		return nil, fmt.Errorf("AI服务错误: %s", aiResp.Error.Message)
	}

	if len(aiResp.Choices) == 0 {
		return nil, fmt.Errorf("AI返回空结果")
	}

	// 解析AI返回的JSON
	var similarQuestions []model.SimilarQuestion
	content := aiResp.Choices[0].Message.Content

	// 清理可能的markdown代码块和非JSON内容
	cleanContent := s.cleanJSONResponse(content)

	fmt.Printf("AI返回内容: %s\n", cleanContent)

	if err := json.Unmarshal([]byte(cleanContent), &similarQuestions); err != nil {
		fmt.Printf("AI响应JSON解析失败: %v，使用模拟数据\n", err)
		return s.getMockSimilarQuestions(questionText, knowledgePoint, limit)
	}

	// 验证返回的数据结构
	if len(similarQuestions) == 0 {
		fmt.Printf("AI返回空数组，使用模拟数据\n")
		return s.getMockSimilarQuestions(questionText, knowledgePoint, limit)
	}

	fmt.Printf("成功获取 %d 道同类题目\n", len(similarQuestions))
	return similarQuestions, nil
}

// cleanJSONResponse 清理AI返回的JSON响应
func (s *AIService) cleanJSONResponse(content string) string {
	// 移除可能的markdown代码块标记和其他非JSON内容
	cleaned := content

	// 移除 ```json 和 ``` 标记
	if len(cleaned) >= 7 && cleaned[:7] == "```json" {
		cleaned = cleaned[7:]
	}
	if len(cleaned) >= 3 && cleaned[len(cleaned)-3:] == "```" {
		cleaned = cleaned[:len(cleaned)-3]
	}

	// 移除开头的换行和空格
	for len(cleaned) > 0 && (cleaned[0] == '\n' || cleaned[0] == ' ' || cleaned[0] == '\r') {
		cleaned = cleaned[1:]
	}

	// 移除结尾的换行和空格
	for len(cleaned) > 0 && (cleaned[len(cleaned)-1] == '\n' || cleaned[len(cleaned)-1] == ' ' || cleaned[len(cleaned)-1] == '\r') {
		cleaned = cleaned[:len(cleaned)-1]
	}

	return cleaned
}

// getMockSimilarQuestions 返回模拟的相似题目（用于测试或API不可用时）
func (s *AIService) getMockSimilarQuestions(questionText, knowledgePoint string, limit int) ([]model.SimilarQuestion, error) {
	fmt.Printf("使用模拟数据，原题: %s, 知识点: %s\n", questionText, knowledgePoint)

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
