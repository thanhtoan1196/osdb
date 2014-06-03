package main

import (
	"fmt"
	"path"

	"github.com/docopt/docopt.go"
	"github.com/oz/osdb"
)

// Get an anonymous client connected to OSDB.
func getClient() (client *osdb.Client, err error) {
	if client, err = osdb.NewClient(); err != nil {
		return
	}
	if err = client.LogIn("", "", ""); err != nil {
		return
	}
	return
}

func Get(file string, lang string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	fmt.Printf("- Getting subtitles for file: %s\n", path.Base(file))
	subs, err := client.FileSearch(file, []string{lang})
	if err != nil {
		return err
	}
	if best := subs.Best(); best != nil {
		dest := file[0:len(file)-len(path.Ext(file))] + ".srt"
		fmt.Printf("- Downloading to: %s\n", dest)
		// FIXME check if dest exists instead of overwriting. O:)
		return client.DownloadTo(best, dest)
	}
	return fmt.Errorf("No subtitles found!")
}

func main() {
	usage := `OSDB, an OpenSubtitles client.

Usage:
	osdb get [--language=<lang>] <file>
	osdb -h | --help
	osdb --version

Options:
	--language=<lang>	Language when searching subtitles [default: ENG].
`
	arguments, err := docopt.Parse(usage, nil, true, "OSDB 0.1a", false)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}
	lang := "ENG"
	if arguments["--language"] != nil {
		lang = arguments["--language"].(string)
	}

	if arguments["get"] == true {
		if err = Get(arguments["<file>"].(string), lang); err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}
}
