package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/alufers/owm-patrons/common"

	"github.com/google/go-github/v19/github"
	"golang.org/x/oauth2"
)

var config = common.PopulatedConfig

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		common.WriteJSON(w, 400, common.J{
			"message": "This request must be a POST",
		})
		return
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_KEY")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	ref, _, err := client.Git.GetRef(ctx, config.RepoAuthor, config.RepoName, "heads/"+config.Branch)
	if err != nil {
		fmt.Println("An error has occured when fetching ref", err.Error())
		common.WriteJSON(w, 500, common.J{
			"message": "An error has occured when fetching ref",
		})
		return
	}
	_, _, err = client.Git.UpdateRef(ctx, config.BotUsername, config.RepoName, ref, true)
	if err != nil {
		fmt.Println("An error has occured when updating ref", err.Error())
		common.WriteJSON(w, 500, common.J{
			"message": "An error has occured when updating ref",
		})
		return
	}

	fileContent, _, _, err := client.Repositories.GetContents(ctx, config.BotUsername, config.RepoName, config.PatronsFilePath, &github.RepositoryContentGetOptions{
		Ref: config.Branch,
	})

	if err != nil {
		fmt.Println("An error has occured fetching patrons.json contents", err.Error())
		common.WriteJSON(w, 500, common.J{
			"message": "An error has occured fetching patrons.json contents",
		})
		return
	}
	jsonData, err := fileContent.GetContent()
	if err != nil {
		panic(err)
	}
	patrons := common.J{}
	err = json.Unmarshal([]byte(jsonData), &patrons)
	if err != nil {
		fmt.Println("Failed to parse patrons.json", err.Error())
		common.WriteJSON(w, 500, common.J{
			"message": "Failed to parse patrons.json",
		})
		return
	}

	patrons["patrons"] = append(patrons["patrons"].([]interface{}), common.J{
		"username": "rnickson",
		"color":    "magenta",
		"tier":     "socjalizm krul",
		"badge":    "lewak",
	})
	newData, _ := json.MarshalIndent(patrons, "", "  ")
	// commitName := "otwarty-bot-pullujacy"
	// commitEmail := "ddd@bjorn.ml"
	commitMessage := "ddd"
	branch := config.Branch
	sha := fileContent.GetSHA()
	_, _, err = client.Repositories.UpdateFile(ctx, config.BotUsername, config.RepoName, config.PatronsFilePath, &github.RepositoryContentFileOptions{
		// Author: &github.CommitAuthor{
		// 	Name:  &commitName,
		// 	Email: &commitEmail,
		// },
		Branch:  &branch,
		Content: newData,
		Message: &commitMessage,
		SHA:     &sha,
	})
	if err != nil {
		fmt.Println("Failed to write patrons.json", err.Error())
		common.WriteJSON(w, 500, common.J{
			"message": "Failed to write patrons.json",
		})
		return
	}
	title := "Add new patrons"
	head := "otwarty-bot-pullujacy:" + config.Branch
	base := config.Branch
	body := `
		witam proszę mi zmergować tego pull requesta
		zmerguj że tego pulla człowieku
		> może grzeczniej
		poprosze feelfrelinuksie tego merga
		> wkurwił mnie pan ostro więc Panu nie dam
		ty kurwa biedaku jebany. merguj tego pierdolonego pulla
		oszukujesz uczciwego kontrybutora
		> po moim trupie dawaj napierdalamy się o niego kto ma więcej commitów
		Gdzie masz profil? Wyślę @alufers
		> Co @alufers Kurwa dawaj kontrybucje
	`
	client.PullRequests.Create(ctx, config.RepoAuthor, config.RepoName, &github.NewPullRequest{
		Title: &title,
		Head:  &head,
		Base:  &base,
		Body:  &body,
	})
	common.WriteJSON(w, 200, common.J{
		"message": "Success!",
		"data":    patrons,
	})
	//client.Repositories()
}
