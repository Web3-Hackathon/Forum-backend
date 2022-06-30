import json 


def get_config():
    f = open("./config.json")
    data = json.load(f)
    f.close()
    return data 