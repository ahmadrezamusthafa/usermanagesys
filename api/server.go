package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"path"
	"runtime"
	"tokopedia.se.training/Project1/usermanagesys/api/configuration"
	"tokopedia.se.training/Project1/usermanagesys/api/gnsq"
	"tokopedia.se.training/Project1/usermanagesys/api/service"
	"tokopedia.se.training/Project1/usermanagesys/helper/processor"
)

var SERVER Server

type Server struct {
	Configuration *configuration.Configuration
	NsqModule     *gnsq.GNSQModule
	UserService   *service.UserService
	MainService   *service.MainService
}

func NewServer(config *configuration.Configuration,
	nsqModule *gnsq.GNSQModule,
	userService *service.UserService,
	mainService *service.MainService,
) *Server {
	SERVER = Server{
		Configuration: config,
		NsqModule:     nsqModule,
		UserService:   userService,
		MainService:   mainService,
	}

	return &SERVER
}

func (a *Server) NewRouter() *mux.Router {

	_, runningFile, _, _ := runtime.Caller(1)
	frontendPath := path.Join(path.Dir(runningFile), "../../frontend")
	fmt.Println(frontendPath)

	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir(frontendPath))))

	router.HandleFunc("/user/get_users", a.UserService.GetUsers).Methods("GET")
	router.HandleFunc("/user/get_users_paging", a.UserService.GetUsersPaging).Methods("GET")
	router.HandleFunc("/publish", a.MainService.PublishNSQ).Methods("GET")
	router.HandleFunc("/get_visitor_count", a.MainService.GetVisitorCount).Methods("GET")

	return router
}

func (s *Server) Run() {
	log.Println("Setup HTTP Server")
	var port = processor.ExtractServerAddressPort(s.Configuration.App.BackEndAddress)
	log.Println("Starting server at http://localhost:" + port + "/")
	log.Println("Starting web at http://localhost:" + port + "/web/")

	log.Println("Initialize Router")
	router := s.NewRouter()
	corsMiddleware := cors.New(cors.Options{
		AllowedMethods: []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"}, //http request yang diizinkan
		AllowedHeaders: []string{"Content-Type", "X-CSRF-Token", "Authorization", ""},
		Debug:          false,
	})

	log.Println("Listen using middleware")
	err := http.ListenAndServe(":"+port, corsMiddleware.Handler(router))
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}
