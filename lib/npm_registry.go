package lib

import (
	"encoding/json"
	"fmt"
	"github.com/sethmcl/gofrosty/vendor/nsemver"
	"net/http"
	"regexp"
	"strings"
)

// NpmRegistryClient interacts with a remote npm registry server
type NpmRegistryClient struct {
	RootURL   string
	AuthToken string
}

type npmInfo struct {
	Versions map[string]interface{} `json:"versions"`
}

// NewNpmRegistryClient create NpmRegistryClient
func NewNpmRegistryClient(url string, token string) *NpmRegistryClient {
	return &NpmRegistryClient{
		RootURL:   url,
		AuthToken: token,
	}
}

// Get send GET request
func (n *NpmRegistryClient) Get(url string) (*http.Response, error) {
	res, err := n.Request("GET", url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 && res.StatusCode <= 500 {
		return nil, fmt.Errorf("HTTP %d response", res.StatusCode)
	}

	return res, nil
}

// Request send request
func (n *NpmRegistryClient) Request(method string, url string) (*http.Response, error) {
	GetContext().Info("%s %s", method, url)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if n.IsRegistryURL(url) && n.AuthToken != "" {
		ctx := GetContext()
		ctx.Debug("Attaching NPM token (%s) to request", ctx.NpmAuthToken)
		auth := fmt.Sprintf("Bearer %s", n.AuthToken)
		req.Header.Add("authorization", auth)
	}

	return client.Do(req)
}

// ListModuleVersions return slice of available versions for a given module
func (n *NpmRegistryClient) ListModuleVersions(name string) ([]string, error) {
	encodeName := strings.Replace(name, "/", "%2f", -1)
	url := fmt.Sprintf("%s/%s", n.RootURL, encodeName)
	res, err := n.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("Unable to fetch module versions for %s", name)
	}

	info := &npmInfo{}
	err = json.NewDecoder(res.Body).Decode(info)
	if err != nil {
		return nil, err
	}

	var versionList []string
	// GetContext().Info("%s versions: %s", name, info.Versions)
	for version := range info.Versions {
		versionList = append(versionList, version)
	}

	return versionList, nil

}

// GetAPIURL returns API url
func (n *NpmRegistryClient) GetAPIURL(cmd string) string {
	return fmt.Sprintf("%s/-/%s", n.RootURL, cmd)
}

// GenTarURL return url to tar file hosted by registry
func (n *NpmRegistryClient) GenTarURL(name string, version string) string {
	// Scoped package
	// Example: @foo/biz => https://registry.npmjs.org/@foo/biz/-/biz-1.0.0.tgz
	if strings.Contains(name, "/") {
		parts := strings.Split(name, "/")
		return fmt.Sprintf("%s/%s/-/%s-%s.tgz", n.RootURL, name, parts[1], version)
	}

	// Global package
	// Example: bar => https://registry.npmjs.org/bar/-/bar-1.0.0.tgz
	return fmt.Sprintf("%s/%s/-/%s-%s.tgz", n.RootURL, name, name, version)
}

// GetTarURL returns url to tar file hosted by registry. Will resolve version if specified as range.
func (n *NpmRegistryClient) GetTarURL(name string, version string) (string, error) {
	ctx := GetContext()

	candidates, err := ctx.NpmRegistry.ListModuleVersions(name)
	if err != nil {
		ctx.DumpStack()
		return "", fmt.Errorf("Cannot list module versions for %s :: %s", name, err.Error())
	}

	ver, err := nsemver.MatchLatest(version, candidates)
	if err != nil {
		ctx.DumpStack()
		return "", err
	}

	return n.GenTarURL(name, ver), nil
}

// IsRegistryURL returns true if this URL points to the npm registry
func (n *NpmRegistryClient) IsRegistryURL(url string) bool {
	match, _ := regexp.MatchString(fmt.Sprintf("^%s", n.RootURL), url)
	return match
}
