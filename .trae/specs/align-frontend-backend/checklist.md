# Checklist

## Type Definitions
- [x] PaperAccount interface matches backend PaperAccount struct fields
- [x] Position interface matches backend PaperPosition struct fields
- [x] Alert interface matches backend Alert struct fields
- [x] Subscription interface matches backend subscription response fields
- [x] Strategy interface matches backend marketplace response fields
- [x] PortfolioReport interface matches backend BacktestReport struct fields
- [x] All decimal types are handled as strings in TypeScript

## API Services
- [x] Auth API service has register() method matching POST /api/v1/auth/register
- [x] Auth API service has login() method matching POST /api/v1/auth/login
- [x] Paper API service has getAccount() method matching GET /api/v1/paper/account
- [x] Paper API service has resetAccount() method matching POST /api/v1/paper/account/reset
- [x] Paper API service has createOrder() method matching POST /api/v1/paper/orders
- [x] Paper API service has getPositions() method matching GET /api/v1/paper/positions
- [x] Alert API service has CRUD methods matching /api/v1/alerts endpoints
- [x] Subscription API service has getSubscription() method
- [x] Subscription API service has createCheckoutSession() method
- [x] Kline API service has getLatest() method
- [x] Kline API service has triggerBackfill() method
- [x] Marketplace API service has listStrategies() method
- [x] Marketplace API service has purchaseStrategy() method

## Error Handling
- [x] Error response type defined with `error: string` field
- [x] HTTP 400 errors handled with error message display
- [x] HTTP 401 errors trigger token removal and redirect
- [x] HTTP 403 errors handled appropriately
- [x] HTTP 404 errors handled with user-friendly message
- [x] HTTP 409 (conflict) errors handled for registration
- [x] HTTP 500 errors handled with generic error message

## Integration
- [x] All existing components updated to use new API services
- [x] All existing stores updated to use new API services
- [x] No direct axios calls remain in components
- [x] All API calls have proper TypeScript types
- [x] Build succeeds without type errors
