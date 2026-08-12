# Encode Monkey values as JSON.

profile = {
  "name": "Hammed",
  "active": true,
  "ratio": 0.75,
  "scores": [10, 20, 30],
  "metadata": null
};

println(json_encode(profile));
