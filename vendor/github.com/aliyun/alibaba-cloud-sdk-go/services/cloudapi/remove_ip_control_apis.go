package cloudapi

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// RemoveIpControlApis invokes the cloudapi.RemoveIpControlApis API synchronously
// api document: https://help.aliyun.com/api/cloudapi/removeipcontrolapis.html
func (client *Client) RemoveIpControlApis(request *RemoveIpControlApisRequest) (response *RemoveIpControlApisResponse, err error) {
	response = CreateRemoveIpControlApisResponse()
	err = client.DoAction(request, response)
	return
}

// RemoveIpControlApisWithChan invokes the cloudapi.RemoveIpControlApis API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/removeipcontrolapis.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RemoveIpControlApisWithChan(request *RemoveIpControlApisRequest) (<-chan *RemoveIpControlApisResponse, <-chan error) {
	responseChan := make(chan *RemoveIpControlApisResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.RemoveIpControlApis(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// RemoveIpControlApisWithCallback invokes the cloudapi.RemoveIpControlApis API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/removeipcontrolapis.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RemoveIpControlApisWithCallback(request *RemoveIpControlApisRequest, callback func(response *RemoveIpControlApisResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *RemoveIpControlApisResponse
		var err error
		defer close(result)
		response, err = client.RemoveIpControlApis(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// RemoveIpControlApisRequest is the request struct for api RemoveIpControlApis
type RemoveIpControlApisRequest struct {
	*requests.RpcRequest
	StageName     string `position:"Query" name:"StageName"`
	IpControlId   string `position:"Query" name:"IpControlId"`
	SecurityToken string `position:"Query" name:"SecurityToken"`
	GroupId       string `position:"Query" name:"GroupId"`
	ApiIds        string `position:"Query" name:"ApiIds"`
}

// RemoveIpControlApisResponse is the response struct for api RemoveIpControlApis
type RemoveIpControlApisResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateRemoveIpControlApisRequest creates a request to invoke RemoveIpControlApis API
func CreateRemoveIpControlApisRequest() (request *RemoveIpControlApisRequest) {
	request = &RemoveIpControlApisRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "RemoveIpControlApis", "apigateway", "openAPI")
	return
}

// CreateRemoveIpControlApisResponse creates a response to parse from RemoveIpControlApis response
func CreateRemoveIpControlApisResponse() (response *RemoveIpControlApisResponse) {
	response = &RemoveIpControlApisResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}