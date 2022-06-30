import pymysql
import pymysql.cursors
import time

from dbutils.persistent_db import PersistentDB
from config.config import get_config

from sessions import session
import datetime

config = get_config()


class Database_Forum:
    def __init__(self):
        self.db = None
        self.connect()

    def connect(self):
        try:
            self.db = PersistentDB(
                creator=pymysql,
                host=config["database"]["host"],
                user=config["database"]["username"],
                password=config["database"]["password"],
                database="crypto_forum",
                autocommit=True,
                cursorclass=pymysql.cursors.DictCursor,
            )
        except Exception as c:
            print(c)
            print("Could not connect to MySQL server. Check details!")
            exit(1)
    
    def get_users(self) -> dict:
        with self.db.connection() as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT * FROM users")
            users = cursor.fetchall()
            

        if users is None:
            return {"status": False, "message": "No users have been found"}

        return {"status": True, "message": "Fetched users", "users": users}
    
    def get_user_by_id(self, user_id):
        with self.db.connection() as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT * FROM users where id = %s", (user_id,))
            user = cursor.fetchone()
        
        if user is None:
            return None 
        else:
            return "Found"

    def ban_user(self, user_id, ban_reason, ban_expiry):
        if self.get_user_by_id(user_id) is None:
            return {"status": False, "message": "User has not been found"}
        else:
            ban_expiry = datetime.datetime.fromtimestamp(ban_expiry).strftime('%Y-%m-%d %H:%M:%S')
            with self.db.connection() as conn:
                cursor = conn.cursor()
                cursor.execute("UPDATE users SET banned = %s, ban_reason = %s, ban_expiry = %s WHERE id = %s", (1, ban_reason, ban_expiry, user_id,))
                conn.commit()
            
            return {"message": "Successfully banned the user!", "status": True}
    
    def unban_user(self, user_id):
        if self.get_user_by_id(user_id) is None:
            return {"status": False, "message": "User has not been found"}
        else:
            with self.db.connection() as conn:
                cursor = conn.cursor()
                cursor.execute("UPDATE users SET banned = %s, ban_reason = %s WHERE id = %s", (0, "", user_id))
                conn.commit()
            
            return {"message": "Successfully unbanned the user!", "status": True}
    
    def mute_user(self, user_id, mute_reason, mute_expiry):
        if self.get_user_by_id(user_id) is None:
            return {"status": False, "message": "User has not been found"}
        else:
            mute_expiry = datetime.datetime.fromtimestamp(mute_expiry).strftime('%Y-%m-%d %H:%M:%S')
            with self.db.connection() as conn:
                cursor = conn.cursor()
                cursor.execute("UPDATE users SET muted = %s, mute_reason = %s, mute_expiry = %s WHERE id = %s", (1, mute_reason, mute_expiry, user_id,))
                conn.commit()
            
            return {"message": "Successfully muted the user!", "status": True}
    
    def unmute_user(self, user_id):
        if self.get_user_by_id(user_id) is None:
            return {"status": False, "message": "User has not been found"}
        else:
            with self.db.connection() as conn:
                cursor = conn.cursor()
                cursor.execute("UPDATE users SET muted = %s, mute_reason = %s WHERE id = %s", (0, "", user_id))
                conn.commit()
            
            return {"message": "Successfully unmuted the user!", "status": True}

    def get_threads(self):
        with self.db.connection() as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT * FROM threads")
            threads = cursor.fetchall()
            

        if threads is None:
            return {"status": False, "message": "No threads have been found"}

        return {"status": True, "message": "Fetched threads", "threads": threads}
    
    def delete_thread(self, thread_id):
        with self.db.connection() as conn:
            cursor = conn.cursor()
            cursor.execute("UPDATE threads SET hidden = %s WHERE id = %s", (1,thread_id))
            conn.commit()
        
        return {"message": "Successfully deleted a thread", "status": True}
    
    def ip_ban(self, ip_address, ban_expiry, admin_session, ban_reason):
        by_who, status = session.check_session(session_key=admin_session)
        timestamp = datetime.datetime.fromtimestamp(ban_expiry).strftime('%Y-%m-%d %H:%M:%S')
        with self.db.connection() as conn:
            cursor = conn.cursor()
            cursor.execute("INSERT INTO ip_bans (ip_address, ban_expiry, banned_by, ban_reason) VALUES (%s, %s, %s, %s)", (ip_address, timestamp, by_who["username"], ban_reason))
            conn.commit()
        
        return {"message": "Successfully banned an IP address", "status": True}

    def ip_unban(self, ban_id):
        with self.db.connection() as conn:
            cursor = conn.cursor()
            cursor.execute("DELETE FROM ip_bans WHERE id=%s", (ban_id,))
            conn.commit()
        
        return {"message": "Successfully unbanned an IP address", "status": True}
    
    def get_ip_bans(self):
        with self.db.connection() as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT * FROM ip_bans")
            ip_bans = cursor.fetchall()
            

        if ip_bans is None:
            return {"status": False, "message": "No ip_bans have been found"}

        return {"status": True, "message": "Fetched ip_bans", "ip_bans": ip_bans}
    

    def get_thread_sections(self):
        with self.db.connection() as conn:
            cursor = conn.cursor()
            cursor.execute("SELECT * FROM thread_sections")
            thread_sections = cursor.fetchall()
            

        if thread_sections is None:
            return {"status": False, "message": "No thread sections have been found"}

        return {"status": True, "message": "Fetched thread sections", "thread_sections": thread_sections}
    
    def create_section(self, category, parent, name):
        with self.db.connection() as conn:
            cursor = conn.cursor()
            cursor.execute("INSERT INTO thread_sections (category, parent, name) VALUES (%s, %s, %s)", (category, parent, name))
            conn.commit()
        
        return {"message": "Successfully created a thread section", "status": True}