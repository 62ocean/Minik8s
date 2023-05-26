//
// Created by 12609 on 2023/2/4.
//

#include "kv_http_server.h"
#include <stdio.h>

using namespace httplib;
using namespace std;

int main() {
    KVHttpServer server;
    server.init();
}

void KVHttpServer::init() {
    server.Get("/", Get_Test);
    // server.Get("/getLocal",Get_GetLocal);
    // server.Get("/getGlobal",Get_GetGlobal);
    // server.Post("/set",Post_Set);
    // server.Delete("/delete",Delete_Delete);
    // server.Get("/",Get_Color);
    // server.Get("/getDst",Get_Dst);
    server.listen(ip, port);
}
void KVHttpServer::Get_Test(const httplib::Request& req, httplib::Response& resp) {
    printf("get request!\n");
    system("head -1 /proc/self/cgroup|cut -d/ -f3 > ip.txt");
    FILE *f = fopen("./ip.txt","r");
    char s[20];
    fgets(s, 10, f);
    fclose(f);
    resp.set_content(s, "text/plain");
    
}

// void KVHttpServer::Get_Color(const httplib::Request& req, httplib::Response& resp){
//     string key = req.get_param_value("id");
//     string value;
//     printf("receive message id: %s\n",key.data());
//     if(key=="000000003"){
//         value = "red";
//     }else if(key=="000000004"){
//         value = "blue";
//     }else if(key=="000000020"){
//         value = "green";
//     }else{
//         value = "invalid id";
//     }
//     resp.set_content(value, "text/plain");
// }

// void KVHttpServer::Get_Dst(const httplib::Request& req, httplib::Response& resp){
//     string key = req.get_param_value("id");
//     string value;
//     value = db.readDataFromDatabase(key);
//     int pos = value.find(",");
//     string ret = value.substr(0,pos+1);
//     resp.set_content(ret, "text/plain");
// }

// void KVHttpServer::Get_GetLocal(const httplib::Request& req, httplib::Response& resp){
//     string key = req.get_param_value("key");
//     string value;
//     // LrcKV.GetLocal(key,value);
//     resp.set_content(value, "text/plain");
// }
// void KVHttpServer::Get_GetGlobal(const httplib::Request& req, httplib::Response& resp){
//     // TODO: 等实现Master之后
// }
// void KVHttpServer::Post_Set(const httplib::Request& req, httplib::Response& resp){
//     string key = req.get_param_value("key");
//     string value = req.get_param_value("value");
//     // LrcKV.Put(key,value);
//     resp.set_content("ok\n", "text/plain");
// }
// void KVHttpServer::Delete_Delete(const httplib::Request& req, httplib::Response& resp){
//     string key = req.get_param_value("key");
//     // LrcKV.Delete(key);
//     resp.set_content("ok\n", "text/plain");
// }
