import type { AppError } from "./types.js";

export function errorResponse(code: string, message: string): AppError {
  return {
    state: "error",
    data: null,
    error: {
      code,
      message,
    },
    meta: {
      timestamp: new Date().toISOString(),
    },
  };
}
