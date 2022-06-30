from pydantic import BaseModel

class Authenticate(BaseModel):
    username: str
    password: str 

class BanUser(BaseModel):
    id: int
    ban_reason: str 
    ban_expiry: int 

class DeleteThread(BaseModel):
    id: int

class BanIP(BaseModel):
    ip_address: str
    ban_expiry: int
    ban_reason: str

class UnbanIP(BaseModel):
    id: int 

class CreateThreadSection(BaseModel):
    category: str 
    parent: str
    name: str