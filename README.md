# Neon bridge - selfhosted homepage dashboard

A modern homelab dashboard built to monitor your selfhosted(and not just) services. The intention of this dashboard is to provide as much as possible valueable information to user. It's not just stats, but also different kind of alerts and messages from services helping you to see which ones needs your attention more. I tried to keep semi-transparent type of UI to leverage background effects. In one of future releases I will allow to manually edit background image/video to make it look even better.

<img width="2009" height="931" alt="Image" src="https://github.com/user-attachments/assets/1f487a7c-1503-4bfc-83ea-a7537862ac1b" />

## 🚀 Features

- **Modern UI**: Beautiful glassmorphism design with multiple themes
- **Theme System**: 3 themes out of box
- **Real-time Updates**: Live system stats and service monitoring
- **Real-time Alerts**: Receive real time alerts from your service instances
- **Responsive Design**: Works perfectly on desktop, tablet, and mobile
- **Plugin Architecture**: Easy to add new widgets and services
- **TypeScript**: Full type safety throughout the application

## Supported widgets

- AdGuard Home
- Sonarr
- Radarr
- Lidarr
- Transmission
- and much more to come

## Installation

The easiest way to get started is with Docker Compose:

```bash
git clone https://github.com/avdept/neon-bridge
cd neon-bridge
docker compose up -d
```

Then visit **http://localhost:3200** to access your dashboard!

That's it! The Docker setup will automatically:

- Build and start the Go backend server
- Build and serve the frontend application
- Set up the database
- Handle all the networking

## 🛠️ Development

### Prerequisites

- Node.js 20.19+ or 22.12+ (Note: Current version 20.16.0 shows warnings but works)
- npm or yarn
- Go lang

### Setup

#### 1. Clone and Install Dependencies

```bash
git clone https://github.com/avdept/neon-bridge
cd neon-bridge
npm install
```

#### 2. Setup Backend (Go Server)

The dashboard requires a Go backend server to fetch data from your services.

```bash
# Navigate to server directory
cd server

# Install Go dependencies
go mod tidy

# Build the server
go build -o neon-bridge-server

# Create environment file (optional)
cp .env.example .env
# Edit .env with your database settings if needed

# Start the server
./neon-bridge-server
```

The server will start on `http://localhost:8080` by default.

#### 3. Setup Frontend

```bash
# In the root directory
npm run dev
```

The frontend will start on `http://localhost:5173` and automatically connect to the backend server.

#### 4. Production Build

```bash
# Build frontend
npm run build

# Build backend
cd server && go build -o dashboard-server

# Start both services (you may want to use a process manager like PM2 or systemd)
./server/dashboard-server & npm run preview
```

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build

## 🎯 View Your Dashboard

Once both the backend server (`http://localhost:8080`) and frontend dev server (`http://localhost:5173`) are running, open http://localhost:5173 to see your homelab dashboard!

### Backend API Endpoints

The Go server provides several API endpoints:

- `GET /api/v1/system/stats` - System statistics
- `GET /api/v1/widgets` - Widget configurations
- `POST /api/v1/widgets` - Create new widgets
- `GET /api/v1/proxy/{service}/{widget_id}` - Proxy requests to services (AdGuard, Sonarr, etc.)

### Adding Widgets

Use the dashboard's built-in widget management UI to add and configure widgets for your services like:

- AdGuard Home
- Sonarr
- Radarr
- Lidarr
- Transmission
