package render

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type TemplateRenderer struct {
	config       TemplateRendererConfig
	baseTemplate *template.Template
}

type TemplateRendererConfig struct {
	TemplateType           string
	TemplatesPath          string
	ComponentTemplatesPath string
	PageTemplatesPath      string
	OutputPath             string
	BaseTemplateName       string
}

type templateData struct {
	Page        any
	LastUpdated string
}

var funcMap = template.FuncMap{
	"join": strings.Join,
	"title": func(s string) string {
		if s == "" {
			return ""
		}
		r := []rune(s)
		r[0] = unicode.ToUpper(r[0])
		return string(r)
	},
}

func New(config TemplateRendererConfig) (*TemplateRenderer, error) {
	baseTemplate, err := template.New("").Funcs(funcMap).ParseFiles(filepath.Join(config.TemplatesPath, config.BaseTemplateName+"."+config.TemplateType))
	if err != nil {
		return nil, err
	}

	_, err = baseTemplate.ParseGlob(filepath.Join(config.TemplatesPath, config.ComponentTemplatesPath, "*."+config.TemplateType))
	if err != nil {
		return nil, err
	}

	return &TemplateRenderer{config: config, baseTemplate: baseTemplate}, nil
}

func (r *TemplateRenderer) RenderToFile(templateName string, data any, outputName string) error {
	pageTemplate, err := r.baseTemplate.Clone()
	if err != nil {
		return err
	}

	templatePath := filepath.Join(
		r.config.TemplatesPath,
		r.config.PageTemplatesPath,
		templateName+"."+r.config.TemplateType,
	)
	_, err = pageTemplate.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	outputFileName := outputName + ".html"
	outputPath := filepath.Join(r.config.OutputPath, outputFileName)

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	templateData := templateData{
		Page:        data,
		LastUpdated: time.Now().Format("2006-01-02"),
	}

	return pageTemplate.ExecuteTemplate(file, r.config.BaseTemplateName, templateData)
}

func (r *TemplateRenderer) CopyStaticFiles(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
