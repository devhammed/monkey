# Decode JSON into Monkey values.

payload = json_decode('{"user":{"name":"Hammed"},"ratio":0.75,"roles":["admin","author"]}');

println(payload["user"]["name"]);
println(payload["ratio"]);
println(payload["roles"]);
