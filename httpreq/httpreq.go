package httpreq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo"
)

type contentType struct {
	MIMEApplicationJSON                  string
	MIMEApplicationJSONCharsetUTF8       string
	MIMEApplicationJavaScript            string
	MIMEApplicationJavaScriptCharsetUTF8 string
	MIMEApplicationXML                   string
	MIMEApplicationXMLCharsetUTF8        string
	MIMETextXML                          string
	MIMETextXMLCharsetUTF8               string
	MIMEApplicationForm                  string
	MIMEApplicationFormUTF8              string
	MIMEApplicationProtobuf              string
	MIMEApplicationMsgpack               string
	MIMETextHTML                         string
	MIMETextHTMLCharsetUTF8              string
	MIMETextPlain                        string
	MIMETextPlainCharsetUTF8             string
	MIMEMultipartForm                    string
	MIMEOctetStream                      string
}

const (
	charsetUTF8 = "charset=UTF-8"
)

type Header struct {
	ContentType string
	Token       string
}

var ContentType = contentType{
	MIMEApplicationJSON:                  "application/json",
	MIMEApplicationJSONCharsetUTF8:       "application/json" + "; " + charsetUTF8,
	MIMEApplicationJavaScript:            "application/javascript",
	MIMEApplicationJavaScriptCharsetUTF8: "application/javascript" + "; " + charsetUTF8,
	MIMEApplicationXML:                   "application/xml",
	MIMEApplicationXMLCharsetUTF8:        "application/xml" + "; " + charsetUTF8,
	MIMETextXML:                          "text/xml",
	MIMETextXMLCharsetUTF8:               "text/xml" + "; " + charsetUTF8,
	MIMEApplicationForm:                  "application/x-www-form-urlencoded",
	MIMEApplicationFormUTF8:              "application/x-www-form-urlencoded" + "; " + charsetUTF8,
	MIMEApplicationProtobuf:              "application/protobuf",
	MIMEApplicationMsgpack:               "application/msgpack",
	MIMETextHTML:                         "text/html",
	MIMETextHTMLCharsetUTF8:              "text/html" + "; " + charsetUTF8,
	MIMETextPlain:                        "text/plain",
	MIMETextPlainCharsetUTF8:             "text/plain" + "; " + charsetUTF8,
	MIMEMultipartForm:                    "multipart/form-data",
	MIMEOctetStream:                      "application/octet-stream",
}

func POST(token, url string, param interface{}, v interface{}) error {
	b, err := json.Marshal(param)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	httpReq.Header.Set("Content-Type", "application/json")
	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := (&http.Client{}).Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg, err := GetErrorMessage(resp)
		if err != nil {
			return &echo.HTTPError{Code: resp.StatusCode, Message: fmt.Sprintf("[%s]%s", resp.Status, err.Error())}

		}
		return &echo.HTTPError{Code: resp.StatusCode, Message: fmt.Sprintf("[%s]%s", resp.Status, errMsg)}
	}
	if v != nil {
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&v); err != nil {
			return err
		}
	}

	return nil
}

func GET(token, url string, v interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("HTTP New Request Error: %s", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP Request Error: %s", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg, err := GetErrorMessage(resp)
		if err != nil {
			return &echo.HTTPError{Code: resp.StatusCode, Message: fmt.Sprintf("[%s]%s", resp.Status, err.Error())}

		}
		return &echo.HTTPError{Code: resp.StatusCode, Message: fmt.Sprintf("[%s]%s", resp.Status, errMsg)}
	}

	if v != nil {
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&v); err != nil {
			return err
		}
	}

	return nil
}

func NewPost(url string, param []byte, header *Header, transport *http.Transport) ([]byte, error) {

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(param))
	if header != nil {
		httpReq.Header.Set("Content-Type", header.ContentType)
		if header.Token != "" {
			httpReq.Header.Set("Authorization", "Bearer "+header.Token)
		}
	} else {
		httpReq.Header.Set("Content-Type", ContentType.MIMEApplicationFormUTF8)
	}
	var client *http.Client
	if transport != nil {
		client = &http.Client{Transport: transport}
	} else {
		client = &http.Client{}
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg, err := GetErrorMessage(resp)
		if err != nil {
			return nil, &echo.HTTPError{Code: resp.StatusCode, Message: fmt.Sprintf("[%s]%s", resp.Status, err.Error())}

		}
		return nil, &echo.HTTPError{Code: resp.StatusCode, Message: fmt.Sprintf("[%s]%s", resp.Status, errMsg)}
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

func NewGet(url string, header *Header, transport *http.Transport) ([]byte, error) {

	httpReq, err := http.NewRequest("GET", url, nil)
	if header != nil {
		httpReq.Header.Set("Content-Type", header.ContentType)
		if header.Token != "" {
			httpReq.Header.Set("Authorization", "Bearer "+header.Token)
		}
	} else {
		httpReq.Header.Set("Content-Type", ContentType.MIMEApplicationFormUTF8)
	}

	var client *http.Client
	if transport != nil {
		client = &http.Client{Transport: transport}
	} else {
		client = &http.Client{}
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg, err := GetErrorMessage(resp)
		if err != nil {
			return nil, &echo.HTTPError{Code: resp.StatusCode, Message: fmt.Sprintf("[%s]%s", resp.Status, err.Error())}

		}
		return nil, &echo.HTTPError{Code: resp.StatusCode, Message: fmt.Sprintf("[%s]%s", resp.Status, errMsg)}
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respData, nil

}

type FileItem struct {
	Key      string //image_content
	FileName string //test.jpg
	Content  []byte //[]byte
}

func NewPostFile(url string, paramTexts map[string]interface{}, paramFile FileItem) ([]byte, error) {
	// if paramFiles ==nil {
	// 	return NewPost(url,paramTexts,header,transport)
	// }

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for k, v := range paramTexts {
		bodyWriter.WriteField(k, v.(string))
	}
	fileWriter, err := bodyWriter.CreateFormFile(paramFile.Key, paramFile.FileName)
	if err != nil {
		return nil, err
	}
	fileWriter.Write(paramFile.Content)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errMsg, err := GetErrorMessage(resp)
		if err != nil {
			return nil, &echo.HTTPError{Code: resp.StatusCode, Message: fmt.Sprintf("[%s]%s", resp.Status, err.Error())}

		}
		return nil, &echo.HTTPError{Code: resp.StatusCode, Message: fmt.Sprintf("[%s]%s", resp.Status, errMsg)}
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respData, nil
}

func GetErrorMessage(resp *http.Response) (string, error) {
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var v struct {
		Result  interface{} `json:"result"`
		Success bool        `json:"success"`
		Error   struct {
			Code    int         `json:"code,omitempty"`
			Details interface{} `json:"details,omitempty"`
			Message string      `json:"message,omitempty"`
		} `json:"error"`
	}
	if err := json.Unmarshal(respData, &v); err != nil {
		return "", err
	}

	return fmt.Sprint(v.Error.Details), nil
}
