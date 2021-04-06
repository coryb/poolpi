Import("env")
import platform
import os
import requests
import urllib.request
import re
import tarfile
import io
tool = env.get("PROJECT_DIR") + '/mklittlefs/mklittlefs'
if not os.path.exists(tool):
    platform = platform.system().lower()
    r = requests.get(
        "https://github.com/earlephilhower/mklittlefs/releases/latest"
    )
    links = re.findall(r'/[^"]+[.]tar[.]gz', r.text)
    for l in links:
        if platform in l:
            print("Fetching URL: " + "https://github.com" + l)
            stream = urllib.request.urlopen("https://github.com" + l)
            tar = tarfile.open(fileobj=stream, mode="r|gz")
            tar.extractall()
            tar.close()
            break

env.Replace(MKSPIFFSTOOL=tool)
