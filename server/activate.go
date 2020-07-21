package main

import (
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

const (
	botUserName    = "reportpostbot"
	botDisplayName = "Report Post Bot"
)

// OnActivate register the plugin command
func (p *Plugin) OnActivate() error {
	botUserID, err := p.ensureBotExists()
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot user")
	}
	p.botUserID = botUserID
	p.router = p.InitAPI()
	return nil
}

func (p *Plugin) ensureBotExists() (string, error) {
	bot := &model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
	}

	return p.Helpers.EnsureBot(bot)
}
