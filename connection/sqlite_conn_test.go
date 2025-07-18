package connection

import "testing"

func TestSqliteConn_UpdateMessageStatus(t *testing.T) {

	conn := SqliteConn{DBPath: "../notifies.db"}
	conn.SetMessageId("123456789", 1, SENT)

}
