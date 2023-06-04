package main

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"log"
	"net"
	"net/http"
	"strings"
)

var IP string

func main() {
	//get ip
	ip, err := GetOutBoundIP()
	if err != nil {
		fmt.Println(ip)
	}
	IP = ip
	// HTTP server
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	nodeWS := new(restful.WebService)
	nodeWS.Path("/").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	nodeWS.Route(nodeWS.GET("").To(Greet))
	wsContainer.Add(nodeWS)

	// run
	server := &http.Server{Addr: ":8888", Handler: wsContainer}
	defer server.Close()
	log.Fatal(server.ListenAndServe())
}

func Greet(request *restful.Request, response *restful.Response) {
	ret := "hello from " + IP + "!"
	response.Write([]byte(ret))
}

func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
