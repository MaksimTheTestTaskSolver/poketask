poketask is a server which will merge a random cat image with pokemon image.

To run the server execute `./poketsk` from the repository directory. It will start a server on `:8080` port

binary is built for `darwin`, if you use another operating system please use provided Dockerfile and make commands `make build` and `make run`

You can use the `localhost:8080/api/v1/pokemon/:pokemonId` to get the image of desired pokemon merged with an image of a random cat.

I don't know at the moment how to properly limit the amount of concurrent requests to the target APIs, so I wrote a simple request limiter as a quick solution to make it in time. Also, there are no tests for the same reason.