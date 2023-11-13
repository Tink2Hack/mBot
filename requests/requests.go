package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Tink2Hack/mBot/data"
	"github.com/Tink2Hack/mBot/env"
)

var Token = ""
var Urls = []string{
	// urls[0] = all unregistered targets
	"https://platform.synack.com/api/targets?filter%5Bprimary%5D=unregistered&filter%5Bsecondary%5D=all&filter%5Bcategory%5D=all&filter%5Bindustry%5D=all&sorting%5Bfield%5D=dateUpdated&sorting%5Bdirection%5D=desc",
	// urls[1] = available missions sorted by price (v2)
	"https://platform.synack.com/api/tasks/v2/tasks?perPage=20&viewed=true&page=1&status=PUBLISHED&sort=AMOUNT&sortDir=desc",
	// urls[2] = QR window
	"https://platform.synack.com/api/targets?filter%5Bprimary%5D=all&filter%5Bsecondary%5D%5B%5D=a&filter%5Bsecondary%5D%5B%5D=l&filter%5Bsecondary%5D%5B%5D=l&filter%5Bsecondary%5D%5B%5D=quality_period&filter%5Bcategory%5D=all&filter%5Bindustry%5D=all&sorting%5Bfield%5D=dateUpdated&sorting%5Bdirection%5D=desc",
	// urls[3] = claimed missions
	"https://platform.synack.com/api/tasks/v2/tasks?perPage=20&viewed=true&page=1&status=CLAIMED",
	// urls[4] = beginning of URL to edit missions
	"https://platform.synack.com/api/tasks/v2/tasks/",
	// urls[5] = authenticate URL
	"https://login.synack.com/api/authenticate",
	// urls[6] = claimed Amount
	"https://platform.synack.com/api/tasks/v2/researcher/claimed_amount",
	// urls[7] = connect to target
	"https://platform.synack.com/api/launchpoint",
}

func SetHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "not a bot (definitely not made by github.com/un4gi)")
	req.Header.Set("Authorization", "Bearer "+Token)
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Referer", "https://platform.synack.com/tasks/user/available")
	req.Header.Set("X-CSRF-Token", "xxxx")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
}

func DoGetRequest(target string) (int, io.ReadCloser) {
	ctx, cancel := context.WithTimeout(context.Background(), 180000*time.Millisecond)
	defer cancel()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("GET", target, nil)
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
		return 0, nil
	}
	SetHeaders(req)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
		return 0, nil
	}
	return resp.StatusCode, resp.Body
}

func DoPostRequest(target string, jsonStr []byte) (int, io.ReadCloser) {
	ctx, cancel := context.WithTimeout(context.Background(), 180000*time.Millisecond)
	defer cancel()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("POST", target, bytes.NewBuffer(jsonStr))
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
	}

	SetHeaders(req)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
	}
	return resp.StatusCode, resp.Body
}

func SetLoginHeaders(req *http.Request, token string, cookie string) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://login.synack.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", token)
	req.Header.Set("Origin", "https://login.synack.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", cookie)
}

func DoLoginGetRequest(target string) (int, io.ReadCloser, http.Header) {
	ctx, cancel := context.WithTimeout(context.Background(), 180000*time.Millisecond)
	defer cancel()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("GET", target, nil)
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
		return 0, nil, nil
	}
	SetHeaders(req)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
		return 0, nil, nil
	}
	return resp.StatusCode, resp.Body, resp.Header
}

func DoLoginPostRequest(target string, jsonStr []byte, token string, cookie string) (int, io.ReadCloser) {
	ctx, cancel := context.WithTimeout(context.Background(), 180000*time.Millisecond)
	defer cancel()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("POST", target, bytes.NewBuffer(jsonStr))
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
	}

	SetLoginHeaders(req, token, cookie)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
	}
	return resp.StatusCode, resp.Body
}

func SetGrantTokenHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://login.synack.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", "xxxx")
	req.Header.Set("Origin", "https://login.synack.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
}

func DoGrantTokenRequest(target string) (int, io.ReadCloser) {
	ctx, cancel := context.WithTimeout(context.Background(), 180000*time.Millisecond)
	defer cancel()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("GET", target, nil)
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
		return 0, nil
	}
	SetGrantTokenHeaders(req)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil && err != context.Canceled && err != io.EOF {
		if err == context.DeadlineExceeded {
			log.Printf(env.ErrorColor, err)
		} else {
			log.Printf(env.ErrorColor, err)
			return 0, nil
		}
	}
	return resp.StatusCode, resp.Body
}

func ConnectToTarget(listing string) (int, io.ReadCloser) {
	type Connection struct {
		ListingID string `json:"listing_id"`
	}

	connection := Connection{
		ListingID: listing,
	}

	json, err := json.Marshal(connection)
	if err != nil {
		log.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180000*time.Millisecond)
	defer cancel()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest(http.MethodPut, Urls[7], bytes.NewBuffer(json))
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
	}

	SetHeaders(req)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
	}
	return resp.StatusCode, resp.Body
}

func VerifyOptimusDownload() bool {
	var connection data.ConnectionStatus

	ctx, cancel := context.WithTimeout(context.Background(), 180000*time.Millisecond)
	defer cancel()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("GET", Urls[7], nil)
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
		return false
	}
	SetHeaders(req)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Printf(env.ErrorColor, err)
		return false
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	if err := json.Unmarshal(body, &connection); err != nil {
		log.Println("Can not unmarshal JSON")
	}

	if connection.Slug == "scz3994tx0" {
		return true
	} else if connection.Slug == "" {
		log.Printf(env.InfoColor, "No target was connected. Connecting to OPTIMUSDOWNLOAD...")
		canConnect, _ := ConnectToTarget("scz3994tx0")
		if canConnect == 200 {
			log.Printf(env.SuccessColor, "Connected to OPTIMUSDOWNLOAD!")
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
