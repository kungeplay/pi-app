package audio

import (
	"bytes"
	"io"
	"log"
	"os/exec"
	"strconv"
)

func SetVolume(volume int) {
	//设置音量大小
	amixerCmd := exec.Command("amixer", "set", "Master", strconv.Itoa(volume)+"%")
	var amixerBuffer bytes.Buffer
	amixerCmd.Stdout = &amixerBuffer
	if err := amixerCmd.Start(); err != nil {
		log.Printf("amixer command start error:%s", err.Error())
		return
	}
	if err := amixerCmd.Wait(); err != nil {
		log.Printf("amixer command wait error:%s", err.Error())
		return
	}
}

// 播放提示音
func PlayAttention(audioFile string) {
	//设置音量大小
	aplayCmd := exec.Command("aplay", audioFile)
	var playBuffer bytes.Buffer
	aplayCmd.Stdout = &playBuffer
	if err := aplayCmd.Start(); err != nil {
		log.Println("aplayCmd command start error:%s", err.Error())
		return
	}
	if err := aplayCmd.Wait(); err != nil {
		log.Println("aplayCmd command wait error:%s", err.Error())
		return
	}
}

// 播放文字音频
func PlayText(content string) {
	log.Println("audio content:" + content)

	echoCmd := exec.Command("echo", content)
	piperCmd := exec.Command("/home/liujiakun/Data/thirdSoft/piper/piper", "--model", "/home/liujiakun/Data/thirdSoft/piper/voices/zh_CN-huayan-medium.onnx", "--output-raw")
	aplayCmd := exec.Command("aplay", "-r", "22050", "-f", "S16_LE", "-t", "raw", "-")

	r1, w1 := io.Pipe()
	defer r1.Close()
	defer w1.Close()
	echoCmd.Stdout = w1
	piperCmd.Stdin = r1

	r2, w2 := io.Pipe()
	defer r2.Close()
	defer w2.Close()
	piperCmd.Stdout = w2
	aplayCmd.Stdin = r2

	var buffer bytes.Buffer
	aplayCmd.Stdout = &buffer

	if err := echoCmd.Start(); err != nil {
		log.Printf("echo command start error:%s", err.Error())
		return
	}
	if err := piperCmd.Start(); err != nil {
		log.Printf("piper command start error:%s", err.Error())
		return
	}
	if err := aplayCmd.Start(); err != nil {
		log.Printf("aplay command start error:%s", err.Error())
		return
	}
	if err := echoCmd.Wait(); err != nil {
		log.Printf("echo command wait error:%s", err.Error())
		return
	}
	if err := w1.Close(); err != nil {
		log.Printf("pipe 1 close error:%s", err.Error())
		return
	}
	if err := piperCmd.Wait(); err != nil {
		log.Printf("piper command wait error:%s", err.Error())
		return
	}
	if err := w2.Close(); err != nil {
		log.Printf("pipe 2 close error:%s", err.Error())
		return
	}
	if err := aplayCmd.Wait(); err != nil {
		log.Printf("aplay command wait error:%s", err.Error())
		return
	}

	log.Println("audio play finished " + buffer.String())
}
