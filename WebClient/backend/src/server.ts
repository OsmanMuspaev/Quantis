import express from "express";
import cookieParser from "cookie-parser";
import { connectRedis } from "./redis/redis.client";
import { sessionMiddleware } from "./middlewares/session.middleware";
import loginCheckRoute from "./routes/login.check.route";
import loginStartRoute from "./routes/login.start.route";
import loginVerifyRoute from "./routes/login.verify.route";
import logoutRoute from "./routes/logout.route";
import authRefreshRoute from "./routes/auth.refresh.route";

async function bootstrap() {
  await connectRedis();

  const app = express();
  app.use(express.json());
  app.use(cookieParser());
  app.use(sessionMiddleware);

  app.use(loginCheckRoute);
  app.use(loginStartRoute);
  app.use(loginVerifyRoute);
  app.use(logoutRoute);
  app.use(authRefreshRoute);

  app.get("/health", (_, res) => {
    res.send("OK");
  });

  app.listen(8080, () => {
    console.log("WebClient backend running on :8080");
  });
}

bootstrap();
