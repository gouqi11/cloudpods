// Copyright 2019 Yunion
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

package hcso

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go/auth/aksk"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/cloudprovider"
	"yunion.io/x/onecloud/pkg/util/httputils"
	"yunion.io/x/pkg/errors"
)

type akClient struct {
	client *http.Client
	aksk   aksk.SignOptions
}

func (self *akClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Del("Accept")
	if req.Method == string(httputils.GET) || req.Method == string(httputils.DELETE) || req.Method == string(httputils.PATCH) {
		req.Header.Del("Content-Length")
	}
	aksk.Sign(req, self.aksk)
	return self.client.Do(req)
}

func (self *SHuaweiClient) getAkClient() *akClient {
	return &akClient{
		client: self.getDefaultClient(),
		aksk: aksk.SignOptions{
			AccessKey: self.accessKey,
			SecretKey: self.accessSecret,
		},
	}
}

func (self *SHuaweiClient) getDefaultClient() *http.Client {
	if self.httpClient != nil {
		return self.httpClient
	}
	self.httpClient = self.cpcfg.AdaptiveTimeoutHttpClient()
	ts, _ := self.httpClient.Transport.(*http.Transport)
	self.httpClient.Transport = cloudprovider.GetCheckTransport(ts, func(req *http.Request) (func(resp *http.Response), error) {
		service, method, path := strings.Split(req.URL.Host, ".")[0], req.Method, req.URL.Path
		respCheck := func(resp *http.Response) {
			if resp.StatusCode == 403 {
				if self.cpcfg.UpdatePermission != nil {
					self.cpcfg.UpdatePermission(service, fmt.Sprintf("%s %s", method, path))
				}
			}
		}
		if self.cpcfg.ReadOnly {
			if req.Method == "GET" {
				return respCheck, nil
			}
			return nil, errors.Wrapf(cloudprovider.ErrAccountReadOnly, "%s %s", req.Method, req.URL.Path)
		}
		return respCheck, nil
	})
	return self.httpClient
}

func (self *SHuaweiClient) request(method httputils.THttpMethod, url string, query url.Values, params map[string]interface{}) (jsonutils.JSONObject, error) {
	client := self.getAkClient()
	if len(query) > 0 {
		url = fmt.Sprintf("%s?%s", url, query.Encode())
	}
	var body jsonutils.JSONObject = nil
	if len(params) > 0 {
		body = jsonutils.Marshal(params)
	}
	header := http.Header{}
	if len(self.projectId) > 0 {
		header.Set("X-Project-Id", self.projectId)
	}
	_, resp, err := httputils.JSONRequest(client, context.Background(), method, url, header, body, self.debug)
	if err != nil {
		if e, ok := err.(*httputils.JSONClientError); ok && e.Code == 404 {
			return nil, errors.Wrapf(cloudprovider.ErrNotFound, err.Error())
		}
		return nil, err
	}
	return resp, err
}

func (self *SHuaweiClient) getAkList(regionId, resource string, query url.Values) (jsonutils.JSONObject, error) {
	uri := fmt.Sprintf("https://iam-pub.cn-south-278.aipcc-gz.com/v3.0/OS-CREDENTIAL/credentials")
	return self.request(httputils.GET, uri, query, nil)
}

// func (self *SHuaweiClient) createAkList(regionId, resource string, query url.Values) (jsonutils.JSONObject, error) {
// 	uri := fmt.Sprintf("https://iam-pub.%s.aipcc-gz.com/v3.0/OS-CREDENTIAL/credentials")
// 	return self.request(httputils.GET, uri, query, nil)
// }
