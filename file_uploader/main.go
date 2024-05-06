package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// content is a struct which contains a file's name, its type and its data.
type content struct {
    fname string
    ftype string
    fdata []byte
}

func sendPostRequest(url string, files ...content) ([]byte, error) {
    var (
        buf = new(bytes.Buffer)
        w   = multipart.NewWriter(buf)
    )

    for _, f := range files {
        part, err := w.CreateFormFile(f.ftype, filepath.Base(f.fname))
        if err != nil {
            return []byte{}, err
        }

        _, err = part.Write(f.fdata)
        if err != nil {
            return []byte{}, err
        }
    }

    err := w.Close()
    if err != nil {
        return []byte{}, err
    }

    req, err := http.NewRequest("POST", url, buf)
    if err != nil {
        return []byte{}, err
    }
    req.Header.Add("Content-Type", w.FormDataContentType())

    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        return []byte{}, err
    }
    defer res.Body.Close()

    cnt, err := io.ReadAll(res.Body)
    if err != nil {
        return []byte{}, err
    }
    return cnt, nil
}


// Creates a new file upload http request with optional extra params
func sendPostRequestFile(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}


func main() {
	// // Example usage
	// cnt, err := sendPostRequest("http://localhost:3100/api/fileUpload", content{
	// 	fname: "file.txt",
	// 	ftype: "file",
	// 	fdata: []byte("Hello, World!"),
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// println(string(cnt))

	// ---

	path, _ := os.Getwd()
	path += "/test.pdf"
	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := sendPostRequestFile("http://localhost:3100/api/fileUpload", extraParams, "file", "test.pdf")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
                resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)

		fmt.Println(body)
	}	
}
