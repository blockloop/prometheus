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
	"sort"

	"github.com/go-kit/kit/log/level"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/tsdb/labels"
)

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

	ap := h.tsdb().Appender()

	// rollback the appends if it is not committed successfully
	committed := false
	defer func() {
		if !committed {
			ap.Rollback()
		}
	}()
	for _, ts := range timeseries {
		lbls := make(labels.Labels, len(ts.Labels))
		for i, l := range ts.Labels {
			lbls[i] = labels.Label{
				Name:  l.GetName(),
				Value: l.GetValue(),
			}
		}
		// soring guarantees hash consistency
		sort.Sort(lbls)

		var ref uint64
		var err error
		for _, s := range ts.Samples {
			if ref == 0 {
				ref, err = ap.Add(lbls, s.GetTimestamp(), s.GetValue())
			} else {
				err = ap.AddFast(ref, s.GetTimestamp(), s.GetValue())
			}
			if err != nil {
				level.Error(h.logger).Log("msg", "failure while writing to store", "err", err)
			}
		}
	}

	if err := ap.Commit(); err != nil {
		level.Error(h.logger).Log("msg", "failure trying to commit write to store", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	committed = true
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
