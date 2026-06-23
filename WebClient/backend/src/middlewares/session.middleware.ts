import { Request, Response, NextFunction } from "express";

declare global {
  namespace Express {
    interface Request {
      sessionId?: string | null;
    }
  }
}

export function sessionMiddleware(
  req: Request,
  _: Response,
  next: NextFunction
) {
  req.sessionId = req.cookies?.session_id || null;
  next();
}
