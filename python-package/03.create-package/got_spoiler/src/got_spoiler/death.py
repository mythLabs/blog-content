import random

def can(name):
    res = bool(random.getrandbits(1))
    if res == True:
        ep = random.randrange(1,9)
        print(f"{name} is going to die in episode {ep}")
    else:
        print(f"By the grace of old gods, {name} will survive")