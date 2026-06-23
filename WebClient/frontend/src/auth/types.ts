export type AuthState =
  | "unknown"
  | "anonymous"
  | "need_verify"
  | "authorized";

export type AuthStatusResponse = {
  status: "pending" | "access_denied" | "approved" | "anonymous";
};

export type VerifyResponse = {
  status: "authorized";
};

export type ApiError = {
  statusCode: number;
  message?: string;
};
