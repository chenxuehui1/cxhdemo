package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"


	"github.com/chenxuehui1/cxhdemo/service"
	"github.com/chenxuehui1/cxhdemo/service/process"
)

var opt level.Option

var cliType = flag.String("clitype", "master", "client tyep")
var bindPort = flag.Int("bindPort", 7946, "gossip bindport")
var advertisePort = flag.Int("advertisePort", 7946, "gossip advertisePort")
var hostname = flag.String("peerName", "cicd-jenkins-srv-agent.devops.svc", "peer hostname")
var knownPeers = flag.String("knownPeers", "jenkins-watchmen.devops.svc:7946", "gossip need seed node")
var logLevel = flag.String("logLevel", "info", "output log level")


func main() {
	flag.Parse()

		if os.Getenv("LOG_LEVEL") != "" {
		*logLevel = os.Getenv("LOG_LEVEL")
	}

	setLogLevel(*logLevel)

	logger := log.NewLogfmtLogger(os.Stdout)
	logger = level.NewFilter(logger, opt)

	switch *cliType {
	case "cxhdemo":
		fmt.Println("this is cxhdemo")
		pr, err := service.CreatePeerCxhdemo(*bindPort, *advertisePort, *hostname, *knownPeers)
		if err != nil {
			panic("init memberlist err==>" + err.Error())
		}
		pr.Join()

		//捕获系统退出信号，退出前执行Leave()安全退出
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-stop
		level.Debug(logger).Log("msg", "Receive system kill signal peer leave")
		pr.Leave()
		return
	default:
		level.Error(logger).Log("clitype is", *cliType)
	}

}
func init() {
	http.Handle("/metrics", promhttp.Handler())
	//log.Fatal(http.ListenAndServe(":8092", nil))

}

func setLogLevel(logLevel string) error {
	switch logLevel {
	case "debug":
		opt = level.AllowDebug()
	case "info":
		opt = level.AllowInfo()
	case "warn":
		opt = level.AllowWarn()
	case "error":
		opt = level.AllowError()
	default:
		return fmt.Errorf("unrecognized log level %q", logLevel)
	}
	//service.Opt = op
	process.Opt = opt
	return nil
}
