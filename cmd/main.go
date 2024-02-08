package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/api"
	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/helper"
	"github.com/System-Analysis-and-Design-2023-SUT/Server/internal/settings"
	models "github.com/System-Analysis-and-Design-2023-SUT/Server/models/queue"
	logging "github.com/System-Analysis-and-Design-2023-SUT/Server/pkg/logger"
	"github.com/hashicorp/memberlist"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/pflag"
)

var logger *logging.Logger

func init() {
	var err error
	logger, err = logging.NewLogger("sad server", false)
	if err != nil {
		log.Fatal("could not initialize main logger")
	}
}

var settingsPath string
var nodeName string

func main() {
	log.Print("Server is Starting ...")
	pflag.StringVar(&settingsPath, "settings", "/opt/server/settings.yml", "Path to settings file")
	pflag.StringVar(&nodeName, "nodeName", "1", "Name of node")
	pflag.Parse()

	var st settings.Settings

	err := cleanenv.ReadConfig(settingsPath, &st)
	if err != nil {
		logger.FatalS("Could not read settings", "error", err.Error())
	}
	_, err = st.IsValid()
	if err != nil {
		logger.FatalS("Setting file is not valid", "error", err.Error())
	}

	gossopingServer := setupGossopingServers(&st)
	go func() {
		runGossopingServer(gossopingServer, st.Global.GossopingPort, "gossoping_server")
	}()

	helper, err := helper.NewHelper(gossopingServer)
	if err != nil {
		logger.FatalS("Could not create helper", "error", err.Error())
	}
	q := models.NewQueue()
	s := models.NewSubscriber()

	internalAPIServer := setupHTTPServer(&st, helper, q, s)
	go func() {
		runHTTPServer(internalAPIServer, st.Global.APIPort, "api_server")
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGHUP, syscall.SIGINT)
	<-shutdown

	logger.Info("Shutting Down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := gossopingServer.Leave(time.Second * 5); err != nil {
		logger.Fatal("Could not shutdown gossoping server gracefully", "error", err.Error())
	}

	if err := internalAPIServer.Shutdown(ctx); err != nil {
		logger.Fatal("Could not shutdown internal api server gracefully", "error", err.Error())
	}

	fmt.Println("Shutting Down server...")

}

func removeSuffix(s string) string {
	re := regexp.MustCompile(`-\d+$`)
	return re.ReplaceAllString(s, "")
}

func setupGossopingServers(settings *settings.Settings) *memberlist.Memberlist {
	logger.InfoS("Initializing gossoping server.")
	logger.Infof("memberlist_server Starting listening on port %d.", settings.Global.MemberlistPort)

	config := memberlist.DefaultLocalConfig()
	config.BindPort = settings.Global.MemberlistPort
	config.BindAddr = settings.Replica.BindAddress

	config.Name = nodeName
	list, err := memberlist.Create(config)
	if err != nil {
		logger.Fatalf("Error initializing Cluster node with error %v", err)
		return nil
	}
	list.LocalNode().Meta = []byte(strconv.Itoa(settings.Global.APIPort))

	for _, host := range settings.Replica.Hostname {
		for id := 1; id <= settings.Replica.MemberCount; id++ {
			nodeName := fmt.Sprintf("%s-%d:8081", removeSuffix(host), id)
			_, err := list.Join([]string{nodeName})
			if err != nil {
				logger.Infof("Error joining Cluster node %s with error %v", nodeName, err)
			} else {
				logger.Infof("Connected to %s", nodeName)
			}
		}
	}
	return list
}

func runGossopingServer(ml *memberlist.Memberlist, port int, serverName string) {
	logger.Infof("%s Starting listening on port %d.", serverName, port)

	var address = fmt.Sprintf("%s:%d", ml.LocalNode().Addr.String(), port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatalf("Cannot start the cluster node: %v", err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read data from the connection
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	// Process the received data
	data := buffer[:n]
	fmt.Printf("Received data: %s\n", data)
}

func setupHTTPServer(settings *settings.Settings, helper *helper.Helper, q *models.Queue, s *models.Subscriber) *http.Server {
	logger.InfoS("Initializing http server.")

	apiServer, err := api.NewAPIServer(settings, helper, q, s)
	if err != nil {
		logger.FatalS("Could not initialize API Server", "error", err.Error())
	}

	// InterCommunication API Server Setup
	apiAddress := fmt.Sprintf(":%d", settings.Global.APIPort)
	APIServer := &http.Server{
		Addr:              apiAddress,
		Handler:           apiServer,
		ReadTimeout:       settings.Global.ReadTimeout,
		ReadHeaderTimeout: settings.Global.ReadHeaderTimeout,
		WriteTimeout:      settings.Global.WriteTimeout,
		IdleTimeout:       settings.Global.IdleTimeout,
		MaxHeaderBytes:    settings.Global.MaxHeaderBytes,
	}

	return APIServer
}

func runHTTPServer(server *http.Server, port int, serverName string) {
	logger.Infof("%s Starting listening on port %d.", serverName, port)
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		logger.Fatal("could not create "+serverName+" server listener", "error", err.Error())
	}
	err = server.Serve(ln)
	if err != nil {
		logger.InfoS("Serving failed", "error", err.Error(), "serverName", serverName)
	}
}
