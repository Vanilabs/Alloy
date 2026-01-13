package database

import 	"github.com/gocql/gocql"


func CreateCassandraEntities(session *gocql.Session) error {
	cql := `
	CREATE TABLE IF NOT EXISTS messages (
	    conversation_id UUID,
	    timestamp timestamp,
	    id UUID,
	    sender_id UUID,
	    text text,
	    attachments text,
	    PRIMARY KEY ((conversation_id), timestamp, id)
	) WITH CLUSTERING ORDER BY (timestamp ASC);
	`
	return session.Query(cql).Exec()
}


// func DropCassandraEntities(session *gocql.Session) error {
//     dropCQL := `DROP TABLE IF EXISTS messages;`
// 	return session.Query(dropCQL).Exec();
// }