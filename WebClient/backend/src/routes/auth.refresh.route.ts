import { Router, Request, Response } from "express";
import { SessionRepository } from "../redis/session.repository";
import { refreshWithAuth } from "../services/auth.service";
import { SESSION_COOKIE_NAME } from "../config/cookies";

const router = Router();

router.post("/auth/refresh", async (req: Request, res: Response) => {
  try {
    const sessionId = req.cookies?.[SESSION_COOKIE_NAME];

    if (!sessionId) {
      return res.status(401).json({ error: "No session" });
    }

    const session = await SessionRepository.get(sessionId);

    if (!session || session.status !== "authorized" || !session.refreshToken) {
      return res.status(401).json({ error: "Not authorized" });
    }

    const refreshResponse = await refreshWithAuth(session.refreshToken);

    await SessionRepository.set(sessionId, {
      status: "authorized",
      accessToken: refreshResponse.access_token,
      refreshToken: refreshResponse.refresh_token,
      userId: session.userId,
    });

    return res.json({ status: "ok" });
  } catch (error) {
    console.error("Auth refresh failed:", error);
    return res.status(500).json({ error: "Refresh failed" });
  }
});

export default router;
