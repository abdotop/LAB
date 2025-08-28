# Real-Time Multiplayer Chess App

A production-ready mobile chess application with real-time multiplayer capabilities built with Go (Fiber) backend and Flutter frontend.

## Features

- **Real-time multiplayer chess** with WebSocket communication
- **User authentication** with JWT tokens
- **Public and private lobbies** for game matchmaking
- **Complete chess rules engine** with move validation
- **Game state persistence** with PostgreSQL
- **Responsive mobile UI** with Flutter
- **Docker deployment** ready

## Architecture

### Backend (Go + Fiber)
- **HTTP API** for authentication, lobbies, and games
- **WebSocket server** for real-time game communication
- **Chess rules engine** with FEN/PGN support
- **JWT authentication** middleware
- **PostgreSQL** database with GORM

### Frontend (Flutter)
- **Riverpod** for state management
- **GoRouter** for navigation
- **WebSocket client** for real-time updates
- **Responsive chess board** with drag-and-drop moves

## Quick Start

### Using Docker (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd LAB
```

2. Start the application:
```bash
docker-compose up -d
```

3. The backend will be available at `http://localhost:8080`

### Manual Setup

#### Backend

1. Navigate to backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. Start PostgreSQL database (or use Docker):
```bash
docker run --name chess-db -e POSTGRES_USER=chess -e POSTGRES_PASSWORD=chess -e POSTGRES_DB=chess -p 5432:5432 -d postgres:15
```

5. Run the backend:
```bash
go run .
```

#### Frontend

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install Flutter dependencies:
```bash
flutter pub get
```

3. Run code generation:
```bash
flutter packages pub run build_runner build
```

4. Start the Flutter app:
```bash
flutter run
```

## API Documentation

### Authentication

- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `GET /api/users/me` - Get current user profile
- `PATCH /api/users/me` - Update user profile

### Lobbies

- `POST /api/lobbies` - Create new lobby
- `GET /api/lobbies/public` - Get public lobbies
- `POST /api/lobbies/:id/join` - Join lobby
- `DELETE /api/lobbies/:id` - Close lobby

### Games

- `GET /api/games/:id` - Get game state
- `GET /api/games/:id/moves` - Get game moves
- `POST /api/games/:id/resign` - Resign game
- `POST /api/games/:id/draw-offer` - Offer draw
- `POST /api/games/:id/draw-respond` - Respond to draw offer

### WebSocket Events

#### Client → Server
- `join` - Join game room
- `move` - Make a move
- `resign` - Resign game
- `offerDraw` - Offer draw

#### Server → Client
- `joined` - Joined game room
- `state` - Game state update
- `moveAccepted` - Move was accepted
- `illegalMove` - Move was rejected
- `gameOver` - Game finished
- `yourTurn` - It's your turn

## Database Schema

### Users
- `id` (UUID, Primary Key)
- `username` (String, Unique)
- `password` (String, Hashed)
- `rating` (Integer, Default: 1200)
- `avatar_url` (String, Optional)
- `fcm_token` (String, Optional)

### Lobbies
- `id` (UUID, Primary Key)
- `type` (Enum: PUBLIC/PRIVATE)
- `creator_id` (UUID, Foreign Key)
- `invite_code` (String, Optional)
- `time_control` (JSONB)
- `status` (Enum: OPEN/MATCHED/CLOSED)

### Games
- `id` (UUID, Primary Key)
- `white_id` (UUID, Foreign Key)
- `black_id` (UUID, Foreign Key)
- `fen` (String, Current position)
- `pgn` (String, Game notation)
- `turn` (String, 'w' or 'b')
- `clocks` (JSONB, Time remaining)
- `result` (String, Optional)
- `termination` (String, Optional)

### Moves
- `id` (UUID, Primary Key)
- `game_id` (UUID, Foreign Key)
- `mover_id` (UUID, Foreign Key)
- `san` (String, Standard Algebraic Notation)
- `from` (String, Source square)
- `to` (String, Target square)
- `promotion` (String, Optional)

## Environment Variables

### Backend (.env)
```env
SERVER_PORT=8080
DATABASE_URL=postgres://chess:chess@localhost:5432/chess?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
REDIS_URL=redis://localhost:6379
WS_ORIGINS=https://*
FCM_SERVER_KEY=your-fcm-server-key
```

## Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend
flutter test
```

## Deployment

The application is Docker-ready with multi-stage builds:

1. **Production build:**
```bash
docker-compose -f docker-compose.prod.yml up -d
```

2. **Environment setup:**
   - Configure database connection
   - Set JWT secret
   - Configure CORS origins
   - Set up HTTPS with reverse proxy

## Development

### Code Generation (Frontend)
When modifying model classes:
```bash
cd frontend
flutter packages pub run build_runner build --delete-conflicting-outputs
```

### Database Migrations (Backend)
The application uses GORM's auto-migration feature. For production, consider using proper migration scripts.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Chess piece Unicode symbols for the board
- Fiber framework for the Go backend
- Riverpod for Flutter state management
- PostgreSQL for data persistence