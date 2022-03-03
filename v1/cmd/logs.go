package cmd

import (
	"encoding/json"
	"github.com/klovercloud-ci/ctl/dependency_manager"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"time"
)

func GetLogs() *cobra.Command {
	return &cobra.Command{
		Use:       "logs",
		Short:     "Get logs by process ID",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := v1.IsUserLoggedIn(); err != nil {
				cmd.Printf("[ERROR]: %v", err.Error())
				return nil
			}
			var processId, page, limit string
			var follow bool
			var apiServerUrl string
			for _, each := range args {
				if strings.Contains(strings.ToLower(each), "page=") {
					strs := strings.Split(strings.ToLower(each), "=")
					if len(strs) > 0 {
						page = strs[1]
					}
				} else if strings.Contains(strings.ToLower(each), "limit=") {
					strs := strings.Split(strings.ToLower(each), "=")
					if len(strs) > 0 {
						limit = strs[1]
					}
				} else if strings.Contains(strings.ToLower(each), "processid=") {
					strs := strings.Split(strings.ToLower(each), "=")
					if len(strs) > 0 {
						processId = strs[1]
					}
				} else if strings.Contains(strings.ToLower(each), "follow") {
					follow = true
				} else if strings.Contains(strings.ToLower(each), "-f") {
					follow = true
				} else if strings.Contains(strings.ToLower(each), "apiserver=") {
					strs := strings.Split(strings.ToLower(each), "=")
					if len(strs) > 1 {
						apiServerUrl = strs[1]
					}
				}
			}
			cfg := v1.GetConfigFile()
			if apiServerUrl == "" {
				if cfg.ApiServerUrl == "" {
					cmd.Println("[ERROR]: Api server url not found!")
					return nil
				}
			} else {
				cfg.ApiServerUrl = apiServerUrl
			}
			err := cfg.Store()
			if err != nil {
				cmd.Println("[ERROR]: ", err.Error())
				return nil
			}
			if page == "" {
				page = "0"
			}
			if limit == "" {
				limit = "10"
			}
			if processId == "" {
				return nil
			}
			return getLogs(cmd, processId, page, limit, follow)
		},
	}
}

func getLogs(cmd *cobra.Command, processId string, page string, limit string, follow bool) error {
	cfg := v1.GetConfigFile()
	pipelineService := dependency_manager.GetPipelineService()
	code, data, err := pipelineService.Logs(cfg.ApiServerUrl+"pipelines/"+processId, page, limit)
	if err != nil {
		cmd.Println("[ERROR]: ", err.Error())
		return nil
	} else if code != 200 {
		cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
		return nil
	} else if data != nil {
		byteBody, _ := json.Marshal(data)
		var result []string
		json.Unmarshal(byteBody, &result)
		for _, logData := range result {
			cmd.Println(logData)
		}
	} else {
		return nil
	}
	if follow {
		if data == nil {
			limit = strconv.Itoa(5)
		} else {
			i, _ := strconv.Atoi(page)
			page = strconv.Itoa(i + 1)
		}
		time.Sleep(time.Millisecond * 500)
		getLogs(cmd, processId, page, limit, follow)
	}
	return nil
}
