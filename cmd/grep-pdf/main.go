package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

type (
	Args struct {
		dir     string
		text    string
		json    bool
		onlyHit bool
		verbose bool
		debug   bool
	}

	SimpleResult struct {
		Path string `json:"path"`
		Text string `json:"text"`
	}

	DetailResult struct {
		Path string `json:"path"`
		Page int32  `json:"page"`
		Line int32  `json:"line"`
		Text string `json:"text"`
	}
)

var (
	args   Args
	appLog = log.New(os.Stdout, "", 0)
)

func init() {
	flag.StringVar(&args.dir, "dir", ".", "directory")
	flag.StringVar(&args.text, "text", "", "search text")
	flag.BoolVar(&args.json, "json", false, "output as JSON")
	flag.BoolVar(&args.onlyHit, "only-hit", true, "show ONLY HIT")
	flag.BoolVar(&args.verbose, "verbose", false, "verbose mode")
	flag.BoolVar(&args.debug, "debug", false, "debug mode")
}

func newSimpleResult(path string, text string) *SimpleResult {
	return &SimpleResult{path, text}
}

func newDetailResult(path string, page, line int32, text string) *DetailResult {
	return &DetailResult{path, page, line, text}
}

func main() {
	flag.Parse()

	if args.text == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if args.dir == "" {
		args.dir = "."
	}

	var (
		rootCtx  = context.Background()
		ctx, cxl = context.WithCancel(rootCtx)
	)
	defer cxl()

	if err := run(ctx); err != nil {
		panic(err)
	}
}

func abs(p string) string {
	v, _ := filepath.Abs(p)
	return v
}

func genErr(procName string, err error) error {
	return fmt.Errorf("%s failed: %w", procName, err)
}

func run(_ context.Context) error {

	rootDir := abs(args.dir)
	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.Contains(filepath.Base(path), "~$") {
			return nil
		}

		if !strings.HasSuffix(path, ".pdf") {
			return nil
		}

		absPath := abs(path)
		file, reader, err := pdf.Open(absPath)
		if err != nil {
			return genErr("os.Open()", err)
		}
		defer file.Close()

		if args.debug {
			appLog.Printf("Document Open: %s", absPath)
		}

		totalPage := reader.NumPage()
		for page := 1; page <= totalPage; page++ {

			p := reader.Page(page)
			if p.V.IsNull() {
				continue
			}

			rows, err := p.GetTextByRow()
			if err != nil {
				return genErr("p.GetPlainText()", err)
			}

			var (
				relPath, _ = filepath.Rel(rootDir, absPath)
				sb         strings.Builder
				found      bool
				count      = 1
			)
			for _, row := range rows {
				sb.Reset()

				for _, word := range row.Content {
					sb.WriteString(word.S)
				}

				text := sb.String()
				if strings.Contains(text, args.text) {
					found = true

					if args.verbose {
						err = outDetail(newDetailResult(relPath, int32(page), int32(count), text))
						if err != nil {
							return genErr("outDetail", err)
						}
					} else {
						err = outSimple(newSimpleResult(relPath, "HIT"))
						if err != nil {
							return genErr("outSimple(hit)", err)
						}

						return nil
					}
				}

				count++
			}

			if !found {
				if !args.onlyHit {
					err = outSimple(newSimpleResult(relPath, "NO HIT"))
					if err != nil {
						return genErr("outSimple(no hit)", err)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func outSimple(result *SimpleResult) error {
	var (
		message = fmt.Sprintf("%s: %s", result.Path, result.Text)
		err     error
	)

	if args.json {
		message, err = toJson(result)
		if err != nil {
			return genErr("toJson(result)", err)
		}
	}

	appLog.Println(message)

	return nil
}

func outDetail(result *DetailResult) error {
	var (
		message = fmt.Sprintf("%s (%3d,%3d): %q", result.Path, result.Page, result.Line, result.Text)
		err     error
	)

	if args.json {
		message, err = toJson(result)
		if err != nil {
			return genErr("toJson(result)", err)
		}
	}

	appLog.Println(message)

	return nil
}

func toJson(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", genErr("json.Marshal(result)", err)
	}

	return string(b), nil
}
