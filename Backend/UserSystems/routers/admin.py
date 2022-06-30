from fastapi import APIRouter, HTTPException, Request
from routers.middleware import middleware
from models.admin import Authenticate, BanUser, DeleteThread, BanIP, UnbanIP, CreateThreadSection
from database.database import Database
from database.database_forum import Database_Forum

db = Database()
db_forum = Database_Forum()

router = APIRouter(
    prefix="/admin",
    tags=["public"],
)


@router.post("/authenticate")
async def authenticate(auth: Authenticate):
    return db.login(auth.username, auth.password)


@router.get("/session")
async def check_session(request: Request):
    message, status = middleware(request, needs_admin=False)

    if status is False:
        return message
    else:
        return message

@router.get("/logs")
async def fetch_logs(request: Request):
    message, status = middleware(request, needs_admin=True)

    if status is False:
        return message
    else:
        return db.get_logs()

@router.get("/users/fetch")
async def fetch_users(request: Request):
    message, status = middleware(request, needs_admin=False)

    if status is False:
        return message
    else:
        return db_forum.get_users()

@router.post("/users/ban")
async def ban_user(ban: BanUser, request: Request):
    message, status = middleware(request, needs_admin=False)

    if status is False:
        return message
    else:
        return db_forum.ban_user(ban.id, ban.ban_reason, ban.ban_expiry)

@router.post("/users/unban")
async def unban_user(ban: BanUser, request: Request):
    message, status = middleware(request, needs_admin=False)

    if status is False:
        return message
    else:
        return db_forum.unban_user(ban.id)



@router.post("/users/mute")
async def mute_user(ban: BanUser, request: Request):
    message, status = middleware(request, needs_admin=False)

    if status is False:
        return message
    else:
        return db_forum.mute_user(ban.id, ban.ban_reason, ban.ban_expiry)

@router.post("/users/unmute")
async def unmute_user(ban: BanUser, request: Request):
    message, status = middleware(request, needs_admin=False)

    if status is False:
        return message
    else:
        return db_forum.unmute_user(ban.id)

@router.get("/threads/fetch")
async def fetch_users(request: Request):
    message, status = middleware(request, needs_admin=False)

    if status is False:
        return message
    else:
        return db_forum.get_threads()

@router.post("/threads/delete")
async def fetch_users(thread: DeleteThread, request: Request):
    message, status = middleware(request, needs_admin=False)

    if status is False:
        return message
    else:
        return db_forum.delete_thread(thread.id)


@router.post("/users/ipban")
async def fetch_users(ban: BanIP, request: Request):
    message, status = middleware(request, needs_admin=False)

    if status is False:
        return message
    else:
        return db_forum.ip_ban(ban.ip_address, ban.ban_expiry, request.headers.get("Authorization"), ban.ban_reason)

@router.post("/users/ipunban")
async def fetch_users(ban: UnbanIP, request: Request):
    message, status = middleware(request, needs_admin=False)

    if status is False:
        return message
    else:
        return db_forum.ip_unban(ban.id)

@router.get("/users/fetch_ipbans")
async def fetch_users(request: Request):
    message, status = middleware(request, needs_admin=False)

    if status is False:
        return message
    else:
        return db_forum.get_ip_bans()

@router.get("/threads/fetch_sections")
async def fetch_users(request: Request):
    message, status = middleware(request, needs_admin=True)

    if status is False:
        return message
    else:
        return db_forum.get_thread_sections()

@router.post("/threads/create_section")
async def fetch_users(thread: CreateThreadSection, request: Request):
    message, status = middleware(request, needs_admin=True)

    if status is False:
        return message
    else:
        return db_forum.create_section(thread.category, thread.parent, thread.name)

