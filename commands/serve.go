package cmd

import (
	"fmt"
	"net/http"

	command "github.com/daijinru/mango-packages-command"
	httpService "github.com/daijinru/mango-runner/http"
)

func NewServiceRPC() *command.Command {
  cmd := &command.Command{
    Use: "serve",
    Args: command.ExactArgs(1),
    RunE: func(cmd *command.Command, args[]string) error {
      return nil
    },
  }
  cmd.AddCommand(NewServiceRpcStart())
  return cmd
}

func NewServiceRpcStart() *command.Command {
  return &command.Command{
    Use: "start",
    Args: command.ExactArgs(1),
    RunE: func(cmd *command.Command, args[]string) error {
      ciService := &httpService.CiService{}
      http.HandleFunc("/pipeline/create", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
          ciService.CreatePipeline(w, r)
        default:
          http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
        }
      })
      http.HandleFunc("/pipeline/status", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
          ciService.ReadPipelineStatus(w, r)
        default:
          http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
        }
      })
      http.HandleFunc("/pipeline/stdout", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
          ciService.ReadPipeline(w, r)
        default:
          http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
        }
      })
      http.HandleFunc("/pipeline/list", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
          ciService.ReadPipelines(w, r)
        default:
          http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
        }
      })
      fmt.Println("🌏 Now listening at port: " + args[0])
      fmt.Println(http.ListenAndServe(":" + args[0], nil))
      return nil
    },
  }
}
