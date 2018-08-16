package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mp4096/clerk"
)

const (
	ok                       int = 0
	noConfigFileSpecified    int = 3
	invalidCommand           int = 4
	couldNotOpenConfigFile   int = 5
	couldNotSendOrPreview    int = 6
	couldNotFindMarkdownFile int = 7
)

var a clerk.Action = clerk.Nothing
var send = false
var fs = flag.NewFlagSet("clerk", flag.ExitOnError)
var configFilename = ""
var mdFilename = ""

func init() {
	fs.BoolVar(&send, "s", false, "Send emails; dry run otherwise")
	fs.StringVar(&configFilename, "c", "", "Config filename")
	fs.StringVar(&mdFilename, "m", "", "Markdown filename; if not specified, will try to use the latest file starting with an ISO timestamp")
}

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "approve":
		a = clerk.Approve
	case "distribute":
		a = clerk.Distribute
	case "help":
		printHelp()
		return
	default:
		fmt.Printf("%q is not valid command\n", os.Args[1])
		os.Exit(invalidCommand)
	}
	fs.Parse(os.Args[2:])

	if len(configFilename) == 0 {
		fmt.Println("Config file not specified")
		os.Exit(noConfigFileSpecified)
	}

	c := new(clerk.Config)
	errConf := c.ImportFromFile(configFilename)
	if errConf != nil {
		fmt.Println("Error opening config file")
		fmt.Println(errConf)
		os.Exit(couldNotOpenConfigFile)
	}

	fmt.Println("Hello,", c.Author.Name)

	if len(mdFilename) == 0 {
		fmt.Println("Markdown file not specified - will try to use the latest one")
		mdFilenameAuto, err := clerk.FindLatestFile()
		if err != nil {
			fmt.Println(err)
			os.Exit(couldNotFindMarkdownFile)
		}
		fmt.Printf("Found file '%s'\n", mdFilenameAuto)
		mdFilename = mdFilenameAuto
	}

	if err := clerk.ProcessFile(mdFilename, a, send, c); err != nil {
		fmt.Println("Error while sending email or opening preview")
		fmt.Println(err)
		os.Exit(couldNotSendOrPreview)
	}
}

func printHelp() {
	fmt.Println("usage: clerk <command> [<args>]")
	fmt.Println("Available commands are: ")
	fmt.Println(" approve     Send a request for approval")
	fmt.Println(" distribute  Distribute the transcripts")
	fmt.Println(" help        Print this message")
	fmt.Println("Call clerk <command> -h for more information about specific command usage")
}
