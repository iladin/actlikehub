package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type GitHubWorkflow struct {
	Jobs map[string]struct {
		RunsOn string                   `yaml:"runs-on"`
		Steps  []map[string]interface{} `yaml:"steps"`
	} `yaml:"jobs"`
}

type GitLabCI struct {
	Stages []string                          `yaml:"stages"`
	Jobs   map[string]map[string]interface{} `yaml:"jobs"`
}

func convertGitHubToGitLab(githubWorkflow GitHubWorkflow) map[string]interface{} {
	gitlabCI := map[string]interface{}{
		"stages": []string{},
	}

	for jobName, job := range githubWorkflow.Jobs {
		// Add the job name to the stages list
		gitlabCI["stages"] = append(gitlabCI["stages"].([]string), jobName)

		steps := []string{}
		for _, step := range job.Steps {
			if script, exists := step["run"]; exists {
				if scriptStr, ok := script.(string); ok { // Type assertion to string
					steps = append(steps, scriptStr)
				} else {
					fmt.Printf("Warning: 'run' step in job '%s' is not a string, skipping\n", jobName)
				}
			}
		}

		// Add the job directly to the root level
		gitlabCI[jobName] = map[string]interface{}{
			"stage":  jobName,
			"script": steps,
		}
	}

	return gitlabCI
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <github_workflow.yaml>")
		return
	}

	filePath := os.Args[1]
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var githubWorkflow GitHubWorkflow
	if err := yaml.Unmarshal(data, &githubWorkflow); err != nil {
		fmt.Println("Error parsing YAML:", err)
		return
	}

	gitlabCI := convertGitHubToGitLab(githubWorkflow)
	output, err := yaml.Marshal(gitlabCI)
	if err != nil {
		fmt.Println("Error generating GitLab CI YAML:", err)
		return
	}

	fmt.Println(string(output))
}
