package main

import (
	"github.com/klovercloud-ci/ctl/config"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"time"
)

func main(){
	config.InitEnvironmentVariables()
	cli()
}
func cli(){
	cmd := &cobra.Command{
		Use:          "ctl",
		Short:        "Cli to use klovercloud-ci apis!",
		Version:      "v1",
		SilenceUsage: true,
	}
	cmd.AddCommand(GetLogs())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func GetLogs()*cobra.Command {
	return &cobra.Command{
		Use:       "logs",
		Short:     "get logs by processId",
		ValidArgs: []string{},
		RunE: func(cmd *cobra.Command, args []string) error {

			var processId,page,limit string
			var follow bool
			for _,each:=range args{
				if strings.Contains(strings.ToLower(each),"page"){
					strs:=strings.Split(strings.ToLower(each),"=")
					if len(strs)>0{
						page=strs[1]
					}
				}
				if strings.Contains(strings.ToLower(each),"limit"){
					strs:=strings.Split(strings.ToLower(each),"=")
					if len(strs)>0{
						limit=strs[1]
					}
				}
				if strings.Contains(strings.ToLower(each),"processid"){
					strs:=strings.Split(strings.ToLower(each),"=")
					if len(strs)>0{
						processId=strs[1]
					}
				}
				if strings.Contains(strings.ToLower(each),"follow"){
					strs:=strings.Split(strings.ToLower(each),"=")
					if len(strs)>0{
						if strings.ToLower(strs[1])=="true"{
							follow=true
						}else {
							follow=false
						}
					}
				}
				if strings.Contains(strings.ToLower(each),"-f"){
					follow=true
				}
			}

			if page==""{page="0"}
			if limit==""{limit="10"}
			if processId==""{
				return nil
			}
				return getLogs(cmd, processId, page, limit,follow)
		},
	}
}

func getLogs(cmd *cobra.Command, processId string, page string, limit string, follow bool) error {
	pipelineService := GetPipelineService()
	code, data, err := pipelineService.Logs(config.ApiServerUrl+"pipelines/"+processId, page, limit)
	if err != nil {
		cmd.Println("[ERROR]: ", err.Error())
		return nil
	}
	if code != 200 {
		cmd.Println("[ERROR]: ", "Something wen wrong! StatusCode: ", code)
		return nil
	}
	if data!=nil {
		cmd.Println(data)
	}
	if follow{
		if data==nil{
			limit= strconv.Itoa( 5)
		}else{
			i,_:=strconv.Atoi(page)
			page= strconv.Itoa(i + 1)
		}
		time.Sleep(time.Millisecond*500)
		getLogs(cmd,processId,page,limit,follow)
	}
	return nil
}
