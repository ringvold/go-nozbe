package nozbe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Nozbe service constants
const (
	NozbeAPI = "https://webapp.nozbe.com/api"
)

var dlog = log.New(os.Stderr, "[nozbe] ", log.LstdFlags)
var client = &http.Client{
	Timeout: time.Second * 30,
}

type Nozbe interface {
	GetProjects() []Project
	CreateAction(name string, params map[string]string) Action
}

// structures ///////////////////////////

type NozbeClient struct {
	APIToken string
	username string
	password string
}

type Action struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	NameShow    string    `json:"name_show"`
	Done        bool      `json:"done"`
	DoneTime    time.Time `json:"done_time"`
	ProjectID   string    `json:"project_id"`
	ProjectName string    `json:"project_name"`
	ContextID   string    `json:"context_id"`
	ContextName string    `json:"context_name"`
	ContextIcon string    `json:"context_icon"`
	Next        string    `json:"next"`
}

type Login struct {
	Key string
}

type Project struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Body     string `json:"body,omitempty"`
	BodyShow string `json:"body_show,omitempty"`
	Count    string `json:"count,omitempty"`
}

type CreateResponse struct {
	Response string `json:"response"`
}

// functions ////////////////////////////

// OpenSession opens a session using an existing API token.
func OpenSession(apiToken string) NozbeClient {
	return NozbeClient{APIToken: apiToken}
}

// // NewSession creates a new session by retrieving a user's API token.
func NewSession(username, password string) (session NozbeClient, err error) {
	session.username = username
	session.password = password

	data, err := session.request("/login", nil)
	if err != nil {
		return session, err
	}

	var login Login
	err = decodeLogin(data, &login)
	if err != nil {
		return session, err
	}

	session.username = ""
	session.password = ""
	session.APIToken = login.Key

	return session, nil
}

func (client *NozbeClient) GetProjects() ([]Project, error) {
	var projects []Project
	data, err := client.request("/projects", nil)
	if err != nil {
		return projects, err
	}

	err = decodeProjects(data, &projects)
	return projects, err
}

func (client *NozbeClient) CreateAction(name string, params map[string]string) (Action, error) {
	var action Action
	path := fmt.Sprintf("/newaction/name-%s", name)
	data, err := client.request(path, params)
	if err != nil {
		return action, err
	}
	var create CreateResponse
	err = decodeCreateResponse(data, &create)
	dlog.Println(create.Response)
	action.ID = create.Response
	return action, err
}

// support //////////////////////////////

func (session *NozbeClient) request(path string, params map[string]string) ([]byte, error) {
	requestURL := NozbeAPI + path

	if params != nil {
		for key, value := range params {
			requestURL += fmt.Sprintf("/%s-%s", key, value)
		}
	}

	if session.APIToken != "" {
		requestURL += fmt.Sprintf("/key-%s", session.APIToken)
	} else {
		requestURL += fmt.Sprintf("/email-%s/password-%s", session.username, session.password)
	}

	dlog.Printf("GETing from URL: %s", requestURL)
	req, err := http.NewRequest("GET", requestURL, nil)

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return content, fmt.Errorf(resp.Status)
	}

	return content, nil
}

func decodeProjects(data []byte, projects *[]Project) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(projects)
	if err != nil {
		return err
	}
	return nil
}

func decodeAction(data []byte, action *Action) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(action)
	if err != nil {
		return err
	}
	return nil
}

func decodeCreateResponse(data []byte, create *CreateResponse) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(create)
	if err != nil {
		return err
	}
	return nil
}

func decodeLogin(data []byte, login *Login) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(login)
	if err != nil {
		return err
	}
	return nil
}
