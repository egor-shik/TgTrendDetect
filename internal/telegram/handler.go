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
  processMessage(update.Message, e)
  return nil
 })

 
 dispatcher.OnEditChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateEditChannelMessage) error {
  processMessage(update.Message, e)
  return nil
 })

 gaps := updates.New(updates.Config{
  Handler: dispatcher,
 })

 return gaps
}

func processMessage(msgClass tg.MessageClass, e tg.Entities) {
 msg, ok := msgClass.(*tg.Message)
 if !ok {
  return 
 }

 post := NewPost(msg, e)
 if post.Message == "" && post.Views == 0 {

  fmt.Printf("[Channel %s (%d)] Post #%d | Media Content\n---\n", post.ChanTitle, post.ChanID, post.PostID)
  return
 }

 fmt.Printf("[Channel %s (%d)] Post #%d | Views: %d | Reactions: %d\nText: %s\n---\n",
  post.ChanTitle, post.ChanID, post.PostID, post.Views, post.Reactions, post.Message)
}

type IncomingPost struct {
 Message   string
 PostID    int
 ChanID    int64
 ChanTitle string
 Date      time.Time
 Reactions int
 Views     int
}

func NewPost(msg *tg.Message, e tg.Entities) *IncomingPost {
 var chID int64
 var chTitle string

 if peer, ok := msg.PeerID.(*tg.PeerChannel); ok {
  chID = peer.ChannelID
  if channel, found := e.Channels[chID]; found {
   chTitle = channel.Title
  }
 }

 if chTitle == "" {
  chTitle = "Unknown"
 }

 return &IncomingPost{
  Message:   msg.Message,
  PostID:    msg.ID,
  ChanID:    chID,
  ChanTitle: chTitle,
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