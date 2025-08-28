package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string     `json:"username" gorm:"unique;not null"`
	Password  string     `json:"-" gorm:"not null"`
	Rating    int        `json:"rating" gorm:"default:1200"`
	AvatarURL *string    `json:"avatarUrl,omitempty"`
	FCMToken  *string    `json:"fcmToken,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type LobbyType string
type LobbyStatus string

const (
	LobbyTypePublic  LobbyType = "PUBLIC"
	LobbyTypePrivate LobbyType = "PRIVATE"
	
	LobbyStatusOpen    LobbyStatus = "OPEN"
	LobbyStatusMatched LobbyStatus = "MATCHED"
	LobbyStatusClosed  LobbyStatus = "CLOSED"
)

type TimeControl struct {
	Initial   int `json:"initial"`   // seconds
	Increment int `json:"increment"` // seconds
}

type Lobby struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Type        LobbyType    `json:"type" gorm:"not null"`
	CreatorID   uuid.UUID    `json:"creatorId" gorm:"type:uuid;not null"`
	Creator     User         `json:"creator" gorm:"foreignKey:CreatorID"`
	InviteCode  *string      `json:"inviteCode,omitempty"`
	TimeControl TimeControl  `json:"timeControl" gorm:"type:jsonb"`
	Status      LobbyStatus  `json:"status" gorm:"default:'OPEN'"`
	CreatedAt   time.Time    `json:"createdAt"`
}

type GameResult string
type Termination string

const (
	GameResultWhiteWins GameResult = "1-0"
	GameResultBlackWins GameResult = "0-1"
	GameResultDraw      GameResult = "1/2-1/2"
	
	TerminationCheckmate     Termination = "CHECKMATE"
	TerminationStalemate     Termination = "STALEMATE"
	TerminationThreefold     Termination = "THREEFOLD"
	TerminationFiftyMove     Termination = "FIFTY_MOVE"
	TerminationInsufficient  Termination = "INSUFFICIENT"
	TerminationTime          Termination = "TIME"
	TerminationResign        Termination = "RESIGN"
	TerminationAgreement     Termination = "AGREEMENT"
)

type GameClocks struct {
	WhiteMs int `json:"wMs"`
	BlackMs int `json:"bMs"`
}

type Game struct {
	ID           uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WhiteID      uuid.UUID    `json:"whiteId" gorm:"type:uuid;not null"`
	BlackID      uuid.UUID    `json:"blackId" gorm:"type:uuid;not null"`
	White        User         `json:"white" gorm:"foreignKey:WhiteID"`
	Black        User         `json:"black" gorm:"foreignKey:BlackID"`
	FEN          string       `json:"fen" gorm:"not null;default:'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1'"`
	PGN          string       `json:"pgn" gorm:"default:''"`
	Turn         string       `json:"turn" gorm:"not null;default:'w'"` // 'w' or 'b'
	Clocks       GameClocks   `json:"clocks" gorm:"type:jsonb"`
	Result       *GameResult  `json:"result,omitempty"`
	Termination  *Termination `json:"termination,omitempty"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

type Move struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GameID    uuid.UUID `json:"gameId" gorm:"type:uuid;not null"`
	Game      Game      `json:"game" gorm:"foreignKey:GameID"`
	MoverID   uuid.UUID `json:"moverId" gorm:"type:uuid;not null"`
	Mover     User      `json:"mover" gorm:"foreignKey:MoverID"`
	SAN       string    `json:"san" gorm:"not null"`        // Standard Algebraic Notation
	From      string    `json:"from" gorm:"not null"`       // e.g., "e2"
	To        string    `json:"to" gorm:"not null"`         // e.g., "e4"
	Promotion *string   `json:"promotion,omitempty"`        // 'q', 'r', 'b', 'n'
	CreatedAt time.Time `json:"createdAt"`
}