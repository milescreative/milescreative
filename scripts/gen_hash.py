import bcrypt
password = b"gofdAs-ximtis-wybri7" # DO NOT REMOVE THE b
hashed = bcrypt.hashpw(password, bcrypt.gensalt())
print(f'The hashed password is: {hashed.decode()}')

docker_interpolation= hashed.decode().replace("$", "$$")
print(f'The hashed password for an docker env is: {docker_interpolation}')
