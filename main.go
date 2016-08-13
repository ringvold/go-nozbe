package nozbe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	Timeout: time.Second * 10,
}

// structures ///////////////////////////

// Session represents an active connection to the Nozbe API.
type Session struct {
	APIToken string
	username string
	password string
}

// {
//   "id": "e6f437e4805",
//   "name": "test",
//   "name_show": "\ttest",
//   "done": 0,
//   "done_time": "0",
//   "time": "60",
//   "project_id": "014f30c20b",
//   "project_name": "Alfred Workflow for Nozbe",
//   "context_id": "5acebca56",
//   "context_name": "Computer",
//   "context_icon": "icon-47.png",
//   "next": "next"
// }

type Action struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	NameShow    string    `json:"name_show"`
	Done        bool      `json:"done"`
	DoneTime    time.Time `json:"done_time"`
	ProjectID   int       `json:"project_id"`
	ProjectName string    `json:"project_name"`
	ContextID   int       `json:"context_id"`
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

// functions ////////////////////////////

// OpenSession opens a session using an existing API token.
func OpenSession(apiToken string) Session {
	return Session{APIToken: apiToken}
}

// // NewSession creates a new session by retrieving a user's API token.
func NewSession(username, password string) (session Session, err error) {
	session.username = username
	session.password = password

	data, err := session.get("/login", nil)
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

func (session *Session) GetProjects() ([]Project, error) {
	var projects []Project
	// params := map[string]string{}
	data, err := session.get("/projects", nil)
	if err != nil {
		return projects, err
	}

	err = decodeProjects(data, &projects)
	return projects, err
}

// support //////////////////////////////

func (session *Session) request(method string, requestURL string, body io.Reader) ([]byte, error) {

	if session.APIToken != "" {
		requestURL += fmt.Sprintf("/key-%s", session.APIToken)
	} else {
		requestURL += fmt.Sprintf("/email-%s/password-%s", session.username, session.password)

	}

	req, err := http.NewRequest(method, requestURL, body)

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

func (session *Session) get(path string, params map[string]string) ([]byte, error) {
	requestURL := NozbeAPI + path

	if params != nil {
		for key, value := range params {
			requestURL += fmt.Sprintf("/%s-%s", key, value)
		}
	}

	dlog.Printf("GETing from URL: %s", requestURL)
	return session.request("GET", requestURL, nil)
}

func decodeProjects(data []byte, projects *[]Project) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(projects)
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
