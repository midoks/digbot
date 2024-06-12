package cmd

import (
	"fmt"
	// "net"
	"time"
	// "strconv"
	// "sync"
	"runtime"
	// "syscall"

	"digbot/core/slog"
	"digbot/core/gonmap"

	"github.com/urfave/cli"

)

var Scan = cli.Command{
	Name:        "scan",
	Usage:       "This command scan ip all port service",
	Description: `scan ip port`,
	Action:      ScanRun,
	Flags: []cli.Flag{
		stringFlag("ip, i", "", "ip address"),
	},
}

func ScanRun(c *cli.Context) error {
	slog.Println(slog.INFO, "当前环境为：", runtime.GOOS)
	slog.Println(slog.INFO, "所有扫描任务已下发完毕")

	ip := c.String("ip")

	start := time.Now()

	nmap := gonmap.New()
	status, response := nmap.ScanTimeout(ip, 80, 20 * time.Second)
	fmt.Println(ip, 80,status, response)

	// status, response = nmap.ScanTimeout(ip, 80, 200000)
	// fmt.Println(ip, 80,status, response)
	
	time.Sleep(10)

	end := time.Since(start)

	
	fmt.Println("花费的时间：", end)
	return nil
}


