// To parse this JSON data, do
//
//     final logModel = logModelFromJson(jsonString);

import 'dart:convert';

LogModel logModelFromJson(String str) => LogModel.fromJson(json.decode(str));

String logModelToJson(LogModel data) => json.encode(data.toJson());

class LogModel {
  LogModel({
    this.error,
    this.logs,
    this.message,
  });

  bool error;
  List<Log> logs;
  String message;

  factory LogModel.fromJson(Map<String, dynamic> json) => LogModel(
        error: json["error"],
        logs: List<Log>.from(json["logs"].map((x) => Log.fromJson(x))),
        message: json["message"],
      );

  Map<String, dynamic> toJson() => {
        "error": error,
        "logs": List<dynamic>.from(logs.map((x) => x.toJson())),
        "message": message,
      };
}

class Log {
  Log({
    this.logId,
    this.host,
    this.username,
    this.status,
    this.ts,
    this.raw,
    this.machineId,
  });

  String logId;
  String host;
  String username;
  String status;
  int ts;
  String raw;
  String machineId;

  factory Log.fromJson(Map<String, dynamic> json) => Log(
        logId: json["log_id"],
        host: json["host"],
        username: json["username"],
        status: json["status"],
        ts: json["ts"],
        raw: json["raw"],
        machineId: json["machine_id"],
      );

  Map<String, dynamic> toJson() => {
        "log_id": logId,
        "host": host,
        "username": username,
        "status": status,
        "ts": ts,
        "raw": raw,
        "machine_id": machineId,
      };
}
