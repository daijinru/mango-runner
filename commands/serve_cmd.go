package cmd

import (
	command "github.com/daijinru/mango-packages-command"
	httpService "github.com/daijinru/mango-runner/http"
)

func NewServiceRPC() *command.Command {
	cmd := &command.Command{
		Use:  "serve",
		Args: command.ExactArgs(1),
		RunE: func(cmd *command.Command, args []string) error {
			return nil
		},
	}
	cmd.AddCommand(NewServiceHttpStart())
	return cmd
}

func NewServiceHttpStart() *command.Command {
	return &command.Command{
		Use:  "start",
		Args: command.ExactArgs(1),
		RunE: func(cmd *command.Command, args []string) error {
			config := &httpService.ServiceHttpConfig{
				Port: args[0],
			}
			httpService.NewServiceHttpStart(config)
			return nil
		},
	}
}
