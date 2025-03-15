package record

import (
	"encoding/binary"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"os"
	"testing"
	"time"
)

func TestRec(t *testing.T) {
	RecByPortAudio("test.aiff", 5)
}
func TestPortAudio(t *testing.T) {
	portaudio.Initialize()
	defer portaudio.Terminate()
	apis, err := portaudio.HostApis()
	if err != nil {
		panic(err)
	}
	for _, api := range apis {
		fmt.Println(api.Name + " " + api.Type.String())
	}
}

func TestRec2(t *testing.T) {
	f, err := os.Create("testRec2.wav")
	chk(err)
	afterSig := time.After(5 * time.Second)
	defer func() {
		chk(f.Close())
	}()

	portaudio.Initialize()
	time.Sleep(1)
	defer portaudio.Terminate()
	in := make([]int16, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(in), in)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
loop:
	for {
		chk(stream.Read())
		chk(binary.Write(f, binary.LittleEndian, in))
		select {
		case <-afterSig:
			break loop
		default:
		}
	}
	chk(stream.Stop())
}
