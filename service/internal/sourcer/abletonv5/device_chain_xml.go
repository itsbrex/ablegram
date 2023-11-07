package abletonv5

// XmlTrackDeviceChain is a struct that represents the XML structure of a track's device chain.
type XmlTrackDeviceChain struct {
	Mixer       XmlMixer             `xml:"Mixer"`
	DeviceChain XmlActualDeviceChain `xml:"DeviceChain"`
}

// XmlActualDeviceChain is a struct that represents the XML structure of a track's actual device chain.
type XmlActualDeviceChain struct {
	Mixer   XmlMixer      `xml:"Mixer"`
	Devices XmlDeviceList `xml:"Devices"`
}

type XmlDeviceList struct {
	Reverb          []XmlReverbDevice          `xml:"Reverb"`
	Delay           []XmlDelayDevice           `xml:"Delay"`
	MidiArpeggiator []XmlMidiArpeggiatorDevice `xml:"MidiArpeggiator"`
	MidiChord       []XmlMidiChordDevice       `xml:"MidiChord"`
}

func (dl *XmlDeviceList) GetCount() int {
	var c int
	c += len(dl.Reverb)
	c += len(dl.Delay)
	c += len(dl.MidiArpeggiator)

	return c
}
