package main

import (
	"SharepointBot/config"
	"SharepointBot/db"
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/imroc/req/v3"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var SCOPE = "https://graph.microsoft.com/Files.Read.All https://graph.microsoft.com/Sites.Read.All"

type OAUTH2CallbackBody struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
	Scope        string `json:"scope"`
	GrantType    string `json:"grant_type"`
}

type MicrosoftOUATH2Response struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SharepointResponse struct {
	OdataContext  string `json:"@odata.context"`
	OdataNextLink string `json:"@odata.nextLink"`
	Value         []struct {
		OdataEtag            string    `json:"@odata.etag"`
		CreatedDateTime      time.Time `json:"createdDateTime"`
		ETag                 string    `json:"eTag"`
		Id                   string    `json:"id"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
		WebUrl               string    `json:"webUrl"`
		CreatedBy            struct {
			User struct {
				Email       string `json:"email"`
				Id          string `json:"id"`
				DisplayName string `json:"displayName"`
			} `json:"user"`
		} `json:"createdBy"`
		LastModifiedBy struct {
			User struct {
				Email       string `json:"email"`
				Id          string `json:"id"`
				DisplayName string `json:"displayName"`
			} `json:"user"`
		} `json:"lastModifiedBy"`
		ParentReference struct {
			Id     string `json:"id"`
			SiteId string `json:"siteId"`
		} `json:"parentReference"`
		ContentType struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"contentType"`
	} `json:"value"`
}

type SharepointNotificationResponse struct {
	OdataContext         string    `json:"@odata.context"`
	OdataEtag            string    `json:"@odata.etag"`
	CreatedDateTime      time.Time `json:"createdDateTime"`
	ETag                 string    `json:"eTag"`
	Id                   string    `json:"id"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	WebUrl               string    `json:"webUrl"`
	CreatedBy            struct {
		User struct {
			Email       string `json:"email"`
			Id          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"user"`
	} `json:"createdBy"`
	LastModifiedBy struct {
		User struct {
			Email       string `json:"email"`
			Id          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"user"`
	} `json:"lastModifiedBy"`
	ParentReference struct {
		Id     string `json:"id"`
		SiteId string `json:"siteId"`
	} `json:"parentReference"`
	ContentType struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"contentType"`
	FieldsOdataContext string `json:"fields@odata.context"`
	Fields             struct {
		OdataEtag                string    `json:"@odata.etag"`
		Title                    string    `json:"Title"`
		ModerationStatus         int       `json:"_ModerationStatus"`
		Body                     string    `json:"Body"`
		Expires                  time.Time `json:"Expires"`
		Id                       string    `json:"id"`
		ContentType              string    `json:"ContentType"`
		Modified                 time.Time `json:"Modified"`
		Created                  time.Time `json:"Created"`
		AuthorLookupId           string    `json:"AuthorLookupId"`
		EditorLookupId           string    `json:"EditorLookupId"`
		UIVersionString          string    `json:"_UIVersionString"`
		Attachments              bool      `json:"Attachments"`
		Edit                     string    `json:"Edit"`
		LinkTitleNoMenu          string    `json:"LinkTitleNoMenu"`
		LinkTitle                string    `json:"LinkTitle"`
		ItemChildCount           string    `json:"ItemChildCount"`
		FolderChildCount         string    `json:"FolderChildCount"`
		ComplianceFlags          string    `json:"_ComplianceFlags"`
		ComplianceTag            string    `json:"_ComplianceTag"`
		ComplianceTagWrittenTime string    `json:"_ComplianceTagWrittenTime"`
		ComplianceTagUserId      string    `json:"_ComplianceTagUserId"`
	} `json:"fields"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type Embed struct {
	Author struct {
		Name    string `json:"name"`
		URL     string `json:"url"`
		IconURL string `json:"icon_url"`
	} `json:"author"`
	Title       string       `json:"title"`
	URL         string       `json:"url"`
	Description string       `json:"description"`
	Color       int          `json:"color"`
	Fields      []EmbedField `json:"fields"`
	Thumbnail   struct {
		URL string `json:"url"`
	} `json:"thumbnail"`
	Image struct {
		URL string `json:"url"`
	} `json:"image"`
	Footer struct {
		Text    string `json:"text"`
		IconURL string `json:"icon_url"`
	} `json:"footer"`
}

type WebhookBody struct {
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Content   string  `json:"content"`
	Embeds    []Embed `json:"embeds"`
}

type DiscordWebhookResponse struct {
	Type         int    `json:"type"`
	Content      string `json:"content"`
	Mentions     []any  `json:"mentions"`
	MentionRoles []any  `json:"mention_roles"`
	Attachments  []any  `json:"attachments"`
	Embeds       []struct {
		Type   string `json:"type"`
		URL    string `json:"url"`
		Color  int    `json:"color"`
		Fields []struct {
			Name   string `json:"name"`
			Value  string `json:"value"`
			Inline bool   `json:"inline"`
		} `json:"fields"`
		Thumbnail struct {
			URL      string `json:"url"`
			ProxyURL string `json:"proxy_url"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			Flags    int    `json:"flags"`
		} `json:"thumbnail"`
	} `json:"embeds"`
	Timestamp       time.Time `json:"timestamp"`
	EditedTimestamp any       `json:"edited_timestamp"`
	Flags           int       `json:"flags"`
	Components      []any     `json:"components"`
	ID              string    `json:"id"`
	ChannelID       string    `json:"channel_id"`
	Author          struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Avatar        any    `json:"avatar"`
		Discriminator string `json:"discriminator"`
		PublicFlags   int    `json:"public_flags"`
		Flags         int    `json:"flags"`
		Bot           bool   `json:"bot"`
		GlobalName    any    `json:"global_name"`
		Clan          any    `json:"clan"`
	} `json:"author"`
	Pinned          bool   `json:"pinned"`
	MentionEveryone bool   `json:"mention_everyone"`
	Tts             bool   `json:"tts"`
	WebhookID       string `json:"webhook_id"`
}

func (server *httpImpl) SendNotificationToWebhook(webhook string, editing bool, notification db.SharepointNotification) string {
	if !editing {
		webhook += "?wait=true"
	}

	if len([]rune(notification.Description)) > 4096 {
		notification.Description = string([]rune(notification.Description)[0:4093]) + "..."
	}

	request := req.C().DevMode().R()

	createdOn := time.Unix(int64(notification.CreatedOn), 0)
	created := createdOn.Format("02. 01. 2006 ob 15.04")

	modifiedOn := time.Unix(int64(notification.ModifiedOn), 0)
	modified := modifiedOn.Format("02. 01. 2006 ob 15.04")

	body := WebhookBody{
		Username:  "Intranet",
		AvatarURL: "",
		Content:   "Novo obvestilo na intranetu",
		Embeds: []Embed{
			{
				Author: struct {
					Name    string `json:"name"`
					URL     string `json:"url"`
					IconURL string `json:"icon_url"`
				}{Name: notification.CreatedBy, URL: "", IconURL: ""},
				Title:       notification.Name,
				Description: notification.Description,
				Color:       15258703,
				URL:         fmt.Sprintf("https://gimnazijabezigrad.sharepoint.com/Lists/ObvAkt/DispForm.aspx?ID=%s", notification.ID),
				Fields: []EmbedField{
					{
						Name:   "Ustvarjeno",
						Value:  fmt.Sprintf("`%s`", created),
						Inline: true,
					},
					{
						Name:   "Nazadnje spremenjeno",
						Value:  fmt.Sprintf("`%s`", modified),
						Inline: true,
					},
					{
						Name:   "Nazadnje spremenil",
						Value:  fmt.Sprintf("`%s`", notification.ModifiedBy),
						Inline: true,
					},
				},
				Thumbnail: struct {
					URL string `json:"url"`
				}{URL: "https://www.gimb.org/wp-content/uploads/2017/01/logo.png"},
			},
		},
	}

	request.SetBodyJsonMarshal(body)

	var resp *req.Response
	var err error
	if editing {
		resp, err = request.Patch(webhook)
	} else {
		resp, err = request.Post(webhook)
	}
	if resp == nil || err != nil {
		server.logger.Errorw("failure while sending message to discord webhook", "err", err)
		return ""
	}
	server.logger.Infow("Discord responded with status code", "statusCode", resp.StatusCode)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		server.logger.Errorw("error while sending message to Discord", "body", resp.String())
	}

	if editing {
		return ""
	}

	var unmarshal DiscordWebhookResponse
	err = resp.Unmarshal(&unmarshal)
	if err != nil {
		server.logger.Errorw("could not unmarshal response", "body", resp.String())
		return ""
	}

	return unmarshal.ID
}

func (server *httpImpl) GetSharepointNotificationsGoroutine(accessToken string) {
	server.logger.Infow("getting Sharepoint notifications")

	client := req.C()

	client.Headers = make(http.Header)
	client.Headers.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	nextLink := "https://graph.microsoft.com/v1.0/sites/root/lists/54521912-06dd-4ccc-8edb-8173c9629fd8/items"
	for nextLink != "" {
		res, err := client.R().Get(nextLink)
		if err != nil {
			server.logger.Errorw("error getting all Sharepoint items", "err", err)
			break
		}

		var response SharepointResponse
		err = res.UnmarshalJson(&response)
		if err != nil {
			server.logger.Errorw("error parsing Microsoft response", "err", err)
			break
		}

		nextLink = response.OdataNextLink
		if nextLink != "" {
			server.logger.Infow("got next page on Sharepoint", "page", nextLink)
		}

		for _, v := range response.Value {
			notificationDb, noterr := server.db.GetSharepointNotification(v.Id)
			if (noterr == nil && notificationDb.ModifiedOn == int(v.LastModifiedDateTime.Unix())) || (noterr != nil && !errors.Is(noterr, sql.ErrNoRows)) {
				if err != nil {
					server.logger.Errorw("error retrieving Sharepoint notification", "id", v.Id, "webUrl", v.WebUrl, "err", err)
				}
				continue
			}

			res, err = client.R().Get(fmt.Sprintf("https://graph.microsoft.com/v1.0/sites/root/lists/54521912-06dd-4ccc-8edb-8173c9629fd8/items/%s", v.Id))
			if err != nil {
				server.logger.Errorw("error getting a Sharepoint notification", "id", v.Id, "err", err)
				break
			}

			var notificationResponse SharepointNotificationResponse
			err = res.UnmarshalJson(&notificationResponse)
			if err != nil {
				server.logger.Errorw("error parsing Sharepoint notification response", "id", v.Id, "err", err)
				break
			}

			// ne posodabljaj za vsak drek
			if noterr == nil && int(notificationResponse.Fields.Modified.Unix()) == notificationDb.ModifiedOn {
				continue
			}

			opt := &md.Options{}
			converter := md.NewConverter("", true, opt)
			markdown, err := converter.ConvertString(notificationResponse.Fields.Body)
			if err != nil {
				server.logger.Errorw("error parsing Sharepoint HTML", "err", err)
				break
			}

			// ker discord je pač retarded
			r := regexp.MustCompile(`\[(?P<URL>.*)]\(.*\)`)
			res := r.FindAllStringSubmatch(markdown, -1)
			for _, l := range res {
				if len(l) < 2 {
					continue
				}
				markdown = strings.ReplaceAll(markdown, l[0], l[1])
			}

			notificationResponse.Fields.Body = markdown

			expires := int(notificationResponse.Fields.Expires.Unix())
			if expires < 0 {
				expires = 0
			}

			if errors.Is(noterr, sql.ErrNoRows) {
				server.logger.Infow("creating new notification", "id", v.Id)

				not := db.SharepointNotification{
					ID:             notificationResponse.Id,
					Name:           notificationResponse.Fields.Title,
					Description:    notificationResponse.Fields.Body,
					CreatedOn:      int(notificationResponse.Fields.Created.Unix()),
					ModifiedOn:     int(notificationResponse.Fields.Modified.Unix()),
					CreatedBy:      notificationResponse.CreatedBy.User.DisplayName,
					ModifiedBy:     notificationResponse.LastModifiedBy.User.DisplayName,
					MessageIDs:     "[]",
					ExpiresOn:      expires,
					HasAttachments: notificationResponse.Fields.Attachments,
				}

				ids := make([]string, 0)
				for _, webhook := range server.config.Webhooks {
					id := server.SendNotificationToWebhook(webhook, false, not)
					if id == "" {
						continue
					}
					ids = append(ids, fmt.Sprintf("%s/messages/%s", webhook, id))
				}

				marshal, err := json.Marshal(ids)
				if err != nil {
					server.logger.Errorw("error while marshalling message IDs", "err", err)
					continue
				}
				not.MessageIDs = string(marshal)

				err = server.db.InsertSharepointNotification(not)
				if err != nil {
					server.logger.Errorw("error inserting Sharepoint notification", "id", v.Id, "notification", notificationResponse, "not", not, "err", err)
					continue
				}
			} else {
				server.logger.Infow("updating an existing notification", "id", v.Id)

				notificationDb.ModifiedOn = int(notificationResponse.Fields.Modified.Unix())
				notificationDb.ModifiedBy = notificationResponse.LastModifiedBy.User.DisplayName
				notificationDb.ExpiresOn = expires
				notificationDb.Name = notificationResponse.Fields.Title
				notificationDb.Description = notificationResponse.Fields.Body

				err := server.db.UpdateSharepointNotification(notificationDb)
				if err != nil {
					server.logger.Errorw("error updating Sharepoint notification", "id", v.Id, "notification", notificationResponse, "not", notificationDb, "err", err)
					continue
				}

				var unmarshal []string
				err = json.Unmarshal([]byte(notificationDb.MessageIDs), &unmarshal)
				if err != nil {
					server.logger.Errorw("error unmarshalling message IDs", "err", err)
					continue
				}

				for _, webhook := range unmarshal {
					server.SendNotificationToWebhook(webhook, true, notificationDb)
				}
			}
		}
	}
}

func (server *httpImpl) SharepointGoroutine() {
	server.logger.Infow("starting Sharepoint goroutine")

	for {
		if server.config.MicrosoftOAUTH2RefreshToken == "" {
			server.logger.Infow("no Microsoft OAUTH2 refresh token was found")
			server.MicrosoftOAUTH2URL()
			server.MicrosoftOAUTH2Callback()
			return // konča gorutino, avtomatično znova zažene program
		}

		client := req.C()

		body := map[string]string{
			"client_id":     server.config.MicrosoftOAUTH2ClientID,
			"client_secret": server.config.MicrosoftOAUTH2Secret,
			"refresh_token": server.config.MicrosoftOAUTH2RefreshToken,
			"scope":         SCOPE,
			"grant_type":    "refresh_token",
		}

		res, err := client.R().SetFormData(body).Post("https://login.microsoftonline.com/organizations/oauth2/v2.0/token")
		if err != nil {
			server.logger.Errorw("error getting token", "err", err)
			break
		}

		var response MicrosoftOUATH2Response
		err = res.UnmarshalJson(&response)
		if err != nil {
			server.logger.Errorw("error parsing Microsoft response", "err", err)
			break
		}

		accessToken := response.AccessToken
		refreshToken := response.RefreshToken

		server.config.MicrosoftOAUTH2RefreshToken = refreshToken
		err = config.SaveConfig(server.config)
		if err != nil {
			server.logger.Errorw("error saving config", "err", err)
			break
		}

		server.GetSharepointNotificationsGoroutine(accessToken)

		server.logger.Infow("ran Sharepoint goroutine")
		time.Sleep(time.Hour)
	}

	server.logger.Infow("exiting Sharepoint goroutine")
}

func (server *httpImpl) MicrosoftOAUTH2URL() {
	fmt.Printf("Obiščite stran in avtorizirajte session: https://login.microsoftonline.com/organizations/oauth2/v2.0/authorize?client_id=%s&response_type=code&response_mode=query&scope=offline_access %s\n", server.config.MicrosoftOAUTH2ClientID, SCOPE)
}

func (server *httpImpl) MicrosoftOAUTH2Callback() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Microsoft code: ")
	code, err := reader.ReadString('\n')
	if err != nil {
		server.logger.Fatalw("error reading input", "err", err)
	}

	client := req.C()

	body := map[string]string{
		"client_id":     server.config.MicrosoftOAUTH2ClientID,
		"client_secret": server.config.MicrosoftOAUTH2Secret,
		"code":          code,
		"scope":         SCOPE,
		"grant_type":    "authorization_code",
	}

	res, err := client.R().SetFormData(body).Post("https://login.microsoftonline.com/organizations/oauth2/v2.0/token")
	if err != nil {
		server.logger.Fatalw("error getting token", "err", err)
		return
	}

	var response MicrosoftOUATH2Response
	err = res.UnmarshalJson(&response)
	if err != nil {
		server.logger.Fatalw("error unmarshalling token", "err", err)
		return
	}

	server.config.MicrosoftOAUTH2RefreshToken = response.RefreshToken
	err = config.SaveConfig(server.config)
	if err != nil {
		server.logger.Fatalw("error saving token", "err", err)
		return
	}

	server.logger.Infow("token received successfully")
}
