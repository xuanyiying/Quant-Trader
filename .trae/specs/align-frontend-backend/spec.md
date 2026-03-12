# Frontend-Backend API Alignment Spec

## Why
The frontend and backend code have inconsistencies in API endpoints, request/response formats, data structures, and error handling. This misalignment causes potential runtime errors and poor user experience. Aligning them ensures seamless integration and reliable data exchange.

## What Changes
- Add missing API endpoints in frontend that exist in backend
- Fix data type mismatches between frontend TypeScript interfaces and backend Go structs
- Align error handling patterns between frontend and backend
- Add missing request/response type definitions
- Update API response handling to match backend response format

## Impact
- Affected specs: API Integration, Data Types, Error Handling
- Affected code: 
  - Frontend: `src/types/data.ts`, `src/types/market.ts`, `src/api/axios.ts`
  - Backend: All `api/*.go` handlers

## ADDED Requirements

### Requirement: Complete API Endpoint Coverage
The frontend SHALL implement API calls for all backend endpoints.

#### Backend Endpoints to Align:
| Endpoint | Method | Frontend Status |
|----------|--------|-----------------|
| `/api/v1/auth/register` | POST | Missing |
| `/api/v1/auth/login` | POST | Missing |
| `/api/v1/paper/account` | GET | Partial |
| `/api/v1/paper/account/reset` | POST | Missing |
| `/api/v1/paper/orders` | POST | Missing |
| `/api/v1/paper/positions` | GET | Partial |
| `/api/v1/portfolio/report` | GET | Partial |
| `/api/v1/alerts` | GET/POST/DELETE | Partial |
| `/api/v1/subscription` | GET | Partial |
| `/api/v1/subscription/checkout` | POST | Missing |
| `/api/v1/klines/latest` | GET | Missing |
| `/api/v1/klines/backfill` | POST | Missing |
| `/api/v1/marketplace` | GET | Partial |
| `/api/v1/marketplace/:id/purchase` | POST | Partial |

### Requirement: Data Type Alignment
The frontend TypeScript interfaces SHALL match backend Go struct field names and types.

#### Scenario: PaperAccount alignment
- **GIVEN** backend returns `{"balance": "100000.00", "initial_balance": "100000.00"}`
- **WHEN** frontend receives the response
- **THEN** TypeScript interface should have matching fields

#### Current Mismatches to Fix:

**PaperAccount (backend model/paper.go):**
```go
type PaperAccount struct {
    ID            int64
    UserID        int64
    Balance       decimal.Decimal
    InitialBalance decimal.Decimal
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```
**Frontend missing:** `initial_balance`, `created_at`, `updated_at`

**Position (backend model/paper.go):**
```go
type PaperPosition struct {
    ID        int64
    AccountID int64
    Symbol    string
    Side      string
    Qty       decimal.Decimal
    AvgPrice  decimal.Decimal
}
```
**Frontend missing:** `id`, `account_id`, `side`

**Alert (backend model/alert.go):**
```go
type Alert struct {
    ID            int64
    UserID        int64
    Symbol        string
    ConditionType string
    TargetValue   decimal.Decimal
    Status        string
    CreatedAt     time.Time
    TriggeredAt   *time.Time
}
```
**Frontend missing:** `status`, `created_at`, `triggered_at`

**Subscription (backend api/subscription.go response):**
```go
// Response fields:
// tier_name, max_symbols, status, expires_at
```
**Frontend missing:** `max_symbols`, `status`, `expires_at`

### Requirement: Error Response Format Alignment
The frontend SHALL handle all backend error response formats consistently.

#### Backend Error Response Format:
```json
{"error": "error message"}
```

#### HTTP Status Codes to Handle:
| Status Code | Backend Usage | Frontend Handling |
|-------------|---------------|-------------------|
| 400 | Bad Request | Partial |
| 401 | Unauthorized | Implemented |
| 403 | Forbidden | Missing |
| 404 | Not Found | Missing |
| 409 | Conflict (email exists) | Missing |
| 500 | Internal Error | Partial |

### Requirement: API Service Layer
The frontend SHALL have a centralized API service layer with typed methods.

#### Scenario: API service structure
- **WHEN** frontend needs to make API calls
- **THEN** use centralized service methods with proper typing

## MODIFIED Requirements

### Requirement: TypeScript Interface Updates
Update existing interfaces to match backend exactly.

**Strategy interface:**
```typescript
// Current
export interface Strategy {
  id: number;
  name: string;
  description: string;
  price: number;
  author: string;
  is_subscribed: boolean;
  type?: string;
  performance_metrics?: unknown;
  subscriber_count?: number;
}

// Should match backend response:
export interface Strategy {
  id: number;
  name: string;
  description: string;
  price: string;  // decimal from backend
  author: string;
  metrics: Record<string, unknown>;  // JSON from backend
}
```

**PortfolioReport interface:**
```typescript
// Should match backend model.BacktestReport
export interface PortfolioReport {
  final_balance: string;
  total_return: number;
  sharpe_ratio: number;
  max_drawdown: number;
  win_rate: number;
  total_trades: number;
  equity_curve: EquityPoint[];
}

export interface EquityPoint {
  timestamp: string;
  equity: string;
  drawdown: number;
}
```

## REMOVED Requirements

### Requirement: None
No features are being removed, only aligned or added.
