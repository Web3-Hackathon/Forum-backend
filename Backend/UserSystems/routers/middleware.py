from fastapi import Request
from sessions import session 

def middleware(request: Request, needs_admin=False):
    session_key = request.headers.get("Authorization")
    print(session_key)
    if session_key is None:
        return {"message": "Invalid Authorization header!", "status": False}, False
    else:
        if len(session_key) < 1:
            return {"message": "Invalid Authorization header!", "status": False}, False
        else:
            session_message, session_bool = session.check_session(session_key)

            if session_bool is False:
                return session_message, False 
            
            if needs_admin is True:
                if not (session_message["is_admin"]):
                    return {"message": "Unauthorized request!", "status": False}, False 
            
            return {"message": "Check Passed", "status": True}, True 