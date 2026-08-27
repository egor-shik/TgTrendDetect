package config 

type Config struct {
    AppID int 
    AppHash string 
    Phone string 
    SessionFilePath string 
}

func NewConfig(appID int, appHash, phone, sessionFilePath string) *Config {
    return &Config{
    AppID:  appID,
    AppHash:  appHash,
    Phone: phone,
    SessionFilePath: sessionFilePath,
    }
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