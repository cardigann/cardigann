package indexer

import (
	"net/http"
	"testing"
	"time"

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
    test:
      path: /profile.php

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

const exampleDefinitionWithMultiRow = `
---
  site: example
  links:
    - http://www.example.org

  caps:
    categories:
      2: Audio
      3: Other

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
    test:
      path: /profile.php

  search:
    path: torrents.php
    inputs:
      $raw: "search={{ .Query.Keywords }}&cat=0"
    rows:
      selector: table.results tbody tr:not(.dateheader)
      after: 1
      dateheaders:
        selector: .dateheader
        filters:
          - name: regexp
            args: "^Added on (.+?)$"
          - name: dateparse
            args: Monday, Jan 02, 2006
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
`

const exampleSearchPageWithDateHeadersAndMultiRow = `
<html>
<body>
  <table class="results">
    <tbody>
      <tr class="dateheader">
        <td colspan="5">Added on Thursday, Aug 25, 2016</td>
      </tr>
      <tr>
        <td rowspan="2"><a href="category.php?id=2">Sound</a></td>
      </tr>
      <tr>
        <td><a href="details.php?1_archive">Llama llama 1</a></td>
        <td><a href="/download/1_archive.torrent">Download</a></td>
        <td>4GB</td>
        <td>2006-01-02 15:04:05</td>
      </tr>
      <tr class="dateheader">
        <td colspan="5">Added on Thursday, Aug 20, 2016</td>
      </tr>
      <tr>
        <td rowspan="2"><a href="category.php?id=2">Sound</a></td>
      </tr>
      <tr>
        <td><a href="details.php?2_archive">Llama llama 2</a></td>
        <td><a href="/download/2_archive.torrent">Download</a></td>
        <td>4GB</td>
        <td>2006-01-02 15:04:05</td>
      </tr>
      <tr>
        <td rowspan="2"><a href="category.php?id=3">Other</a></td>
      </tr>
      <tr>
        <td><a href="details.php?3_archive">Llama llama 3</a></td>
        <td><a href="/download/3_archive.torrent">Download</a></td>
        <td>4GB</td>
        <td>2006-01-02 15:04:05</td>
      </tr>
    </tbody>
  </table>
</body>
</html>
`

// registerResponder wraps httpmock.RegisterResponder and fixes bug with Request assignment
func registerResponder(method, url string, f func(req *http.Request) (*http.Response, error)) {
	httpmock.RegisterResponder(method, url, func(innerreq *http.Request) (*http.Response, error) {
		resp, err := f(innerreq)
		if err != nil {
			return nil, err
		}
		resp.Request = innerreq
		return resp, nil
	})
}

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

	var loggedIn bool

	registerResponder("GET", "https://example.org/", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, "Ok"), nil
	})

	registerResponder("GET", "https://example.org/profile.php", func(req *http.Request) (*http.Response, error) {
		if !loggedIn {
			resp := httpmock.NewStringResponse(http.StatusTemporaryRedirect, "")
			resp.Header.Set("Location", "/login.php")
			resp.Request = req
			return resp, nil
		}
		return httpmock.NewStringResponse(http.StatusOK, "Ok"), nil
	})

	registerResponder("GET", "https://example.org/login.php", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, exampleLoginPage), nil
	})

	registerResponder("POST", "https://example.org/login.php", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, exampleLoginErrorPage), nil
	})

	r := NewRunner(def, conf)
	err = r.login()

	if err == nil || err.Error() != "Login failed" {
		t.Fatalf("Expected 'Login failed', got %#v", err)
	}

	registerResponder("POST", "https://example.org/login.php", func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(http.StatusOK, "Success!")

		if pwd := req.FormValue("llamas_password"); pwd != "mypassword" {
			t.Fatalf("Incorrect password %q was provided", pwd)
		}

		return resp, nil
	})

	loggedIn = true

	err = r.login()
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

	var loggedIn bool

	registerResponder("GET", "https://example.org/", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, ""), nil
	})

	registerResponder("GET", "https://example.org/profile.php", func(req *http.Request) (*http.Response, error) {
		if !loggedIn {
			resp := httpmock.NewStringResponse(http.StatusTemporaryRedirect, "")
			resp.Header.Set("Location", "/login.php")
			return resp, nil
		}
		return httpmock.NewStringResponse(http.StatusOK, ""), nil
	})

	registerResponder("GET", "https://example.org/login.php", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, exampleLoginPage), nil
	})

	registerResponder("POST", "https://example.org/login.php", func(req *http.Request) (*http.Response, error) {
		loggedIn = true
		return httpmock.NewStringResponse(http.StatusOK, "Success"), nil
	})

	registerResponder("GET", "https://example.org/torrents.php", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, exampleSearchPage), nil
	})

	results, err := r.Search(torznab.Query{"t": "tv-search", "q": "llamas", "cat": []int{torznab.CategoryAudio_Foreign.ID}})
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

func TestIndexerDefinitionRunner_SearchWithMultiRow(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	def, err := ParseDefinition([]byte(exampleDefinitionWithMultiRow))
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

	var loggedIn bool

	registerResponder("GET", "https://example.org/", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, ""), nil
	})

	registerResponder("GET", "https://example.org/profile.php", func(req *http.Request) (*http.Response, error) {
		if !loggedIn {
			resp := httpmock.NewStringResponse(http.StatusOK, "")
			resp.Header.Set("Refresh", "1; /login.php")
			return resp, nil
		}
		return httpmock.NewStringResponse(http.StatusOK, ""), nil
	})

	registerResponder("GET", "https://example.org/login.php", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, exampleLoginPage), nil
	})

	registerResponder("POST", "https://example.org/login.php", func(req *http.Request) (*http.Response, error) {
		loggedIn = true
		return httpmock.NewStringResponse(http.StatusOK, "Success"), nil
	})

	registerResponder("GET", "https://example.org/torrents.php", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, exampleSearchPageWithDateHeadersAndMultiRow), nil
	})

	results, err := r.Search(torznab.Query{
		"t":   "tv-search",
		"q":   "llamas",
		"cat": []int{torznab.CategoryAudio.ID},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 result, got %d", len(results))
	}

	if results[1].Title != "Llama llama 2" {
		t.Fatalf("Expected row 2 to have title of %q, got %q",
			"Llama llama 2",
			results[1].Title)
	}

	expectedDate := time.Date(2016, time.August, 20, 0, 0, 0, 0, time.UTC)
	if !results[1].PublishDate.Equal(expectedDate) {
		t.Fatalf("Expected row 2 to have publish date of %q, got %q",
			expectedDate.String(), results[1].PublishDate)
	}
}
