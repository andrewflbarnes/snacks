package main

import (
	"fmt"
	"os"

	"github.com/andrewflbarnes/snacks/internal/judy"
	"github.com/andrewflbarnes/snacks/internal/loris"
)

var (
	version string = "undefined"
)

func main() {
	if len(os.Args) < 2 {
		helpText("No options provided")
		return
	}

	first := os.Args[1]
	switch first {
	case "-v":
		versionText()
	case "-h":
		helpText("")
	case "judy":
		judy.Judy()
	case "loris":
		loris.Loris()
	default:
		helpText(fmt.Sprintf("Unsupported option \"%s\"", first))
	}
}

func versionText() {
	fmt.Println(version)
}

func helpText(msg string) {
	if len(msg) > 0 {
		fmt.Println(msg)
		fmt.Println()
	}
	fmt.Println("Usage:")
	fmt.Println("  attack:  snacks <attack type> [<attack option>, ...] <target>")
	fmt.Println("  help:    snacks -h")
	fmt.Println("  version: snacks -v")
	fmt.Println()
	fmt.Println("Attack types:")
	fmt.Println("  judy:    RUDY style attacks defaulting to application/json")
	fmt.Println(`  loris:   Slow Loris style attacks, defaulting to "x-snacks: slowloris"`)
	fmt.Println()
	fmt.Println("Attack options:")
	fmt.Println("  -v:      Enable debug logging")
	fmt.Println("  -vv:     Enable trace logging")
	fmt.Println("  -j:      Enable JSON logging")
	fmt.Println("  For a full list of attack specific options see the help for each attack e.g.")
	fmt.Println("  snacks loris -h")
	fmt.Println()
	fmt.Println("Target:")
	fmt.Println("  The URL of the endpoints to target e.g.")
	fmt.Println("  http://locallhost/vulnerable/endpoint")
	fmt.Println()
	fmt.Println("When not specified the URL the below defaults are used")
	fmt.Println("  scheme:  http")
	fmt.Println("  host:    localhost")
	fmt.Println("  port:    80")
	fmt.Println("  path:    /")
}
