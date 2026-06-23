#include "db.h"

DB::DB(const std::string& conninfo) : conninfo(conninfo), conn(nullptr) {
    ensureConnection();
}

DB::~DB() {
    if (conn) {
        PQfinish(static_cast<PGconn*>(conn));
    }
}

// Ensures the PostgreSQL connection is alive. Reconnects if the previous connection was lost.
void DB::ensureConnection() {
    if (!conn || PQstatus(static_cast<PGconn*>(conn)) != CONNECTION_OK) {
        if (conn) PQfinish(static_cast<PGconn*>(conn));
        
        conn = PQconnectdb(conninfo.c_str());
        
        if (PQstatus(static_cast<PGconn*>(conn)) != CONNECTION_OK) {
            std::string errMsg = PQerrorMessage(static_cast<PGconn*>(conn));
            PQfinish(static_cast<PGconn*>(conn));
            conn = nullptr;
            throw std::runtime_error("Failed to connect to database: " + errMsg);
        }
    }
}
