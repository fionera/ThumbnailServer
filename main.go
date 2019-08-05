package main

import (
    "context"
    "flag"
    "fmt"
    "github.com/sirupsen/logrus"
    "net/http"
    "os"
    "os/signal"
    "runtime/pprof"
)

var srcFolder string

var httpServer *http.Server
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
    flag.Parse()
    srcFolder = flag.Arg(0)

    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            logrus.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    ctx, cancel := context.WithCancel(context.Background())

    go listenCtrlC(ctx, cancel)

    httpServer = &http.Server{
        Addr: ":8081",
    }

    http.HandleFunc("/", handle)
    logrus.Info("Starting on Port 8081")
    err := httpServer.ListenAndServe()
    if err != nil {
        logrus.Fatal(err)
    }
}

func listenCtrlC(ctx context.Context, cancel context.CancelFunc) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    <-c
    httpServer.Shutdown(ctx)
    cancel()
    fmt.Fprintln(os.Stderr, "Press ^C again to exit instantly.")
    <-c
    fmt.Fprintln(os.Stderr, "\nKilled!")
    os.Exit(255)
}

func handle(w http.ResponseWriter, r *http.Request) {
    err, tp := ParseThumbnailPath(r.URL.Path)

    switch err {
    case ErrNoThumbnail:
        handleDefault(w, r)
    case nil:
        handleThumbnail(w, r, tp)
    default:
        logrus.Error(err)
        handleError(w, r, err)
    }
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
    _, err := fmt.Fprint(w, "Hello World!")
    if err != nil {
        logrus.Error(err)
    }
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
    w.Write([]byte(err.Error()))
}

func handleThumbnail(w http.ResponseWriter, r *http.Request, tp *ThumbnailPath) {
    err := CreateThumbnail(tp, w)
    if err != nil {
        logrus.Error(err)
        handleError(w, r, err)
    }
}
