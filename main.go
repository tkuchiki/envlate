package main

import (
	"bytes"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
)

func getFp(filename string) (*os.File, error) {
	var f *os.File
	var err error
	stdinfi, err := os.Stdin.Stat()
	if err != nil {
		return f, err
	}

	if stdinfi.Mode()&os.ModeNamedPipe == 0 {
		f, err = os.Open(filename)
		if err != nil {
			return f, err
		}
	} else {
		f = os.Stdin
	}

	return f, nil
}

func getEnvMap() map[string]string {
	envs := make(map[string]string)

	for _, e := range os.Environ() {
		env := strings.SplitN(e, "=", 2)
		envs[env[0]] = env[1]
	}

	return envs
}

func renderTemplate(tplb []byte, envs map[string]string, unsetErr bool) (string, error) {
	tpl, err := template.New("").Parse(string(tplb))
	if err != nil {
		return "", err
	}

	option := "zero"
	if unsetErr {
		option = "error"
	}

	var txt bytes.Buffer
	err = tpl.Option("missingkey="+option).Execute(&txt, envs)
	if err != nil {
		re := regexp.MustCompile(`template: :([0-9]+):([0-9]+): executing .+: map has no entry for key "(.+)"`)
		group := re.FindStringSubmatch(err.Error())
		if len(group) != 4 {
			return "", err
		}
		err = fmt.Errorf("line %s char %s : no entry for key `%s`", group[1], group[2], group[3])
	}
	return txt.String(), err
}

var (
	file  = kingpin.Flag("file", "Template file").Short('f').String()
	unset = kingpin.Flag("unset-error", "Treat unset variables as an error").Short('u').Bool()
)

func main() {
	kingpin.CommandLine.Help = "Expand environment variables in template (the templates use Go text/template syntax)"
	kingpin.Version("0.1.0")
	kingpin.Parse()

	f, ferr := getFp(*file)
	if ferr != nil {
		log.Fatal(ferr)
	}
	defer f.Close()

	b, berr := ioutil.ReadAll(f)
	if berr != nil {
		log.Fatal(berr)
	}

	envs := getEnvMap()

	txt, terr := renderTemplate(b, envs, *unset)
	if terr != nil {
		log.Fatal(terr)
	}

	fmt.Print(txt)
}
