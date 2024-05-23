import 'package:path_provider/path_provider.dart';
import 'package:sqflite_common_ffi/sqflite_ffi.dart';
import 'dart:io';
import 'package:path/path.dart' as path;

class SCNDatabase {
  static SCNDatabase? instance = null;

  final Database _db;

  SCNDatabase._(this._db) {}

  static create() async {
    var docPath = await getApplicationDocumentsDirectory();
    var dbpath = path.join(docPath.absolute.path, 'scn.db');

    if (Platform.isWindows || Platform.isLinux) {
      sqfliteFfiInit();
    }

    var db = await databaseFactoryFfi.openDatabase(dbpath,
        options: OpenDatabaseOptions(
          version: 1,
          onCreate: (db, version) async {
            initDatabase(db);
          },
          onUpgrade: (db, oldVersion, newVersion) async {
            upgradeDatabase(db, oldVersion, newVersion);
          },
        ));

    return instance = SCNDatabase._(db);
  }

  static void initDatabase(Database db) async {
    await db.execute('CREATE TABLE requests (id INTEGER PRIMARY KEY, timestamp DATETIME, name TEXT, url TEXT, response_code INTEGER, response TEXT, status TEXT)');

    await db.execute('CREATE TABLE logs (id INTEGER PRIMARY KEY, timestamp DATETIME, level TEXT, text TEXT, additional TEXT)');

    await db.execute('CREATE TABLE messages (message_id INTEGER PRIMARY KEY, receive_timestamp DATETIME, channel_id TEXT, timestamp TEXT, data JSON)');
  }

  static void upgradeDatabase(Database db, int oldVersion, int newVersion) {
    // ...
  }
}
