#include "db.h"

void DB::pushNotification(
    const std::string& userId, 
    const std::string& type, 
    const std::string& title, 
    const std::string& message, 
    const crow::json::wvalue& payload
) {
    ensureConnection();

    std::string payloadStr = payload.dump();
    if (payloadStr == "null") payloadStr = "{}";

    const char* paramValues[5];
    paramValues[0] = userId.c_str();
    paramValues[1] = type.c_str();
    paramValues[2] = title.c_str();
    paramValues[3] = message.c_str();
    paramValues[4] = payloadStr.c_str();

    const char* query = 
        "INSERT INTO notifications (user_id, type, title, message, payload) "
        "VALUES ($1, $2, $3, $4, $5::jsonb)";

    PGresult* res = PQexecParams(
        static_cast<PGconn*>(conn),
        query,
        5,
        NULL,
        paramValues,
        NULL,
        NULL,
        0 
    );

    if (PQresultStatus(res) != PGRES_COMMAND_OK) {
        std::cerr << "pushNotification error: " << PQerrorMessage(static_cast<PGconn*>(conn)) << "\n";
    }

    PQclear(res);
}

// Soft-deletes notifications by marking them as sent.
void DB::markNotificationsAsSent(const std::vector<int>& ids, std::string userId) {
    if (ids.empty()) return;
    ensureConnection();

    std::string pg_array = "{";
    for (size_t i = 0; i < ids.size(); ++i) {
        pg_array += std::to_string(ids[i]);
        if (i < ids.size() - 1) pg_array += ",";
    }
    pg_array += "}";

    const char* paramValues[2];
    paramValues[0] = pg_array.c_str();
    paramValues[1] = userId.c_str();

    const char* sql = "UPDATE notifications SET is_sent_tg = TRUE WHERE id = ANY($1::int[]) AND user_id = $2";

    PGresult* res = PQexecParams(static_cast<PGconn*>(conn), sql, 2, nullptr, paramValues, nullptr, nullptr, 0);

    if (PQresultStatus(res) != PGRES_COMMAND_OK) {
        CROW_LOG_ERROR << "markNotificationsAsSent error: " << PQerrorMessage(static_cast<PGconn*>(conn));
    }

    PQclear(res);
}

// Fetches all unsent notifications for a user.
std::vector<crow::json::wvalue> DB::getUnsentNotifications(std::string userId) {
    ensureConnection();
    std::vector<crow::json::wvalue> notifications;
    
    const char* paramValues[1];
    paramValues[0] = userId.c_str();

    const char* sql = "SELECT id, type, title, message, payload FROM notifications "
                      "WHERE user_id = $1 AND is_sent_tg = FALSE";

    PGresult* res = PQexecParams(static_cast<PGconn*>(conn), sql, 1, nullptr, paramValues, nullptr, nullptr, 0);

    if (PQresultStatus(res) == PGRES_TUPLES_OK) {
        int rows = PQntuples(res);
        for (int i = 0; i < rows; i++) {
            crow::json::wvalue n;
            n["id"] = std::stoi(PQgetvalue(res, i, 0));
            n["type"] = PQgetvalue(res, i, 1);
            n["title"] = PQgetvalue(res, i, 2);
            n["message"] = PQgetvalue(res, i, 3);
            
            std::string payloadStr = PQgetvalue(res, i, 4);
            n["payload"] = crow::json::load(payloadStr);
            
            notifications.push_back(std::move(n));
        }
    }
    PQclear(res);
    return notifications;
}
