#pragma once
#include <string>
#include <vector>
#include <unordered_set>

struct UserContext {
    std::string userId;
    bool blocked = false;
    std::unordered_set<std::string> permissions;
    std::unordered_set<std::string> roles;
};