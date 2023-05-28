//
// Created by 12609 on 2023/2/4.
//

#ifndef LRCPHXKV_KV_HTTP_SERVER_H
#define LRCPHXKV_KV_HTTP_SERVER_H

#include <iostream>
#include <stdio.h>

#include "httplib.h"
// #include "../fdb/fdb.h"
// #include "include/nlohmann/json.hpp"

#define IP "142.16.1.11"
#define PORT 8181

// using json = nlohmann::json;

class KVHttpServer {
   private:
    httplib::Server server;
    std::string ip;
    int port;

    // static void Get_GetLocal(const httplib::Request& req, httplib::Response& resp);
    // static void Get_GetGlobal(const httplib::Request& req, httplib::Response& resp);
    // static void Post_Set(const httplib::Request& req, httplib::Response& resp);
    // static void Delete_Delete(const httplib::Request& req, httplib::Response& resp);
    // static void Get_Color(const httplib::Request& req, httplib::Response& resp);
    // static void Get_Dst(const httplib::Request& req, httplib::Response& resp);
    static void Get_Test(const httplib::Request& req, httplib::Response& resp);

   public:
    KVHttpServer(std::string ip, int port) : ip(ip), port(port) {}
    KVHttpServer() {
        ip = IP;
        port = PORT;
    }
    ~KVHttpServer() {}
    void init();
    // static FDB db;
};

#endif  // LRCPHXKV_KV_HTTP_SERVER_H
