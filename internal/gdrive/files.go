package gdrive

import (
	"fmt"
	"strings"

	"google.golang.org/api/drive/v3"
)

// FileInfo represents a file or folder in Google Drive
type FileInfo struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	MimeType     string      `json:"mimeType"`
	WebViewLink  string      `json:"webViewLink"`
	ModifiedTime string      `json:"modifiedTime"`
	Size         int64       `json:"size"`
	Parents      []string    `json:"parents,omitempty"`
	Owners       []OwnerInfo `json:"owners,omitempty"`
}

// OwnerInfo represents an owner of a Drive file
type OwnerInfo struct {
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
}

// ListOptions controls which files ListFiles returns
type ListOptions struct {
	// Parent is the folder ID whose direct children are listed. Required.
	Parent string
	// Query filters by name substring (Drive "name contains").
	Query string
	// MaxResults caps the number of files returned; 0 means unlimited.
	MaxResults int64
}

// listPageSize is the Drive API maximum page size for files.list.
const listPageSize = 1000

// BuildListQuery builds the Drive API "q" expression for the given options.
// Trashed files are always excluded.
func BuildListQuery(opts ListOptions) string {
	terms := []string{
		fmt.Sprintf("'%s' in parents", escapeQueryValue(opts.Parent)),
		"trashed = false",
	}
	if opts.Query != "" {
		terms = append(terms, fmt.Sprintf("name contains '%s'", escapeQueryValue(opts.Query)))
	}
	return strings.Join(terms, " and ")
}

// escapeQueryValue escapes a value for use inside single quotes in a Drive
// API query expression.
func escapeQueryValue(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	return strings.ReplaceAll(s, `'`, `\'`)
}

// ListFiles returns the direct children of the folder given in opts.Parent,
// following nextPageToken until MaxResults is reached or pages are exhausted.
// Shared drive contents are included.
func ListFiles(svc *drive.Service, opts ListOptions) ([]FileInfo, error) {
	if opts.Parent == "" {
		return nil, fmt.Errorf("parent folder ID is required")
	}

	result := []FileInfo{}
	pageToken := ""
	for {
		pageSize := int64(listPageSize)
		if opts.MaxResults > 0 {
			remaining := opts.MaxResults - int64(len(result))
			if remaining <= 0 {
				break
			}
			if remaining < pageSize {
				pageSize = remaining
			}
		}

		call := svc.Files.List().
			Q(BuildListQuery(opts)).
			Fields("nextPageToken, files(id, name, mimeType, webViewLink, modifiedTime, size)").
			OrderBy("folder,name").
			SupportsAllDrives(true).
			IncludeItemsFromAllDrives(true).
			PageSize(pageSize)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		resp, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("unable to list files: %v", err)
		}

		for _, f := range resp.Files {
			result = append(result, fromDriveFile(f))
		}

		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	return result, nil
}

// GetFile returns metadata for a single file or folder. Shared drive items
// are supported.
func GetFile(svc *drive.Service, fileID string) (*FileInfo, error) {
	f, err := svc.Files.Get(fileID).
		Fields("id, name, mimeType, webViewLink, parents, modifiedTime, size, owners(emailAddress, displayName)").
		SupportsAllDrives(true).
		Do()
	if err != nil {
		return nil, fmt.Errorf("unable to get file: %v", err)
	}

	info := fromDriveFile(f)
	return &info, nil
}

func fromDriveFile(f *drive.File) FileInfo {
	info := FileInfo{
		ID:           f.Id,
		Name:         f.Name,
		MimeType:     f.MimeType,
		WebViewLink:  f.WebViewLink,
		ModifiedTime: f.ModifiedTime,
		Size:         f.Size,
		Parents:      f.Parents,
	}
	for _, o := range f.Owners {
		if o == nil {
			continue
		}
		info.Owners = append(info.Owners, OwnerInfo{
			EmailAddress: o.EmailAddress,
			DisplayName:  o.DisplayName,
		})
	}
	return info
}
