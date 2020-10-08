package config

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/pborman/getopt/v2"
	"io"
	"net/http"
	"os"
	"strings"
)

type Creds struct {
	User     string `json:"u"`
	Password string `json:"p"`
}

type Cfg struct {
	LogLocationErr   string
	LogLocationInfo  string
	LogLocationWarn  string
	LogLocationDebug string
	TOInsecure       bool
	TOUser           string
	TOPass           string
	TOUrl            string
	Tr               *http.Transport
}

type ToResponse struct {
	Response string `json:"response"`
}

func Dclose(c io.Closer) {
	if err := c.Close(); err != nil {

	}
}

func ErrCheck(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (cfg Cfg) ErrorLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationErr) }
func (cfg Cfg) WarningLog() log.LogLocation { return log.LogLocation(cfg.LogLocationWarn) }
func (cfg Cfg) InfoLog() log.LogLocation    { return log.LogLocation(cfg.LogLocationInfo) }
func (cfg Cfg) DebugLog() log.LogLocation   { return log.LogLocation(cfg.LogLocationDebug) }
func (cfg Cfg) EventLog() log.LogLocation   { return log.LogLocation(log.LogLocationNull) } // event logging not used.

func GetCfg() (Cfg, error) {
	var err error
	logLocationDebugPtr := getopt.StringLong("log-location-debug", 'd', "", "Where to log debugs. May be a file path, stdout, stderr, or null, default ''")
	logLocationErrorPtr := getopt.StringLong("log-location-error", 'e', "stderr", "Where to log errors. May be a file path, stdout, stderr, or null, default stderr")
	logLocationInfoPtr := getopt.StringLong("log-location-info", 'i', "stderr", "Where to log info. May be a file path, stdout, stderr, or null, default stderr")
	logLocationWarnPtr := getopt.StringLong("log-location-warning", 'w', "stderr", "Where to log warnings. May be a file path, stdout, stderr, or null, default stderr")
	toInsecurePtr := getopt.BoolLong("traffic-ops-insecure", 'I', "[true | false] ignore certificate errors from Traffic Ops")
	toUserPtr := getopt.StringLong("traffic-ops-user", 'u', "", "Traffic Ops username. Required.")
	toPassPtr := getopt.StringLong("traffic-ops-passowrd", 'p', "", "Traffic Ops Password. Required")
	toUrlPtr := getopt.StringLong("traffic-ops-url", 'U', "", "Traffic ops base URL. Required.")
	helpPtr := getopt.BoolLong("help", 'h', "Print usage information and exit")
	getopt.ParseV2()

	logLocationDebug := *logLocationDebugPtr
	logLocationError := *logLocationErrorPtr
	logLocationInfo := *logLocationInfoPtr
	logLocationWarn := *logLocationWarnPtr
	toInsecure := *toInsecurePtr
	toURL := *toUrlPtr
	toUser := *toUserPtr
	toPass := *toPassPtr
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: toInsecure}}
	help := *helpPtr

	if help {
		Usage()
		return Cfg{}, nil
	}

	missingStr := "Missing required argument"
	usageStr := "\nBasic usage: ToDnssecRefresh --traffic-ops-url=myurl --traffic-ops-user=myuser --traffic-ops-password=mypass"
	if strings.TrimSpace(toURL) == "" {
		return Cfg{}, errors.New(missingStr + " --traffic-ops-url" + usageStr)
	}
	if strings.TrimSpace(toUser) == "" {
		return Cfg{}, errors.New(missingStr + " --traffic-ops-user" + usageStr)
	}
	if strings.TrimSpace(toPass) == "" {
		return Cfg{}, errors.New(missingStr + " --traffic-ops-password" + usageStr)
	}

	cfg := Cfg{
		LogLocationDebug: logLocationDebug,
		LogLocationErr:   logLocationError,
		LogLocationInfo:  logLocationInfo,
		LogLocationWarn:  logLocationWarn,
		TOInsecure:       toInsecure,
		Tr:               tr,
		TOUrl:            toURL,
		TOUser:           toUser,
		TOPass:           toPass,
	}

	if err = log.InitCfg(cfg); err != nil {
		return Cfg{}, errors.New("Initializing loggers: " + err.Error() + "\n")
	}

	return cfg, nil
}

func PrintConfig(cfg Cfg) {
	log.Debugf("TOUrl: %s\n", cfg.TOUrl)
	log.Debugf("TOUser: %s\n", cfg.TOUser)
	log.Debugf("TOPass: Pass len: %d\n", len(cfg.TOPass))
	log.Debugf("TOInsecure: %t\n", cfg.TOInsecure)
	log.Debugf("LogLocationDebug: %s\n", cfg.LogLocationDebug)
	log.Debugf("LogLocationErr: %s\n", cfg.LogLocationErr)
	log.Debugf("LogLocationInfo: %s\n", cfg.LogLocationInfo)
	log.Debugf("LogLocationWarn: %s\n", cfg.LogLocationWarn)
}

func Usage() {
	fmt.Println("\t  --log-location-debug=[value] | -d [value], Where to log debugs. May be a file path, stdout, stderr, or null, default stderr")
	fmt.Println("\t  --log-location-error=[value] | -e [value], Where to log errors. May be a file path, stdout, stderr, or null, default stderr")
	fmt.Println("\t  --log-location-info=[value] | -i [value], Where to log info. May be a file path, stdout, stderr, or null, default stderr")
	fmt.Println("\t  --log-location-warning=[value] | -w [value], Where to log warnings. May be a file path, stdout, stderr, or null, default stderr")
	fmt.Println("\t  --traffic-ops-url=[url] | -u [url], Traffic Ops URL. Must be the full URL, including the scheme. Required.")
	fmt.Println("\t  --traffic-ops-user=[username] | -U [username], Traffic Ops username. Required.")
	fmt.Println("\t  --traffic-ops-password=[password] | -P [password], Traffic Ops password. Required.")
	fmt.Println("\t  --help | -h, Print usage information and exit")
}
