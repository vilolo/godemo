import requests
import json

url = 'http://api-robot.mozi.local/v1/user/login'
data = {
    'account': 'mozi',
    'password': '000000',
    'device_id': '1',
    'os_ver': '1'
}

def send_post(url, data):
    return requests.post(url=url, data=data).json()
    # json.dumps(result, indent=4)

b = send_post(url, data)
#b = json.loads(a)
print(b["data"]["token"])
