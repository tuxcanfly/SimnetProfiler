// simnet network profiler

package main

import (
	"github.com/conformal/btcrpcclient"
	"github.com/conformal/btcutil"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func TimeStampUnix() int32 {
	return int32(time.Now().Unix())
}

func main() {

	// Only override the handlers for notifications you care about.
	// Also note most of the handlers will only be called if you register
	// for notifications.  See the documentation of the btcrpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := btcrpcclient.NotificationHandlers{
		OnAccountBalance: func(account string, balance btcutil.Amount, confirmed bool) {
			log.Printf("New balance for account %s: %v", account,
				balance)
		},
	}

	// Connect to local btcwallet RPC server using websockets.
	certHomeDir := btcutil.AppDataDir("btcwallet", false)
	certs, err := ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &btcrpcclient.ConnConfig{
		Host:         "localhost:18554",
		Endpoint:     "ws",
		User:         "rpcuser",
		Pass:         "rpcpass",
		Certificates: certs,
	}
	client, err := btcrpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}

	client.NotifyNewTransactions(true)

	StartTimeStamp := time.Now()
	data, err := client.ListTransactionsCount("", 1)

	if err != nil {
		log.Printf("ListTransactionsCount RPC Error: %s", err)

	} else {
		txid := data[0]
		log.Printf("At Time: %s", StartTimeStamp)
		log.Printf("Transactions Count: %s \n", txid)

	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			log.Printf("Client shutdown: %v, stopping proc and exiting..", sig)

			time.AfterFunc(time.Second*1, func() {
				log.Println("Client shutting down...")
				client.Shutdown()
				log.Println("Client shutdown complete.")
			})
		}
		os.Exit(1)
	}()

	client.WaitForShutdown()
}
