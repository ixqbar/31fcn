package spider

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strings"
)

type TF31User struct {
	Name string
	Password string
}

func (obj *TF31User) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var content string
	if err := d.DecodeElement(&content, &start); err != nil {
		return err
	}

	if len(content) > 0 && strings.Index(content, ",") >= 0 {
		userInfo := strings.Split(content, ",")
		obj.Name = userInfo[0]
		obj.Password = userInfo[1]
		return nil
	}

	return errors.New("error f31user config")
}

type TSpiderConfig struct {
	F31User TF31User `xml:"f3cn"`
	TargetURL string `xml:"url"`
	PullOnStartup bool `xml:"task>startup"`
	Schedule string `xml:"task>schedule"`
	RedisServerAddress string `xml:"redis_server"`
}

var GTSpiderConfig *TSpiderConfig

func ParseXmlConfig(path string) (*TSpiderConfig, error) {
	if len(path) == 0 {
		return nil, errors.New("not found configure xml file")
	}

	n, err := GetFileSize(path)
	if err != nil || n == 0 {
		return nil, errors.New("not found configure xml file")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	GTSpiderConfig = &TSpiderConfig{}

	data := make([]byte, n)

	m, err := f.Read(data)
	if err != nil {
		return nil, err
	}

	if int64(m) != n {
		return nil, errors.New(fmt.Sprintf("expect read configure xml file size %d but result is %d", n, m))
	}

	err = xml.Unmarshal(data, &GTSpiderConfig)
	if err != nil {
		return nil, err
	}

	Logger.Printf("read config %+v", GTSpiderConfig)

	return GTSpiderConfig, nil
}
