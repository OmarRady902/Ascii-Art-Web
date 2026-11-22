package handler

import (
	"ascii-art-web/banner"
	"ascii-art-web/render"
	"ascii-art-web/renderpage"
	"net/http"
	"strings"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	renderpage.RenderPage(w, renderpage.PageData{})
}

func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	bannerFile := r.FormValue("banner")

	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, "\\n", "\n")

	if strings.TrimSpace(text) == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	for _, ch := range text {
		if ch != '\n' && (ch < 32 || ch > 126) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	if bannerFile == "" || (bannerFile != "standard" && bannerFile != "shadow" && bannerFile != "thinkertoy") {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	content, err := banner.LoadBannerFile("banners/" + bannerFile + ".txt")
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if len(strings.Split(content, "\n")) < 855 {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	charMap := banner.BuildCharacterMap(content)

	var result strings.Builder

	result.Grow(len(text) * 64)

	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			result.WriteString("\n")
			continue
		}
		result.WriteString(render.RenderText(line, charMap))
	}
	

	renderpage.RenderPage(w, renderpage.PageData{
		Result: result.String(),
	})
}
