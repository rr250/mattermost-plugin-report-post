package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/postreport", p.postReport).Methods("POST")
	return r
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) postReport(w http.ResponseWriter, req *http.Request) {
	body, err1 := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err1 != nil {
		p.API.LogError("can't read body", err1)
	}

	if err2 := req.Body.Close(); err2 != nil {
		p.API.LogError("can't read body", err2)
	}
	log.Println(string(body))
	var reportpost ReportPost
	if err3 := json.Unmarshal(body, &reportpost); err3 != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err4 := json.NewEncoder(w).Encode(err3); err4 != nil {
			p.API.LogError("can't unmarshall body", err4)
		}
	}
	log.Println(reportpost.ID)
	configuration := p.getConfiguration()
	log.Println(configuration)
	postModel := &model.Post{
		UserId:    p.botUserID,
		ChannelId: configuration.ChannelID,
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{
				{
					Text: "Report Alert:-\n\tReported: " +
						reportpost.ReportedName +
						"\n\tReported ID: " + reportpost.ReportedID +
						"\n\tReported Channel ID: " + reportpost.ChannelID +
						"\n\tReported Username: " + reportpost.ReportedUserName +
						"\n\tReported By: " + reportpost.ReportedBy +
						"\n\tReported By ID: " + reportpost.ReportedByID +
						"\n\tReported Text ID: " + reportpost.ReportedTextID +
						"\n\nReported Text:-\n" + reportpost.ReportedText + "\n ",
				},
			},
		},
	}

	_, err5 := p.API.CreatePost(postModel)
	if err5 != nil {
		p.API.LogError("failed to create post", err5)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err6 := json.NewEncoder(w).Encode(err5); err6 != nil {
			p.API.LogError("can't create post", err6)
		}
	}
}
