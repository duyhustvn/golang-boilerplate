package server

import (
	"boilerplate/internal/handlers"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/redis/go-redis/v9"
)

// Run the https server
func Run(app *handlers.App) {
	// tlsConfig := &tls.Config{
	// 	PreferServerCipherSuites: true,
	// 	CipherSuites: []uint16{
	// 		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	// 		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	// 		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	// 		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	// 	},
	// 	MinVersion: tls.VersionTLS12,
	// }

	subscribe := app.Pubsub.Subscribe(context.Background(), app.Cfg.Redis.Channel)

	// Close the subscription when we are done
	defer subscribe.(*redis.PubSub).Close()

	go func() {
		fmt.Println("Start listen message in thread")
		for {
			msg, err := subscribe.(*redis.PubSub).ReceiveMessage(context.Background())
			if err != nil {
				panic(err)
			}

			fmt.Println("received message: ", msg.Payload)
		}
	}()

	// ch := subscribe.(*redis.PubSub).Channel()
	// for msg := range ch {
	// 	fmt.Println(msg.Channel, msg.Payload)
	// }

	if err := app.Pubsub.Publish(context.Background(), app.Cfg.Redis.Channel, "Ping"); err != nil {
		fmt.Println("Publish error", err)
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", app.Cfg.Server.Port),
	}

	runHTTP := func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Println((fmt.Sprintf("Listening on port: %s ...", app.Cfg.Server.Port)))

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("ListenAndServe error: ", err)
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go runHTTP(wg)
	wg.Wait()
}
