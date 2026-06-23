import { Router, Request, Response } from "express";
import { SessionRepository } from "../redis/session.repository";
import { SESSION_COOKIE_NAME } from "../config/cookies";

const router = Router();

router.post("/logout", async (req: Request, res: Response) => {
  const sessionId = req.cookies?.[SESSION_COOKIE_NAME];

  if (sessionId) {
    await SessionRepository.delete(sessionId);
  }

  res.clearCookie(SESSION_COOKIE_NAME, {
    httpOnly: true,
    sameSite: "lax",
    path: "/",
  });

  return res.json({ status: "ok" });
});

export default router;
