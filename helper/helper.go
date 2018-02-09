package helper

import (
	"html"
	"log"
)

func LinkTag(url string) string {
	return `<link type="text/css" rel="stylesheet" href="` + html.EscapeString(url) + `"></link>`
}

func ScriptTag(url string) string {
	return `<script type="text/javascript" src="` + html.EscapeString(url) + `"></script>`
}

func AssetTag(kind, url string) string {
	var buf string
	if kind == "css" {
		buf = LinkTag(url)
	} else if kind == "js" {
		buf = ScriptTag(url)
	} else {
		log.Println("go-webpack: unsupported asset kind: " + kind)
		buf = ""
	}
	return buf
}
