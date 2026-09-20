package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Config struct {
	Format   string
	DryRun   bool
	MaxChars int
	Files    []string
	APIKey   string
}

type JevRequest struct {
	Model     string              `json:"model"`
	State     map[string]string   `json:"state"`
	Questions map[string]Question `json:"questions"`
}

type Question struct {
	Type         string      `json:"type"`
	Instructions string      `json:"instructions"`
	Criteria     []Criterion `json:"criteria,omitempty"`
}

type Criterion struct {
	Score       int    `json:"score"`
	Description string `json:"description"`
}

type JevResponse struct {
	Answers map[string]Answer `json:"answers"`
}

type Answer struct {
	Type       string  `json:"type"`
	Score      float64 `json:"score,omitempty"`
	Noul       float64 `json:"noul,omitempty"`
	Confidence float64 `json:"confidence"`
}

var cfg Config

var rootCmd = &cobra.Command{
	Use:   "eval-diff [files...]",
	Short: "Send code/diff content to TypeSafe Jev API for severity evaluation",
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Files = args
		cfg.APIKey = os.Getenv("TYPESAFE_API_KEY")

		loadEnvFile(".env")
		if cfg.APIKey == "" {
			cfg.APIKey = os.Getenv("TYPESAFE_API_KEY")
		}

		if err := validateEnv(&cfg); err != nil {
			logError(cfg.Format, err.Error())
			os.Exit(1)
		}

		content, err := readInput(cfg)
		if err != nil {
			logError(cfg.Format, err.Error())
			os.Exit(1)
		}

		if strings.TrimSpace(content) == "" {
			handleOutput(cfg, 1.0, 1.0, 0.0, nil)
			os.Exit(0)
		}

		content = truncateContent(content, cfg.MaxChars)

		if cfg.DryRun {
			printDryRun(cfg, len(content))
			os.Exit(0)
		}

		rawBody, err := sendAPIRequest(cfg.APIKey, content)
		if err != nil {
			logError(cfg.Format, err.Error())
			os.Exit(1)
		}

		jevResp, err := parseAPIResponse(rawBody)
		if err != nil {
			logError(cfg.Format, err.Error())
			os.Exit(1)
		}

		scoreAns := jevResp.Answers["impact_score"]
		breakingAns := jevResp.Answers["has_breaking_changes"]

		handleOutput(cfg, scoreAns.Score, scoreAns.Confidence, breakingAns.Noul, rawBody)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&cfg.Format, "format", "f", "standard", "Output format: 'standard', 'score', or 'full'")
	rootCmd.Flags().BoolVarP(&cfg.DryRun, "dry-run", "d", false, "Print prepared payload stats without invoking Jev API")
	rootCmd.Flags().IntVarP(&cfg.MaxChars, "max-chars", "m", 100000, "Maximum character threshold for truncating input content")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func loadEnvFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
	}
}

func validateEnv(cfg *Config) error {
	if cfg.APIKey == "" && !cfg.DryRun {
		return fmt.Errorf("TYPESAFE_API_KEY environment variable is not set")
	}
	return nil
}

func readInput(cfg Config) (string, error) {
	if len(cfg.Files) > 0 {
		return readFiles(cfg.Files)
	}
	return readStdin()
}

func readFiles(files []string) (string, error) {
	var builder strings.Builder
	for _, filePath := range files {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("Failed to read file %s: %w", filePath, err)
		}
		builder.Write(data)
		builder.WriteString("\n")
	}
	return builder.String(), nil
}

func readStdin() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", fmt.Errorf("No file arguments passed and no input piped via stdin")
	}

	stdinBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("Failed to read standard input: %w", err)
	}
	return string(stdinBytes), nil
}

func truncateContent(content string, maxChars int) string {
	if len(content) > maxChars {
		return content[:maxChars] + "\n... [truncated]"
	}
	return content
}

func buildPayload(content string) JevRequest {
	return JevRequest{
		Model: "jev-latest",
		State: map[string]string{
			"git_diff": content,
		},
		Questions: map[string]Question{
			"impact_score": {
				Type:         "score",
				Instructions: "Evaluate the breaking risk and operational severity of this code change on a 1-5 scale.",
				Criteria: []Criterion{
					{Score: 1, Description: "Docs only, comments, formatting, or zero-impact non-code changes"},
					{Score: 2, Description: "Minor internal change, refactor, or additive feature with zero API change"},
					{Score: 3, Description: "Moderate logic change, dependency update, or potential performance shift"},
					{Score: 4, Description: "High risk schema, core logic change, or API deprecation warning"},
					{Score: 5, Description: "Breaking change, API contract violation, or destructive database update"},
				},
			},
			"has_breaking_changes": {
				Type:         "noul",
				Instructions: "Does this content contain breaking API alterations or destructive schema changes?",
			},
		},
	}
}

func sendAPIRequest(apiKey, content string) ([]byte, error) {
	reqBody := buildPayload(content)
	jsonPayload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize request payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.typesafe.ai/v1/systemone", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("Failed to construct HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func parseAPIResponse(body []byte) (JevResponse, error) {
	var resp JevResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, fmt.Errorf("Failed to parse JSON response: %w", err)
	}
	return resp, nil
}

func printDryRun(cfg Config, charLen int) {
	fmt.Printf("Dry-Run Prepared Payload:\n")
	fmt.Printf("- Source: %d file(s) or stdin\n", len(cfg.Files))
	fmt.Printf("- Character Length: %d / %d max\n", charLen, cfg.MaxChars)
}

func handleOutput(cfg Config, score, confidence, breakingNoul float64, rawBody []byte) {
	setOutput("impact_score", fmt.Sprintf("%.2f", score))
	setOutput("confidence", fmt.Sprintf("%.2f", confidence))
	setOutput("is_breaking", fmt.Sprintf("%t", breakingNoul > 0.8))

	switch cfg.Format {
	case "score":
		fmt.Printf("%.2f\n", score)
	case "full":
		if len(rawBody) > 0 {
			fmt.Println(string(rawBody))
		} else {
			fmt.Println("{}")
		}
	default:
		fmt.Printf("Evaluation Summary:\n")
		fmt.Printf("- Impact Score: %.2f / 5.0\n", score)
		fmt.Printf("- Confidence: %.2f\n", confidence)
		fmt.Printf("- Breaking Change Probability: %.2f\n", breakingNoul)
	}
}

func setOutput(name, value string) {
	githubOutput := os.Getenv("GITHUB_OUTPUT")
	if githubOutput != "" {
		f, err := os.OpenFile(githubOutput, os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			defer f.Close()
			f.WriteString(fmt.Sprintf("%s=%s\n", name, value))
		}
	}
}

func logError(format, msg string) {
	if format == "standard" {
		fmt.Printf("::error::%s\n", msg)
	}
}
