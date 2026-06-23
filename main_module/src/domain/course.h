#pragma once
#include <string>

struct Course {
    int id = 0;
    std::string title;
    std::string description;
    std::string author_id;
    bool is_deleted = false;
};
