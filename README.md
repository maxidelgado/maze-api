# Maze API

This API exposes the logic for both creating a maze (with its corresponding quadrants, spots and paths) and
the endpoints to start a game and move around an existing maze, until you arrive to the exit spot.

### How to run

Just run the following command

```sh
$ docker-compose up -d
```

Note: to clear the database, you can run:
```sh
$ docker-compose restart
```

### How to test

First you need to create a maze as follows:
```sh 
$ curl --location --request POST 'localhost:3000/api/v1/mazes' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "center": [
          0,
          0
      ],
      "spots": [
          {
              "name": "entrance",
              "gold_amount": 0,
              "coordinate": [1,1]
          },
          {
              "name": "exit",
              "gold_amount": 50,
              "coordinate": [-1,1]
          }
      ],
      "paths": [
          {
              "origin": [1,1],
              "destiny": [-1,1]
          }
      ]
  }'
```
_IMPORTANT: playable maze must contain an entrance and exit spots, and they both
must be connected by any path_

Then, you can start a new game by using the maze_id, returned from the previous request:
```sh 
$ curl --location --request POST 'localhost:3000/api/v1/games' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "maze_id": "17b11a39-a969-4f0a-bf31-4af77a959cc4"
  }'
```

Now you can start moving around the maze based on the game info, for example:
```json
{
    "id": "7f2ced2a-9c4e-47d7-b9fd-19c6c5518d5a",
    "entrance": "[1,1]",
    "exit": "[-1,1]",
    "minimum_distance": 6,
    "player_stats": {
        "total_gold": 0,
        "distance_covered": 0,
        "current_spot": "[1,1]",
        "movements": null,
        "allowed_movements": [
            "[1,-1]"
        ]
    },
    "start_date": "2020-10-30T00:14:10.984025909-03:00",
    "end_date": "0001-01-01T00:00:00Z"
}
```

You should use the "allowed_movements" field to perform your next movement around the maze, for example:
```sh
$ curl --location --request PUT 'localhost:3000/api/v1/games/7f2ced2a-9c4e-47d7-b9fd-19c6c5518d5a/move' \
--header 'Content-Type: application/json' \
--data-raw '{
    "spot": "[1,-1]"
}'
```

Repeat the last step until you arrived to the exit spot

_Note: more examples about the other CRUD operations [here](example.rest)_

_Note 2: unit testing has low coverage and was added as a demonstration_