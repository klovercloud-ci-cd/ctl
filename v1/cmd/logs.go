package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/klovercloud-ci/ctl/dependency_manager"
	v1 "github.com/klovercloud-ci/ctl/v1"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"time"
)

func GetLogs() *cobra.Command {
	command := cobra.Command{
		Use:       "logs",
		Short:     "Get logs by process ID",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := v1.GetConfigFile()
			if cfg.Token == "" {
				cmd.Println("[ERROR]: %v", "user is not logged in")
				return nil
			}
			var processId, page, limit string
			var follow, skipSsl bool
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
				} else if strings.Contains(strings.ToLower(each), "--skipssl") {
					skipSsl = true
				}
			}
			if apiServerUrl == "" {
				if cfg.ApiServerUrl == "" {
					cmd.Println("[ERROR]: Api server url not found!")
					return nil
				}
			} else {
				cfg.ApiServerUrl = apiServerUrl
				err := cfg.Store()
				if err != nil {
					cmd.Println("[ERROR]: ", err.Error())
					return nil
				}
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
			return getLogs(cmd, cfg.ApiServerUrl, cfg.Token, processId, page, limit, follow, skipSsl, 0)
		},
		DisableFlagParsing: true,
	}
	command.SetUsageTemplate("Usage: \n" +
		"  cli logs [processid=PROCESS_ID] [page=PAGE_NUMBER] [limit=LIMIT_NUMBER] [follow | -f] [apiserver=APISERVER_URL] [--skipssl] \n" +
		"  cli help logs \n" +
		"\nOptions: \n" +
		"  --skipssl\t" + "Ignore SSL certificate errors \n" +
		"  help\t" + "Show this screen. \n")
	return &command
}

func getLogs(cmd *cobra.Command, apiServerUrl, token, processId, page, limit string, follow, skipSsl bool, skip int64) error {
	pipelineService := dependency_manager.GetPipelineService()
	code, data, err := pipelineService.Token(token).SkipSsl(skipSsl).Logs(apiServerUrl+"pipelines/"+processId, page, limit)
	if err != nil {
		cmd.Println("[ERROR]: ", err.Error())
		return nil
	} else if code != 200 {
		cmd.Println("[ERROR]: ", "Something went wrong! StatusCode: ", code)
		return nil
	} else if data != nil {
		byteBody, _ := json.Marshal(data)
		var result []string
		err := json.Unmarshal(byteBody, &result)
		if err != nil {
			return err
		}
		for i := skip; i < int64(len(result)); i++ {
			if result[i] != "" {
				if strings.HasSuffix(strings.ToLower(result[i]), "step started") {
					cmd.Println("\n")
					color.Set(color.FgHiCyan, color.Bold)
					fmt.Println(center(result[i], 110))
					color.Unset()
					cmd.Println("\n")
				} else {
					cmd.Println(result[i])
				}
			}
		}
		lim, _ := strconv.Atoi(limit)
		if len(result) < lim {
			skip = int64(len(result))
		} else {
			skip = 0
		}
	} else if data == nil {
		p, _ := strconv.Atoi(page)
		page = strconv.Itoa(p - 1)
	}
	if follow {
		if skip == 0 {
			p, _ := strconv.Atoi(page)
			page = strconv.Itoa(p + 1)
		}
		time.Sleep(time.Millisecond * 500)
		err := getLogs(cmd, apiServerUrl, token, processId, page, limit, follow, skipSsl, skip)
		if err != nil {
			return err
		}
	}
	return nil
}

func center(s string, w int) string {
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
}
