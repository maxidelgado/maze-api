# Maze API

This API exposes the logic for both creating a maze (with its corresponding quadrants, spots and paths) and
the endpoints to start a game and move around an existing maze, until you arrive to the exit spot.

### How to run

Just run the following command

```bash
$ docker-compose up -d
```

Note: to clear the database, you can run:
```bash
$ docker-compose restart
```

### How to use

#### Create a maze

You can create a maze as follows:
```bash 
$ curl --location --request POST 'localhost:3000/api/v1/mazes' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "name": "amazing maze",
      "center": [ // Optional
          0,
          0
      ],
      "spots": [ // Optional
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
      "paths": [ // Optional
          {
              "origin": [1,1],
              "destiny": [-1,1]
          }
      ]
  }'
```
_IMPORTANT: a playable maze must contain an entrance and exit spots, and they both
must be connected by any path_

#### Update a maze

You can:
    - Move quadrants
    - Add spots
    - Add paths
    
```bash
$ curl --location --request PUT 'localhost:3000/api/v1/mazes/96d9a144-ac8d-497c-bc5a-248012d7687d' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "center": [-1,4],
      "spots": [
          {
              "name": "another spot",
              "gold_amount": 0,
              "coordinate": [2,-2]
          }
      ],
      "paths": [
           {
              "origin": [1,1],
              "destiny": [2,-2]
          }
      ]
  }'
```

#### Get an existing maze

If you want to get maze details, you can get it by an id or by using the name:

By id:
```bash
$ curl --location --request GET 'localhost:3000/api/v1/mazes/96d9a144-ac8d-497c-bc5a-248012d7687d'
```

By name: (will return any maze matching the name)
```bash
$ curl --location --request GET 'localhost:3000/api/v1/mazes?name=test'
```

#### Delete a maze

```bash
$ curl --location --request DELETE 'localhost:3000/api/v1/mazes/96d9a144-ac8d-497c-bc5a-248012d7687d'
```

#### Delete a spot or a path

Delete spot:
```bash
$ curl --location --request DELETE 'localhost:3000/api/v1/mazes/96d9a144-ac8d-497c-bc5a-248012d7687d/spot' \
--header 'Content-Type: application/json' \
--data-raw '[1,1]'
```

Delete path:
```bash
$ curl --location --request DELETE 'localhost:3000/api/v1/mazes/96d9a144-ac8d-497c-bc5a-248012d7687d/path' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "origin": [-1,1],
      "destiny": [-1,-1]
  }'
```

#### Create a game

You can create a new game by providing the id of the selected maze, and the name of the game:

```bash
$ curl --location --request POST 'localhost:3000/api/v1/games' \
--header 'Content-Type: application/json' \
--data-raw '{
    "maze_id": "2b267c65-107a-42e2-8343-b9a53dcd8492",
    "name": "amazing game"
}'
```

You will be positioned at the entrance of the maze, so you can start moving from here by using the "allowed_movements" field. For example:
```json
{
    "id": "a4b4abde-ac4a-4ce6-a1c9-e66cd7717b54",
    "name": "test 2",
    "minimum_distance": 24.60112615949154,
    "player_stats": {
        "total_gold": 0,
        "distance_covered": 0,
        "current_spot": "[1,1]",
        "allowed_movements": [
            {
                "key": "[1,-1]",
                "name": "spot 1"
            },
            {
                "key": "[9,2]",
                "name": "spot 2"
            },
            {
                "key": "[-7,5]",
                "name": "spot 3"
            }
        ]
    },
    "start_date": "2020-10-31T00:18:06.392391141-03:00",
    "end_date": "0001-01-01T00:00:00Z"
}
```
So you can move to (1,-1), (9,2) and (-7,5)

#### Moving

You can perform a movement to the next spot as follows:
```bash
$ curl --location --request PUT 'localhost:3000/api/v1/games/a4b4abde-ac4a-4ce6-a1c9-e66cd7717b54/move' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "spot": "[9,2]"
  }'
```
From here you should repeat until you arrive to the "exit" spot.

#### Delete a game

```bash
$ curl --location --request DELETE 'localhost:3000/api/v1/games/96d9a144-ac8d-497c-bc5a-248012d7687d'
```

#### Get a game

If you want to get game details, you can get it by an id or by using the name:

By id:
```bash
$ curl --location --request GET 'localhost:3000/api/v1/games/96d9a144-ac8d-497c-bc5a-248012d7687d'
```

By name: (will return any game matching the name)
```bash
$ curl --location --request GET 'localhost:3000/api/v1/games?name=test'
```

_Note: more examples about the other CRUD operations [here](examples)_

_Note 2: unit testing has low coverage and was added as a demonstration_