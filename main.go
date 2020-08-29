package nginx-rtmp-stat

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	RTMP struct {
		NginxVersion     string `xml:"nginx_version"`
		NginxRTMPVersion string `xml:"nginx_rtmp_version"`
		Built            string `xml:"built"`
		PID              int    `xml:"pid"`
		Uptime           int    `xml:"uptime"`
		Accepted         int    `xml:"naccepted"`
		BwIn             int    `xml:"bw_in"`
		BytesIn          int    `xml:"bytes_in"`
		BytesOut         int    `xml:"bytes_out"`
		Server           Server `xml:"server"`
	}
	Server struct {
		Application []Application `xml:"application"`
	}
	Application struct {
		Name string `xml:"name"`
		Live Live   `xml:"live"`
	}
	Live struct {
		Clients int      `xml:"nclients"`
		Stream  []Stream `xml:"stream"`
	}
	Stream struct {
		Name     string   `xml:"name"`
		Time     int      `xml:"time"`
		BwIn     int      `xml:"bw_in"`
		BytesIn  int      `xml:"bytes_in"`
		BwOut    int      `xml:"bw_out"`
		BytesOut int      `xml:"bytes_out"`
		BwAudio  int      `xml:"bw_audio"` // Might be bitrate
		BwVideo  int      `xml:"bw_video"`
		Client   []Client `xml:"client"`
		Meta     Meta     `xml:"meta"`
	}
	Client struct {
		ID         int           `xml:"id"`
		Address    string        `xml:"address"` // IP address
		Time       int           `xml:"time"`
		FlashVer   string        `xml:"flashver"`
		Dropped    int           `xml:"dropped"`
		AVSync     int           `xml:"avsync"`
		Timestamp  int           `xml:"timestamp"`
		Active     BoolIfPresent `xml:"active"`
		Publishing BoolIfPresent `xml:"publishing"`
	}
	Meta struct {
		Video Video `xml:"video"`
		Audio Audio `xml:"audio"`
	}
	Video struct {
		Width     int    `xml:"width"`
		Height    int    `xml:"height"`
		FrameRate int    `xml:"frame_rate"`
		Codec     string `xml:"codec"`
		Profile   string `xml:"profile"`
		Level     string `xml:"level"`
	}
	Audio struct {
		Codec      string `xml:"codec"`
		Profile    string `xml:"profile"`
		Channels   int    `xml:"channels"`
		SampleRate int    `xml:"sample_rate"`
	}
	BoolIfPresent bool
)

func (c *BoolIfPresent) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	*c = true
	return nil
}
func FromURL(url string) (*RTMP, error) {
	res, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("failed to get stats: %w", err)
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("failed to read body: %w", err)
		return nil, err
	}
	rtmp := &RTMP{}
	err = xml.Unmarshal(b, &rtmp)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal xml: %w", err)
		return nil, err
	}
	return rtmp, nil
}
