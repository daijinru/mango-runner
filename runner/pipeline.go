package runner

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/daijinru/mango-runner/utils"
)

type Pipeline struct {
	Tag       string
	FilePath  string
	Filename  string
	Directory string

	File *os.File
}

// PipelineArgs Tag is the project name, Path is typically the cwd provided by the Workspace.
type PipelineArgs struct {
	Tag  string
	Path string
}

func NewPipeline(args *PipelineArgs) (*Pipeline, error) {
	dir, err := url.JoinPath(args.Path, "/pipelines/")
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	tag := args.Tag
	if tag == "" {
		tag = utils.GenerateUUIDFileName()
	}

	now := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s", tag, now)
	filePath := filepath.Join(dir, filename+".txt")

	return &Pipeline{
		Tag:       tag,
		FilePath:  filePath,
		Filename:  filename,
		Directory: dir,
	}, nil
}

func (pip *Pipeline) IfLogFileExists() (*os.File, error) {
	_, err := os.Stat(pip.FilePath)
	if err != nil && os.IsNotExist(err) {
		_, err := os.Create(pip.FilePath)
		if err != nil {
			return nil, err
		}
	}

	if pip.File == nil {
		pip.File, err = os.OpenFile(pip.FilePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	}

	return pip.File, nil
}

func (pip *Pipeline) WriteInfo(text string) error {
	file, err := pip.IfLogFileExists()
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("[%s] %s", utils.TimeNow(), text))

	return err
}

func (pip *Pipeline) WriteError(cause error, command string) error {
	file, err := pip.IfLogFileExists()
	if err != nil {
		return err
	}
	errorMsg := fmt.Sprintf("[Failed] [Tag:%s] [Filename:%s] [Error:%s] [script:%s]", pip.Tag, pip.Filename, cause.Error(), command)
	_, err = file.WriteString(errorMsg)

	return err
}

func (pip *Pipeline) CloseFile() error {
	if pip.File != nil {
		err := pip.File.Close()
		return err
	}
	return nil
}

func (pip *Pipeline) ReadDir() ([]string, error) {
	files, err := os.ReadDir(pip.Directory)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filenames = append(filenames, file.Name())
	}

	filteredFilenames := make([]string, 0)
	for _, filename := range filenames {
		if strings.Contains(filename, pip.Tag) {
			filteredFilenames = append(filteredFilenames, filename)
		}
	}

	sort.Slice(filteredFilenames, func(i, j int) bool {
		file1 := filepath.Join(pip.Directory, filenames[i])
		file2 := filepath.Join(pip.Directory, filenames[j])
		info1, _ := os.Stat(file1)
		info2, _ := os.Stat(file2)
		return info1.ModTime().Before(info2.ModTime())
	})
	return filteredFilenames, nil
}

func (pip *Pipeline) Read(tag string, timestamp string) string {
	filename := fmt.Sprintf("%s_%s", tag, timestamp)
	filePath := filepath.Join(pip.Directory, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(content)
}

func (pip *Pipeline) ReadFile(filename string) string {
	filePath := filepath.Join(pip.Directory, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(content)
}

func (pip *Pipeline) Callback(urlStr string, newQueries ...string) error {
	decodedURL, err := url.QueryUnescape(urlStr)
	if err != nil {
		return err
	}

	parsedURL, err := url.Parse(decodedURL)
	if err != nil {
		return err
	}

	queries := parsedURL.Query()
	for i := 0; i < len(newQueries); i += 2 {
		key := newQueries[i]
		value := newQueries[i+1]
		queries.Add(key, value)
	}

	parsedURL.RawQuery = queries.Encode()

	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodyStr := string(body)
	pip.WriteInfo(bodyStr)
	return nil
}
