package main

import (
	"flag"
	"fmt"
	"os"

	ck "github.com/mp4096/clerk"
)

var a ck.Action = ck.NOTHING
var send = false
var fs = flag.NewFlagSet("clerk", flag.ExitOnError)
var configFilename = ""
var mdFilename = ""

func init() {
	fs.BoolVar(&send, "s", false, "Send emails; dry run otherwise")
	fs.StringVar(&configFilename, "c", "", "Config filename")
	fs.StringVar(&mdFilename, "m", "", "Markdown filename; use latest if not specified")
}

func main() {

	// Program flow:
	// - process input arguments :check:
	// - import configuration :check:
	// - read in and process the markdown file
	// - if not in dry run, read in the credentials
	// - if dry run, describe your actions and bail out
	// - IF APPROVE:
	//   - create an email object
	//   - send this email
	// - IF DISTRIBUTE:
	//   - create an email object (general transcript)
	//   - send this email
	//   - analyze individual recipients
	//   - create N customised email objects
	//   - send these emails

	if len(os.Args) == 1 {
		fmt.Println("usage: clerk <command> [<args>]")
		fmt.Println("Available commands are: ")
		fmt.Println(" approve     Send a request for approval")
		fmt.Println(" distribute  Distribute the transcripts")
		return
	}

	switch os.Args[1] {
	case "approve":
		a = ck.APPROVE
	case "distribute":
		a = ck.DISTRIBUTE
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
	fs.Parse(os.Args[2:])

	if len(configFilename) == 0 {
		fmt.Println("Error: config file not specified")
		os.Exit(3)
	}

	// TODO: Find latest Markdown file if none was specified
	c := new(ck.Config)
	errConf := c.ImportFromFile(configFilename)
	if errConf != nil {
		fmt.Println("Error opening config file")
		fmt.Println(errConf)
		os.Exit(3)
	}

	fmt.Println("Hello,", c.Author.Name)

	if err := ck.ProcessFile(mdFilename, a, send, c); err != nil {
		fmt.Println("Error sending email")
		fmt.Println(err)
		os.Exit(4)
	}
}
