import pymysql
import pymysql.cursors
import time

from argon2 import PasswordHasher

from dbutils.persistent_db import PersistentDB
from config.config import get_config

from sessions import session 


config = get_config()

class Database:
    def __init__(self):
        self.db = None
        self.connect()
        self.password_hasher = PasswordHasher()

    def connect(self):
        try:
            self.db = PersistentDB(
                creator=pymysql,
                host=config["database"]["host"],
                user=config["database"]["username"],
                password=config["database"]["password"],
                database=config["database"]["database"],
                autocommit=True,
                cursorclass=pymysql.cursors.DictCursor,
            )
        except Exception as c:
            print(c)
            print("Could not connect to MySQL server. Check details!")
            exit(1)

    def update_hashed_password(self, username, password):
        with self.db.connection() as conn:
            cursor = conn.cursor()
            cursor.execute("UPDATE users set password = %s where username = %s", (password, username))
    
    def login(self, username, password):
        with self.db.connection() as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT username, password, permission FROM users where username = %s", (username,))
            user = cursor.fetchone()

            if user is None:
                return {"message": "Invalid login details!", "status": False}
            else:
                user = dict(user)
                hashed_password = user["password"]

                try:
                    if self.password_hasher.verify(hashed_password, password):
                        ses_object, key = session.create_session(username, bool(user["permission"]))
                        return {"message": "You have successfully logged in!", "session": key}
                except Exception as e:
                    return {"message": "Invalid login details!", "status": False}
    
    def get_logs(self) -> dict:
        with self.db.connection() as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT * FROM logs")
            logs = cursor.fetchall()
            

        if logs is None:
            return {"status": False, "message": "No logs have been found"}

        return {"status": True, "message": "Fetched logs", "logs": logs}