# Decode JSON into Monkey values.

payload = json_decode('{"user":{"name":"Hammed"},"ratio":0.75,"roles":["admin","author"]}');

print(payload["user"]["name"]);
print(payload["ratio"]);
print(payload["roles"]);
