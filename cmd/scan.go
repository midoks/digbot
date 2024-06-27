package cmd

import (
	"fmt"
	"net"
	"time"
	"strconv"
	"sync"
	"runtime"
	"net/url"
	"strings"
	// "syscall"

	"digbot/app"
	"digbot/core/slog"
	"digbot/core/slog/stdio/chinese"
	"digbot/core/gonmap"
	"digbot/core/scanner"
	"digbot/lib/color"
	"digbot/lib/misc"
	// "digbot/lib/uri"
	"github.com/lcvvvv/simplehttp"

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

var (
	disableKey       = []string{"MatchRegexString", "Service", "ProbeName", "Response", "Cert", "Header", "Body", "IP"}
	importantKey     = []string{"ProductName", "DeviceType"}
	varyImportantKey = []string{"Hostname", "FingerPrint", "ICP"}
)

func getTimeout(i int) time.Duration {
	switch {
	case i > 10000:
		return time.Millisecond * 200
	case i > 5000:
		return time.Millisecond * 300
	case i > 1000:
		return time.Millisecond * 400
	default:
		return time.Millisecond * 500
	}
}

func outputOpenResponse(addr net.IP, port int) {
	//输出结果
	protocol := gonmap.GuessProtocol(port)
	target := fmt.Sprintf("%s://%s:%d", protocol, addr.String(), port)
	URL, _ := url.Parse(target)
	outputHandler(URL, "response is empty", map[string]string{
		"IP":   URL.Hostname(),
		"Port": strconv.Itoa(port),
	})
}

func outputNmapFinger(URL *url.URL, resp *gonmap.Response) {
	if responseFilter(resp.Raw) == true {
		return
	}
	finger := resp.FingerPrint
	m := misc.ToMap(finger)
	m["Response"] = resp.Raw
	m["IP"] = URL.Hostname()
	m["Port"] = URL.Port()
	//补充归属地信息
	// if app.Setting.CloseCDN == false {
	// 	result, _ := cdn.Find(URL.Hostname())
	// 	m["Addr"] = result
	// }
	outputHandler(URL, finger.Service, m)
}

func outputHandler(URL *url.URL, keyword string, m map[string]string) {
	m = misc.FixMap(m)
	if respRaw := m["Response"]; respRaw != "" {
		if m["Service"] == "http" || m["Service"] == "https" {
			m["Digest"] = strconv.Quote(getHTTPDigest(respRaw))
		} else {
			m["Digest"] = strconv.Quote(getRawDigest(respRaw))
		}
	}
	m["Length"] = strconv.Itoa(len(m["Response"]))
	sourceMap := misc.CloneMap(m)
	for _, keyword := range disableKey {
		delete(m, keyword)
	}
	for key, value := range m {
		if key == "FingerPrint" {
			continue
		}
		m[key] = misc.StrRandomCut(value, 24)
	}
	fingerPrint := color.StrMapRandomColor(m, true, importantKey, varyImportantKey)
	fingerPrint = misc.FixLine(fingerPrint)
	format := "%-30v %-" + strconv.Itoa(misc.AutoWidth(color.Clear(keyword), 26+color.Count(keyword))) + "v %s"
	printStr := fmt.Sprintf(format, URL.String(), keyword, fingerPrint)
	slog.Println(slog.DATA, printStr)

	if jw := app.Setting.OutputJson; jw != nil {
		sourceMap["URL"] = URL.String()
		sourceMap["Keyword"] = keyword
		jw.Push(sourceMap)
	}
	if cw := app.Setting.OutputCSV; cw != nil {
		sourceMap["URL"] = URL.String()
		sourceMap["Keyword"] = keyword
		delete(sourceMap, "Header")
		delete(sourceMap, "Cert")
		delete(sourceMap, "Response")
		delete(sourceMap, "Body")
		sourceMap["Digest"] = strconv.Quote(sourceMap["Digest"])
		for key, value := range sourceMap {
			sourceMap[key] = chinese.ToUTF8(value)
		}
		cw.Push(sourceMap)
	}
}


func getHTTPDigest(s string) string {
	var length = 24
	var digestBuf []rune
	_, body := simplehttp.SplitHeaderAndBody(s)
	body = chinese.ToUTF8(body)
	for _, r := range []rune(body) {
		buf := []byte(string(r))
		if len(digestBuf) == length {
			return string(digestBuf)
		}
		if len(buf) > 1 {
			digestBuf = append(digestBuf, r)
		}
	}
	return string(digestBuf) + misc.StrRandomCut(body, length-len(digestBuf))
}

func responseFilter(strArgs ...string) bool {
	var match = app.Setting.Match
	var notMatch = app.Setting.NotMatch

	if match != "" {
		for _, str := range strArgs {
			//主要结果中包含关键则，则会显示
			if strings.Contains(str, app.Setting.Match) == true {
				return false
			}
		}
	}

	if notMatch != "" {
		for _, str := range strArgs {
			//主要结果中包含关键则，则会显示
			if strings.Contains(str, app.Setting.NotMatch) == true {
				return true
			}
		}
	}
	return false
}


func getRawDigest(s string) string {
	var length = 24
	if len(s) < length {
		return s
	}
	var digestBuf []rune
	for _, r := range []rune(s) {
		if len(digestBuf) == length {
			return string(digestBuf)
		}
		if 0x20 <= r && r <= 0x7E {
			digestBuf = append(digestBuf, r)
		}
	}
	return string(digestBuf) + misc.StrRandomCut(s, length-len(digestBuf))
}
func outputUnknownResponse(addr net.IP, port int, response string) {
	if responseFilter(response) == true {
		return
	}
	//输出结果
	target := fmt.Sprintf("unknown://%s:%d", addr.String(), port)
	URL, _ := url.Parse(target)
	outputHandler(URL, "无法识别该协议", map[string]string{
		"Response": response,
		"IP":       URL.Hostname(),
		"Port":     strconv.Itoa(port),
	})
}

func generatePortScanner(wg *sync.WaitGroup) *scanner.PortClient {
	PortConfig := scanner.DefaultConfig()
	PortConfig.Threads = app.Setting.Threads
	PortConfig.Timeout = getTimeout(len(app.Setting.Port))
	if app.Setting.ScanVersion == true {
		PortConfig.DeepInspection = true
	}
	client := scanner.NewPortScanner(PortConfig)
	client.HandlerClosed = func(addr net.IP, port int) {
		//nothing
	}
	client.HandlerOpen = func(addr net.IP, port int) {
		outputOpenResponse(addr, port)
	}
	client.HandlerNotMatched = func(addr net.IP, port int, response string) {
		outputUnknownResponse(addr, port, response)
	}
	client.HandlerMatched = func(addr net.IP, port int, response *gonmap.Response) {
		URLRaw := fmt.Sprintf("%s://%s:%d", response.FingerPrint.Service, addr.String(), port)
		URL, _ := url.Parse(URLRaw)
		// if appfinger.SupportCheck(URL.Scheme) == true {
			// pushURLTarget(URL, response)
			// return
		// }
		outputNmapFinger(URL, response)
	}
	client.HandlerError = func(addr net.IP, port int, err error) {
		slog.Println(slog.DEBUG, "PortScanner Error: ", fmt.Sprintf("%s:%d", addr.String(), port), err)
	}
	client.Defer(func() {
		wg.Done()
	})
	return client
}

var (
	PortScanner   *scanner.PortClient
)

func ScanRun(c *cli.Context) error {

	var wg = &sync.WaitGroup{}
	wg.Add(1)

	PortScanner = generatePortScanner(wg)

	wg.Wait()


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


