import { Router, Request, Response } from "express";
import { SessionRepository } from "../redis/session.repository";
import { generateToken } from "../utils/token";
import { SESSION_COOKIE_NAME } from "../config/cookies";

const router = Router();

router.get("/login", async (req: Request, res: Response) => {
  const type = req.query.type as string | undefined;

  if (!type || !["github", "yandex", "code"].includes(type)) {
    return res.status(400).send("Invalid login type");
  }

  let sessionId = req.cookies?.[SESSION_COOKIE_NAME];

  if (!sessionId) {
    sessionId = generateToken(16);

    res.cookie(SESSION_COOKIE_NAME, sessionId, {
      httpOnly: true,
      sameSite: "lax",
      path: "/",
    });
  }

  const entryToken = generateToken(16);

  await SessionRepository.set(sessionId, {
    status: "pending",
    entryToken,
  });

  if (type === "code") {
    return res.json({ code: entryToken });
  }

  return res.redirect("/verify");
});

export default router;
