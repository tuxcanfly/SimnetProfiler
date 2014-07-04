// simnet network profiler

package main

import (
	"fmt"
	"github.com/conformal/btcjson"
	"github.com/conformal/btcrpcclient"
	"github.com/conformal/btcutil"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)




func Serialize(item btcjson.ListTransactionsResult) string {
	amt := fmt.Sprintf("%.8f", item.Amount)
	fee := fmt.Sprintf("%.8f", item.Fee)
	return item.Address + "," + amt + "," + fee
}

connected := make(chan struct{}

func main() {

	// based off of btcwebsocket example for reference
	ntfnHandlers := btcrpcclient.NotificationHandlers{		
		OnBtcdConnected: func(conn bool) {
			if conn && !firstConn {
				firstConn = true
				connected <- struct{}{}
			}
		},
	}

	// Connect to local btcwallet RPC server using websockets.
	certHomeDir := btcutil.AppDataDir("btcd", false)
	certs, err := ioutil.ReadFile(filepath.Join(certHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &btcrpcclient.ConnConfig{
		Host:         "localhost:18556",
		Endpoint:     "ws",
		User:         "rpcuser",
		Pass:         "rpcpass",
		Certificates: certs,
	}
	client, err := btcrpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}
	
	
    // Create the wallet.
	if err := client.CreateEncryptedWallet("walletpass"); err != nil {
		if err := a.cmd.Process.Kill(); err != nil {
			log.Printf("Cannot kill wallet process after failed "+
				"wallet creation: %v", err)
		}
		if err := a.Cleanup(); err != nil {
			log.Printf("Cannot remove actor directory after "+
				"failed wallet creation: %v", err)
		}
		return err
	}

	client.NotifyNewTransactions(true)

	for {
		log.Println("LN 55")
		StartTime := time.Now().Unix()
		data, err := client.ListTransactionsCount("", 100)
		
		log.Println("LN 59")
		// Client
		if err != nil {
			log.Printf("ListTransactionsCount RPC Error: %s", err)
			break
		} else {
            txn := data[0]
			for {
				log.Println("LN 67")
				TxnArray, err := client.ListTransactionsCount("", 100)
				if err != nil {
					log.Printf("ListTransactionsCount RPC Error: %s", err)
					break
				}
                log.Println("LN 73")
				log.Println("Checking Transactions")
				var TxnFound bool = true

				for i := range TxnArray {
					if Serialize(TxnArray[i]) == Serialize(txn) {
						TxnFound = true
						break
					} else {

						TxnFound = false

					}

				}

				if TxnFound == false {
					EndTime := time.Now().Unix()
					delta := float64((StartTime - EndTime))

					if delta == 0 {
						delta = float64(0.00000001) // at least a satoshi
					}
					tps := 100.0 / delta
					log.Printf("Transactions per second: %f", tps)
					break
				}
			}
		}
	}

	// shutdown with ctrl-c
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
