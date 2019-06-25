package Search

import (
	"fmt"
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"time"
)

const MessageQuery string = "SELECT M.guid,M.text,M.date,M.is_from_me,H.id,H.person_centric_id as appleId,M.ROWID FROM MESSAGE AS M LEFT JOIN HANDLE AS H ON  M.handle_id = H.ROWID WHERE `TEXT` LIKE ?"
const ChatQuery string = "SELECT CH.chat_id,H.ID,H.person_centric_id FROM chat_handle_join AS CH LEFT JOIN handle AS H ON CH.handle_id = H.ROWID LEFT JOIN chat_message_join AS CM ON CH.chat_id = CM.chat_id WHERE CM.message_id = ?"
const ChatMessageQuery = "SELECT M.ROWID,M.guid,M.text,M.date,M.is_from_me FROM MESSAGE AS M JOIN chat_message_join AS CM ON CM.MESSAGE_ID = M.ROWID WHERE CM.CHAT_ID = ? AND M.DATE > ? AND M.DATE < ? ORDER BY date ASC"
type Query struct {
	Db string
	searchStmt * sqlite3.Stmt
	enrichStmt * sqlite3.Stmt
	searchChatStmt * sqlite3.Stmt

}


func (this*Query) open() (* sqlite3.Conn, bool){
	connection, err := sqlite3.Open(this.Db)
	if err != nil {
		fmt.Printf("Failed to access the database with error %s\n", err)
		return nil,true
	}
	connection.BusyTimeout(5 * time.Second)

	return connection,false
}
func (this *Query)GetChatMessages(input* Message, radiusInSec int) ([]Message,bool)  {
	enriched :=input.chat==nil
	if !enriched {
		enriched = this.Enrich(input)
	}
	if !enriched{
		return nil,false
	} else {
		radiusInNs := int64(radiusInSec)*time.Second.Nanoseconds()
		center := int64(input.dateInt)
		connection, err := this.open()
		if err{
			return nil, false
		}
		defer connection.Close()
		var stmtErr error
		if this.searchChatStmt != nil{
			this.searchChatStmt.Reset()
			stmtErr = this.searchChatStmt.Bind(input.rowId, center-radiusInNs, center+radiusInNs)
		} else{
			this.searchChatStmt, stmtErr = connection.Prepare(ChatMessageQuery, input.rowId)
		}
		if stmtErr != nil{
			fmt.Printf("Failed to prepare the query statement with error %s\n", stmtErr)
			return nil, false
		}
		messages := make([]Message, 10, 10)
		defer this.searchChatStmt.Close()
		for{
			hasRow, err := this.searchChatStmt.Step()
			if err != nil{
				fmt.Printf("Error while stepping through results with message %s\n", err)
				return nil, false
			}
			if !hasRow{
				// The query is finished
				break
			}
			rowId := getInt(this.searchStmt, 0)
			guid := getText(this.searchStmt, 1)
			text := getText(this.searchStmt, 2)
			date := getInt(this.searchStmt, 3)
			fromMe := getInt(this.searchStmt, 4)

			msg := NewMessage(rowId, guid, text, date, fromMe, input.to)
			msg.chat = input.chat
			messages = append(messages, msg)
		}
		return messages, true
	}
}

func (this*Query) Enrich(message* Message) bool{
	connection, err := this.open()
	if err {
		return false
	}
	defer connection.Close()
	var stmtErr error
	if this.enrichStmt != nil {
		this.enrichStmt.Reset()
		stmtErr = this.enrichStmt.Bind(message.rowId)
	}else {
		this.enrichStmt,stmtErr = connection.Prepare(ChatQuery,message.rowId)
	}
	if stmtErr != nil {
		fmt.Printf("Failed to prepare the query statement with error %s\n", stmtErr)
		return false
	}
	chat := Chat{participants: make(map[string]participant)}

	defer this.enrichStmt.Close()
	for  {
		hasRow, err := this.enrichStmt.Step()
		if err != nil {
			fmt.Printf("Error while stepping through results with message %s\n", err)
			return false
		}
		if !hasRow {
			// The query is finished
			break
		}
		//we are setting this value n row time but who cares :D

		chat.id = getInt(this.enrichStmt, 0)

		handleId := getText(this.enrichStmt, 1)
		personId := getText(this.enrichStmt, 2)
		chat.addParticipant(personId,handleId)
	}
	message.chat = &chat
	return true
}

func (this*Query)Search(query string) []Message {
	connection, err :=this.open()
	if err {
		return nil
	}
	defer connection.Close()
	var stmtErr error
	likeQuery := fmt.Sprintf("%%%s%%",query)
	if this.searchStmt != nil {
		this.searchStmt.Reset()
		stmtErr = this.searchStmt.Bind(likeQuery)
	}else {
		this.searchStmt,stmtErr = connection.Prepare(MessageQuery,likeQuery)
	}
	if stmtErr != nil {
		fmt.Printf("Failed to prepare the query statement with error %s\n", stmtErr)
		return nil
	}
	defer this.searchStmt.Close()
	messages :=make([]Message,0,10)
	for  {
		hasRow, err := this.searchStmt.Step()
		if err != nil {
			fmt.Printf("Error while stepping through results with message %s\n", err)
			return nil
		}
		if !hasRow {
			// The query is finished
			break
		}
		guid:= getText(this.searchStmt, 0)
		text:= getText(this.searchStmt, 1)
		date:= getInt(this.searchStmt,2)
		fromMe := getInt(this.searchStmt,3)
		handleId := getText(this.searchStmt, 4)
		personId := getText(this.searchStmt, 5)
		rowId := getInt(this.searchStmt,6)

		p := participant{id:personId,handles:[]string{handleId}}
		message := NewMessage(rowId,guid,text,date,fromMe,p)
		messages = append(messages, message)
	}
	return messages
}
func getText(stmt * sqlite3.Stmt, column int) string {
	value, ok, err := stmt.ColumnText(column)

	if err != nil {
		fmt.Printf("error while retriving text from stmt at index %d with error %s",column, err)
		return ""
	}
	if !ok {
		// The column was NULL
		return ""
	}
	return value
}
func getInt(stmt * sqlite3.Stmt, column int) int {
	value, ok, err := stmt.ColumnInt(column)

	if err != nil {
		fmt.Printf("error while retriving text from stmt at index %d with error %s",column, err)
		return 0
	}
	if !ok {
		// The column was NULL
		return 0
	}
	return value
}