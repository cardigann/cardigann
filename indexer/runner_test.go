package indexer

import (
	"net/http"
	"testing"

	"github.com/cardigann/cardigann/config"
	"github.com/cardigann/cardigann/torznab"
	"github.com/jarcoal/httpmock"
)

const exampleDefinition2 = `
---
  site: example
  links:
    - http://www.example.org

  caps:
    categories:
      2:  Audio

    modes:
      search: q

  login:
    path: /login.php
    form: form
    inputs:
      username: "{{ .Config.username }}"
      llamas_password: "{{ .Config.password }}"
    error:
      selector: .loginerror a

  search:
    path: torrents.php
    inputs:
      $raw: "search={{ .Query.Keywords }}&cat=0"
    rows:
      selector: table.results tbody tr
    fields:
      category:
        selector: td:nth-child(1) a
        attribute: href
        filters:
          - name: querystring
            args: id
      title:
        selector: td:nth-child(2) a
      details:
        selector: td:nth-child(2) a
        attribute: href
      download:
        selector: td:nth-child(3) a
        attribute: href
      size:
        selector: td:nth-child(4)
      date:
        selector: td:nth-child(5)
        filters:
          - name: dateparse
            args: 2006-01-02 15:04:05
      seeders:
        selector: td:nth-child(6)
        filters:
          - name: regexp
            args: "^(\\d+) seeders"
      leechers:
        selector: td:nth-child(7)
        filters:
          - name: regexp
            args: "^(\\d+) leechers"
`

const exampleLoginPage = `
<html>
<body>
  <form method="post">
    <input type="text" name="username"></input>
    <input type="text" name="llamas_password"></input>
    <input type="submit" value="submit"></input>
  </form>
</body>
</html>
`

const exampleLoginErrorPage = `
<html>
<body>
  <div class="loginerror">
    <a href="">Login <strong>failed</strong></a>
    <strong>Forgotten password</strong
  </div>
</body>
</html>
`

const exampleSearchPage = `
<html>
<body>
  <table class="results">
    <tbody>
      <tr>
        <td><a href="category.php?id=2">Sound</a></td>
        <td><a href="details.php?mma_llama_309960_archive">Llama llama</a></td>
        <td><a href="/download/mma_llama_309960/mma_llama_309960_archive.torrent">Download</a></td>
        <td>4GB</td>
        <td>2006-01-02 15:04:05</td>
        <td>12 seeders</td>
        <td>100 leechers</td>
      </tr>
    </tbody>
  </table>
</body>
</html>
`

func TestIndexerDefinitionRunner_Login(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	def, err := ParseDefinition([]byte(exampleDefinition2))
	if err != nil {
		t.Fatal(err)
	}

	conf := &config.ArrayConfig{
		"example": map[string]string{
			"username": "myusername",
			"password": "mypassword",
			"url":      "https://example.org/",
		},
	}

	httpmock.RegisterResponder("GET", "https://example.org/login.php", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(http.StatusOK, exampleLoginPage)
		resp.Request = req
		return resp, nil
	})

	httpmock.RegisterResponder("POST", "https://example.org/login.php", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(http.StatusOK, exampleLoginErrorPage)
		resp.Request = req
		return resp, nil
	})

	r := NewRunner(def, conf)
	err = r.Login()

	if err == nil || err.Error() != "Login failed" {
		t.Fatalf("Expected 'Login failed', got %#v", err)
	}

	httpmock.RegisterResponder("POST", "https://example.org/login.php", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(http.StatusOK, "Success!")
		resp.Request = req

		if pwd := req.FormValue("llamas_password"); pwd != "mypassword" {
			t.Fatalf("Incorrect password %q was provided", pwd)
		}

		return resp, nil
	})

	err = r.Login()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIndexerDefinitionRunner_Search(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	def, err := ParseDefinition([]byte(exampleDefinition2))
	if err != nil {
		t.Fatal(err)
	}

	conf := &config.ArrayConfig{
		"example": map[string]string{
			"username": "myusername",
			"password": "mypassword",
			"url":      "https://example.org/",
		},
	}

	r := NewRunner(def, conf)

	httpmock.RegisterResponder("GET", "https://example.org/torrents.php", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(http.StatusOK, exampleSearchPage)
		resp.Request = req
		return resp, nil
	})

	results, err := r.Search(torznab.Query{"t": "tv-search", "q": "llamas", "cat": []int{torznab.CategoryAudio.ID}})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Link != "https://example.org/download/mma_llama_309960/mma_llama_309960_archive.torrent" {
		t.Fatal("Incorrect download link")
	}

	if results[0].Seeders != 12 {
		t.Fatal("Incorrect seeders count")
	}

	if results[0].Peers != 112 {
		t.Fatal("Incorrect peers count")
	}
}
