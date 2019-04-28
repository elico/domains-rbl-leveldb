package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/twinj/uuid"
)

func insert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
	case "POST":
	case "HEAD":
	default:
		// Unsupportted method PUT etc..
		JSONOptions(w, r)
		return
	}

	testhost := r.FormValue("host")
	testurl := r.FormValue("url")

	w.Header().Set("Host-Var-Length", strconv.Itoa(len(testhost)))
	w.Header().Set("Url-Var-Length", strconv.Itoa(len(testurl)))

	switch {
	case (len(testhost) > 0):
		_, err := url.Parse("http://" + testhost + "/")
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed host varialble: %s", testhost))
			return
		}
	case (len(testurl) > 0):
		parsedURL, err := url.Parse(testurl)
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed url varialble: %s", testurl))
			return
		}
		testhost = parsedURL.Hostname()
	default:
		JSONError(w, r, fmt.Errorf("host or url varialbles are required"))
		return
	}
	w.Header().Set("X-Test-Domain", testhost)

	err := database.Put([]byte(testhost), []byte("1"), nil)
	if err != nil {
		JSONError(w, r, err)
		return

	}

	newmap := make(map[string]interface{})
	newmap["Insert-Domain"] = testhost

	json.NewEncoder(w).Encode(newmap)

}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
	case "POST":
	case "HEAD":
	default:
		// Unsupportted method PUT etc..
		JSONOptions(w, r)
		return
	}
	testhost := r.FormValue("host")
	testurl := r.FormValue("url")
	w.Header().Set("Host-Var-Length", strconv.Itoa(len(testhost)))
	w.Header().Set("Url-Var-Length", strconv.Itoa(len(testurl)))

	switch {
	case (len(testhost) > 0):
		_, err := url.Parse("http://" + testhost + "/")
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed host varialble: %s", testhost))
			return
		}
	case (len(testurl) > 0):
		parsedURL, err := url.Parse(testurl)
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed url varialble: %s", testurl))
			return
		}
		testhost = parsedURL.Hostname()
	default:
		JSONError(w, r, fmt.Errorf("host or url varialbles are required"))
		return
	}
	w.Header().Set("X-Test-Domain", testhost)

	err := database.Delete([]byte(testhost), nil)
	if err != nil {
		JSONError(w, r, err)
		return

	}

	newmap := make(map[string]interface{})
	newmap["Delete-Domain"] = testhost

	json.NewEncoder(w).Encode(newmap)
}

func demoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
	case "POST":
	case "HEAD":
	default:
		// Unsupportted method PUT etc..
		JSONOptions(w, r)
		return
	}
	testhost := r.FormValue("host")
	testurl := r.FormValue("url")
	w.Header().Set("Host-Var-Length", strconv.Itoa(len(testhost)))
	w.Header().Set("Url-Var-Length", strconv.Itoa(len(testurl)))

	switch {
	case (len(testhost) > 0):
		_, err := url.Parse("http://" + testhost + "/")
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed host varialble: %s", testhost))
			return
		}
	case (len(testurl) > 0):
		parsedURL, err := url.Parse(testurl)
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed url varialble: %s", testurl))
			return
		}
		testhost = parsedURL.Hostname()
	default:
		JSONError(w, r, fmt.Errorf("host or url varialbles are required"))
		return
	}
	w.Header().Set("X-Test-Domain", testhost)
}

func search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
	case "POST":
	case "HEAD":
	default:
		// Unsupportted method PUT etc..
		JSONOptions(w, r)
		return
	}
	testhost := r.FormValue("host")
	testurl := r.FormValue("url")
	w.Header().Set("Host-Var-Length", strconv.Itoa(len(testhost)))
	w.Header().Set("Url-Var-Length", strconv.Itoa(len(testurl)))

	switch {
	case (len(testhost) > 0):
		_, err := url.Parse("http://" + testhost + "/")
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed host varialble: %s", testhost))
			return
		}
	case (len(testurl) > 0):
		parsedURL, err := url.Parse(testurl)
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed url varialble: %s", testurl))
			return
		}
		testhost = parsedURL.Hostname()
	default:
		JSONError(w, r, fmt.Errorf("host or url varialbles are required"))
		return
	}

	response := make(map[string]interface{})
	response["OK"] = "1"

	_, err := database.Get([]byte(testhost), nil)

	switch {
	case err == leveldb.ErrNotFound:
	case err != nil:
		JSONError(w, r, fmt.Errorf(err.Error()))
		return
	default:
		// There is a row and there is no error aka blacklisted
		response["blacklisted"] = "true"
		response["OK"] = "0"
		w.Header().Set("X-Vote", "BLOCK")
	}

	json.NewEncoder(w).Encode(response)
}

func batchInsert(w http.ResponseWriter, r *http.Request) {
	u := uuid.NewV4()

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "PUT":
	case "POST":
	default:
		// Unsupportted method PUT etc..
		JSONOptions(w, r)
		return
	}

	var bodyBytes []byte

	// w.Header().Set("Host-Var-Length", strconv.Itoa(len(testhost)))
	// w.Header().Set("Url-Var-Length", strconv.Itoa(len(testurl)))

	switch r.Method {
	case "PUT":
		if r.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(r.Body)
		}
	case "POST":
		file, _, err := r.FormFile("listfile")
		if err != nil {
			fmt.Println(err)
		} else {
			bodyBytes, _ = ioutil.ReadAll(file)
		}
	default:
		JSONBatchOptions(w, r)
		return
	}

	bodyString := string(bodyBytes)
	if len(bodyString) == 0 {
		http.Error(w, "Empty body detected", 400)
		return
	}

	fmt.Println(u, "BatchInsert Request file body size => ", len(bodyString))

	lines := strings.Split(string(bodyString), "\n")
	batch := new(leveldb.Batch)

	fmt.Println(u, "BatchInsert Request body lines count => ", len(lines))

	fmt.Println(u, ": Started a BatchInsert update from => ", r.RemoteAddr)
	significantProgress := len(lines) / 20
	for i, line := range lines {
		if len(line) > 1 {
			batch.Put([]byte(line), []byte("1"))

		}
		if i > 0 && significantProgress > 0 && (i%significantProgress) == 0 {
			fmt.Println(u, ": Added to the batch => "+strconv.Itoa(i))
		}
	}

	err := database.Write(batch, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		fmt.Println(err)
		return
	}
	fmt.Println(u, ": Fininshed adding to the batch => ", len(lines))

	newmap := make(map[string]interface{})
	newmap["Body-Size"] = len(bodyString)
	newmap["Lines-Count"] = len(lines)

	json.NewEncoder(w).Encode(newmap)
}

func batchDelete(w http.ResponseWriter, r *http.Request) {
	u := uuid.NewV4()

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "PUT":
	case "POST":
	default:
		// Unsupportted method PUT etc..
		JSONOptions(w, r)
		return
	}

	var bodyBytes []byte

	// w.Header().Set("Host-Var-Length", strconv.Itoa(len(testhost)))
	// w.Header().Set("Url-Var-Length", strconv.Itoa(len(testurl)))

	switch r.Method {
	case "PUT":
		if r.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(r.Body)
		}
	case "POST":
		file, _, err := r.FormFile("listfile")
		if err != nil {
			fmt.Println(err)
		} else {
			bodyBytes, _ = ioutil.ReadAll(file)
		}
	default:
		JSONBatchOptions(w, r)
		return
	}

	bodyString := string(bodyBytes)
	if len(bodyString) == 0 {
		http.Error(w, "Empty body detected", 400)
		return
	}

	fmt.Println(u, "BatchDelete Request file body size => ", len(bodyString))

	lines := strings.Split(string(bodyString), "\n")
	batch := new(leveldb.Batch)

	fmt.Println(u, "BatchDelete body lines count => ", len(lines))

	fmt.Println(u, ": Started a BatchDelete update from => ", r.RemoteAddr)
	significantProgress := len(lines) / 20
	for i, line := range lines {
		if len(line) > 1 {
			batch.Delete([]byte(line))

		}
		if i > 0 && significantProgress > 0 && (i%significantProgress) == 0 {
			fmt.Println(u, ": Added to the batch => "+strconv.Itoa(i))
		}
	}

	err := database.Write(batch, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		fmt.Println(err)
		return
	}

	newmap := make(map[string]interface{})
	newmap["Body-Size"] = len(bodyString)
	newmap["Lines-Count"] = len(lines)

	json.NewEncoder(w).Encode(newmap)
}
