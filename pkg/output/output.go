package output

import (
	"fmt"
	"os"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
	BgRed   = "\033[41m"
	BgGreen = "\033[42m"
)

func Banner() {
	fmt.Print(Cyan + Bold)
	fmt.Println(`
__________-------____                 ____-------__________
          \------____-------___--__---------__--___-------____------/
           \//////// / / / / / \   _-------_   / \ \ \ \ \ \\\\\\\\/
             \////-/-/------/_/_| /___   ___\ |_\_\------\-\-\\\\/
               --//// / /  /  //|| (O)\ /(O) ||\\  \  \ \ \\\\--
                    ---__/  // /| \_  /V\  _/ |\ \\  \__---
                         -//  / /\_ ------- _/\ \  \\-
                           \_/_/ /\---------/\ \_\_/
                               ----\   |   /----
                                    | -|- |
                                   /   |   \
                                   ---- \___|
`)
	fmt.Println(Reset)
	fmt.Println()
	fmt.Print(Cyan + Bold)
	fmt.Println("    ██╗  ██╗███████╗██╗  ██╗██████╗ ██╗      ██████╗ ██╗████████╗")
	fmt.Println("    ██║ ██╔╝██╔════╝╚██╗██╔╝██╔══██╗██║     ██╔═══██╗██║╚══██╔══╝")
	fmt.Println("    █████╔╝ █████╗   ╚███╔╝ ██████╔╝██║     ██║   ██║██║   ██║   ")
	fmt.Println("    ██╔═██╗ ██╔══╝   ██╔██╗ ██╔═══╝ ██║     ██║   ██║██║   ██║   ")
	fmt.Println("    ██║  ██╗███████╗██╔╝ ██╗██║     ███████╗╚██████╔╝██║   ██║   ")
	fmt.Println("    ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝     ╚══════╝ ╚═════╝ ╚═╝   ╚═╝")
	fmt.Println(Reset)
	fmt.Println(Dim + "    ┌─────────────────────────────────────────────────────────────────────────┐")
	fmt.Println("    │" + Reset + Magenta + "  Kernel Exploit Suggester & Auto-Exploiter v1.0" + Dim + "                      │")
	fmt.Println("    │" + Reset + Yellow + "  Author: " + Reset + White + "past3l@mileniumsec" + Dim + "                                          │")
	fmt.Println("    │" + Reset + Blue + "  GitHub: " + Reset + Cyan + "github.com/past3l" + Dim + "                                            │")
	fmt.Println("    │                       │")
	fmt.Println("    └─────────────────────────────────────────────────────────────────────────┘" + Reset)
	fmt.Println()
}

func Info(format string, a ...any) {
	fmt.Printf(Cyan+"[*] "+Reset+format+"\n", a...)
}

func Success(format string, a ...any) {
	fmt.Printf(Green+Bold+"[+] "+Reset+Green+format+Reset+"\n", a...)
}

func Error(format string, a ...any) {
	fmt.Printf(Red+"[-] "+Reset+format+"\n", a...)
}

func Warn(format string, a ...any) {
	fmt.Printf(Yellow+"[!] "+Reset+format+"\n", a...)
}

func Fatal(format string, a ...any) {
	fmt.Fprintf(os.Stderr, Red+Bold+"[FATAL] "+Reset+format+"\n", a...)
	os.Exit(1)
}

func RootBanner() {
	fmt.Println()
	fmt.Println(BgGreen + Black + Bold + "                                                                              " + Reset)
	fmt.Println(BgGreen + Black + Bold + "    ██████╗  ██████╗  ██████╗ ████████╗    ███████╗██╗  ██╗███████╗██╗     ██╗ " + Reset)
	fmt.Println(BgGreen + Black + Bold + "    ██╔══██╗██╔═══██╗██╔═══██╗╚══██╔══╝    ██╔════╝██║  ██║██╔════╝██║     ██║ " + Reset)
	fmt.Println(BgGreen + Black + Bold + "    ██████╔╝██║   ██║██║   ██║   ██║       ███████╗███████║█████╗  ██║     ██║ " + Reset)
	fmt.Println(BgGreen + Black + Bold + "    ██╔══██╗██║   ██║██║   ██║   ██║       ╚════██║██╔══██║██╔══╝  ██║     ██║ " + Reset)
	fmt.Println(BgGreen + Black + Bold + "    ██║  ██║╚██████╔╝╚██████╔╝   ██║       ███████║██║  ██║███████╗███████╗███████╗ " + Reset)
	fmt.Println(BgGreen + Black + Bold + "    ╚═╝  ╚═╝ ╚═════╝  ╚═════╝    ╚═╝       ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝ " + Reset)
	fmt.Println(BgGreen + Black + Bold + "                                                                              " + Reset)
	fmt.Println()
	fmt.Printf(Green + Bold + "    UID: %d | GID: %d | EUID: %d\n" + Reset, os.Getuid(), os.Getgid(), os.Geteuid())
	fmt.Println(Green + "    You are now root! 👑" + Reset)
	fmt.Println()
}

const Black = "\033[30m"
