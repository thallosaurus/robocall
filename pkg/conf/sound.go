package conf

import (
	"fmt"
	"log"
	"os"
)

type SampleEntry struct {
	SoundName string
	File      *os.File
}

func FromSoundName(name string) SampleEntry {
	file, err := os.OpenFile(fmt.Sprintf("/var/lib/asterisk/sounds/en/%s.gsm", name), os.O_RDONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	return SampleEntry{
		SoundName: name,
		File:      file,
	}
}
