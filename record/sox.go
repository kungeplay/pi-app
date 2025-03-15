package record

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
)

func RecBySox(fileName string, seconds int) error {
	log.Printf("begin rec by sox fileName:%s,duration seconds: %d", fileName, seconds)
	//设置音量大小
	recCmd := exec.Command("rec", "-r", "16k", "-c", "1", "-b", "16", fileName, "trim", "0", strconv.Itoa(seconds))
	var recBuffer bytes.Buffer
	recCmd.Stdout = &recBuffer
	if err := recCmd.Start(); err != nil {
		log.Printf("rec command start error:%s", err.Error())
		return err
	}
	if err := recCmd.Wait(); err != nil {
		log.Printf("rec command wait error:%s", err.Error())
		return err
	}
	log.Printf("end rec by sox fileName:%s,duration seconds: %d", fileName, seconds)
	return nil
}
