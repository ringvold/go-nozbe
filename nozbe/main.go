/*

The nozbe command will display a user's Toggl account information.

Usage:
    nozbe projects API_TOKEN

The API token can be retrieved from a user's account information page at toggl.com.

*/
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ringvold/go-nozbe"
)

func main() {
	if len(os.Args) != 3 {
		println("usage:")
		println(os.Args[0], "projects", "API_TOKEN")
		println(os.Args[0], "create-action", "API_TOKEN")
		return
	}
	lastArg := os.Args[len(os.Args)-1]
	client := nozbe.OpenSession(lastArg)

	// PROJECTS
	if os.Args[1] == "projects" {
		if len(os.Args) != 3 {
			println("usage:", os.Args[0], "projects", "API_TOKEN")
			return
		}
		projects, err := client.GetProjects()
		if err != nil {
			println("error:", err)
			return
		}
		data, err := json.MarshalIndent(&projects, "", "    ")
		println("projects:", string(data))
	}

	// CREATE ACTION
	if os.Args[1] == "create-action" {
		if len(os.Args) != 3 {
			println("usage:", os.Args[0], "create-action", "API_TOKEN")
			return
		}
		params := map[string]string{
			"project_id": "014f30c20b",
			"next":       "true",
		}
		projects, err := client.CreateAction("testnamelol", params)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		data, err := json.MarshalIndent(&projects, "", "    ")
		println("Response:", string(data))
	}

}
