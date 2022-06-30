from fastapi import FastAPI 
from fastapi.middleware.cors import CORSMiddleware
from routers import admin
from database.database import Database
from sessions import session
import json 

db = Database()
app = FastAPI(title="Crypto API", redoc_url=None, docs_url=None, description="Celestine Public API")


app.include_router(admin.router)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
