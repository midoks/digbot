package cmd

import (
	"fmt"
	"net"
	"time"
	"strconv"
	"sync"
	"runtime"

	"github.com/urfave/cli"
)

var Scan = cli.Command{
	Name:        "scan",
	Usage:       "This command scan ip all port service",
	Description: `scan ip all port`,
	Action:      ScanRun,
	Flags: []cli.Flag{
		stringFlag("ip, i", "", "ip address"),
	},
}

// taskschan 中存储要扫描的端口，
// reschan 存储开放的端口号
// exitchan 存储当前的goroutine是否完成的状态
// wgscan 同步goroutine
func IpPortScan(ip string, taskschan chan int, reschan chan int, exitchan chan bool, wgscan *sync.WaitGroup) {
	defer func() {
		// fmt.Println("任务完成")
		exitchan <- true
		wgscan.Done()
	}()
	// fmt.Println("开始任务")
	for {
		port, ok := <-taskschan
		if !ok {
			break
		}
		_, err := net.DialTimeout("tcp", ip+":"+strconv.Itoa(port), time.Second)
		if err == nil {
			reschan <- port
			fmt.Println("开放的端口", port)
		}else {
			// fmt.Println("关闭的端口", port)
		}
	}
}

// func ScanRun(c *cli.Context) error {

// 	fmt.Println(c.String("ip"))
// 	ip := c.String("ip")
// 	var res []int = make([]int, 0)

// 	start := time.Now()
// 	for i := 1; i < 65536; i++ {
// 		// address := "127.0.0.1:8080"
// 		address := fmt.Sprintf("%s:%d", ip, i)
// 		fmt.Println(address)

// 		_, err := net.DialTimeout("tcp", address, time.Second)
// 		if err == nil {
// 			fmt.Printf("端口 %s 开放\n", address)
// 			res = append(res, i)
// 		} else {
// 			fmt.Printf("端口 %s 关闭\n", address)
// 		}
// 		// time.Sleep(1 * time.Second)
// 	}
// 	end := time.Since(start)
// 	fmt.Println("花费的总时间", end)
// 	fmt.Println("开放的端口", res)
	
// 	return nil
// }

func ScanRun(c *cli.Context) error {

	// defaultports := [...]int{21, 22, 23, 25, 80, 443, 8080,
	// 	110, 135, 139, 445, 389, 489, 587, 1433, 1434,
	// 	1521, 1522, 1723, 2121, 3306, 3389, 4899, 5631,
	// 	5632, 5800, 5900, 7071, 43958, 65500, 4444, 8888,
	// 	6789, 4848, 5985, 5986, 8081, 8089, 8443, 10000,
	// 	6379, 7001, 7002}

	var defaultports []int = make([]int, 0)
	for i := 1; i < 65536; i++ {
		defaultports = append(defaultports, i)
	}
	// fmt.Println(defaultports)

	ip := c.String("ip")

	taskschan := make(chan int, len(defaultports))
	reschan := make(chan int, len(defaultports))
	goroutine_nums := runtime.NumCPU()
	exitchan := make(chan bool, goroutine_nums)
	var wgp sync.WaitGroup
	for _, value := range defaultports {
		taskschan <- value
	}
    //向taskschan中写完数据时，就需要关闭taskschan,否则goroutine会一直认为该channel会写入数据，会一直等待
	close(taskschan)

	start := time.Now()
    //开启四个goroutine执行扫描任务
	for i := 0; i < goroutine_nums; i++ {
		wgp.Add(1)
		go IpPortScan(ip, taskschan, reschan, exitchan, &wgp)
	}
	wgp.Wait()
    //判断4个goroutine是否都执行完了，当他们都执行完了写入到reschan，rechan才可以被关闭
	for i := 0; i < goroutine_nums; i++ {
		<-exitchan
	}
	end := time.Since(start)

	close(exitchan)
	close(reschan)
	for {
		openport, ok := <-reschan
		if !ok {
			break
		}
		fmt.Println("开放的端口：", openport)
	}
	fmt.Println("花费的时间：", end)
	return nil
}


