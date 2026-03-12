export interface ApiError {
  message: string;
  status: number;
  code?: string;
}

export interface ErrorResponse {
  error: string;
}

export const ErrorCodes = {
  UNAUTHORIZED: 'UNAUTHORIZED',
  FORBIDDEN: 'FORBIDDEN',
  NOT_FOUND: 'NOT_FOUND',
  CONFLICT: 'CONFLICT',
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  INTERNAL_ERROR: 'INTERNAL_ERROR',
} as const;

export type ErrorCode = (typeof ErrorCodes)[keyof typeof ErrorCodes];

export class AppError extends Error implements ApiError {
  message: string;
  status: number;
  code?: string;

  constructor(message: string, status: number, code?: string) {
    super(message);
    this.message = message;
    this.status = status;
    this.code = code;
    this.name = 'AppError';
  }
}

export class ValidationError extends AppError {
  constructor(message: string) {
    super(message, 400, ErrorCodes.VALIDATION_ERROR);
    this.name = 'ValidationError';
  }
}

export class UnauthorizedError extends AppError {
  constructor(message: string = 'Unauthorized') {
    super(message, 401, ErrorCodes.UNAUTHORIZED);
    this.name = 'UnauthorizedError';
  }
}

export class ForbiddenError extends AppError {
  constructor(message: string = 'Forbidden') {
    super(message, 403, ErrorCodes.FORBIDDEN);
    this.name = 'ForbiddenError';
  }
}

export class NotFoundError extends AppError {
  constructor(message: string = 'Not found') {
    super(message, 404, ErrorCodes.NOT_FOUND);
    this.name = 'NotFoundError';
  }
}

export class ConflictError extends AppError {
  constructor(message: string) {
    super(message, 409, ErrorCodes.CONFLICT);
    this.name = 'ConflictError';
  }
}

export class InternalError extends AppError {
  constructor(message: string = 'Internal server error') {
    super(message, 500, ErrorCodes.INTERNAL_ERROR);
    this.name = 'InternalError';
  }
}
