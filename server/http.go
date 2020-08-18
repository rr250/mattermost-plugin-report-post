package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/getreason", p.getReason).Methods("POST")
	r.HandleFunc("/reason", p.reportWithPredefinedReason).Methods("POST")
	r.HandleFunc("/customreason", p.reportWithCustomReason).Methods("POST")
	r.HandleFunc("/getcustomreason", p.getReportWithCustomReason).Methods("POST")
	r.HandleFunc("/deletepost", p.deletePost).Methods("POST")
	return r
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) getReason(w http.ResponseWriter, req *http.Request) {
	body, err1 := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err1 != nil {
		p.API.LogError("can't read body", err1)
	}

	if err2 := req.Body.Close(); err2 != nil {
		p.API.LogError("can't read body", err2)
	}
	var postDetails PostDetails
	if err3 := json.Unmarshal(body, &postDetails); err3 != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err4 := json.NewEncoder(w).Encode(err3); err4 != nil {
			p.API.LogError("can't unmarshall body", err4)
		}
	}

	post, err7 := p.API.GetPost(postDetails.PostID)
	if err7 != nil {
		p.API.LogError("can't fetch post", err7)
	}

	postModel := &model.Post{
		UserId:    p.botUserID,
		ChannelId: post.ChannelId,
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{
				{
					Text: "Why do you want to report post?",
					Actions: []*model.PostAction{
						{
							Integration: &model.PostActionIntegration{
								URL: fmt.Sprintf("/plugins/%s/reason", manifest.ID),
								Context: model.StringInterface{
									"action":       "reason",
									"reportpostid": postDetails.PostID,
								},
							},
							Name: "SELECT",
							Type: model.POST_ACTION_TYPE_SELECT,
							Options: []*model.PostActionOptions{
								{
									Text:  "SPAM",
									Value: "Spam",
								},
								{
									Text:  "INAPPROPRIATE",
									Value: "Inappropriate",
								},
								{
									Text:  "HARASSMENT",
									Value: "Harassment",
								},
								{
									Text:  "HATE SPEECH",
									Value: "Hate Speech",
								},
								{
									Text:  "HATE SPEECH",
									Value: "Hate Speech",
								},
								{
									Text:  "MOCKING",
									Value: "Mocking",
								},
								{
									Text:  "ABUSIVE",
									Value: "Abusive",
								},
							},
						},
						{
							Integration: &model.PostActionIntegration{
								URL: fmt.Sprintf("/plugins/%s/getcustomreason", manifest.ID),
								Context: model.StringInterface{
									"action":       "getcustomreason",
									"reportpostid": postDetails.PostID,
								},
							},
							Type: model.POST_ACTION_TYPE_BUTTON,
							Name: "CUSTOM REASON",
						},
					},
				},
			},
		},
	}

	p.API.SendEphemeralPost(postDetails.CurrentUserID, postModel)
}

func (p *Plugin) reportWithPredefinedReason(w http.ResponseWriter, req *http.Request) {
	request := model.PostActionIntegrationRequestFromJson(req.Body)
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
	p.postReport(request.UserId, request.Context["reportpostid"].(string), request.Context["selected_option"].(string))
	p.API.DeleteEphemeralPost(request.UserId, request.PostId)
}

func (p *Plugin) getReportWithCustomReason(w http.ResponseWriter, req *http.Request) {
	request := model.PostActionIntegrationRequestFromJson(req.Body)
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
	dialogRequest := model.OpenDialogRequest{
		TriggerId: request.TriggerId,
		URL:       fmt.Sprintf("/plugins/%s/customreason", manifest.ID),
		Dialog: model.Dialog{
			Title:       "Why do you want to report this post?",
			CallbackId:  model.NewId(),
			SubmitLabel: "Report",
			State:       request.Context["reportpostid"].(string),
			Elements: []model.DialogElement{
				{
					DisplayName: "Reason",
					Name:        "reason",
					Type:        "text",
					SubType:     "text",
					Default:     " ",
				},
			},
		},
	}
	if pErr := p.API.OpenInteractiveDialog(dialogRequest); pErr != nil {
		p.API.LogError("Failed opening interactive dialog " + pErr.Error())
	}
	p.API.DeleteEphemeralPost(request.UserId, request.PostId)
}

func (p *Plugin) reportWithCustomReason(w http.ResponseWriter, req *http.Request) {
	request := model.SubmitDialogRequestFromJson(req.Body)
	reason := request.Submission["reason"].(string)
	reportpostID := request.State
	p.postReport(request.UserId, reportpostID, reason)
}

func (p *Plugin) postReport(currentUserID string, reportpostID string, reason string) {
	configuration := p.getConfiguration()
	post, err1 := p.API.GetPost(reportpostID)
	if err1 != nil {
		p.API.LogError("failed to get post", err1)
	}
	currentUser, err2 := p.API.GetUser(currentUserID)
	if err2 != nil {
		p.API.LogError("failed to get current user", err2)
	}
	reportedUser, err3 := p.API.GetUser(post.UserId)
	if err2 != nil {
		p.API.LogError("failed to get reported user", err3)
	}
	channel, err8 := p.API.GetChannel(post.ChannelId)
	if err8 != nil {
		p.API.LogError("failed to get channel", err8)
	}
	reportpost := ReportPost{
		ID:               model.NewId(),
		ReportedBy:       currentUser.GetFullName(),
		ReportedByID:     currentUserID,
		CreatedAt:        time.Now(),
		ReportedName:     reportedUser.GetFullName(),
		ReportedID:       reportedUser.Id,
		ChannelID:        post.ChannelId,
		ChannelName:      channel.Name,
		ReportedUserName: reportedUser.Username,
		ReportedEmail:    reportedUser.Email,
		ReportedText:     post.Message,
		ReportedTextID:   post.Id,
		Reason:           reason,
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
						"\n\tReason: " + reportpost.Reason +
						"\n\tReported Text ID: " + reportpost.ReportedTextID +
						"\n\nReported Text:-\n" + reportpost.ReportedText,
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
	}
}

func writePostActionIntegrationResponseOk(w http.ResponseWriter, response *model.PostActionIntegrationResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response.ToJson())
}
