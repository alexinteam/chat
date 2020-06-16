// DB-related functions.

package chat

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	// Global connection to DB
	db *sql.DB

	// Users
	stmtGetUserById       *sql.Stmt
	stmtGetUserByUsername *sql.Stmt
	// Authentication
	stmtMakeSession        *sql.Stmt
	stmtGetUserBySession   *sql.Stmt
	stmtDeleteSession      *sql.Stmt
	stmtDeleteUserSessions *sql.Stmt
	// Rooms
	stmtGetAllRooms *sql.Stmt
	// Messages
	stmtInsertMessage         *sql.Stmt
	stmtGetRoomMessagesByUser *sql.Stmt
)

func prepareStmt(db *sql.DB, query string) *sql.Stmt {
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal("Could not prepare '" + query + "': " + err.Error())
	}
	return stmt
}

func initStmts() {
	// Users
	stmtGetUserById = prepareStmt(db, `
        SELECT id, full_name, username, email
        FROM auth_user
        WHERE id = $1
    `)
	stmtGetUserByUsername = prepareStmt(db, `
        SELECT id, full_name, username, email, password
        FROM auth_user
        WHERE username = $1
    `)

	// Authentication
	stmtMakeSession = prepareStmt(db, `
        INSERT INTO auth_session
        (key, user_id, create_date, expire_date)
        VALUES
        ($1, $2, CURRENT_TIMESTAMP, $3)
    `)
	stmtGetUserBySession = prepareStmt(db, `
        SELECT u.id, u.full_name, u.username, u.email
        FROM auth_session AS s
        LEFT JOIN auth_user AS u ON u.id = s.user_id
        WHERE s.key = $1
            AND s.expire_date > CURRENT_TIMESTAMP
    `)
	stmtDeleteSession = prepareStmt(db, `
        DELETE FROM auth_session
        WHERE key = $1
    `)
	stmtDeleteUserSessions = prepareStmt(db, `
        DELETE FROM auth_session
        WHERE user_id = $1
    `)

	// Rooms
	stmtGetAllRooms = prepareStmt(db, `
        SELECT id, name
        FROM room
    `)

	// Messages in room: for the user, from the user or broadcast
	stmtGetRoomMessagesByUser = prepareStmt(db, `
        SELECT *
        FROM (
            SELECT
                m.id, 'message', m.text, m.send_date,
                us.id, us.username, us.full_name, us.email,
                ur.id, ur.username, ur.full_name, ur.email
            FROM message AS m
            -- Sender
            LEFT JOIN auth_user AS us ON us.id = m.sender_id
            -- Recipient
            LEFT JOIN auth_user AS ur ON ur.id = m.recipient_id
            WHERE
                (
                    (m.recipient_id = $1 OR m.recipient_id IS NULL)
                    OR m.sender_id = $1
                )
                AND m.room_id = $2
            ORDER BY m.send_date DESC
            LIMIT $3
        ) AS tmp
        ORDER BY send_date ASC
    `)

	stmtInsertMessage = prepareStmt(db, `
        INSERT INTO message
        (room_id, sender_id, recipient_id, text, send_date)
        VALUES
        ($1, $2, $3, $4, $5)
    `)
}

func dbConnect(dbUser string, dbPass string, dbName string) *sql.DB {
	var err error

	dbConnection := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName)
	db, err = sql.Open("postgres", dbConnection)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("DB ping failed:", err)
	}

	return db
}
