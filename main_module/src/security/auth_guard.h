#pragma once
#include <crow.h>
#include "jwt.h"

// Extracts and validates the JWT from the request.
// Returns 200 on success, 401 on invalid token, 403 if user is blocked.
inline int authGuard(
    const crow::request& req,
    UserContext& ctx
) {
    try {
        ctx = parseAndVerifyJWT(req);

        if (ctx.blocked) {
            return 403;
        }

        return 200;
    } catch (const std::exception&) {
        return 401;
    }
}