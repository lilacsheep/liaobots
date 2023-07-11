package liaobots

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
)

const (
	defaultAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67"
	userInfoUrl  = "https://liaobots.com/api/user"
	modelsUrl    = "https://liaobots.com/api/models"
	chatUrl      = "https://liaobots.com/api/chat"
	recommendUrl = "https://liaobots.com/api/recommend"
)

type Client struct {
	Token  string
	models []Model
}

func (c *Client) cli() *gclient.Client {
	return gclient.New()
}

func (c *Client) GetResponse(url string, req interface{}) (string, error) {
	cli := c.cli().ContentJson().SetAgent(defaultAgent)
	cli = cli.SetHeaderMap(g.MapStrStr{
		"Origin":      "https://liaobots.com",
		"Referer":     "https://liaobots.com",
		"X-Auth-Code": c.Token,
		"Authority":   "liaobots.com",
	})
	response, err := cli.Post(context.Background(), url, req)
	if err != nil {
		return "", err
	}
	defer response.Close()
	return response.ReadAllString(), nil
}

func (c *Client) UserInfo() (*UserResponse, error) {
	resp, err := c.GetResponse(userInfoUrl, &UserReq{Authcode: c.Token})
	if err != nil {
		return nil, err
	}
	var info = &UserResponse{}
	err = gjson.New(resp).Scan(info)
	return info, err
}

func (c *Client) Models() (*ModelsResponse, error) {
	resp, err := c.GetResponse(modelsUrl, &ModelReq{})
	if err != nil {
		return nil, err
	}
	var (
		info   = &ModelsResponse{}
		models []Model
	)
	err = gjson.New(resp).Scan(&models)
	if err != nil {
		return nil, err
	}
	info.Data = models
	return info, err
}

func (c *Client) Chat(req *ChatReq) (string, error) {
	// var (
	// 	req = &request.ChatReq{
	// 		ConversationID: uuid.NewV4().String(),
	// 		Model:          model,
	// 		Messages:       messages,
	// 	}
	// )
	if req.Model.ID == "" {
		return "", nil
	}
	resp, err := c.GetResponse(chatUrl, req)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (c *Client) Recommend(messages []Message) error {
	req := &RecommendReq{
		Messages: messages,
		AuthCode: c.Token,
	}
	_, err := c.GetResponse(recommendUrl, req)
	return err
}

func (c *Client) GetModel(id string) (*Model, error) {
	for _, model := range c.models {
		if model.ID == id {
			return &model, nil
		}
	}
	return nil, nil
}

func NewClient(token string) (*Client, error) {
	cli := &Client{Token: token}
	resp, err := cli.Models()
	if err != nil {
		return nil, err
	}
	cli.models = resp.Data
	return cli, nil
}