package main

import (
	"fmt"
	"log"
	"os"
	"pi-app/aigc"
	"pi-app/audio"
	"pi-app/record"
	"pi-app/stt"
	"pi-app/wake"
	"strconv"
	"strings"
	"time"
)

const (
	AUDIO_VOLUME        = 50
	COMMAND_DURATION    = 4
	COMMON_RES          = "../resources/common.res"
	DEFAULT_SENSITIVITY = 0.3
	HOTWORD_MODEL_FILE  = "../resources/xiaoyixiaoyi.pmdl"
)

func main() {
	recChannel := make(chan int, 1)
	//0.5
	hotWordModelFile := HOTWORD_MODEL_FILE
	sensitivity := DEFAULT_SENSITIVITY
	if len(os.Args) > 1 {
		hotWordModelFile = os.Args[1]
	}
	if len(os.Args) > 2 {
		if sensitivityTemp, err := strconv.ParseFloat(os.Args[2], 32); err != nil {
			log.Printf("os.Args[2] to sensitivity convert error:%s", os.Args[3])
			return
		} else {
			sensitivity = sensitivityTemp
		}
	}
	go wake.ListenerHotWord(COMMON_RES, hotWordModelFile, float32(sensitivity), recChannel)
	for {
		if _, ok := <-recChannel; ok {
			audio.SetVolume(AUDIO_VOLUME)
			audio.PlayAttention("../resources/phone_call.wav")
			doChat(COMMAND_DURATION)
		} else {
			log.Println("Channel has been closed.")
			break
		}
	}
}

func doChat(commandDuration int) {
	lastChatTime := time.Now()
	for {
		command, err := getUserCommand(commandDuration)
		if err != nil {
			log.Printf("getUserCommand error %s", err.Error())
		}
		if strings.Contains(command, "安静") {
			log.Println("chat abort,see you again")
			return
		}
		if len(command) <= 2 {
			if time.Now().Compare(lastChatTime.Add(time.Duration(2*commandDuration)*time.Second)) > 0 {
				log.Println("chat end,see you again")
				return
			}
			continue
		}
		chatResp := aigc.Chat(command)
		audio.PlayText(chatResp)
		audio.PlayAttention("../resources/game_finish.wav")
		lastChatTime = time.Now()
	}
}

func getUserCommand(listenSeconds int) (string, error) {
	recFileName := "/tmp/" + os.Args[0] + "-" + time.Now().Format("20060102150405") + ".wav"
	if err := record.RecBySox(recFileName, listenSeconds); err != nil {
		log.Printf("rec audio file error %s", err.Error())
	}
	asr, err := stt.Asr(recFileName)
	if err != nil {
		log.Println("speech to text error " + err.Error())
		return "", err
	}
	if len(asr) == 0 {
		return "", fmt.Errorf("speech to text empty")
	}
	return asr[0], nil
}
