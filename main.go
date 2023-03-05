package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/saml"
	log "github.com/sirupsen/logrus"
)

var (
	SERVER   string
	USERNAME string
	PASSWORD string
	TERM     string
)

func init() {
	os.Mkdir("Download", os.ModePerm)

	flag.StringVar(&SERVER, "s", "", "Server Address")
	flag.StringVar(&USERNAME, "u", "", "Username")
	flag.StringVar(&PASSWORD, "p", "", "Password")
	flag.StringVar(&TERM, "t", "", "Search Term")
	flag.Parse()
}

func NewClient() *api.SP {
	var sp *api.SP

	authCnfg := &strategy.AuthCnfg{
		SiteURL:  SERVER,
		Username: USERNAME,
		Password: PASSWORD,
	}

	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp = api.NewSP(client)

	return sp
}

func Download(url string) error {
	authCnfg := &strategy.AuthCnfg{
		SiteURL:  SERVER,
		Username: USERNAME,
		Password: PASSWORD,
	}

	client := &gosip.SPClient{AuthCnfg: authCnfg}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating GET request - err: %v", err)
	}

	rsp, err := client.Execute(req)
	if err != nil {
		return fmt.Errorf("error sending GET request - err: %v", err)
	}
	defer rsp.Body.Close()

	path := strings.Split(url, "/")
	out, err := os.Create(fmt.Sprintf("Download/%s", path[len(path)-1]))
	if err != nil {
		return fmt.Errorf("error creating output directory - err: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, rsp.Body)
	if err != nil {
		return fmt.Errorf("error saving file - err: %v", err)
	}

	log.Printf("Download completed successfully - URL: %s", url)

	return nil
}

func main() {
	sp := NewClient()
	if sp == nil {
		log.Fatalln("error performing authentication")
	}

	log.Println("Performing search via API")
	log.Printf("Search term: %s", TERM)
	res, err := sp.Search().PostQuery(&api.SearchQuery{
		QueryText: TERM,
		RowLimit:  1000,
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, r := range res.Results() {
		log.Printf("Document found: %s", r["OriginalPath"])
		Download(r["OriginalPath"])
	}
}
