package aigc

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

type Role = string

const (
	System    Role = "system"
	User      Role = "user"
	Assistant Role = "assistant"
)

type Model = string

const (
	QWEN1_5_1_8B_CHAT Model = "qwen1.5-1.8b-chat"
)

type ChatCache struct {
	Messages   []Message `json:"messages"`
	LatestTime time.Time `json:"latest_time_stamp"`
}

type TongYiQianWenChatReq struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TongYiQianWenChatResp struct {
	Choices           []Choice `json:"choices"`
	Created           int64    `json:"created"`
	Id                string   `json:"id"`
	Model             string   `json:"model"`
	Object            string   `json:"object"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Usage             *Usage   `json:"usage"`
}
type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Logprobs     string  `json:"logprobs"`
	Message      Message `json:"message"`
}

var chatCache ChatCache

//func main() {
//	// 从标准输入流中接收输入数据
//	input := bufio.NewScanner(os.Stdin)
//	fmt.Printf("开始进入交互:\n")
//	// 逐行扫描
//	for input.Scan() {
//		line := input.Text()
//		// 输入bye时 结束
//		if line == "bye" {
//			break
//		}
//		reply := Chat(line)
//		playAudio(reply)
//	}
//}

func InitChatInfo() *ChatCache {
	messages := make([]Message, 0)
	messages = append(messages, Message{
		Role:    System,
		Content: "You are a helpful assistant.",
	})
	chat := ChatCache{
		Messages:   messages,
		LatestTime: time.Now(),
	}
	return &chat
}

func Chat(content string) string {
	if time.Now().Second()-chatCache.LatestTime.Second() > 60 {
		chatCache = *InitChatInfo()
		log.Println("reBegin chat")
	}

	chatCache.Messages = append(chatCache.Messages, Message{
		Role:    User,
		Content: content,
	})

	replyMes, err := queryOpenApi(chatCache.Messages)
	if err != nil {
		panic(err)
	}
	chatCache.Messages = append(chatCache.Messages, *replyMes)
	chatCache.LatestTime = time.Now()
	return replyMes.Content
}

func queryOpenApi(messages []Message) (*Message, error) {
	posturl := "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"
	AUTHORIZATION_KEY := os.Getenv("ALIYUNCS_AUTHORIZATION_KEY") //获取环境变量中查询天气开放平台的key
	req := TongYiQianWenChatReq{
		Model:    QWEN1_5_1_8B_CHAT,
		Messages: messages,
	}
	marshal, _ := json.Marshal(req)
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(marshal))
	if err != nil {
		log.Println("http NewRequest error", err)
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", AUTHORIZATION_KEY)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Println("http client do error", err)
		return nil, err
	}
	defer res.Body.Close()

	resp := &TongYiQianWenChatResp{}
	err = json.NewDecoder(res.Body).Decode(resp)
	if err != nil {
		log.Printf("http client new decoder error %s", err.Error())
		return nil, err
	}
	respStr, _ := json.Marshal(resp)
	log.Println("queryOpenApi resp:" + string(respStr))
	if resp.Usage == nil {
		log.Println("http client resp usage empty")
		return nil, errors.New("http client resp usage error")
	}
	return &resp.Choices[0].Message, nil
}
