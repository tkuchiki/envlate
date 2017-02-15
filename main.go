package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	re      = regexp.MustCompile(`template: :([0-9]+):([0-9]+): executing .+: map has no entry for key "(.+)"`)
	funcMap = template.FuncMap{
		"default": func(def, val string) string {
			if val != "" {
				return val
			}
			return def
		},
	}
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
	option := "zero"
	if unsetErr {
		option = "error"
	}

	t := template.New("").Option("missingkey=" + option).Funcs(funcMap)
	tpl, err := t.Parse(string(tplb))
	if err != nil {
		return "", err
	}

	var txt bytes.Buffer
	err = tpl.Execute(&txt, envs)
	if err != nil {
		group := re.FindStringSubmatch(err.Error())
		if len(group) != 4 {
			return "", err
		}
		err = fmt.Errorf("line %s char %s : no entry for key `%s`", group[1], group[2], group[3])
	}
	return txt.String(), err
}

func stringToFileMode(mode string) (os.FileMode, error) {
	m, err := strconv.ParseInt(mode, 8, 0)
	if err != nil {
		var fmode os.FileMode
		return fmode, err
	}

	return os.FileMode(m), err
}

func writeFile(fpath, data string, perm os.FileMode) error {
	return ioutil.WriteFile(fpath, []byte(data), perm)
}

var (
	file   = kingpin.Flag("file", "Template file").Short('f').String()
	unset  = kingpin.Flag("unset-error", "Treat unset variables as an error").Short('u').Bool()
	output = kingpin.Flag("output", "Write the output to the file rather than to stdout").Short('o').PlaceHolder("FILENAME").String()
	mode   = kingpin.Flag("mode", "File permission").Default("0644").String()
)

func main() {
	kingpin.CommandLine.Help = "Expand environment variables in template (the templates use Go text/template syntax)"
	kingpin.Version("0.1.1")
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

	if *output != "" {
		m, merr := stringToFileMode(*mode)
		if merr != nil {
			log.Fatal(merr)
		}

		werr := writeFile(*output, txt, m)
		if werr != nil {
			log.Fatal(werr)
		}
	} else {
		fmt.Print(txt)
	}
}
