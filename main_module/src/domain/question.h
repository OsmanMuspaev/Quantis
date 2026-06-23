#pragma once
#include <string>
#include <vector>
#include <crow.h>

// Represents a single versioned question.
// Questions use append-only versioning: "updating" a question creates a new version.
struct Question {
    int id = 0;
    int version = 1;
    std::string author_id;
    std::string title;
    std::string content;
    std::vector<std::string> options;
    int correct_option = 0;
    bool is_deleted = false;

    // Serializes for the author/admin view (includes correct_option and is_deleted).
    crow::json::wvalue to_json() const {
        crow::json::wvalue j;
        j["id"] = id;
        j["version"] = version;
        j["author_id"] = author_id;
        j["title"] = title;
        j["content"] = content;
        
        crow::json::wvalue::list optList;
        for (const auto& opt : options) {
            optList.push_back(crow::json::wvalue(opt));
        }
        j["options"] = std::move(optList);
        
        j["correct_option"] = correct_option;
        j["is_deleted"] = is_deleted;
        return j;
    }

    // Serializes for the student view (hides correct_option and is_deleted).
    crow::json::wvalue to_student_json() const {
        crow::json::wvalue j;
        j["id"] = id;
        j["version"] = version;
        j["title"] = title;
        j["content"] = content;
        
        crow::json::wvalue::list optList;
        for (const auto& opt : options) {
            optList.push_back(crow::json::wvalue(opt));
        }
        j["options"] = std::move(optList);
        
        return j;
    }
};
