package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/alufers/owm-patrons/common"

	"github.com/google/go-github/v21/github"
	"golang.org/x/oauth2"
)

//J stands for JSON
type J map[string]interface{}

var config = common.PopulatedConfig

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	serialized, err := json.Marshal(data)
	if err != nil {
		serialized, _ = json.Marshal(err.Error())
	}
	w.Write(serialized)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_KEY")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	fileContent, _, _, err := client.Repositories.GetContents(ctx, config.RepoAuthor, config.RepoName, config.PatronsFilePath, &github.RepositoryContentGetOptions{
		Ref: config.Branch,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, "An error has occured fetching patrons.json contents", err.Error())
		writeJSON(w, 500, J{
			"message": "An error has occured fetching patrons.json contents",
			"error":   err.Error(),
		})
		return
	}
	jsonData, err := fileContent.GetContent()
	if err != nil {
		panic(err)
	}
	patrons := J{}
	err = json.Unmarshal([]byte(jsonData), &patrons)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse patrons.json", err.Error())
		writeJSON(w, 500, J{
			"message": "Failed to parse patrons.json",
		})
		return
	}
	w.Header().Add("Cache-Control", "s-maxage=3600, maxage=0")
	writeJSON(w, 200, patrons)
	//client.Repositories()
}
