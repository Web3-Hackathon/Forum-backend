import random
import string 
import time



session_object = {}

def create_session(user, permission):
    random_string = ''.join(random.choices(string.ascii_lowercase + string.digits, k=16))
    new_session =  {"exp": int(time.time() + 3600), "is_admin": permission, "username": user}
    session_object[random_string] =  new_session
    print(session_object)
    return new_session, random_string


def check_session(session_key):
    session = session_object.get(session_key, None)

    if session is None:
        return {"message": "Invalid Session!"}, False 
    else:
        if int(time.time()) > session["exp"]:
            return {"message": "Invalid Session!"}, False 
        else:
            return session, True 
