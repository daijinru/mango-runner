package httpService

import (
	"fmt"
	"github.com/daijinru/mango-runner/runner"
	"github.com/ttacon/chalk"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		url := r.URL.String()
		err := r.ParseForm()
		if err != nil {
			fmt.Printf("failed to parse request body: %v\n", err)
		}
		params := r.Form

		msg := fmt.Sprintf("Received from IP: %s, URL: %s, Params: %v", ip, url, params)
		fmt.Println(chalk.Dim, runner.AddPrefixMsg(msg))

		next.ServeHTTP(w, r)
	})
}

// RouteConfig ServiceFunc receives ResponseWriter and Request and no need to return
type RouteConfig struct {
	Path        string
	Method      string
	ServiceFunc func(w http.ResponseWriter, r *http.Request)
}

func createRouteHandler(config RouteConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == config.Method {
			config.ServiceFunc(w, r)
		} else {
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

type ServiceHttpConfig struct {
	Port string
}

func NewServiceHttpStart(config *ServiceHttpConfig) {
	ciService := &CiService{}

	routes := []RouteConfig{
		{
			Path:   "/pipeline/create",
			Method: http.MethodPost,
			ServiceFunc: func(w http.ResponseWriter, r *http.Request) {
				ciService.CreatePipeline(w, r)
			},
		},
		{
			Path:   "/pipeline/stdout",
			Method: http.MethodPost,
			ServiceFunc: func(w http.ResponseWriter, r *http.Request) {
				ciService.ReadPipeline(w, r)
			},
		},
		{
			Path:   "/pipeline/list",
			Method: http.MethodPost,
			ServiceFunc: func(w http.ResponseWriter, r *http.Request) {
				ciService.ReadPipelines(w, r)
			},
		},
		{
			Path:   "/service/status",
			Method: http.MethodPost,
			ServiceFunc: func(w http.ResponseWriter, r *http.Request) {
				ciService.ReadServiceStatus(w, r)
			},
		},
		{
			Path:   "/git/clone",
			Method: http.MethodPost,
			ServiceFunc: func(w http.ResponseWriter, r *http.Request) {
				ciService.GitClone(w, r)
			},
		},
		{
			Path:   "/service/monitor",
			Method: http.MethodPost,
			ServiceFunc: func(w http.ResponseWriter, r *http.Request) {
				ciService.ReadMonitor(w, r)
			},
		},
	}

	for _, route := range routes {
		handler := createRouteHandler(route)
		http.Handle(route.Path, loggingMiddleware(handler))
	}

	fmt.Println("🌏 Now listening at port: " + config.Port)
	fmt.Println(http.ListenAndServe(":"+config.Port, nil))
}
