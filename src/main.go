package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Mod struct {
	Name     string `json:"name"`
	FileSize int64  `json:"fileSize"`
	Key      string `json:"key"`
	Optional bool   `json:"optional"`
}

type Repository struct {
	Name            string  `json:"name"`
	ServerAddress   string  `json:"serverAddress"`
	ServerPort      float64 `json:"serverPort"`
	Password        string  `json:"password"`
	BattlEyeEnabled bool    `json:"battlEyeEnabled"`
	Mods            []Mod   `json:"mods"`
}

type Repositories struct {
	UpdateUrl    string       `json:"updateUrl"`
	DeltaUpdates string       `json:"deltaUpdates"`
	Repositories []Repository `json:"repositories"`
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if len(os.Getenv("AFISYNC_SRV_INJ_SOURCE_UPDATE_URL")) == 0 {
		fmt.Fprintf(os.Stderr, "Environment variable \"%v\" is missing.\n", "AFISYNC_SRV_INJ_SOURCE_UPDATE_URL")
		return events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}, nil
	}
	if len(os.Getenv("AFISYNC_SRV_INJ_REPLACE_UPDATE_URL")) == 0 {
		fmt.Fprintf(os.Stderr, "Environment variable \"%v\" is missing.\n", "AFISYNC_SRV_INJ_REPLACE_UPDATE_URL")
		return events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}, nil
	}
	if len(os.Getenv("AFISYNC_SRV_INJ_SOURCE_REPOSITORY_NAME")) == 0 {
		fmt.Fprintf(os.Stderr, "Environment variable \"%v\" is missing.\n", "AFISYNC_SRV_INJ_SOURCE_REPOSITORY_NAME")
		return events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}, nil
	}
	if len(os.Getenv("AFISYNC_SRV_INJ_TARGET_REPOSITORY_NAME")) == 0 {
		fmt.Fprintf(os.Stderr, "Environment variable \"%v\" is missing.\n", "AFISYNC_SRV_INJ_TARGET_REPOSITORY_NAME")
		return events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}, nil
	}
	if len(os.Getenv("AFISYNC_SRV_INJ_TARGET_REPOSITORY_SERVER_ADDRESS")) == 0 {
		fmt.Fprintf(os.Stderr, "Environment variable \"%v\" is missing.\n", "AFISYNC_SRV_INJ_TARGET_REPOSITORY_SERVER_ADDRESS")
		return events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}, nil
	}
	if len(os.Getenv("AFISYNC_SRV_INJ_TARGET_REPOSITORY_PASSWORD")) == 0 {
		fmt.Fprintf(os.Stderr, "Environment variable \"%v\" is missing.\n", "AFISYNC_SRV_INJ_TARGET_REPOSITORY_PASSWORD")
		return events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}, nil
	}
	targetRepositoryServerPort, err := strconv.ParseFloat(os.Getenv("AFISYNC_SRV_INJ_TARGET_REPOSITORY_SERVER_PORT"), 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed parsing server port environment variable: %v\n", err)
		return events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}, nil
	}
	targetRepositoryBattlEyeEnabled, err := strconv.ParseBool(os.Getenv("AFISYNC_SRV_INJ_TARGET_REPOSITORY_BATTL_EYE_ENABLED"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed parsing battl eye enabled environment variable: %v\n", err)
		return events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}, nil
	}

	// JSON Marshal configuration
	const jsonPrefix = ""
	const jsonIndent = "    "

	// Setup HTTP client with preferred transport settings
	tr := &http.Transport{
		ResponseHeaderTimeout: 30 * time.Second,
	}
	client := &http.Client{Transport: tr}

	// Request our repositories source URL
	resp, err := client.Get(os.Getenv("AFISYNC_SRV_INJ_SOURCE_UPDATE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed requesting source repositories: %v\n", err)
		return events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// Trim Byte Order Mark (BOM) if exists
	bomBytes := []byte("\xef\xbb\xbf")
	bomTrimmed := false
	if bytes.Equal(body[0:3], bomBytes) {
		body = bytes.TrimPrefix(body, bomBytes)
		bomTrimmed = true
	}

	// Decode JSON from response body
	var reps Repositories
	err = json.Unmarshal(body, &reps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed decoding source repositories: %v\n", err)
		return events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}, nil
	}

	// Replace UpdateUrl if it matches the one we requested
	if os.Getenv("AFISYNC_SRV_INJ_SOURCE_UPDATE_URL") == reps.UpdateUrl {
		reps.UpdateUrl = os.Getenv("AFISYNC_SRV_INJ_REPLACE_UPDATE_URL")
	}

	// Make a copy from source repository and append it to the reps
	for i := 0; i < len(reps.Repositories); i++ {
		if reps.Repositories[i].Name == os.Getenv("AFISYNC_SRV_INJ_SOURCE_REPOSITORY_NAME") {
			var otsoRepository Repository = reps.Repositories[i]
			otsoRepository.Name = os.Getenv("AFISYNC_SRV_INJ_TARGET_REPOSITORY_NAME")
			otsoRepository.ServerAddress = os.Getenv("AFISYNC_SRV_INJ_TARGET_REPOSITORY_SERVER_ADDRESS")
			otsoRepository.ServerPort = targetRepositoryServerPort
			otsoRepository.Password = os.Getenv("AFISYNC_SRV_INJ_TARGET_REPOSITORY_PASSWORD")
			otsoRepository.BattlEyeEnabled = targetRepositoryBattlEyeEnabled
			reps.Repositories = append(reps.Repositories, otsoRepository)
			break
		}
	}

	// Encode modified repositories and return it
	encodedReps, err := json.MarshalIndent(reps, jsonPrefix, jsonIndent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed encoding repositories: %v\n", err)
		return events.APIGatewayProxyResponse{Body: "Internal server error", StatusCode: 500}, nil
	}

	// Insert Byte Order Mark (BOM) if it was previously trimmed
	if bomTrimmed {
		slices := [][]byte{bomBytes, encodedReps}
		encodedReps = bytes.Join(slices, []byte(""))
	}

	return events.APIGatewayProxyResponse{Body: string(encodedReps), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
