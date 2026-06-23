#pragma once
#include "user_context.h"
#include <crow.h>

// Parses and verifies JWT from the Authorization header.
// Throws std::runtime_error on missing/invalid/expired token.
UserContext parseAndVerifyJWT(const crow::request& req);