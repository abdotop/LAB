package chess

import (
	"errors"
	"strconv"
	"strings"
)

// Piece represents a chess piece
type Piece struct {
	Type  string // 'p', 'r', 'n', 'b', 'q', 'k'
	Color string // 'w', 'b'
}

// Position represents a chess position
type Position struct {
	Board       [8][8]*Piece
	Turn        string // 'w' or 'b'
	Castling    string // KQkq
	EnPassant   string // '-' or square like 'e3'
	Halfmove    int
	Fullmove    int
}

const StartingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// ParseFEN parses a FEN string into a Position
func ParseFEN(fen string) (*Position, error) {
	parts := strings.Fields(fen)
	if len(parts) != 6 {
		return nil, errors.New("invalid FEN: must have 6 parts")
	}
	
	pos := &Position{}
	
	// Parse board
	rows := strings.Split(parts[0], "/")
	if len(rows) != 8 {
		return nil, errors.New("invalid FEN: board must have 8 rows")
	}
	
	for row, rowStr := range rows {
		col := 0
		for _, char := range rowStr {
			if char >= '1' && char <= '8' {
				// Empty squares
				empty := int(char - '0')
				col += empty
			} else {
				// Piece
				if col >= 8 {
					return nil, errors.New("invalid FEN: too many pieces in row")
				}
				
				color := "w"
				if char >= 'a' && char <= 'z' {
					color = "b"
					char = char - 'a' + 'A'
				}
				
				pieceType := strings.ToLower(string(char))
				pos.Board[row][col] = &Piece{Type: pieceType, Color: color}
				col++
			}
		}
		if col != 8 {
			return nil, errors.New("invalid FEN: row must have 8 squares")
		}
	}
	
	// Parse turn
	pos.Turn = parts[1]
	if pos.Turn != "w" && pos.Turn != "b" {
		return nil, errors.New("invalid FEN: turn must be 'w' or 'b'")
	}
	
	// Parse castling rights
	pos.Castling = parts[2]
	
	// Parse en passant
	pos.EnPassant = parts[3]
	
	// Parse halfmove clock
	var err error
	pos.Halfmove, err = strconv.Atoi(parts[4])
	if err != nil {
		return nil, errors.New("invalid FEN: halfmove must be a number")
	}
	
	// Parse fullmove number
	pos.Fullmove, err = strconv.Atoi(parts[5])
	if err != nil {
		return nil, errors.New("invalid FEN: fullmove must be a number")
	}
	
	return pos, nil
}

// ToFEN converts a Position to FEN string
func (p *Position) ToFEN() string {
	var fen strings.Builder
	
	// Board
	for row := 0; row < 8; row++ {
		if row > 0 {
			fen.WriteString("/")
		}
		
		empty := 0
		for col := 0; col < 8; col++ {
			piece := p.Board[row][col]
			if piece == nil {
				empty++
			} else {
				if empty > 0 {
					fen.WriteString(strconv.Itoa(empty))
					empty = 0
				}
				
				char := strings.ToUpper(piece.Type)
				if piece.Color == "b" {
					char = strings.ToLower(char)
				}
				fen.WriteString(char)
			}
		}
		
		if empty > 0 {
			fen.WriteString(strconv.Itoa(empty))
		}
	}
	
	fen.WriteString(" " + p.Turn)
	fen.WriteString(" " + p.Castling)
	fen.WriteString(" " + p.EnPassant)
	fen.WriteString(" " + strconv.Itoa(p.Halfmove))
	fen.WriteString(" " + strconv.Itoa(p.Fullmove))
	
	return fen.String()
}

// Move represents a chess move
type Move struct {
	From      string
	To        string
	Promotion string
}

// ParseMove parses a move string like "e2e4" or "e7e8q"
func ParseMove(moveStr string) (*Move, error) {
	moveStr = strings.TrimSpace(moveStr)
	
	if len(moveStr) < 4 {
		return nil, errors.New("move must be at least 4 characters")
	}
	
	move := &Move{
		From: moveStr[0:2],
		To:   moveStr[2:4],
	}
	
	if len(moveStr) == 5 {
		move.Promotion = string(moveStr[4])
	}
	
	// Validate square notation
	if !isValidSquare(move.From) || !isValidSquare(move.To) {
		return nil, errors.New("invalid square notation")
	}
	
	return move, nil
}

func isValidSquare(square string) bool {
	if len(square) != 2 {
		return false
	}
	
	file := square[0]
	rank := square[1]
	
	return file >= 'a' && file <= 'h' && rank >= '1' && rank <= '8'
}

// SquareToIndices converts a square like "e4" to row, col indices
func SquareToIndices(square string) (int, int) {
	file := int(square[0] - 'a')
	rank := int(square[1] - '1')
	return 7 - rank, file // Row 0 is rank 8, row 7 is rank 1
}

// IndicesToSquare converts row, col indices to square notation
func IndicesToSquare(row, col int) string {
	file := rune('a' + col)
	rank := rune('8' - row)
	return string(file) + string(rank)
}

// IsLegalMove checks if a move is legal in the current position
func (p *Position) IsLegalMove(move *Move) bool {
	fromRow, fromCol := SquareToIndices(move.From)
	toRow, toCol := SquareToIndices(move.To)
	
	// Check if there's a piece at the from square
	piece := p.Board[fromRow][fromCol]
	if piece == nil {
		return false
	}
	
	// Check if it's the right color's turn
	if piece.Color != p.Turn {
		return false
	}
	
	// Check if trying to capture own piece
	targetPiece := p.Board[toRow][toCol]
	if targetPiece != nil && targetPiece.Color == piece.Color {
		return false
	}
	
	// Basic piece movement rules (simplified)
	switch piece.Type {
	case "p": // Pawn
		return p.isLegalPawnMove(fromRow, fromCol, toRow, toCol, piece.Color)
	case "r": // Rook
		return p.isLegalRookMove(fromRow, fromCol, toRow, toCol)
	case "n": // Knight
		return p.isLegalKnightMove(fromRow, fromCol, toRow, toCol)
	case "b": // Bishop
		return p.isLegalBishopMove(fromRow, fromCol, toRow, toCol)
	case "q": // Queen
		return p.isLegalQueenMove(fromRow, fromCol, toRow, toCol)
	case "k": // King
		return p.isLegalKingMove(fromRow, fromCol, toRow, toCol)
	}
	
	return false
}

func (p *Position) isLegalPawnMove(fromRow, fromCol, toRow, toCol int, color string) bool {
	direction := -1
	startRow := 6
	if color == "b" {
		direction = 1
		startRow = 1
	}
	
	// Forward move
	if fromCol == toCol {
		// One square forward
		if toRow == fromRow+direction && p.Board[toRow][toCol] == nil {
			return true
		}
		// Two squares forward from starting position
		if fromRow == startRow && toRow == fromRow+2*direction && p.Board[toRow][toCol] == nil {
			return true
		}
	}
	
	// Diagonal capture
	if abs(fromCol-toCol) == 1 && toRow == fromRow+direction {
		if p.Board[toRow][toCol] != nil {
			return true
		}
		// En passant (simplified - just check if target square matches en passant square)
		if p.EnPassant != "-" {
			epRow, epCol := SquareToIndices(p.EnPassant)
			if toRow == epRow && toCol == epCol {
				return true
			}
		}
	}
	
	return false
}

func (p *Position) isLegalRookMove(fromRow, fromCol, toRow, toCol int) bool {
	// Rook moves horizontally or vertically
	if fromRow != toRow && fromCol != toCol {
		return false
	}
	
	// Check for obstructions
	return p.isPathClear(fromRow, fromCol, toRow, toCol)
}

func (p *Position) isLegalKnightMove(fromRow, fromCol, toRow, toCol int) bool {
	dr := abs(fromRow - toRow)
	dc := abs(fromCol - toCol)
	return (dr == 2 && dc == 1) || (dr == 1 && dc == 2)
}

func (p *Position) isLegalBishopMove(fromRow, fromCol, toRow, toCol int) bool {
	// Bishop moves diagonally
	if abs(fromRow-toRow) != abs(fromCol-toCol) {
		return false
	}
	
	// Check for obstructions
	return p.isPathClear(fromRow, fromCol, toRow, toCol)
}

func (p *Position) isLegalQueenMove(fromRow, fromCol, toRow, toCol int) bool {
	// Queen combines rook and bishop moves
	return p.isLegalRookMove(fromRow, fromCol, toRow, toCol) ||
		p.isLegalBishopMove(fromRow, fromCol, toRow, toCol)
}

func (p *Position) isLegalKingMove(fromRow, fromCol, toRow, toCol int) bool {
	dr := abs(fromRow - toRow)
	dc := abs(fromCol - toCol)
	
	// Normal king move (one square in any direction)
	if dr <= 1 && dc <= 1 {
		return true
	}
	
	// TODO: Implement castling
	return false
}

func (p *Position) isPathClear(fromRow, fromCol, toRow, toCol int) bool {
	dr := sign(toRow - fromRow)
	dc := sign(toCol - fromCol)
	
	r, c := fromRow+dr, fromCol+dc
	for r != toRow || c != toCol {
		if p.Board[r][c] != nil {
			return false
		}
		r += dr
		c += dc
	}
	
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

// MakeMove applies a move to the position and returns a new position
func (p *Position) MakeMove(move *Move) (*Position, error) {
	if !p.IsLegalMove(move) {
		return nil, errors.New("illegal move")
	}
	
	// Create a copy of the position
	newPos := *p
	newBoard := [8][8]*Piece{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if p.Board[i][j] != nil {
				newBoard[i][j] = &Piece{
					Type:  p.Board[i][j].Type,
					Color: p.Board[i][j].Color,
				}
			}
		}
	}
	newPos.Board = newBoard
	
	fromRow, fromCol := SquareToIndices(move.From)
	toRow, toCol := SquareToIndices(move.To)
	
	// Move the piece
	piece := newPos.Board[fromRow][fromCol]
	newPos.Board[fromRow][fromCol] = nil
	newPos.Board[toRow][toCol] = piece
	
	// Handle pawn promotion
	if piece.Type == "p" && move.Promotion != "" {
		piece.Type = move.Promotion
	}
	
	// Switch turn
	if newPos.Turn == "w" {
		newPos.Turn = "b"
	} else {
		newPos.Turn = "w"
		newPos.Fullmove++
	}
	
	// Update halfmove clock (simplified)
	if piece.Type == "p" || p.Board[toRow][toCol] != nil {
		newPos.Halfmove = 0
	} else {
		newPos.Halfmove++
	}
	
	// Reset en passant (simplified)
	newPos.EnPassant = "-"
	
	// Set en passant if pawn moved two squares
	if piece.Type == "p" && abs(fromRow-toRow) == 2 {
		epRow := (fromRow + toRow) / 2
		newPos.EnPassant = IndicesToSquare(epRow, fromCol)
	}
	
	return &newPos, nil
}

// ToSAN converts a move to Standard Algebraic Notation
func (p *Position) ToSAN(move *Move) (string, error) {
	fromRow, fromCol := SquareToIndices(move.From)
	piece := p.Board[fromRow][fromCol]
	if piece == nil {
		return "", errors.New("no piece at from square")
	}
	
	// Simple SAN implementation (not handling all disambiguation cases)
	san := ""
	
	if piece.Type != "p" {
		san += strings.ToUpper(piece.Type)
	}
	
	// Check if capture
	toRow, toCol := SquareToIndices(move.To)
	if p.Board[toRow][toCol] != nil || (piece.Type == "p" && fromCol != toCol) {
		if piece.Type == "p" {
			san += string(move.From[0]) // file of pawn
		}
		san += "x"
	}
	
	san += move.To
	
	// Handle promotion
	if move.Promotion != "" {
		san += "=" + strings.ToUpper(move.Promotion)
	}
	
	// TODO: Add check/checkmate indicators
	
	return san, nil
}

// IsInCheck checks if the king of the given color is in check
func (p *Position) IsInCheck(color string) bool {
	// Find the king
	var kingRow, kingCol int = -1, -1
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			piece := p.Board[row][col]
			if piece != nil && piece.Type == "k" && piece.Color == color {
				kingRow, kingCol = row, col
				break
			}
		}
		if kingRow != -1 {
			break
		}
	}
	
	if kingRow == -1 {
		return false // King not found (shouldn't happen)
	}
	
	// Check if any opponent piece can attack the king
	opponentColor := "b"
	if color == "b" {
		opponentColor = "w"
	}
	
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			piece := p.Board[row][col]
			if piece != nil && piece.Color == opponentColor {
				kingSquare := IndicesToSquare(kingRow, kingCol)
				fromSquare := IndicesToSquare(row, col)
				move := &Move{From: fromSquare, To: kingSquare}
				
				// Temporarily set turn to allow checking opponent moves
				originalTurn := p.Turn
				p.Turn = opponentColor
				if p.IsLegalMove(move) {
					p.Turn = originalTurn
					return true
				}
				p.Turn = originalTurn
			}
		}
	}
	
	return false
}

// GenerateLegalMoves generates all legal moves for the current position
func (p *Position) GenerateLegalMoves() []*Move {
	var moves []*Move
	
	// Generate all pseudo-legal moves
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			piece := p.Board[row][col]
			if piece != nil && piece.Color == p.Turn {
				fromSquare := IndicesToSquare(row, col)
				
				// Generate moves for each piece type
				for toRow := 0; toRow < 8; toRow++ {
					for toCol := 0; toCol < 8; toCol++ {
						toSquare := IndicesToSquare(toRow, toCol)
						move := &Move{From: fromSquare, To: toSquare}
						
						if p.IsLegalMove(move) {
							// Check if move doesn't leave king in check
							newPos, err := p.MakeMove(move)
							if err == nil && !newPos.IsInCheck(p.Turn) {
								moves = append(moves, move)
								
								// Handle pawn promotion
								if piece.Type == "p" && (toRow == 0 || toRow == 7) {
									// Replace the basic move with promotion moves
									moves = moves[:len(moves)-1]
									for _, promotion := range []string{"q", "r", "b", "n"} {
										promMove := &Move{
											From:      fromSquare,
											To:        toSquare,
											Promotion: promotion,
										}
										moves = append(moves, promMove)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	
	return moves
}

// IsCheckmate checks if the current position is checkmate
func (p *Position) IsCheckmate() bool {
	if !p.IsInCheck(p.Turn) {
		return false
	}
	
	moves := p.GenerateLegalMoves()
	return len(moves) == 0
}

// IsStalemate checks if the current position is stalemate
func (p *Position) IsStalemate() bool {
	if p.IsInCheck(p.Turn) {
		return false
	}
	
	moves := p.GenerateLegalMoves()
	return len(moves) == 0
}

// IsInsufficientMaterial checks if the position has insufficient material to mate
func (p *Position) IsInsufficientMaterial() bool {
	// Count pieces
	whitePieces := make(map[string]int)
	blackPieces := make(map[string]int)
	
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			piece := p.Board[row][col]
			if piece != nil {
				if piece.Color == "w" {
					whitePieces[piece.Type]++
				} else {
					blackPieces[piece.Type]++
				}
			}
		}
	}
	
	// King vs King
	if len(whitePieces) == 1 && len(blackPieces) == 1 {
		return true
	}
	
	// King and Bishop/Knight vs King
	if (len(whitePieces) == 2 && len(blackPieces) == 1) ||
		(len(whitePieces) == 1 && len(blackPieces) == 2) {
		for _, pieces := range []map[string]int{whitePieces, blackPieces} {
			if pieces["b"] == 1 || pieces["n"] == 1 {
				return true
			}
		}
	}
	
	return false
}

// ValidateMove validates a move string and returns a structured move
func ValidateMove(fen, moveStr string) (*Move, error) {
	pos, err := ParseFEN(fen)
	if err != nil {
		return nil, err
	}
	
	move, err := ParseMove(moveStr)
	if err != nil {
		return nil, err
	}
	
	if !pos.IsLegalMove(move) {
		return nil, errors.New("illegal move")
	}
	
	return move, nil
}