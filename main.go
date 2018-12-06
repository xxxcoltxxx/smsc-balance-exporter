package main

import (
    "context"
    "encoding/json"
    "errors"
    "flag"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/prometheus/common/version"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"
    "fmt"
    "strings"
)

var addr = flag.String("listen-address", "0.0.0.0:9601", "The address to listen on for HTTP requests.")
var interval = flag.Int("interval", 3600, "Interval (in seconds) for querying balance.")

var balanceGauge *prometheus.GaugeVec

type smscResponse struct {
    Balance string `json:"balance"`
}

type smscConfig struct {
    Login    string
    Password string
}

var credentials = smscConfig{}

func init() {
    balanceGauge = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Subsystem: "balance",
            Name:      "smsc",
            Help:      "Balance for service in smsc account",
        },
        []string{"service"},
    )

    prometheus.MustRegister(balanceGauge)

    flag.Parse()
}

func main() {
    log.Println("Starting Smsc balance exporter", version.Info())
    log.Println("Build context", version.BuildContext())

    if err := readConfig(); err != nil {
        log.Println(err.Error())
        return
    }

    loadBalance()
    go startBalanceUpdater(*interval)
    srv := &http.Server{
        Addr:         *addr,
        WriteTimeout: time.Second * 2,
        ReadTimeout:  time.Second * 2,
        IdleTimeout:  time.Second * 60,

        Handler: nil,
    }

    http.Handle("/metrics", promhttp.Handler())
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "static/index.html")
    })

    go func() {
        log.Fatal(srv.ListenAndServe())
    }()

    log.Printf("Smsc balance exporter has been started at address %s", *addr)
    log.Printf("Exporter will update balance every %d seconds", *interval)

    c := make(chan os.Signal, 1)

    signal.Notify(c, os.Interrupt)
    signal.Notify(c, syscall.SIGTERM)

    <-c

    log.Println("Smsc balance exporter shutdown")
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()

    err := srv.Shutdown(ctx)
    if err != nil {
        log.Fatal(err)
    }

    os.Exit(0)
}

func readConfig() error {
    if login, ok := os.LookupEnv("SMSC_LOGIN"); ok {
        credentials.Login = login
    } else {
        return errors.New("environment \"SMSC_LOGIN\" is not set")
    }

    if password, ok := os.LookupEnv("SMSC_PASSWORD"); ok {
        credentials.Password = password
    } else {
        return errors.New("environment \"SMSC_PASSWORD\" is not set")
    }

    return nil
}

func startBalanceUpdater(i int) {
    for {
        time.Sleep(time.Second * time.Duration(i))
        loadBalance()
    }
}

func securePrintf(format string, args ...interface{}) {
    var message = fmt.Sprintf(format, args...)
    message = strings.Replace(message, credentials.Login, "<smsc-login>", -1)
    message = strings.Replace(message, credentials.Password, "<smsc-password>", -1)
    log.Print(message)
}

func loadBalance() {
    body, err := loadBody()
    if err != nil {
        securePrintf("Error fetching balance: %s", err.Error())
    }

    jsonResponse := smscResponse{}
    if err := json.Unmarshal(body, &jsonResponse); err != nil {
        log.Printf("Error fetching balance: %s", err.Error())
    }
    if b, err := strconv.ParseFloat(jsonResponse.Balance, 2); err != nil {
        log.Printf("Cannot parse balance: %s", err.Error())
    } else {
        balanceGauge.With(prometheus.Labels{"service": credentials.Login}).Set(b)
    }
}

func loadBody() ([]byte, error) {
    client := http.Client{
        Timeout: time.Second * 2,
    }

    req, err := http.NewRequest(http.MethodGet, "https://smsc.ru/sys/balance.php", nil)
    q := req.URL.Query()
    q.Add("login", credentials.Login)
    q.Add("psw", credentials.Password)
    q.Add("fmt", "3")
    req.URL.RawQuery = q.Encode()

    if err != nil {
        securePrintf("Request error: %s", err.Error())
        return []byte{}, err
    }

    res, err := client.Do(req)
    if err != nil {
        securePrintf("Request error: %s", err.Error())
        return []byte{}, err
    }

    defer func() {
        err := res.Body.Close()
        if err != nil {
            securePrintf("Error close response body: %s", err.Error())
        }
    }()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        securePrintf("Error read response: %s", err.Error())
    }

    return body, nil
}
