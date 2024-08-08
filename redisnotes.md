Redis channels themselves don't automatically clear or delete when a user goes offline. The behavior of Redis Pub/Sub channels is as follows:

1. Persistence: 
   - Redis Pub/Sub channels are not persistent by default. They exist as long as there are subscribers.
   - Messages published to a channel are only delivered to current subscribers. If there are no subscribers, the message is simply discarded.

2. Subscriber disconnection:
   - When a subscriber disconnects (user goes offline), they are automatically unsubscribed from all channels they were subscribed to.
   - The channel itself continues to exist as long as there are other subscribers.

3. Channel lifecycle:
   - A channel is created when the first client subscribes to it.
   - The channel remains active as long as there's at least one subscriber.
   - When the last subscriber unsubscribes or disconnects, the channel is effectively removed, but this is an internal Redis operation and doesn't affect future subscriptions to the same channel name.

4. Message handling:
   - Messages are not stored in the channel. They are immediately delivered to current subscribers or discarded if there are no subscribers.

Given these characteristics, here's how you might handle users going offline in a chat application:

1. When a user connects:
   - Subscribe them to their personal channel (e.g., "user:<user_id>").
   - Subscribe them to any group channels they're part of.

2. When a user disconnects:
   - You don't need to manually unsubscribe them, as Redis will handle this automatically.
   - However, you might want to perform some cleanup operations in your application logic.

3. Handling missed messages:
   - Since messages aren't stored in Redis Pub/Sub channels, if you need users to receive messages sent while they were offline, you'll need to implement this separately.
   - Common strategies include:
     - Storing recent messages in Redis lists or other data structures.
     - Using a separate database to store message history.
     - Implementing a message queue system for offline message delivery.

Here's a simple example of how you might handle a user connecting and disconnecting:

```go
package main

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
)

func handleUserConnection(userID string, redisClient *redis.Client) {
    ctx := context.Background()
    pubsub := redisClient.Subscribe(ctx, "user:"+userID)
    defer pubsub.Close()

    fmt.Printf("User %s connected and subscribed\n", userID)

    // Listen for messages
    channel := pubsub.Channel()
    for msg := range channel {
        fmt.Printf("User %s received: %s\n", userID, msg.Payload)
    }

    fmt.Printf("User %s disconnected\n", userID)
    // Perform any necessary cleanup here
}

func main() {
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // Simulate a user connecting
    go handleUserConnection("user1", redisClient)

    // Simulate sending a message
    redisClient.Publish(context.Background(), "user:user1", "Hello, user1!")

    // Keep the program running
    select {}
}
```

In a real application, you'd need to implement additional logic to handle message persistence, offline message delivery, and other features as required by your specific use case.