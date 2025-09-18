package models

import "emby-plugin/utils"

type EmbyEvent struct {
	Date            string           `mapstructure:"Date" json:"Date"`
	Description     string           `mapstructure:"Description" json:"Description"`
	Event           string           `mapstructure:"Event" json:"Event"`
	Item            EmbyItem         `mapstructure:"Item" json:"Item"`
	PlaybackInfo    *PlaybackInfo    `mapstructure:"PlaybackInfo" json:"PlaybackInfo"`
	Server          EmbyServer       `mapstructure:"Server" json:"Server"`
	Session         *EmbySession     `mapstructure:"Session" json:"Session"`
	Title           string           `mapstructure:"Title" json:"Title"`
	User            *EmbyUser        `mapstructure:"User" json:"User"`
	TranscodingInfo *TranscodingInfo `mapstructure:"TranscodingInfo" json:"TranscodingInfo"`
}

type ExternalURL = utils.ExternalURL

type GenreItem struct {
	ID   int    `mapstructure:"Id" json:"Id"`
	Name string `mapstructure:"Name" json:"Name"`
}

type Studio struct {
	ID   int    `mapstructure:"Id" json:"Id"`
	Name string `mapstructure:"Name" json:"Name"`
}

type RemoteTrailer struct {
	URL string `mapstructure:"Url" json:"Url"`
}

type EmbyItem struct {
	BackdropImageTags       []string          `mapstructure:"BackdropImageTags" json:"BackdropImageTags"`
	Bitrate                 int               `mapstructure:"Bitrate" json:"Bitrate"`
	CommunityRating         float64           `mapstructure:"CommunityRating" json:"CommunityRating"`
	Container               string            `mapstructure:"Container" json:"Container"`
	DateCreated             string            `mapstructure:"DateCreated" json:"DateCreated"`
	ExternalUrls            []ExternalURL     `mapstructure:"ExternalUrls" json:"ExternalUrls"`
	FileName                string            `mapstructure:"FileName" json:"FileName"`
	GenreItems              []GenreItem       `mapstructure:"GenreItems" json:"GenreItems"`
	Genres                  []string          `mapstructure:"Genres" json:"Genres"`
	Height                  int               `mapstructure:"Height" json:"Height"`
	ID                      string            `mapstructure:"Id" json:"Id"`
	ImageTags               map[string]string `mapstructure:"ImageTags" json:"ImageTags"`
	IndexNumber             int               `mapstructure:"IndexNumber" json:"IndexNumber"`
	IsFolder                bool              `mapstructure:"IsFolder" json:"IsFolder"`
	MediaType               string            `mapstructure:"MediaType" json:"MediaType"`
	Name                    string            `mapstructure:"Name" json:"Name"`
	OfficialRating          string            `mapstructure:"OfficialRating" json:"OfficialRating"`
	OriginalTitle           string            `mapstructure:"OriginalTitle" json:"OriginalTitle"`
	Overview                string            `mapstructure:"Overview" json:"Overview"`
	ParentBackdropImageTags []string          `mapstructure:"ParentBackdropImageTags" json:"ParentBackdropImageTags"`
	ParentBackdropItemId    string            `mapstructure:"ParentBackdropItemId" json:"ParentBackdropItemId"`
	ParentID                string            `mapstructure:"ParentId" json:"ParentId"`
	ParentIndexNumber       int               `mapstructure:"ParentIndexNumber" json:"ParentIndexNumber"`
	ParentLogoImageTag      string            `mapstructure:"ParentLogoImageTag" json:"ParentLogoImageTag"`
	ParentLogoItemId        string            `mapstructure:"ParentLogoItemId" json:"ParentLogoItemId"`
	Path                    string            `mapstructure:"Path" json:"Path"`
	PremiereDate            string            `mapstructure:"PremiereDate" json:"PremiereDate"`
	PrimaryImageAspectRatio float64           `mapstructure:"PrimaryImageAspectRatio" json:"PrimaryImageAspectRatio"`
	ProductionYear          int               `mapstructure:"ProductionYear" json:"ProductionYear"`
	ProviderIds             map[string]string `mapstructure:"ProviderIds" json:"ProviderIds"`
	RemoteTrailers          []RemoteTrailer   `mapstructure:"RemoteTrailers" json:"RemoteTrailers"`
	RunTimeTicks            int64             `mapstructure:"RunTimeTicks" json:"RunTimeTicks"`
	SeasonID                string            `mapstructure:"SeasonId" json:"SeasonId"`
	SeasonName              string            `mapstructure:"SeasonName" json:"SeasonName"`
	SeriesID                string            `mapstructure:"SeriesId" json:"SeriesId"`
	SeriesName              string            `mapstructure:"SeriesName" json:"SeriesName"`
	SeriesPrimaryImageTag   string            `mapstructure:"SeriesPrimaryImageTag" json:"SeriesPrimaryImageTag"`
	ServerID                string            `mapstructure:"ServerId" json:"ServerId"`
	Size                    int64             `mapstructure:"Size" json:"Size"`
	SortName                string            `mapstructure:"SortName" json:"SortName"`
	Studios                 []Studio          `mapstructure:"Studios" json:"Studios"`
	TagItems                []map[string]any  `mapstructure:"TagItems" json:"TagItems"`
	Taglines                []string          `mapstructure:"Taglines" json:"Taglines"`
	Type                    string            `mapstructure:"Type" json:"Type"`
	Width                   int               `mapstructure:"Width" json:"Width"`
	Album                   string            `mapstructure:"Album" json:"Album"`
}

type TranscodingInfo struct {
	CompletionPercentage float64 `mapstructure:"CompletionPercentage" json:"CompletionPercentage"`
}

type PlaybackInfo struct {
	PlaySessionID  string `mapstructure:"PlaySessionId" json:"PlaySessionId"`
	PlaylistIndex  int    `mapstructure:"PlaylistIndex" json:"PlaylistIndex"`
	PlaylistLength int    `mapstructure:"PlaylistLength" json:"PlaylistLength"`
	PositionTicks  int64  `mapstructure:"PositionTicks" json:"PositionTicks"`
}

type EmbyServer struct {
	ID      string `mapstructure:"Id" json:"Id"`
	Name    string `mapstructure:"Name" json:"Name"`
	Version string `mapstructure:"Version" json:"Version"`
}

type EmbySession struct {
	ApplicationVersion string `mapstructure:"ApplicationVersion" json:"ApplicationVersion"`
	Client             string `mapstructure:"Client" json:"Client"`
	DeviceID           string `mapstructure:"DeviceId" json:"DeviceId"`
	DeviceName         string `mapstructure:"DeviceName" json:"DeviceName"`
	ID                 string `mapstructure:"Id" json:"Id"`
	RemoteEndPoint     string `mapstructure:"RemoteEndPoint" json:"RemoteEndPoint"`
}

type EmbyUser struct {
	ID   string `mapstructure:"Id" json:"Id"`
	Name string `mapstructure:"Name" json:"Name"`
}
