package config

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// IsRemoteURL checks if the path is a remote URL
func IsRemoteURL(path string) bool {
	return isHTTPURL(path) || isGitURL(path)
}

// FetchRemoteConfig fetches a config file from a remote location
// Supports:
// - HTTP/HTTPS URLs: https://example.com/config.yaml
// - Git HTTPS: https://github.com/org/repo.git//path/to/config.yaml?ref=main
// - Git SSH: git@github.com:org/repo.git//path/to/config.yaml?ref=main
func FetchRemoteConfig(remotePath string) ([]byte, error) {
	if !hasValidExtension(remotePath) {
		return nil, fmt.Errorf("unsupported file type: must be .yaml, .yml, or .json")
	}

	// Check Git first as it's more specific
	if isGitURL(remotePath) {
		return fetchFromGit(remotePath)
	}

	// HTTP check is simpler, do it second
	if isHTTPURL(remotePath) {
		return fetchFromHTTP(remotePath)
	}

	return nil, fmt.Errorf("unsupported remote path format")
}

func hasValidExtension(path string) bool {
	// Remove query parameters if present
	if idx := strings.IndexByte(path, '?'); idx != -1 {
		path = path[:idx]
	}
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml" || ext == ".json"
}

func isHTTPURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func isGitURL(path string) bool {
	// Git SSH format - most specific check first
	if strings.HasPrefix(path, "git@") {
		return true
	}
	// Git HTTPS with .git and // separator (Terraform/Checkov convention)
	return (strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://")) &&
		strings.Contains(path, ".git//")
}

func fetchFromHTTP(urlStr string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func fetchFromGit(gitPath string) ([]byte, error) {
	repoURL, filePath, ref, err := parseGitURL(gitPath)
	if err != nil {
		return nil, err
	}

	tmpDir, err := os.MkdirTemp("", "terratags-git-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	cloneOpts := &git.CloneOptions{
		URL:      repoURL,
		Progress: nil,
		Depth:    1,
	}

	if ref != "" {
		cloneOpts.ReferenceName = plumbing.ReferenceName(ref)
		if !strings.HasPrefix(ref, "refs/") {
			cloneOpts.ReferenceName = plumbing.NewBranchReferenceName(ref)
		}
	}

	_, err = git.PlainClone(tmpDir, false, cloneOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to clone git repository: %w", err)
	}

	cleanPath := filepath.Clean(filePath)
	fullPath := filepath.Join(tmpDir, cleanPath)

	// Ensure the resolved path is still within tmpDir using absolute paths
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve file path: %w", err)
	}
	absTmpDir, err := filepath.Abs(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve temp directory: %w", err)
	}

	if !strings.HasPrefix(absFullPath, absTmpDir+string(filepath.Separator)) {
		return nil, fmt.Errorf("invalid file path: outside repository bounds")
	}

	// #nosec G304 -- Path is validated to be within tmpDir bounds above
	return os.ReadFile(fullPath)
}

// parseGitURL parses a Git URL following Terraform/Checkov conventions
// Format: <git-url>//<file-path>?ref=<branch-or-tag>
// Examples:
//   - https://github.com/org/repo.git//config.yaml?ref=main
//   - git@github.com:org/repo.git//path/to/config.yaml?ref=v1.0.0
func parseGitURL(gitPath string) (repoURL, filePath, ref string, err error) {
	// Find the "//" separator
	idx := strings.Index(gitPath, "//")
	if idx == -1 {
		return "", "", "", fmt.Errorf("invalid git URL format, expected: <git-url>//<file-path>")
	}

	// Handle protocol prefixes more efficiently
	start := 0
	if strings.HasPrefix(gitPath, "https://") {
		start = 8 // len("https://")
		idx = strings.Index(gitPath[start:], "//")
		if idx == -1 {
			return "", "", "", fmt.Errorf("invalid git URL format, expected: <git-url>//<file-path>")
		}
		idx += start
	} else if strings.HasPrefix(gitPath, "http://") {
		start = 7 // len("http://")
		idx = strings.Index(gitPath[start:], "//")
		if idx == -1 {
			return "", "", "", fmt.Errorf("invalid git URL format, expected: <git-url>//<file-path>")
		}
		idx += start
	}

	repoURL = gitPath[:idx]
	remainder := gitPath[idx+2:] // +2 to skip the "//"

	// Extract ref parameter if present
	if qIdx := strings.IndexByte(remainder, '?'); qIdx != -1 {
		filePath = remainder[:qIdx]

		query, err := url.ParseQuery(remainder[qIdx+1:])
		if err != nil {
			return "", "", "", fmt.Errorf("failed to parse query parameters: %w", err)
		}
		ref = query.Get("ref")
	} else {
		filePath = remainder
	}

	return repoURL, filePath, ref, nil
}
