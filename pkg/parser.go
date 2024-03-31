package pkg

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"strings"
)

type parser struct {
	config       Config
	originURL    string
	downloadPath string
}

func NewParser(conf Config, originURL string, downloadPath string) *parser {
	return &parser{config: conf, originURL: originURL, downloadPath: downloadPath}
}
func (p *parser) Download() error {
	downloadURL, err := p.Parse()
	if err != nil {
		return err
	}
	out, err := os.Create(p.downloadPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}
	return nil
}

func (p *parser) getConfig() (domElement, error) {
	hostName := p.getDomainName()
	conf, ok := p.config.HostConf[hostName]
	if !ok {
		return domElement{}, fmt.Errorf("config not found for %s", hostName)
	}
	return conf, nil
}

func (p *parser) Parse() (string, error) {
	conf, err := p.getConfig()
	if err != nil {
		return "", err
	}
	resp, err := http.Get(p.originURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	var downloadURL string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if p.isEpub(n) {
			attributes := make(map[string]string)
			for _, a := range n.Attr {
				attributes[a.Key] = a.Val
			}
			var found = true
			for _, field := range conf.Fields {
				if className, ok := attributes[field.Key]; !ok || !strings.Contains(className, field.Contains) {
					found = false
					break
				}
			}
			if found {
				downloadURL = attributes["href"]
			}
		}
		for c := n.FirstChild; c != nil && len(downloadURL) == 0; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if len(downloadURL) == 0 {
		return "", fmt.Errorf("download URL for epub not found")
	}

	pathPrefix := ""
	if !strings.HasPrefix(downloadURL, "/") {
		pathPrefix = "/"
	}
	// Obtain host name from remoteURL
	parts := strings.Split(p.originURL, "/")
	return parts[0] + "//" + parts[2] + pathPrefix + downloadURL, nil
}

func (p *parser) getDomainName() string {
	parts := strings.Split(p.originURL, "/")
	return parts[2]
}

func (p *parser) isEpub(n *html.Node) bool {
	conf, err := p.getConfig()
	if err != nil {
		return false
	}
	return n != nil && n.Type == html.ElementNode && n.Data == "a" && n.FirstChild != nil && n.FirstChild.Type == html.TextNode && strings.EqualFold(strings.TrimSpace(n.FirstChild.Data), conf.Text)
}
