package record

import (
	"encoding/binary"
	"github.com/gordonklaus/portaudio"
	"os"
	"strings"
	"time"
)

func RecByPortAudio(fileName string, seconds int) error {
	afterSig := time.After(time.Duration(seconds) * time.Second)

	if !strings.HasSuffix(fileName, ".aiff") {
		fileName += ".aiff"
	}
	f, err := os.Create(fileName)
	chk(err)

	// form chunk
	_, err = f.WriteString("FORM")
	chk(err)
	chk(binary.Write(f, binary.BigEndian, int32(0))) //total bytes
	_, err = f.WriteString("AIFF")
	chk(err)

	// common chunk
	_, err = f.WriteString("COMM")
	chk(err)
	chk(binary.Write(f, binary.BigEndian, int32(18)))                  //size
	chk(binary.Write(f, binary.BigEndian, int16(1)))                   //channels
	chk(binary.Write(f, binary.BigEndian, int32(0)))                   //number of samples
	chk(binary.Write(f, binary.BigEndian, int16(32)))                  //bits per sample
	_, err = f.Write([]byte{0x40, 0x0e, 0xac, 0x44, 0, 0, 0, 0, 0, 0}) //80-bit sample rate 44100
	chk(err)

	// sound chunk
	_, err = f.WriteString("SSND")
	chk(err)
	chk(binary.Write(f, binary.BigEndian, int32(0))) //size
	chk(binary.Write(f, binary.BigEndian, int32(0))) //offset
	chk(binary.Write(f, binary.BigEndian, int32(0))) //block
	nSamples := 0
	defer func() {
		// fill in missing sizes
		totalBytes := 4 + 8 + 18 + 8 + 8 + 4*nSamples
		_, err = f.Seek(4, 0)
		chk(err)
		chk(binary.Write(f, binary.BigEndian, int32(totalBytes)))
		_, err = f.Seek(22, 0)
		chk(err)
		chk(binary.Write(f, binary.BigEndian, int32(nSamples)))
		_, err = f.Seek(42, 0)
		chk(err)
		chk(binary.Write(f, binary.BigEndian, int32(4*nSamples+8)))
		chk(f.Close())
	}()

	portaudio.Initialize()
	defer portaudio.Terminate()
	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
	for {
		chk(stream.Read())
		chk(binary.Write(f, binary.BigEndian, in))
		nSamples += len(in)
		select {
		case <-afterSig:
			return nil
		default:
		}
	}
	chk(stream.Stop())

	return nil
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
