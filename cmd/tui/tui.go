package tui

import (
	"fmt"
	"time"
	"tunl-cli/cmd/options"
)

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

func PrintConnectionScreen(opt options.Options, publicUrl, clientVer, serverVer string, expireAt int64) {
	fmt.Printf("🚀 " + BlueBold + "tunl started!" + Reset + "\n\n")
	fmt.Printf("Version              %s\n", clientVer)
	fmt.Printf("Server Version       %s\n", serverVer)
	fmt.Printf("Session expired at   %s\n", time.Now().Add(time.Duration(expireAt)).Format("2006-01-02 15:04:05"))
	if opt.BasicAuth != nil {
		fmt.Printf("BasicAuth            %s:%s\n", opt.BasicAuth.Login, opt.BasicAuth.Pass)
	}
	if opt.Monitor {
		fmt.Printf("Web monitor          %s\n", opt.MonitorAddr.ToProtoString())
	}
	fmt.Printf("Forwarding           %s -> %s\n\n", opt.LocalAddr.ToProtoString(), publicUrl)
	fmt.Printf("Docs                 https://github.com/black40x/tunl.online/\n\n")
	fmt.Printf(Yellow + "HTTP Requests: " + Reset + " \n\n")
}

func PrintURL(m, u string, d time.Duration) {
	fmt.Printf(Yellow+"[%.1fms] "+BlueBold+"[%s]"+Reset+" %s\n", float64(d.Nanoseconds())/1e6, m, u)
}

func PrintError(err error) {
	fmt.Printf(RedBold+"Error: "+Reset+"%s\n", err.Error())
}

func PrintInfo(s string) {
	fmt.Printf(BlueBold+"Info: "+Reset+"%s\n", s)
}

func PrintWarning(s string) {
	fmt.Printf(YellowBold+"Warning: "+Reset+"%s\n", s)
}

func PrintLn(s string) {
	fmt.Println(s)
}
