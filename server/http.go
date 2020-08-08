package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/postreport", p.postReport).Methods("POST")
	r.HandleFunc("/deletepost", p.deletePost).Methods("POST")
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
	var reportpost ReportPost
	if err3 := json.Unmarshal(body, &reportpost); err3 != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err4 := json.NewEncoder(w).Encode(err3); err4 != nil {
			p.API.LogError("can't unmarshall body", err4)
		}
	}
	configuration := p.getConfiguration()
	channel, err8 := p.API.GetChannel(reportpost.ChannelID)
	if err8 != nil {
		p.API.LogError("failed to get channel", err8)
	}
	reportpost.ChannelName = channel.Name
	postModel := &model.Post{
		UserId:    p.botUserID,
		ChannelId: configuration.ChannelID,
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{
				{
					Text: "Report Alert:-\n\tReported: " +
						reportpost.ReportedName +
						"\n\tReported User ID: " + reportpost.ReportedID +
						"\n\tReported Username: " + reportpost.ReportedUserName +
						"\n\tReported Email: " + reportpost.ReportedEmail +
						"\n\tReported By: " + reportpost.ReportedBy +
						"\n\tReported By User ID: " + reportpost.ReportedByID +
						"\n\tReported Channel ID: " + reportpost.ChannelID +
						"\n\tReported Channel Name: " + reportpost.ChannelName +
						"\n\tReported Text ID: " + reportpost.ReportedTextID +
						"\n\nReported Text:-\n" + reportpost.ReportedText + "\n ",
					Actions: []*model.PostAction{
						{
							Integration: &model.PostActionIntegration{
								URL: fmt.Sprintf("/plugins/%s/deletepost", manifest.ID),
								Context: model.StringInterface{
									"action":       "deletepost",
									"reportpostid": reportpost.ReportedTextID,
								},
							},
							Type: model.POST_ACTION_TYPE_BUTTON,
							Name: "Delete Post",
						},
					},
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

	if channel.Type == "O" {
		reactions, err9 := p.API.GetReactions(reportpost.ReportedTextID)
		if err9 != nil {
			p.API.LogError("failed to get reactions", err9)
		}
		emojis := []string{"x", "warning", "heavy_multiplication_x", "no_entry_sign", "bangbang"}
		for _, reaction1 := range reactions {
			for j, emoji := range emojis {
				if emoji == reaction1.EmojiName {
					emojis[j] = emojis[len(emojis)-1]
					emojis = emojis[:len(emojis)-1]
					break
				}
			}
		}
		if len(emojis) > 0 {
			reaction := &model.Reaction{
				UserId:    p.botUserID,
				PostId:    reportpost.ReportedTextID,
				EmojiName: emojis[0],
			}
			_, err7 := p.API.AddReaction(reaction)
			if err7 != nil {
				p.API.LogError("failed to add reaction", err7)
			}
		}
		privateChannel, err10 := p.API.GetDirectChannel(reportpost.ReportedID, p.botUserID)
		if err10 != nil {
			p.API.LogError("failed to get private channel", err10)
		}
		postModel = &model.Post{
			UserId:    p.botUserID,
			ChannelId: privateChannel.Id,
			Message:   "Someone reported your post: " + reportpost.ReportedText,
		}
		p.API.CreatePost(postModel)
	}
	postModel = &model.Post{
		UserId:    p.botUserID,
		ChannelId: reportpost.ChannelID,
		Message:   "Reported Successfully",
	}
	p.API.SendEphemeralPost(reportpost.ReportedByID, postModel)
}

func (p *Plugin) deletePost(w http.ResponseWriter, req *http.Request) {
	request := model.PostActionIntegrationRequestFromJson(req.Body)
	reportpostID := request.Context["reportpostid"].(string)
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
	err1 := p.API.DeletePost(reportpostID)
	if err1 != nil {
		p.API.LogError("failed to delete post", err1)
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   fmt.Sprintf("failed to delete post or post already deleted %s", err1),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
	} else {
		postModel := &model.Post{
			UserId:    request.UserId,
			ChannelId: request.ChannelId,
			Message:   fmt.Sprintf("Deleted successfully"),
		}
		p.API.SendEphemeralPost(request.UserId, postModel)
		// configuration := p.getConfiguration()
		// postModel = &model.Post{
		// 	UserId:    p.botUserID,
		// 	ChannelId: configuration.ChannelID,
		// 	Props: model.StringInterface{
		// 		"attachments": []*model.SlackAttachment{
		// 			{
		// 				Text: "Report Alert:-\n\tReported: " +
		// 					reportpost.ReportedName +
		// 					"\n\tReported User ID: " + reportpost.ReportedID +
		// 					"\n\tReported Username: " + reportpost.ReportedUserName +
		// 					"\n\tReported Email: " + reportpost.ReportedEmail +
		// 					"\n\tReported By: " + reportpost.ReportedBy +
		// 					"\n\tReported By User ID: " + reportpost.ReportedByID +
		// 					"\n\tReported Channel ID: " + reportpost.ChannelID +
		// 					"\n\tReported Channel Name: " + reportpost.ChannelName +
		// 					"\n\tReported Text ID: " + reportpost.ReportedTextID +
		// 					"\n\nReported Text:-\n" + reportpost.ReportedText + "\n ",
		// 			},
		// 		},
		// 	},
		// }

		// p.API.UpdatePost(postModel)
	}
}

func writePostActionIntegrationResponseOk(w http.ResponseWriter, response *model.PostActionIntegrationResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response.ToJson())
}
