package conf

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"os"
)

type Config struct {
	Sip       PjsipConfig
	ExtConfig []ExtConfig
	Samples   []SampleEntry
}

type PjsipConfig struct {
	Name            string
	Username        string
	Password        string
	Host            string
	SelectedContext string
}

func (c *Config) ToFile(path string) error {
	json, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(path, json, 0755); err != nil {
		return err
	}

	return nil
}

func FromDefaultFile() Config {
	return FromFile("/opt/robocall/cnf/config.json")
}

func (c *Config) ToDefaultFile() {
	if err := c.ToFile("/opt/robocall/cnf/config.json"); err != nil {
		log.Fatal(err)
	}
}

func FromFile(path string) Config {
	var fbuf bytes.Buffer
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	io.Copy(&fbuf, file)

	var c Config
	json.Unmarshal(fbuf.Bytes(), &c)
	return c
}

func (c *PjsipConfig) GetAsIni() string {
	tmpl := template.Must(template.ParseFiles("/opt/robocall/configs/pjsip.conf.tmpl"))

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, c)

	if err != nil {
		log.Fatal(err)
	}

	return tpl.String()
}

func (c *Config) ApplyConfig() error {
	c.ToDefaultFile()
	pjsip := c.Sip.GetAsIni()
	pjerr := os.WriteFile("/etc/asterisk/pjsip.conf", []byte(pjsip), 0755)

	if pjerr != nil {
		return pjerr
	} else {
		//svcctl.ReloadService()
		//c.ToFile("/opt/robocall/config.json")
		return nil
	}
}
