package main

import (
	"encoding/json"
	"net/http"

	mattermostModel "github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	// All HTTP endpoints of this plugin require a logged in user.
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	// All HTTP endpoints of this plugin require the user to be a System Admin
	if !p.Client.User.HasPermissionTo(userID, mattermostModel.PermissionManageSystem) {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
	}

	p.Client.Log.Info(r.URL.Path)

	switch path := r.URL.Path; path {
	case "/api/v1/legalhold/list":
		p.listLegalHolds(w, r)
		return
	default:
		http.NotFound(w, r)
	}
}

// listLegalHolds serves a list of all LegalHold objects
func (p *Plugin) listLegalHolds(w http.ResponseWriter, r *http.Request) {
	legalHolds, err := p.KVStore.GetAllLegalHolds()
	if err != nil {
		http.Error(w, "an error occurred fetching the legal holds", http.StatusInternalServerError)
		p.Client.Log.Error(err.Error())
		return
	}

	b, jsonErr := json.Marshal(legalHolds)
	if jsonErr != nil {
		http.Error(w, "Error encoding json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		p.API.LogError("failed to write http response", err.Error())
		return
	}
}