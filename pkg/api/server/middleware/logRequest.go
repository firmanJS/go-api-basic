package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gilcrest/go-API-template/pkg/env"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

// APIAudit struct holds the http request and other attributes needed
// for auditing an http request
type APIAudit struct {
	ctx          context.Context
	RequestID    string    `json:"request_id"`
	TimeStarted  time.Time `json:"time_started"`
	TimeFinished time.Time `json:"time_finished"`
	TimeInMillis int       `json:"time_in_millis"`
	request
	response request
}

type request struct {
	Proto            string `json:"protocol"`
	ProtoMajor       int    `json:"protocol_major"`
	ProtoMinor       int    `json:"protocol_minor"`
	Method           string `json:"request_method"`
	Scheme           string `json:"scheme"`
	Host             string `json:"host"`
	Port             string `json:"port"`
	Path             string `json:"path"`
	Header           string `json:"header"`
	Body             string `json:"body"`
	ContentLength    int64  `json:"content_length"`
	TransferEncoding string `json:"transfer_encoding"`
	Close            bool   `json:"close"`
	Trailer          string `json:"trailer"`
	RemoteAddr       string `json:"remote_address"`
	RequestURI       string `json:"request_uri"`
}

// LogRequest wraps several optional logging functions
//   printRequest - sends request output from httputil.DumpRequest to STDERR
//   logRequest - uses logger util to log requests
//   log2DB - logs request to relational database TODO - not yet implemented
func LogRequest(env *env.Env, aud *APIAudit) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := logSwitch(env, aud, r)
			if err != nil {
				http.Error(w, "Unable to log request", http.StatusBadRequest)
				return
			}
			h.ServeHTTP(w, r) // call original
		})
	}
}

// logSwitch determines which, if any, of the logging methods
// you wish to use will be employed.  Using a cache mechanism I haven't
// implemented yet, you will be able to turn on/off these methods on demand
// as of right now, it's doing all of them, which is ridiculous
func logSwitch(env *env.Env, aud *APIAudit, req *http.Request) error {
	// Check cached key:value pair to determine if printing/logging is on
	// for the service
	// TODO - Implement cache - maybe via Redis?
	var err error

	err = SetRequest(aud, req)
	if err != nil {
		return err
	}

	// err = printRequest(req)
	// if err != nil {
	// 	return err
	// }

	err = logRequest(env, aud)
	if err != nil {
		return err
	}

	err = logRequest2Db(env, aud)
	if err != nil {
		return err
	}

	return nil

}

// SetRequest populates several core fields (TimeStarted, ctx and RequestID)
// for the APIAudit struct being passed into the function
func SetRequest(aud *APIAudit, req *http.Request) error {

	var (
		scheme    string
		jsonBytes []byte
		body      []byte
	)

	// set APIAudit TimeStarted to current time in UTC
	loc, _ := time.LoadLocation("UTC")
	aud.TimeStarted = time.Now().In(loc)

	// split host and port out for cleaner logging
	host, port, err := net.SplitHostPort(req.Host)
	if err != nil {
		return err
	}

	// determine if the request is an HTTPS request
	isHTTPS := req.TLS != nil

	if isHTTPS {
		scheme = "https"
	} else {
		scheme = "http"
	}

	// convert the Header map from the request to an array of bytes
	jsonBytes, err = json.Marshal(req.Header)
	if err != nil {
		return err
	}

	// convert the array of bytes to a JSON string
	headerJSON := string(jsonBytes)

	// convert the Trailer map from the request to an array of bytes
	jsonBytes, err = json.Marshal(req.Trailer)
	if err != nil {
		return err
	}

	// convert the array of bytes to a JSON string
	trailerJSON := string(jsonBytes)

	// get byte Array representation of guid from xid package (12 bytes)
	guid := xid.New()

	// use the String method of the guid object to convert byte array to string (20 bytes)
	rID := guid.String()

	body, err = dumpBody(req)
	if err != nil {
		return err
	}

	aud.ctx = req.Context() // retrieve the context from the http.Request
	aud.RequestID = rID
	aud.request.Proto = req.Proto
	aud.request.ProtoMajor = req.ProtoMajor
	aud.request.ProtoMinor = req.ProtoMinor
	aud.request.Method = req.Method
	aud.request.Scheme = scheme
	aud.request.Host = host
	aud.request.Body = string(body)
	aud.request.Port = port
	aud.request.Path = req.URL.Path
	aud.request.Header = headerJSON
	aud.request.ContentLength = req.ContentLength
	aud.request.TransferEncoding = strings.Join(req.TransferEncoding, ",")
	aud.request.Close = req.Close
	aud.request.Trailer = trailerJSON
	aud.request.RemoteAddr = req.RemoteAddr
	aud.request.RequestURI = req.RequestURI

	return nil

}

// PrintRequest wraps the call to httputil.DumpRequest
func printRequest(req *http.Request) error {

	// func DumpRequest(req *http.Request, body bool) ([]byte, error)
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}
	fmt.Println(string(requestDump))
	return nil
}

// drainBody reads all of b to memory and then returns two equivalent
// ReadClosers yielding the same bytes.
//
// It returns an error if the initial slurp of all bytes fails. It does not attempt
// to make the returned ReadClosers have identical error-matching behavior.
// Function lifted straight from net/http/httputil package
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func dumpBody(req *http.Request) ([]byte, error) {
	var err error
	save := req.Body
	save, req.Body, err = drainBody(req.Body)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer

	chunked := len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"

	if req.Body != nil {
		var dest io.Writer = &b
		if chunked {
			dest = httputil.NewChunkedWriter(dest)
		}
		_, err = io.Copy(dest, req.Body)
		if chunked {
			dest.(io.Closer).Close()
			io.WriteString(&b, "\r\n")
		}
	}

	req.Body = save
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func logFormValues(lgr *zap.Logger, req *http.Request) (*zap.Logger, error) {

	var i int

	err := req.ParseForm()
	if err != nil {
		return nil, err
	}

	for key, valSlice := range req.Form {
		for _, val := range valSlice {
			i++
			formValue := fmt.Sprintf("%s: %s", key, val)
			lgr = lgr.With(zap.String(fmt.Sprintf("Form(%d)", i), formValue))
		}
	}

	for key, valSlice := range req.PostForm {
		for _, val := range valSlice {
			i++
			formValue := fmt.Sprintf("%s: %s", key, val)
			lgr = lgr.With(zap.String(fmt.Sprintf("PostForm(%d)", i), formValue))
		}
	}

	return lgr, nil
}

func logRequest(env *env.Env, aud *APIAudit) error {

	//var err error

	logger := env.Logger
	defer env.Logger.Sync()

	logger.Debug("logRequest started")
	defer logger.Debug("logRequest ended")

	// logger, err = logFormValues(logger, req)
	// if err != nil {
	// 	return err
	// }

	logger.Info("Request received",
		zap.String("request_id", aud.RequestID),
		zap.String("protocol", aud.request.Proto),
		zap.Int("proto_major", aud.request.ProtoMajor),
		zap.Int("proto_minor", aud.request.ProtoMinor),
		zap.String("method", aud.request.Method),
		zap.String("scheme", aud.request.Scheme),
		zap.String("host", aud.request.Host),
		zap.String("port", aud.request.Port),
		zap.String("path", aud.request.Path),
		zap.String("header_json", aud.request.Header),
		zap.Int64("Content Length", aud.request.ContentLength),
		zap.String("Transfer-Encoding", aud.request.TransferEncoding),
		zap.Bool("Close", aud.request.Close),
		zap.String("RemoteAddr", aud.request.RemoteAddr),
		zap.String("RequestURI", aud.request.RequestURI),
	)

	return nil
}

// Creates a record in the appUser table using a stored function
func logRequest2Db(env *env.Env, aud *APIAudit) error {

	var rowsInserted int

	// Calls the BeginTx method of the MainDb opened database
	tx, err := env.DS.MainDb.BeginTx(aud.ctx, nil)
	if err != nil {
		return err
	}

	// Prepare the sql statement using bind variables
	stmt, err := tx.PrepareContext(aud.ctx, `select api.log_request
		(
		p_request_id => $1,
		p_request_timestamp => $2,
		p_protocol => $3,
		p_protocol_major => $4,
		p_protocol_minor => $5,
		p_request_method => $6,
		p_scheme => $7,
		p_host => $8,
		p_port => $9,
		p_path => $10,
		p_header => $11,
		p_content_length => $12,
		p_remote_address => $13)`)

	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(aud.ctx,
		aud.RequestID,
		aud.TimeStarted,
		aud.request.Proto,
		aud.request.ProtoMajor,
		aud.request.ProtoMinor,
		aud.request.Method,
		aud.request.Scheme,
		aud.request.Host,
		aud.request.Port,
		aud.request.Path,
		aud.request.Header,
		aud.request.ContentLength,
		aud.request.RemoteAddr)

	if err != nil {
		return err
	}
	defer rows.Close()

	// Iterate through the returned record(s)
	for rows.Next() {
		if err := rows.Scan(&rowsInserted); err != nil {
			return err
		}
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	// If we have successfully written rows to the db
	// we commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}
