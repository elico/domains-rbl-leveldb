package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func jsonPrettyPrint(in []byte) []byte {
	var out bytes.Buffer
	err := json.Indent(&out, in, "", "\t")
	if err != nil {
		return in
	}
	return out.Bytes()
}

func jsonUnsupportedMethod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cacheResponseFor(w, r, 600)
	http.Error(w, "Unsupported Method", 405)
	errorDetails := make(map[string]interface{})
	errorDetails["error"] = 405
	errorDetails["details"] = "Unsupported Method"
	json.NewEncoder(w).Encode(errorDetails)
}

// JSONOptions ---
func JSONOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cacheResponseFor(w, r, 600)
	w.Header().Set("Allow", "OPTIONS, GET, HEAD, POST")
	errorDetails := make(map[string]interface{})
	errorDetails["error"] = 405
	errorDetails["details"] = "Unsupported Method"
	json.NewEncoder(w).Encode(errorDetails)
}

// JSONBatchOptions ---
func JSONBatchOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cacheResponseFor(w, r, 600)
	w.Header().Set("Allow", "OPTIONS, PUT, POST")
	errorDetails := make(map[string]interface{})
	errorDetails["error"] = 405
	errorDetails["details"] = "Unsupported Method"
	json.NewEncoder(w).Encode(errorDetails)
}

// JSONError ---
func JSONError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	switch {
	case strings.Contains("host or url varialbles are required", err.Error()):
	default:
		fmt.Fprintf(os.Stderr, "Error accured: %s\n", err)
	}

	w.Header().Set("Allow", "OPTIONS, GET, HEAD, POST")
	errorDetails := make(map[string]interface{})
	errorDetails["error"] = err
	errorDetails["details"] = err.Error()
	json.NewEncoder(w).Encode(errorDetails)
}
