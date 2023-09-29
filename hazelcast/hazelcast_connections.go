package hazelcast

import (
	"context"
	"fmt"
	"time"

	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/logger"
)

// hazelcast connect client to cluster
func Hazelcast_Connect() *hazelcast.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	config := hazelcast.NewConfig()
	config.Cluster.Network.SetAddresses("localhost:5701")
	config.Cluster.Cloud.Token = "Ds32viNXq8GDgjVCizDNgjiq-q5LlOKnnFfDN9LShMY"

	config.Logger.Level = logger.InfoLevel
	client, err := hazelcast.StartNewClientWithConfig(context.Background(), config)
	if err != nil {
		fmt.Println("Error Connetion to Hazekcast ...")
	}
	fmt.Println(ctx)
	//	defer client.Shutdown(ctx)
	return client

}

// // hazelcast connect client to cluster
// 		config.MapConfig("myMap").SetBackupCount(1)
// config.MapConfig("myMap").SetEvictionPolicy(hazelcast.EvictionPolicyNone)
// config.MapConfig("myMap").SetMaxIdleSeconds(0)
// config.MapConfig("myMap").SetTimeToLiveSeconds(0)
// config.MapConfig("myMap").SetInMemoryFormat(hazelcast.InMemoryFormatObject)
// config.MapConfig("myMap").SetPersisted(true)
// config.MapConfig("myMap").SetWriteDelaySeconds(0)
