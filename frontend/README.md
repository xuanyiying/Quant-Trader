# Quant-Trader Frontend

Modern React-based dashboard for real-time monitoring and strategy management.

---

## 🚀 Features

- **Real-time Market Data**: Live price updates via WebSocket
- **Trading Interface**: Paper trading with instant order execution
- **Portfolio Analytics**: Performance tracking with interactive charts
- **Alert Management**: Create and monitor price alerts
- **Strategy Marketplace**: Browse and subscribe to trading strategies
- **Responsive Design**: Works on desktop and mobile devices

---

## 🛠 Tech Stack

| Technology | Purpose |
|------------|---------|
| React 18 | UI Framework |
| Vite | Build Tool |
| TypeScript | Type Safety |
| TailwindCSS | Styling |
| ECharts | Charts & Visualization |
| Zustand | State Management |
| Axios | HTTP Client |

---

## 🏁 Quick Start

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

```bash
cd frontend
npm install
```

### Development

```bash
npm run dev
```

Access the dashboard at `http://localhost:5173`

### Build

```bash
npm run build
```

---

## 📁 Project Structure

```
frontend/
├── src/
│   ├── api/              # API client (axios)
│   ├── components/       # React components
│   │   ├── AlertsManager.tsx
│   │   ├── Auth.tsx
│   │   ├── Chart.tsx
│   │   ├── Header.tsx
│   │   ├── PortfolioReport.tsx
│   │   ├── SignalFeed.tsx
│   │   ├── StrategyMarketplace.tsx
│   │   └── TradingPanel.tsx
│   ├── constants/        # App constants
│   ├── hooks/            # Custom React hooks
│   ├── store/            # State management
│   ├── types/            # TypeScript types
│   ├── utils/            # Utility functions
│   ├── App.tsx           # Main app component
│   └── main.tsx          # Entry point
├── public/
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── tailwind.config.js
```

---

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the frontend directory:

```env
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
```

---

## 📖 Usage

### Connect to Backend

1. Ensure the backend is running on `http://localhost:8080`
2. Register a new account or login
3. Access paper trading features

### Market Data

- View real-time prices in the chart
- Subscribe to multiple symbols
- Set up price alerts

### Paper Trading

- Create buy/sell orders
- Monitor open positions
- Track portfolio performance

---

## 🔨 Development Commands

| Command | Description |
|---------|-------------|
| `npm run dev` | Start development server |
| `npm run build` | Build for production |
| `npm run lint` | Run ESLint |
| `npm run preview` | Preview production build |

---

## 📄 License

Distributed under the MIT License. See [LICENSE](../LICENSE) for more information.

---

*Quant-Trader - Empowering your strategy with professional infrastructure.*
