// Copyright 2013 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/digitalocean/prometheus/v2/prompb"
	"github.com/digitalocean/prometheus/v2/web/api/v2"
)

const commitChunkSize = 500

// RemoteWrite is an HTTP handler to handle Prometheus remote_write
func (h *Handler) write(w http.ResponseWriter, r *http.Request) {
	timeseries, err := readRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(timeseries) == 0 {
		// nothing was sent so just nop
		return
	}

	api_v2.WriteTimeSeries(timeseries, h.tsdb, h.logger)
}

func readRequest(r *http.Request) ([]prompb.TimeSeries, error) {
	compressed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %+v", err)
	}
	defer r.Body.Close()

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		return nil, fmt.Errorf("failed to snappy.Decode: %+v", err)
	}

	var req prompb.WriteRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		return nil, fmt.Errorf("failed to proto.Unmarshal: %+v", err)
	}
	return req.GetTimeseries(), nil
}
