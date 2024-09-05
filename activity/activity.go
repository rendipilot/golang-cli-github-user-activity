package activity

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/lipgloss"
)

type GitHubActivity struct {
	Type      string `json:"type"`
	Repo      Repo   `json:"repo"`
	CreatedAt string `json:"created_at"`
	Payload   struct {
		Action  string `json:"action"`
		Ref     string `json:"ref"`
		RefType string `json:"ref_type"`
		Commits []struct {
			Message string `json:"message"`
		} `json:"commits"`
	} `json:"payload"`
}

type Repo struct {
	Name string `json:"name"`
}

func FetchGithubActivity(username string) ([]GitHubActivity, int, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/events", username))

	if err != nil {
		return nil, res.StatusCode, err
	}

	if res.StatusCode == 404 {
		return nil, res.StatusCode, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode, err
	}

	var activity []GitHubActivity

	if err := json.NewDecoder(res.Body).Decode(&activity); err != nil {
		return nil, 0, err
	}

	return activity, 200, nil
}

func DisplayActivity(username string, events []GitHubActivity, code int) error {
	if code == 0 || code != 200 {
		return fmt.Errorf("failed to fetch GitHub activity: status code %d", code)
	}

	if len(events) == 0 {
		return fmt.Errorf("not events to display %d", code)
	}

	fmt.Println(lipgloss.NewStyle().Bold(true).Padding(1, 0).Foreground(lipgloss.Color("#FFCC66")).Render(fmt.Sprintf("%s's recent activity", username)))

	for _, event := range events {
		var action string
		switch event.Type {
		case "PushEvent":
			commitCount := len(event.Payload.Commits)
			action = fmt.Sprintf("Pushed %d commit(s) to %s", commitCount, event.Repo.Name)
		case "IssuesEvent":
			action = fmt.Sprintf("%s an issue in %s", event.Payload.Action, event.Repo.Name)
		case "WatchEvent":
			action = fmt.Sprintf("Starred %s", event.Repo.Name)
		case "ForkEvent":
			action = fmt.Sprintf("Forked %s", event.Repo.Name)
		case "CreateEvent":
			action = fmt.Sprintf("Created %s in %s", event.Payload.RefType, event.Repo.Name)
		default:
			action = fmt.Sprintf("%s in %s", event.Type, event.Repo.Name)
		}

		actionStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("#3C3C3C")).
			Render(fmt.Sprintf("- %s", action))
		fmt.Println(actionStyle)
	}

	return nil
}
