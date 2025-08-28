package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"chess-backend/chess"
	"chess-backend/config"
	"chess-backend/middleware"
	"chess-backend/models"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type Hub struct {
	clients    map[*Client]bool
	gameRooms  map[uuid.UUID]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mutex      sync.RWMutex
}

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
	userID uuid.UUID
	gameID uuid.UUID
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// WebSocket message types
const (
	MsgJoin         = "join"
	MsgLeave        = "leave"
	MsgMove         = "move"
	MsgHeartbeat    = "heartbeat"
	MsgOfferDraw    = "offerDraw"
	MsgRespondDraw  = "respondDraw"
	MsgResign       = "resign"
	
	// Server to client
	MsgJoined       = "joined"
	MsgState        = "state"
	MsgMoveAccepted = "moveAccepted"
	MsgIllegalMove  = "illegalMove"
	MsgOpponentJoined = "opponentJoined"
	MsgYourTurn     = "yourTurn"
	MsgGameOver     = "gameOver"
	MsgError        = "error"
)

// Message payloads
type JoinMessage struct {
	GameID uuid.UUID `json:"gameId"`
}

type MoveMessage struct {
	GameID    uuid.UUID `json:"gameId"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Promotion string    `json:"promotion,omitempty"`
}

type JoinedMessage struct {
	GameID uuid.UUID `json:"gameId"`
	YouAre string    `json:"youAre"` // "WHITE" or "BLACK"
}

type StateMessage struct {
	Game GameSnapshot `json:"game"`
}

type GameSnapshot struct {
	ID          uuid.UUID         `json:"id"`
	FEN         string            `json:"fen"`
	PGN         string            `json:"pgn"`
	Turn        string            `json:"turn"`
	Clocks      models.GameClocks `json:"clocks"`
	Players     GamePlayers       `json:"players"`
	LegalMoves  []string          `json:"legalMoves,omitempty"`
}

type GamePlayers struct {
	White models.User `json:"white"`
	Black models.User `json:"black"`
}

type MoveAcceptedMessage struct {
	GameID uuid.UUID         `json:"gameId"`
	SAN    string            `json:"san"`
	FEN    string            `json:"fen"`
	Turn   string            `json:"turn"`
	Clocks models.GameClocks `json:"clocks"`
}

type YourTurnMessage struct {
	GameID uuid.UUID `json:"gameId"`
}

type GameOverMessage struct {
	Result         models.GameResult   `json:"result"`
	Termination    models.Termination  `json:"termination"`
	WinnerUserID   *uuid.UUID         `json:"winnerUserId,omitempty"`
}

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var GameHub = &Hub{
	clients:    make(map[*Client]bool),
	gameRooms:  make(map[uuid.UUID]map[*Client]bool),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	broadcast:  make(chan []byte),
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			log.Printf("Client %s connected", client.userID)
			h.mutex.Unlock()

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				
				// Remove from game room
				if room, exists := h.gameRooms[client.gameID]; exists {
					delete(room, client)
					if len(room) == 0 {
						delete(h.gameRooms, client.gameID)
					}
				}
				
				log.Printf("Client %s disconnected", client.userID)
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) JoinGameRoom(client *Client, gameID uuid.UUID) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	client.gameID = gameID
	
	if h.gameRooms[gameID] == nil {
		h.gameRooms[gameID] = make(map[*Client]bool)
	}
	
	h.gameRooms[gameID][client] = true
}

func (h *Hub) BroadcastToGame(gameID uuid.UUID, message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if room, exists := h.gameRooms[gameID]; exists {
		for client := range room {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
				delete(room, client)
			}
		}
	}
}

func HandleWebSocket(c *websocket.Conn) {
	// Authenticate user from query params
	token := c.Query("token")
	if token == "" {
		c.WriteJSON(Message{
			Type: MsgError,
			Data: ErrorMessage{Code: "E_WS_401", Message: "Authentication required"},
		})
		c.Close()
		return
	}

	claims, err := middleware.ValidateJWT(token)
	if err != nil {
		c.WriteJSON(Message{
			Type: MsgError,
			Data: ErrorMessage{Code: "E_WS_401", Message: "Invalid token"},
		})
		c.Close()
		return
	}

	client := &Client{
		conn:   c,
		send:   make(chan []byte, 256),
		hub:    GameHub,
		userID: claims.UserID,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		c.handleMessage(msg)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case MsgJoin:
		c.handleJoin(msg.Data)
	case MsgLeave:
		c.handleLeave(msg.Data)
	case MsgMove:
		c.handleMove(msg.Data)
	case MsgHeartbeat:
		c.sendMessage(Message{Type: "heartbeat", Data: map[string]string{"status": "ok"}})
	case MsgOfferDraw:
		c.handleOfferDraw(msg.Data)
	case MsgRespondDraw:
		c.handleRespondDraw(msg.Data)
	case MsgResign:
		c.handleResign(msg.Data)
	default:
		c.sendError("E_MSG_400", "Unknown message type")
	}
}

func (c *Client) handleJoin(data interface{}) {
	var joinMsg JoinMessage
	bytes, _ := json.Marshal(data)
	if err := json.Unmarshal(bytes, &joinMsg); err != nil {
		c.sendError("E_JOIN_400", "Invalid join message")
		return
	}

	// Verify user is part of the game
	var game models.Game
	if err := config.DB.Preload("White").Preload("Black").First(&game, joinMsg.GameID).Error; err != nil {
		c.sendError("E_GAME_404", "Game not found")
		return
	}

	if game.WhiteID != c.userID && game.BlackID != c.userID {
		c.sendError("E_GAME_403", "Not a player in this game")
		return
	}

	// Join game room
	c.hub.JoinGameRoom(c, joinMsg.GameID)

	// Send joined confirmation
	youAre := "WHITE"
	if game.BlackID == c.userID {
		youAre = "BLACK"
	}

	c.sendMessage(Message{
		Type: MsgJoined,
		Data: JoinedMessage{GameID: joinMsg.GameID, YouAre: youAre},
	})

	// Send current game state
	c.sendGameState(game)

	// Notify other players
	c.broadcastToGame(joinMsg.GameID, Message{
		Type: MsgOpponentJoined,
		Data: map[string]interface{}{
			"user": c.getUserInfo(),
		},
	})
}

func (c *Client) handleLeave(data interface{}) {
	// Remove from current game room
	if c.gameID != uuid.Nil {
		c.hub.mutex.Lock()
		if room, exists := c.hub.gameRooms[c.gameID]; exists {
			delete(room, c)
		}
		c.hub.mutex.Unlock()
		c.gameID = uuid.Nil
	}
}

func (c *Client) handleMove(data interface{}) {
	var moveMsg MoveMessage
	bytes, _ := json.Marshal(data)
	if err := json.Unmarshal(bytes, &moveMsg); err != nil {
		c.sendError("E_MOVE_400", "Invalid move message")
		return
	}

	// Get game
	var game models.Game
	if err := config.DB.Preload("White").Preload("Black").First(&game, moveMsg.GameID).Error; err != nil {
		c.sendError("E_GAME_404", "Game not found")
		return
	}

	// Check if it's the player's turn
	isWhitesTurn := game.Turn == "w"
	if (isWhitesTurn && game.WhiteID != c.userID) || (!isWhitesTurn && game.BlackID != c.userID) {
		c.sendError("E_TURN_403", "Not your turn")
		return
	}

	// Validate move
	moveStr := moveMsg.From + moveMsg.To + moveMsg.Promotion
	move, err := chess.ValidateMove(game.FEN, moveStr)
	if err != nil {
		c.sendMessage(Message{
			Type: MsgIllegalMove,
			Data: map[string]interface{}{
				"reason": err.Error(),
				"fen":    game.FEN,
			},
		})
		return
	}

	// Apply move
	pos, _ := chess.ParseFEN(game.FEN)
	newPos, err := pos.MakeMove(move)
	if err != nil {
		c.sendError("E_MOVE_500", "Failed to apply move")
		return
	}

	// Generate SAN
	san, _ := pos.ToSAN(move)

	// Update game state
	game.FEN = newPos.ToFEN()
	game.Turn = newPos.Turn
	if game.PGN == "" {
		game.PGN = san
	} else {
		if newPos.Turn == "w" {
			game.PGN += " " + san
		} else {
			game.PGN += " " + san
		}
	}

	// Save move
	moveRecord := models.Move{
		GameID:    game.ID,
		MoverID:   c.userID,
		SAN:       san,
		From:      move.From,
		To:        move.To,
		Promotion: &move.Promotion,
	}
	config.DB.Create(&moveRecord)

	// Check for game over
	if newPos.IsCheckmate() {
		var result models.GameResult
		var termination models.Termination = models.TerminationCheckmate
		var winnerID *uuid.UUID

		if newPos.Turn == "w" {
			// Black is checkmated, White wins
			result = models.GameResultWhiteWins
			winnerID = &game.WhiteID
		} else {
			// White is checkmated, Black wins
			result = models.GameResultBlackWins
			winnerID = &game.BlackID
		}

		game.Result = &result
		game.Termination = &termination

		c.broadcastToGame(game.ID, Message{
			Type: MsgGameOver,
			Data: GameOverMessage{
				Result:       result,
				Termination:  termination,
				WinnerUserID: winnerID,
			},
		})
	} else if newPos.IsStalemate() {
		result := models.GameResultDraw
		termination := models.TerminationStalemate
		game.Result = &result
		game.Termination = &termination

		c.broadcastToGame(game.ID, Message{
			Type: MsgGameOver,
			Data: GameOverMessage{
				Result:      result,
				Termination: termination,
			},
		})
	} else if newPos.IsInsufficientMaterial() {
		result := models.GameResultDraw
		termination := models.TerminationInsufficient
		game.Result = &result
		game.Termination = &termination

		c.broadcastToGame(game.ID, Message{
			Type: MsgGameOver,
			Data: GameOverMessage{
				Result:      result,
				Termination: termination,
			},
		})
	}

	// Save game
	config.DB.Save(&game)

	// Broadcast move accepted
	c.broadcastToGame(game.ID, Message{
		Type: MsgMoveAccepted,
		Data: MoveAcceptedMessage{
			GameID: game.ID,
			SAN:    san,
			FEN:    game.FEN,
			Turn:   game.Turn,
			Clocks: game.Clocks,
		},
	})

	// Send your turn to next player
	if game.Result == nil {
		nextPlayerID := game.WhiteID
		if game.Turn == "w" {
			nextPlayerID = game.WhiteID
		} else {
			nextPlayerID = game.BlackID
		}

		c.sendToPlayer(nextPlayerID, Message{
			Type: MsgYourTurn,
			Data: YourTurnMessage{GameID: game.ID},
		})
	}
}

func (c *Client) handleOfferDraw(data interface{}) {
	// TODO: Implement draw offer logic
}

func (c *Client) handleRespondDraw(data interface{}) {
	// TODO: Implement draw response logic
}

func (c *Client) handleResign(data interface{}) {
	var gameID uuid.UUID
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &gameID)

	var game models.Game
	if err := config.DB.First(&game, gameID).Error; err != nil {
		c.sendError("E_GAME_404", "Game not found")
		return
	}

	if game.WhiteID != c.userID && game.BlackID != c.userID {
		c.sendError("E_GAME_403", "Not a player in this game")
		return
	}

	// Determine winner
	var result models.GameResult
	var winnerID *uuid.UUID
	if game.WhiteID == c.userID {
		result = models.GameResultBlackWins
		winnerID = &game.BlackID
	} else {
		result = models.GameResultWhiteWins
		winnerID = &game.WhiteID
	}

	termination := models.TerminationResign
	game.Result = &result
	game.Termination = &termination
	config.DB.Save(&game)

	c.broadcastToGame(game.ID, Message{
		Type: MsgGameOver,
		Data: GameOverMessage{
			Result:       result,
			Termination:  termination,
			WinnerUserID: winnerID,
		},
	})
}

func (c *Client) sendMessage(msg Message) {
	bytes, _ := json.Marshal(msg)
	select {
	case c.send <- bytes:
	default:
		close(c.send)
	}
}

func (c *Client) sendError(code, message string) {
	c.sendMessage(Message{
		Type: MsgError,
		Data: ErrorMessage{Code: code, Message: message},
	})
}

func (c *Client) sendGameState(game models.Game) {
	pos, _ := chess.ParseFEN(game.FEN)
	legalMoves := pos.GenerateLegalMoves()
	legalMoveStrings := make([]string, len(legalMoves))
	for i, move := range legalMoves {
		legalMoveStrings[i] = move.From + move.To + move.Promotion
	}

	c.sendMessage(Message{
		Type: MsgState,
		Data: StateMessage{
			Game: GameSnapshot{
				ID:     game.ID,
				FEN:    game.FEN,
				PGN:    game.PGN,
				Turn:   game.Turn,
				Clocks: game.Clocks,
				Players: GamePlayers{
					White: game.White,
					Black: game.Black,
				},
				LegalMoves: legalMoveStrings,
			},
		},
	})
}

func (c *Client) broadcastToGame(gameID uuid.UUID, msg Message) {
	bytes, _ := json.Marshal(msg)
	c.hub.BroadcastToGame(gameID, bytes)
}

func (c *Client) sendToPlayer(userID uuid.UUID, msg Message) {
	bytes, _ := json.Marshal(msg)
	
	c.hub.mutex.RLock()
	defer c.hub.mutex.RUnlock()
	
	for client := range c.hub.clients {
		if client.userID == userID {
			select {
			case client.send <- bytes:
			default:
				close(client.send)
			}
		}
	}
}

func (c *Client) getUserInfo() models.User {
	var user models.User
	config.DB.First(&user, c.userID)
	return user
}