// Copyright © 2018 Yoshiki Shibata. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const GitHubAPIURL = "https://api.github.com"

type Asset struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type AssetsInfo struct {
	Assets []Asset `json:"assets"`
}

func retrieveAssets(owner, repo, tag, token string) error {
	assetsInfo, err := getAssetsInfo(owner, repo, tag, token)
	if err != nil {
		return err
	}

	for _, asset := range assetsInfo.Assets {
		if err = retrieveAsset(asset.Name, asset.URL, token); err != nil {
			return err
		}
	}
	return nil
}

func getAssetsInfo(owner, repo, tag, token string) (*AssetsInfo, error) {
	req, err := newRequest("GET", tagURL(owner, repo, tag), token)
	if err != nil {
		return nil, fmt.Errorf("New Request failed: %v\n", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3.text-match+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Do failed: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status = %v\n", resp.Status)
	}

	var assetsInfo AssetsInfo
	if err := json.NewDecoder(resp.Body).Decode(&assetsInfo); err != nil {
		return nil, fmt.Errorf("Decode failed: %v\n", err)
	}
	return &assetsInfo, nil
}

func retrieveAsset(name, url, token string) error {
	fmt.Printf("Retrieving %s ... ", name)

	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("Create(%s) failed: %v\n", name, err)
	}
	defer f.Close()

	req, err := newRequest("GET", url, token)
	if err != nil {
		return fmt.Errorf("New Request failed: %v\n", err)
	}

	req.Header.Set("Accept", "application/octet-stream")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Do failed: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status = %v\n", resp.Status)
	}

	bytes, err := io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("Failed: %v\n", err)
	}

	fmt.Printf("%d bytes\n", bytes)
	return nil
}

func tagURL(owner, repo, tag string) string {
	// GET /repos/:owner/:repo/releases/tags/:tag
	return GitHubAPIURL +
		"/repos/" + owner + "/" + repo + "/releases/tags/" + tag
}

func newRequest(cmd, url string, token string) (*http.Request, error) {
	req, err := http.NewRequest(cmd, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Token "+token)
	return req, nil
}
