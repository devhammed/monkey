# Encode Monkey values as JSON.

profile = {
  "name": "Hammed",
  "active": true,
  "scores": [10, 20, 30],
  "metadata": null
};

println(json_encode(profile));
