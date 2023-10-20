package cmd

import (
	"fmt"
	"net/http"

	command "github.com/daijinru/mango-packages-command"
	httpService "github.com/daijinru/mango-runner/http"
	"github.com/daijinru/mango-runner/runner"
	"github.com/ttacon/chalk"
)

func NewServiceRPC() *command.Command {
	cmd := &command.Command{
		Use:  "serve",
		Args: command.ExactArgs(1),
		RunE: func(cmd *command.Command, args []string) error {
			return nil
		},
	}
	cmd.AddCommand(NewServiceRpcStart())
	return cmd
}

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

func NewServiceRpcStart() *command.Command {
	return &command.Command{
		Use:  "start",
		Args: command.ExactArgs(1),
		RunE: func(cmd *command.Command, args []string) error {
			ciService := &httpService.CiService{}
			handler1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost:
					ciService.CreatePipeline(w, r)
				default:
					http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
				}
			})
			handler2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost:
					ciService.ReadPipelineStatus(w, r)
				default:
					http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
				}
			})
			handler3 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost:
					ciService.ReadPipeline(w, r)
				default:
					http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
				}
			})
			handler4 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost:
					ciService.ReadPipelines(w, r)
				default:
					http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
				}
			})

			http.Handle("/pipeline/create", loggingMiddleware(handler1))
			http.Handle("/pipeline/status", loggingMiddleware(handler2))
			http.Handle("/pipeline/stdout", loggingMiddleware(handler3))
			http.Handle("/pipeline/list", loggingMiddleware(handler4))

			fmt.Println("🌏 Now listening at port: " + args[0])
			fmt.Println(http.ListenAndServe(":"+args[0], nil))
			return nil
		},
	}
}
