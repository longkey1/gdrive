/*
Copyright © 2026 longkey1

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/longkey1/gdrive/internal/gdrive"
	"github.com/spf13/cobra"
)

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Inspect files and folders in Google Drive",
}

// filesListCmd represents the files list command
var filesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the direct children of a folder",
	Long: `List the direct children of a folder.
The folder is given by --parent as either a folder ID or a Google Drive folder URL.
Listing is not recursive. Trashed items are excluded. Shared drive folders are supported.`,
	RunE: runFilesList,
}

// filesGetCmd represents the files get command
var filesGetCmd = &cobra.Command{
	Use:   "get <file-id-or-url>",
	Short: "Show metadata of a file or folder",
	Long: `Show metadata of a file or folder.
Accepts a Drive file ID, or a Google Drive / Docs / Sheets / Slides URL.`,
	Args: cobra.ExactArgs(1),
	RunE: runFilesGet,
}

func runFilesList(cmd *cobra.Command, args []string) error {
	parentInput, _ := cmd.Flags().GetString("parent")
	query, _ := cmd.Flags().GetString("query")
	maxResults, _ := cmd.Flags().GetInt64("max-results")
	formatStr, _ := cmd.Flags().GetString("format")

	format, err := gdrive.ParseOutputFormat(formatStr)
	if err != nil {
		return err
	}
	if maxResults < 0 {
		return fmt.Errorf("--max-results must be 0 or greater")
	}

	parent, err := gdrive.ParseFileInput(parentInput)
	if err != nil {
		return fmt.Errorf("invalid --parent: %w", err)
	}

	cfg := GetConfig()
	if err := cfg.Validate(); err != nil {
		return err
	}

	svc, err := gdrive.NewService(context.Background(), cfg)
	if err != nil {
		return err
	}

	files, err := gdrive.ListFiles(svc.Drive.Service, gdrive.ListOptions{
		Parent:     parent,
		Query:      query,
		MaxResults: maxResults,
	})
	if err != nil {
		return err
	}

	return gdrive.FormatFileList(os.Stdout, files, format)
}

func runFilesGet(cmd *cobra.Command, args []string) error {
	formatStr, _ := cmd.Flags().GetString("format")

	format, err := gdrive.ParseOutputFormat(formatStr)
	if err != nil {
		return err
	}

	fileID, err := gdrive.ParseFileInput(args[0])
	if err != nil {
		return err
	}

	cfg := GetConfig()
	if err := cfg.Validate(); err != nil {
		return err
	}

	svc, err := gdrive.NewService(context.Background(), cfg)
	if err != nil {
		return err
	}

	file, err := gdrive.GetFile(svc.Drive.Service, fileID)
	if err != nil {
		return err
	}

	return gdrive.FormatFile(os.Stdout, file, format)
}

func init() {
	rootCmd.AddCommand(filesCmd)
	filesCmd.AddCommand(filesListCmd)
	filesCmd.AddCommand(filesGetCmd)

	filesListCmd.Flags().String("parent", "", "Folder ID or URL whose children to list (required)")
	_ = filesListCmd.MarkFlagRequired("parent")
	filesListCmd.Flags().StringP("query", "q", "", "Filter by name substring")
	filesListCmd.Flags().Int64P("max-results", "n", 100, "Maximum number of results (0 = unlimited)")
	filesListCmd.Flags().String("format", "text", "Output format (text or json)")

	filesGetCmd.Flags().String("format", "text", "Output format (text or json)")
}
