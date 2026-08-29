package telegram

import (
 "context"
 "fmt"
 "time"

 "github.com/gotd/td/telegram/updates"
 "github.com/gotd/td/tg"
)

func NewUpdateHandler() *updates.Manager {
    dispatcher := tg.NewUpdateDispatcher()
   
    dispatcher.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
     if msg, ok := update.Message.(*tg.Message); ok {
      post := NewPost(msg)
      fmt.Printf("[Channel %d] Post #%d | Views: %d | Reactions: %d\nText: %s\n---\n",
       post.ChanID, post.PostID, post.Views, post.Reactions, post.Message)
     }
     return nil
    })
   
    gaps := updates.New(updates.Config{
     Handler: dispatcher,
     Storage: updates.NewInMemoryStorage(), 
    })
   
    return gaps
   }

type IncomingPost struct {
 Message   string
 PostID    int
 ChanID    int64
 Date      time.Time
 Reactions int
 Views     int
}

func NewPost(msg *tg.Message) *IncomingPost {
 var chID int64
 if peer, ok := msg.PeerID.(*tg.PeerChannel); ok {
  chID = peer.ChannelID
 }

 return &IncomingPost{
  Message:   msg.Message,
  PostID:    msg.ID,
  ChanID:    chID,
  Date:      time.Unix(int64(msg.Date), 0),
  Reactions: totalReaction(msg.Reactions),
  Views:     msg.Views,
 }
}

func totalReaction(reactions tg.MessageReactions) int {
 total := 0
 for _, r := range reactions.Results {
  total += r.Count
 }
 return total
}