package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
)

func CreateApiInfo() (_result *client.Params) {
	params := &client.Params{
		Action:      tea.String("GetAccountAlias"),
		Version:     tea.String("2015-05-01"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		Pathname:    tea.String("/"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	_result = params
	return _result
}

func buildCredential(accessKeyId string, accessKeySecret string) (credentials.Credential, error) {
	config := new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(accessKeyId).
		SetAccessKeySecret(accessKeySecret)

	provider, err := credentials.NewCredential(config)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func createClient(credential credentials.Credential) (result *client.Client, err error) {
	config := &client.Config{
		Credential: credential,
	}
	config.Endpoint = tea.String("ram.aliyuncs.com")
	result, err = client.NewClient(config)
	return
}

func main() {
	var (
		accessKeyID     string
		accessKeySecret string
	)
	flag.StringVar(&accessKeyID, "ak", "", "Aliyun Access Key ID")
	flag.StringVar(&accessKeySecret, "sk", "", "Aliyun Access Key Secret")
	flag.Parse()

	if accessKeyID == "" {
		accessKeyID = os.Getenv("ALIYUN_ACCESS_KEY_ID")
	}
	if accessKeySecret == "" {
		accessKeySecret = os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	}

	if accessKeyID == "" || accessKeySecret == "" {
		panic("missing required credentials. Please provide via -ak/-sk flags or ALIYUN_ACCESS_KEY_ID/ALIYUN_ACCESS_KEY_SECRET environment variables")
	}
	credential, err := buildCredential(accessKeyID, accessKeySecret)

	if err != nil {
		panic(err)
	}

	apiClient, err := createClient(credential)
	if err != nil {
		panic(err)
	}

	params := CreateApiInfo()
	runtime := &service.RuntimeOptions{}
	request := &client.OpenApiRequest{}
	resp, err := apiClient.CallApi(params, request, runtime)
	if err != nil {
		fmt.Println(err.Error())
	}

	if resp["StatusCode"] == 200 {
		fmt.Println("successful connected to Aliabab Cloud")
	} else {
		fmt.Println(resp["body"])
	}
}
