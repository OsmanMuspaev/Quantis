#pragma once
#include <string>
#include <vector>

struct Test {
    int id = 0;
    int course_id = 0;
    std::string title;
    bool is_active = false;
    bool is_deleted = false;
};

struct UserAnswer {
    std::string question_text;
    std::string user_answer_text;
};

struct AttemptDetails {
    std::string user_id;
    std::vector<UserAnswer> answers;
};
