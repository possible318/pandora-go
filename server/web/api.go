package api

import (
	"bytes"
	"encoding/json"
	http "github.com/bogdanfinn/fhttp"
	tlsClient "github.com/bogdanfinn/tls-client"
	tokens "github.com/pandora_go/exts/token"
	reqTypes "github.com/pandora_go/web/typings/req"
)

const (
	apiPrefix = "https://ai.fakeopen.com"
	tokenKey  = ""
)

var (
	jar     = tlsClient.NewCookieJar()
	options = []tlsClient.HttpClientOption{
		tlsClient.WithTimeoutSeconds(360),
		tlsClient.WithClientProfile(tlsClient.Chrome_110),
		tlsClient.WithNotFollowRedirects(),
		tlsClient.WithCookieJar(jar),       // cookie jar
		tlsClient.WithInsecureSkipVerify(), // ssl
	}
	client, _ = tlsClient.NewHttpClient(tlsClient.NewNoopLogger(), options...)
)

type ChatBot struct {
}

func init() {
}

// ListModels 获取模型列表
func ListModels() (any, error) {
	url := apiPrefix + "/api/models"
	response, err := doRequest(url, http.MethodGet, nil, false)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		// Try read response body as JSON
		var errorResponse any
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return nil, err
		}
		return errorResponse, nil
	}
	var res any
	json.NewDecoder(response.Body).Decode(&res)
	return res, nil
}

// ListConversations 获取会话列表
func ListConversations(offset, limit string) (any, error) {
	url := apiPrefix + "/api/conversations?offset=" + offset + "&limit=" + limit
	response, err := doRequest(url, http.MethodGet, nil, false)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		// Try read response body as JSON
		var errorResponse any
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return nil, err
		}
		return errorResponse, err
	}

	print(response.Body)
	var res any
	json.NewDecoder(response.Body).Decode(&res)
	return res, nil
}

// GetConversation 获取会话内容
func GetConversation(id string) (any, error) {
	url := apiPrefix + "/api/conversation/" + id
	response, err := doRequest(url, http.MethodGet, nil, false)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		var errorResponse any
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return nil, err
		}
		return errorResponse, nil
	}
	var res any
	json.NewDecoder(response.Body).Decode(&res)
	return res, nil
}

// ClearConversations 清空会话列表
func ClearConversations() (any, error) {
	jsonBytes, _ := json.Marshal(reqTypes.ConversationRequest{
		IsVisible: false,
	})

	url := apiPrefix + "/api/conversations"
	response, err := doRequest(url, http.MethodDelete, jsonBytes, false)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		// Try read response body as JSON
		var errorResponse any
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return nil, err
		}
		return errorResponse, nil
	}

	var res any
	json.NewDecoder(response.Body).Decode(&res)
	return res, nil
}

// DelConversation  删除会话
func DelConversation(id string, params reqTypes.ConversationRequest) (any, error) {

	url := apiPrefix + "/api/conversation/" + id

	jsonBytes, _ := json.Marshal(params)

	response, err := doRequest(url, http.MethodPatch, jsonBytes, false)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		// Try read response body as JSON
		var errorResponse any
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return nil, err
		}
		return errorResponse, nil
	}

	var res any
	json.NewDecoder(response.Body).Decode(&res)
	return res, nil
}

// SetConversationTitle  重命名
func SetConversationTitle(id string, params reqTypes.ConversationRequest) (any, error) {
	jsonBytes, _ := json.Marshal(params)

	url := apiPrefix + "/api/conversation/" + id
	response, err := doRequest(url, http.MethodPatch, jsonBytes, false)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		// Try read response body as JSON
		var errorResponse any
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return nil, err
		}
		return errorResponse, nil
	}

	var res any
	json.NewDecoder(response.Body).Decode(&res)
	return res, nil
}

// GenConversationTitle 生成会话标题
func GenConversationTitle(id string, params reqTypes.GenerateTitleRequest) (any, error) {
	jsonBytes, _ := json.Marshal(params)

	url := apiPrefix + "/api/conversation/gen_title/" + id

	response, err := doRequest(url, http.MethodPost, jsonBytes, false)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		// Try read response body as JSON
		var errorResponse any
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return nil, err
		}
		return errorResponse, nil
	}

	var res any
	json.NewDecoder(response.Body).Decode(&res)
	return res, nil
}

// Conversation 发送消息
func Conversation(chatParam any) (*http.Response, error) {
	jsonBytes, _ := json.Marshal(chatParam)
	url := apiPrefix + "/api/conversation"

	response, _ := doRequest(url, http.MethodPost, jsonBytes, true)

	return response, nil
}

// doRequest 发送请求
func doRequest(url, method string, params []byte, sse bool) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(params))
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Pandora/1.0.3 Safari/537.36")
	req.Header.Set("Authorization", "Bearer "+tokens.GetToken(tokenKey))
	if sse {
		req.Header.Set("Accept", "text/event-stream")
	} else {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "*/*")
	}
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	return res, err
}
