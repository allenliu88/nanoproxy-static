package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/allenliu88/nanoproxy-static/logging"
	"github.com/go-logr/zapr"
	"knative.dev/pkg/signals"
)

var passthruRequestHeaderKeys = [...]string{
	"Accept",
	"Accept-Encoding",
	"Accept-Language",
	"Cache-Control",
	"Cookie",
	"Referer",
	"User-Agent",
}

var passthruResponseHeaderKeys = [...]string{
	"Content-Encoding",
	"Content-Language",
	"Content-Type",
	"Cache-Control", // TODO: Is this valid in a response?
	"Date",
	"Etag",
	"Expires",
	"Last-Modified",
	"Location",
	"Server",
	"Vary",
	"Transfer-Encoding",
}

var (
	verbose = kingpin.Flag("verbose", "Verbose log flag").Short('v').Default("false").Bool()
	port    = kingpin.Flag("port", "HTTP port").Short('p').Default("8080").Int()
	target  = kingpin.Flag("target", "Target host or host:port").Short('t').Required().String()

	// Root Context
	ctx = signals.NewContext()
	// Logging
	baseLogger = zapr.NewLogger(logging.NewLogger(ctx, "nanoproxy-static"))
	logger     = logging.IgnoreDebugEvents(baseLogger)
)

func main() {
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	handler := http.DefaultServeMux

	handler.HandleFunc("/", handleFunc)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Info("Starting to listen and serve...")
	s.ListenAndServe()
	select {}
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	logger.Info(fmt.Sprintf("--> %v %v\n", r.Method, r.URL))

	// Construct filtered header to send to origin server
	hh := http.Header{}
	for _, hk := range passthruRequestHeaderKeys {
		if hv, ok := r.Header[hk]; ok {
			hh[hk] = hv
		}
	}

	// Host
	hh.Add("Host", *target)

	// Construct request to send to origin server
	rr := http.Request{
		Method: r.Method,
		URL:    r.URL,
		Header: hh,
		Body:   r.Body,
		// TODO: Is this correct for a 0 value?
		//       Perhaps a 0 may need to be reinterpreted as -1?
		ContentLength: r.ContentLength,
		Close:         r.Close,
	}

	// Forward request to origin server
	rr.URL.Scheme = "http"
	rr.URL.Host = *target
	logger.Info(fmt.Sprintf("--> %v %v\n", rr.Method, rr.URL))
	resp, err := http.DefaultTransport.RoundTrip(&rr)
	if err != nil {
		// TODO: Passthru more error information
		http.Error(w, "Could not reach origin server", 500)
		return
	}
	defer resp.Body.Close()

	if *verbose {
		logger.Info(fmt.Sprintf("--> %+v\n", rr.Header))
		logger.Info(fmt.Sprintf("<-- %v %v %+v\n", resp.Status, resp.ContentLength, resp.Header))
	} else {
		logger.Info(fmt.Sprintf("<-- %v\n", resp.Status))
	}

	// Transfer filtered header from origin server -> client
	respH := w.Header()
	for _, hk := range passthruResponseHeaderKeys {
		if hv, ok := resp.Header[hk]; ok {
			respH[hk] = hv
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Transfer response from origin server -> client
	if resp.ContentLength > 0 {
		// (Ignore I/O errors, since there's nothing we can do)
		io.CopyN(w, resp.Body, resp.ContentLength)
	} else /**if resp.Close**/ { // TODO: Is this condition right? No, fixed by allen.liu@2024-7-10 remove the [if resp.Close], resp.ContentLength may be -1 represent Unkonw
		// Copy until EOF or some other error occurs
		for {
			if _, err := io.Copy(w, resp.Body); err != nil {
				break
			}
		}
	}
}
