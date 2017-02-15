package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
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
	kingpin.Version("0.1.2")
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

	t := NewTemplate(*unset)

	txt, terr := t.render(b, envs)
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
