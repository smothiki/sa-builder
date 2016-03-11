package storage

import (
	"fmt"
	"os"

	"github.com/deis/sa-builder/pkg/gitreceive/git"
)

// SlugBuilderInfo contains all of the object storage related information needed to pass to a slug builder
type SlugBuilderInfo struct {
	pushKey string
	pushURL string
	tarKey  string
	tarURL  string
	slugKey string
	slugURL string
}

// NewSlugBuilderInfo creates and populates a new SlugBuilderInfo based on the given data
func NewSlugBuilderInfo(appName, slugName string, gitSha *git.SHA) *SlugBuilderInfo {
	s3Endpoint := "http://" + os.Getenv("DEIS_BUILDER_SERVICE_HOST") + ":3000"
	tarKey := fmt.Sprintf("home/%s/tar", slugName)
	// this is where workflow tells slugrunner to download the slug from, so we have to tell slugbuilder to upload it to here
	pushKey := fmt.Sprintf("home/%s:git-%s/push", appName, gitSha.Short())
	slugKey := fmt.Sprintf("home/%s:git-%s/slug", appName, gitSha.Short())

	return &SlugBuilderInfo{
		pushKey: pushKey,
		pushURL: fmt.Sprintf("%s/git/%s", s3Endpoint, pushKey),
		tarKey:  tarKey,
		tarURL:  fmt.Sprintf("%s/git/%s", s3Endpoint, tarKey),
		slugKey: slugKey,
		slugURL: fmt.Sprintf("%s/git/%s", s3Endpoint, slugKey),
	}
}

func (s SlugBuilderInfo) PushKey() string { return s.pushKey }
func (s SlugBuilderInfo) PushURL() string { return s.pushURL }
func (s SlugBuilderInfo) TarKey() string  { return s.tarKey }
func (s SlugBuilderInfo) TarURL() string  { return s.tarURL }
func (s SlugBuilderInfo) SlugURL() string { return s.slugURL }
