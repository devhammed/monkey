# Decode JSON into Monkey values.

payload = json_decode('{"user":{"name":"Hammed"},"roles":["admin","author"]}');

println(payload["user"]["name"]);
println(payload["roles"]);
