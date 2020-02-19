package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var envVars = []string{
	"ZNC_VERSION",
	"ZNC_LISTENER_PORT",
	"ZNC_LISTENER_IPV4",
	"ZNC_LISTENER_IPV6",
	"ZNC_LISTENER_SSL",
	"ZNC_LOADMODULE",
	"ZNC_USER_NAME",
	"ZNC_USER_PASS",
	"ZNC_USER_PASS_SALT",
	"ZNC_USER_NICK",
	"ZNC_USER_ALTNICK",
	"ZNC_USER_IDENT",
	"ZNC_USER_LOADMODULE",
	"ZNC_NETWORK_NAME",
	"ZNC_NETWORK_LOADMODULE",
	"ZNC_NETWORK_SERVER",
	"ZNC_NETWORK_PORT",
	"ZNC_NETWORK_SSL",
	"ZNC_NETWORK_CHANS",
}

func tmplRender(ts []string, i interface{}) string {
	t := strings.Join(ts, "\n")
	b := bytes.Buffer{}
	tmpl, _ := template.New("test").Parse(t)
	tmpl.Execute(&b, i)

	return b.String()
}

type zncConf struct {
	Version  string
	Listener zncListener
	Modules  []string
	User     zncUser
}

func (c zncConf) String() string {
	ts := []string{
		"Version = {{.Version}}",
		"{{.Listener}}",
		"{{- range .Modules }}",
		"LoadModule = {{. -}}",
		"{{end}}",
		"",
		"{{.User}}",
	}
	return tmplRender(ts, c)
}

type zncListener struct {
	Port string
	IPV4 string
	IPV6 string
	SSL  string
}

func (l zncListener) String() string {
	ts := []string{
		"<Listener l>",
		"        Port = {{.Port}}",
		"        IPv4 = {{.IPV4}}",
		"        IPv6 = {{.IPV6}}",
		"        SSL = {{.SSL}}",
		"</Listener>",
	}
	return tmplRender(ts, l)
}

type zncUser struct {
	Name    string
	Pass    string
	Nick    string
	AltNick string
	Ident   string
	Modules []string
	Network zncNetwork
}

func (u zncUser) String() string {
	ts := []string{
		"<User {{.Name}}>",
		"        Pass       = {{.Pass}}",
		"        Admin      = true",
		"        Nick       = {{.Nick}}",
		"        AltNick    = {{.AltNick}}",
		"        Ident      = {{.Ident}}",
		"        {{- range .Modules }}",
		"        LoadModule = {{. -}}",
		"        {{end}}",
		"",
		"{{.Network}}",
		"</User>",
	}
	return tmplRender(ts, u)
}

type zncNetwork struct {
	Name    string
	Modules []string
	Server  string
	Port    string
	SSL     string
	Chans   []string
}

func (n zncNetwork) String() string {
	ts := []string{
		"        <Network {{.Name}}>",
		"                {{- range .Modules }}",
		"                LoadModule = {{. -}}",
		"                {{end}}",
		"                Server     = {{.Server}} {{.Port}}",
		"",
		"                {{ range .Chans}}",
		"                <Chan {{.}}>",
		"                </Chan>",
		"                {{- end}}",
		"        </Network>",
	}
	return tmplRender(ts, n)
}

func hashPass(pass, salt string) string {
	ps := pass + salt
	s := sha256.Sum256([]byte(ps))

	return fmt.Sprintf("sha256#%x#%s#", s, salt)
}

func main() {
	envs := map[string]string{}
	for _, i := range envVars {
		v := os.Getenv(i)

		if v == "" {
			fmt.Fprintf(os.Stderr, "%s environment variable not set\n", i)
			os.Exit(1)
		}

		envs[i] = v
	}

	conf := zncConf{
		Version: envs["ZNC_VERSION"],
		Listener: zncListener{
			Port: envs["ZNC_LISTENER_PORT"],
			IPV4: envs["ZNC_LISTENER_IPV4"],
			IPV6: envs["ZNC_LISTENER_IPV6"],
			SSL:  envs["ZNC_LISTENER_SSL"],
		},
		Modules: strings.Split(envs["ZNC_LOADMODULE"], " "),
		User: zncUser{
			Name:    envs["ZNC_USER_NAME"],
			Pass:    hashPass(envs["ZNC_USER_PASS"], envs["ZNC_USER_PASS_SALT"]),
			Nick:    envs["ZNC_USER_NICK"],
			AltNick: envs["ZNC_USER_ALTNICK"],
			Ident:   envs["ZNC_USER_IDENT"],
			Modules: strings.Split(envs["ZNC_USER_LOADMODULE"], " "),
			Network: zncNetwork{
				Name:    envs["ZNC_NETWORK_NAME"],
				Modules: strings.Split(envs["ZNC_NETWORK_LOADMODULE"], " "),
				Server:  envs["ZNC_NETWORK_SERVER"],
				Port:    envs["ZNC_NETWORK_PORT"],
				SSL:     envs["ZNC_NETWORK_SSL"],
				Chans:   strings.Split(envs["ZNC_NETWORK_CHANS"], " "),
			},
		},
	}

	fmt.Println(conf)
}
