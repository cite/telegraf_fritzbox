package fritzbox

// fritzbox.go

// Copyright 2019 Stefan Förster, original code in main.go taken from
// https://github.com/ndecker/fritzbox_exporter, original copyright
// notice:
// Copyright 2016 Nils Decker
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	upnp "github.com/ndecker/fritzbox_exporter/fritzbox_upnp"
	"log"
)

type Fritzbox struct {
	Host string
	Port uint16
}

type Metric struct {
	Service string
	Action  string
	Result  string
	Name    string
}

var metrics = []*Metric{
	{
		Service: "urn:schemas-upnp-org:service:WANCommonInterfaceConfig:1",
		Action:  "GetTotalPacketsReceived",
		Result:  "TotalPacketsReceived",
		Name:    "packets_received",
	},
	{
		Service: "urn:schemas-upnp-org:service:WANCommonInterfaceConfig:1",
		Action:  "GetTotalPacketsSent",
		Result:  "TotalPacketsSent",
		Name:    "packets_sent",
	},
	{
		Service: "urn:schemas-upnp-org:service:WANCommonInterfaceConfig:1",
		Action:  "GetAddonInfos",
		Result:  "TotalBytesReceived",
		Name:    "bytes_received",
	},
	{
		Service: "urn:schemas-upnp-org:service:WANCommonInterfaceConfig:1",
		Action:  "GetAddonInfos",
		Result:  "TotalBytesSent",
		Name:    "bytes_sent",
	},
	{
		Service: "urn:schemas-upnp-org:service:WANCommonInterfaceConfig:1",
		Action:  "GetCommonLinkProperties",
		Result:  "PhysicalLinkStatus",
		Name:    "link_status",
	},
	{
		Service: "urn:schemas-upnp-org:service:WANIPConnection:1",
		Action:  "GetStatusInfo",
		Result:  "ConnectionStatus",
		Name:    "connection_status",
	},
	{
		Service: "urn:schemas-upnp-org:service:WANIPConnection:1",
		Action:  "GetStatusInfo",
		Result:  "Uptime",
		Name:    "uptime",
	},
}

func (s *Fritzbox) Description() string {
	return "a demo plugin"
}

func (s *Fritzbox) SampleConfig() string {
	return `
  ## Host and Port for FRITZ!Box UPnP service
  host = fritz.box
  port = 49000
`
}

func (s *Fritzbox) Gather(acc telegraf.Accumulator) error {
	var host string
	var port uint16
	if s.Host == "" {
		host = "fritz.box"
	} else {
		host = s.Host
	}
	if s.Port == 0 {
		port = 49000
	} else {
		port = s.Port
	}

	root, err := upnp.LoadServices(host, port)
	if err != nil {
		return fmt.Errorf("fritzbox: unable to load services: %v", err)
	}

	// remember what we already called
	var last_service string
	var last_method string
	var result upnp.Result
	fields := make(map[string]interface{})

	for _, m := range metrics {
		if m.Service != last_service || m.Action != last_method {
			service, ok := root.Services[m.Service]
			if !ok {
				// TODO
				log.Println("W! Cannot find defined service %s", m.Service)
				continue
			}
			action, ok := service.Actions[m.Action]
			if !ok {
				// TODO
				log.Println("W! Cannot find defined action %s on service %s", m.Action)
				continue
			}

			result, err = action.Call()
			if err != nil {
				log.Println("E! Unable to call action %s on service %s: %v", m.Action, m.Service, err)
				continue
			}

			// save service and action
			last_service = m.Service
			last_method = m.Action
		}

		fields[m.Name] = result[m.Result]
	}
	acc.AddFields("fritzbox", fields, map[string]string{"host": host})

	return nil
}

func init() {
	inputs.Add("fritzbox", func() telegraf.Input { return &Fritzbox{} })
}
