package wake

import (
	"bytes"
	"encoding/binary"
	"github.com/brentnd/go-snowboy"
	"github.com/gordonklaus/portaudio"
	"log"
	"time"
)

// Sound represents a sound stream implementing the io.Reader interface
// that provides the microphone data.
type Sound struct {
	stream *portaudio.Stream
	data   []int16
}

// Init initializes the Sound's PortAudio stream.
func (s *Sound) Init() {
	inputChannels := 1
	outputChannels := 0
	sampleRate := 16000
	s.data = make([]int16, 1024)

	// initialize the audio recording interface
	err := portaudio.Initialize()
	if err != nil {
		log.Printf("Error initialize audio interface: %s", err.Error())
		return
	}

	// open the sound input stream for the microphone
	stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(s.data), s.data)
	if err != nil {
		log.Printf("Error open default audio stream: %s", err.Error())
		return
	}

	err = stream.Start()
	if err != nil {
		log.Printf("Error on stream start: %s", err.Error())
		return
	}

	s.stream = stream
}

// Close closes down the Sound's PortAudio connection.
func (s *Sound) Close() {
	s.stream.Close()
	portaudio.Terminate()
}

// Read is the Sound's implementation of the io.Reader interface.
func (s *Sound) Read(p []byte) (int, error) {
	s.stream.Read()

	buf := &bytes.Buffer{}
	for _, v := range s.data {
		binary.Write(buf, binary.LittleEndian, v)
	}

	copy(p, buf.Bytes())
	return len(p), nil
}

func ListenerHotWord(commonRes string, hotWordModelFile string, sensitivity float32, recChannel chan int) {
	// open the mic
	mic := &Sound{}
	mic.Init()
	defer mic.Close()

	// open the snowboy detector
	d := snowboy.NewDetector(commonRes)
	defer d.Close()

	// set the handlers
	d.HandleFunc(snowboy.NewHotword(hotWordModelFile, sensitivity), func(string) {
		log.Println("hit the hot word!")
		recChannel <- 1
	})

	d.HandleSilenceFunc(300*time.Millisecond, func(string) {
		log.Println("Silence detected")
	})

	// display the detector's expected audio format
	sr, nc, bd := d.AudioFormat()
	log.Printf("sample rate=%d, num channels=%d, bit depth=%d\n", sr, nc, bd)

	// start detecting using the microphone
	d.ReadAndDetect(mic)
}
