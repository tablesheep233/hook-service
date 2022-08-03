package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/tablesheep233/hook-service/config"
	"github.com/tablesheep233/hook-service/logger"
)

var lock sync.Mutex

func responseError(msg string, writer http.ResponseWriter) {
	response := ResponseMsg{Msg: msg}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(500)
	if err := json.NewEncoder(writer).Encode(&response); err != nil {
		logger.Error.Println("write response err", err)
	}
}

func hook(writer http.ResponseWriter, request *http.Request) {

	if request.Method != "POST" {
		logger.Error.Println("Request Method Error, only support post")
		responseError("Request Method Error, only support post", writer)
		return
	}
	logger.Info.Printf("receive request %s: %s", request.Method, request.URL)

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		logger.Error.Println("read body error", err)
		responseError("read body error", writer)
		return
	}
	logger.Info.Println("request: ", string(body))
	var releaseArg ReleaseArg
	if err := json.Unmarshal(body, &releaseArg); err != nil {
		logger.Error.Println("parse body error", err)
		responseError("parse body error", writer)
		return
	}

	if releaseArg.Action == "published" {
		shellPath := config.Config.Scripts[releaseArg.Repo.Name]
		if shellPath == "" {
			logger.Error.Println("cannot find script", err)
			responseError("cannot find script", writer)
			return
		}
		shellPath = fmt.Sprintf("%s release_tag=\"%s\"", shellPath, releaseArg.Release.Tag_name)
		logger.Info.Println(shellPath)
		lock.Lock()
		defer lock.Unlock()
		{
			start := time.Now()
			command := exec.Command("bash", "-c", shellPath)

			rc, err := command.StdoutPipe()
			if err != nil {
				logger.Error.Println("start script failed", err)
				responseError("start script failed", writer)
				return
			}

			err = command.Start()
			if err != nil {
				logger.Error.Println("start script failed", err)
				responseError("start script failed", writer)
				return
			}
			logger.Info.Println("Start execute the shell, Process Pid:", command.Process.Pid)

			reader := bufio.NewReader(rc)
			for {
				s, err := reader.ReadString('\n')
				if err != nil || io.EOF == err {
					break
				}
				logger.Info.Print(s)
			}

			err = command.Wait()
			if err != nil {
				logger.Error.Println("execute script failed", err)
				responseError("execute script failed", writer)
				return
			}
			end := time.Now()
			logger.Info.Printf("script execute success, PID:%v, cost:%v", command.ProcessState.Pid(), end.Sub(start))
		}
	} else {
		logger.Info.Println("ignore action ", releaseArg.Action)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	if err := json.NewEncoder(writer).Encode(&ResponseMsg{Msg: "ok"}); err != nil {
		logger.Error.Println("write response err", err)
	}
}

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hooks", hook)

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", config.Config.Port),
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error.Println("server start failed:", err)
	}
}
