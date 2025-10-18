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
	
	if isGitURL(remotePath) {
		return fetchFromGit(remotePath)
	}
	
	if isHTTPURL(remotePath) {
		return fetchFromHTTP(remotePath)
	}
	
	return nil, fmt.Errorf("unsupported remote path format")
}

func hasValidExtension(path string) bool {
	// Remove query parameters
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml" || ext == ".json"
}

func isHTTPURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func isGitURL(path string) bool {
	// Git SSH format
	if strings.HasPrefix(path, "git@") {
		return true
	}
	// Git HTTPS with .git and // separator (Terraform/Checkov convention)
	if (strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://")) && 
	   strings.Contains(path, ".git//") {
		return true
	}
	return false
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
	
	return os.ReadFile(filepath.Join(tmpDir, filepath.Clean(filePath)))
}

// parseGitURL parses a Git URL following Terraform/Checkov conventions
// Format: <git-url>//<file-path>?ref=<branch-or-tag>
// Examples:
//   - https://github.com/org/repo.git//config.yaml?ref=main
//   - git@github.com:org/repo.git//path/to/config.yaml?ref=v1.0.0
func parseGitURL(gitPath string) (repoURL, filePath, ref string, err error) {
	// Find the position after the protocol (if present)
	start := 0
	if strings.HasPrefix(gitPath, "https://") {
		start = len("https://")
	} else if strings.HasPrefix(gitPath, "http://") {
		start = len("http://")
	}
	
	// Find the "//" separator after the protocol
	idx := strings.Index(gitPath[start:], "//")
	if idx == -1 {
		return "", "", "", fmt.Errorf("invalid git URL format, expected: <git-url>//<file-path>")
	}
	
	// Split at the actual separator position
	repoURL = gitPath[:start+idx]
	remainder := gitPath[start+idx+2:] // +2 to skip the "//"
	
	// Extract ref parameter
	if strings.Contains(remainder, "?") {
		pathAndQuery := strings.SplitN(remainder, "?", 2)
		filePath = pathAndQuery[0]
		
		query, err := url.ParseQuery(pathAndQuery[1])
		if err != nil {
			return "", "", "", fmt.Errorf("failed to parse query parameters: %w", err)
		}
		ref = query.Get("ref")
	} else {
		filePath = remainder
	}
	
	return repoURL, filePath, ref, nil
}
