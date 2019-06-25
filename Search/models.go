package Search

import (
	"fmt"
	"time"
)

type Message struct{
	rowId int
	guid string
	txt string
	dateInt int64
	date time.Time
	fromMe bool
	chat * Chat
	to participant
}
func NewMessage(rowId int, guid string,txt string,date int, fromMe int, participant participant) Message{
	startOfTime,_ := time.Parse("2006-01-02", "2001-01-01")
	m:= Message{rowId:rowId,guid:guid,txt:txt,dateInt:int64(date),fromMe:fromMe==1,to:participant}
	//date is evaluted starting from startOfTime so we have to convert it to regular epoch
	m.date = time.Unix(0, m.dateInt + int64(startOfTime.UnixNano()))
	return m
	}

//a map is used to store participants just to make the search easier/faster
type Chat struct {
	id int
	participants map[string]participant
}

func (c*Chat)addParticipant(id string, handle string)  {
	p,exist := c.participants[id]
	if !exist {
		c.participants[id] = participant{id:id,handles:[]string{handle}}
	}else{
		handles :=append(p.handles,handle)
		c.participants[id]= participant{id:id,handles:handles}
	}
}

type participant struct {
	id string
	handles []string
}

func (m Message) String() string{
	val := fmt.Sprintf("%s\t%s\n",m.txt,m.date.UTC().Format("2006-01-02 15:04"))
	return val
}







