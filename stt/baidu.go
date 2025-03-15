package stt

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const API_URL = "http://vop.baidu.com/server_api"
const DEV_PID = "1537" //普通话 输入法模型
const FORMAT = "wav"   // support pcm/wav/amr 格式，极速版额外支持m4a 格式
const RATE = "16000"

// refer https://cloud.baidu.com/doc/SPEECH/s/Jlbxdezuf and https://github.com/Baidu-AIP/speech-demo
func Asr(filePath string) ([]string, error) {
	token, err := queryAccessToken()
	if err != nil {
		log.Printf("query access token error:%s", err.Error())
		return nil, err
	}
	asrUrl := API_URL + "?dev_pid=" + DEV_PID + "&token=" + token + "&cuid=123456"
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("open file error:%s", err.Error())
		return nil, err
	}
	res, err := http.Post(asrUrl, "audio/"+FORMAT+";rate="+RATE, file)
	if err != nil {
		log.Printf("http post error:%s", err.Error())
		return nil, err
	}
	defer res.Body.Close()
	resp := &AsrResp{}
	if err = json.NewDecoder(res.Body).Decode(resp); err != nil {
		log.Printf("http client new decoder error:%s", err.Error())
		return nil, err
	}
	marshal, _ := json.Marshal(resp)
	log.Printf("speech to text result:%s", string(marshal))
	if resp.ErrNo != 0 {
		return nil, fmt.Errorf("speech to text failed:%s", resp.ErrMsg)
	}
	log.Println(resp.Result)
	return resp.Result, nil
}

type AsrResp struct {
	CorpusNo string   `json:"corpus_no"`
	ErrMsg   string   `json:"err_msg"`
	ErrNo    int      `json:"err_no"`
	Result   []string `json:"result"`
	Sn       string   `json:"sn"`
}

func queryAccessToken() (string, error) {
	apiKey := os.Getenv("BAIDU_API_KEY")
	secretKey := os.Getenv("BAIDU_SECRET_KEY")
	url := "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=" + apiKey + "&client_secret=" + secretKey
	res, err := http.Get(url)
	if err != nil {
		log.Printf("http get error:%s", err.Error())
		return "", err
	}
	defer res.Body.Close()
	resp := &TokenResp{}
	if err = json.NewDecoder(res.Body).Decode(resp); err != nil {
		log.Printf("http client new decoder error %s", err.Error())
		return "", err
	}
	marshal, _ := json.Marshal(resp)
	log.Printf("query access token resp:%s", string(marshal))
	if resp.AccessToken == "" {
		return "", fmt.Errorf("query access token is empty:%s", string(marshal))
	}
	return resp.AccessToken, nil
}

type TokenResp struct {
	RefreshToken  string `json:"refresh_token"`
	ExpiresIn     int    `json:"expires_in"`
	SessionKey    string `json:"session_key"`
	AccessToken   string `json:"access_token"`
	Scope         string `json:"scope"`
	SessionSecret string `json:"session_secret"`
}
