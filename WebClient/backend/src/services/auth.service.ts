const AUTH_URL = process.env.AUTH_SERVICE_URL || "http://auth:8081";

type AuthVerifyResponse =
  | {
      status: "approved";
      access_token: string;
      refresh_token: string;
      user_id: string;
    }
  | {
      status: "access_denied";
    }
  | {
      status: "pending";
    };

type AuthRefreshResponse = {
  access_token: string;
  refresh_token: string;
};

export async function verifyWithAuth(
  entryToken: string,
  code: string
): Promise<AuthVerifyResponse> {
  const res = await fetch(`${AUTH_URL}/login/verify`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      entry_token: entryToken,
      code,
    }),
  });

  if (!res.ok) {
    throw new Error("Auth service error");
  }

  return res.json();
}

export async function refreshWithAuth(
  refreshToken: string
): Promise<AuthRefreshResponse> {
  const res = await fetch(`${AUTH_URL}/refresh`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      refresh_token: refreshToken,
    }),
  });

  if (!res.ok) {
    throw new Error("Auth refresh failed");
  }

  return res.json();
}
