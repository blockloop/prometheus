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

	"github.com/go-kit/kit/log/level"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/digitalocean/prometheus/v2/prompb"
	"github.com/prometheus/tsdb"
	"github.com/prometheus/tsdb/labels"
)

func (h *Handler) read(w http.ResponseWriter, r *http.Request) {
	query, err := parseRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matchers := toMatchers(query)
	if len(matchers) == 0 {
		http.Error(w, "missing query matcher", http.StatusBadRequest)
		return
	}

	q, err := h.tsdb().Querier(query.GetStartTimestampMs(), query.GetEndTimestampMs())
	if err != nil {
		level.Error(h.logger).Log("msg", "failed to create querier", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer q.Close()

	set, err := q.Select(matchers...)
	if err != nil {
		level.Error(h.logger).Log("msg", "failed to execute query", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ts := toTimeseries(set)

	if err := writeResponse(w, ts); err != nil {
		level.Error(h.logger).Log("msg", "failed to write query", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseRequest(r *http.Request) (*prompb.Query, error) {
	compressed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %+v", err.Error())
	}

	buf, err := snappy.Decode(nil, compressed)
	if err != nil {
		return nil, fmt.Errorf("snappy decode failed: %+v", err)
	}

	var req prompb.ReadRequest
	if err := proto.Unmarshal(buf, &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proto: %+v", err)
	}

	if len(req.GetQueries()) != 1 {
		return nil, fmt.Errorf("exactly one query must be sent. Got %d", len(req.GetQueries()))
	}

	return req.GetQueries()[0], nil
}

func toPbLabels(in labels.Labels) []prompb.Label {
	res := make([]prompb.Label, len(in))
	for i, l := range in {
		res[i] = prompb.Label{
			Name:  l.Name,
			Value: l.Value,
		}
	}
	return res
}

func toMatchers(query *prompb.Query) []labels.Matcher {
	ms := query.GetMatchers()
	if ms == nil {
		return nil
	}

	mt := make([]labels.Matcher, len(ms))
	for i, m := range ms {
		switch m.GetType() {
		case prompb.LabelMatcher_EQ:
			mt[i] = labels.NewEqualMatcher(m.GetName(), m.GetValue())
		case prompb.LabelMatcher_NEQ:
			mt[i] = labels.Not(labels.NewEqualMatcher(m.GetName(), m.GetValue()))
		case prompb.LabelMatcher_RE:
			mt[i] = labels.NewMustRegexpMatcher(m.GetName(), m.GetValue())
		case prompb.LabelMatcher_NRE:
			mt[i] = labels.Not(labels.NewMustRegexpMatcher(m.GetName(), m.GetValue()))
		default:
			continue
		}
	}
	return mt

}

func toTimeseries(set tsdb.SeriesSet) []*prompb.TimeSeries {
	ts := make([]*prompb.TimeSeries, 0)

	for set.Next() {
		series := set.At()
		res := &prompb.TimeSeries{
			Labels:  toPbLabels(series.Labels()),
			Samples: []prompb.Sample{},
		}
		it := series.Iterator()
		for it.Next() {
			t, v := it.At()
			res.Samples = append(res.Samples, prompb.Sample{
				Timestamp: t,
				Value:     v,
			})
		}

		ts = append(ts, res)
	}
	return ts
}

func writeResponse(w http.ResponseWriter, ts []*prompb.TimeSeries) error {
	marshaled, err := proto.Marshal(&prompb.ReadResponse{
		Results: []*prompb.QueryResult{{Timeseries: ts}},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal repsonse: %+v", err)
	}

	enc := snappy.Encode(nil, marshaled)

	if _, err := w.Write(enc); err != nil {
		return fmt.Errorf("failed to write response body: %+v", err)
	}
	return nil
}
