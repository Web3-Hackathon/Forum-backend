from argon2 import PasswordHasher
 
pw = PasswordHasher()

print(pw.hash("darko1234"))