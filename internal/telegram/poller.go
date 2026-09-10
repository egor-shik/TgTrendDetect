package telegram

import (
 "context"
 "fmt"
 "time"

 "github.com/gotd/td/tg"
)

type ChannelPoller struct {
 client *Client
}

func NewChannelPoller(client *Client) *ChannelPoller {
 return &ChannelPoller{client: client}
}

func (p *ChannelPoller) StartPolling(ctx context.Context, channels []*tg.Channel, interval time.Duration) {
 ticker := time.NewTicker(interval)
 defer ticker.Stop()

 lastSeenPost := make(map[int64]int)

 fmt.Printf("[Poller] Started tracking %d channels (interval: %v)\n", len(channels), interval)

 for {
  select {
  case <-ctx.Done():
   return
  case <-ticker.C:
   for _, ch := range channels {
    p.checkChannel(ctx, ch, lastSeenPost)
   }
  }
 }
}

func (p *ChannelPoller) checkChannel(ctx context.Context, ch *tg.Channel, lastSeen map[int64]int) {
    api := p.client.tgClient.API()
   
    history, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
     Peer:  &tg.InputPeerChannel{ChannelID: ch.ID, AccessHash: ch.AccessHash},
     Limit: 5, 
    })
    if err != nil {
     return
    }
   
    var messages []tg.MessageClass
    switch h := history.(type) {
    case *tg.MessagesMessages:
     messages = h.Messages
    case *tg.MessagesMessagesSlice:
     messages = h.Messages
    case *tg.MessagesChannelMessages:
     messages = h.Messages
    }
   
    for _, m := range messages {
     msg, ok := m.(*tg.Message)
     if !ok {
      continue 
     }
   
     if !msg.Post {
      continue
     }
   
     lastID := lastSeen[ch.ID]
   
     if lastID != 0 && msg.ID > lastID {
      post := NewPostFromPoller(msg, ch.Title)
      fmt.Printf("\nPOST DETECTED!\n[Channel: %s | ID: %d]\nViews: %d | Reactions: %d\nText: %s\n---\n",
       post.ChanTitle, post.PostID, post.Views, post.Reactions, post.Message)
     }
     lastSeen[ch.ID] = msg.ID
     break
    }
   }

func NewPostFromPoller(msg *tg.Message, chanTitle string) *IncomingPost {
 var chID int64
 if peer, ok := msg.PeerID.(*tg.PeerChannel); ok {
  chID = peer.ChannelID
 }

 return &IncomingPost{
  Message:   msg.Message,
  PostID:    msg.ID,
  ChanID:    chID,
  ChanTitle: chanTitle,
  Date:      time.Unix(int64(msg.Date), 0),
  Reactions: totalReaction(msg.Reactions),
  Views:     msg.Views,
 }
}