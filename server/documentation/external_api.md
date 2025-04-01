# Project Kallaxy external - 3rd party API resources

## Open Library API (books)

## The Movie DB (movies and tv shows)

?language=en-US
?language=fr-FR



https://api.themoviedb.org/3/search/movie  
Search for movies by their original, translated and alternative titles.

query parameters : 
* query ?q= (string) REQUIRED
* include_adult (bool)
* language (string)
* primary_release_year (string)
* region region=FR
* year


https://api.themoviedb.org/3/search/tv  
Search for TV shows by their original, translated and also known as names.
query parameters :
* query ?q=
* first_air_date_year (int32)
* include_adult
* language
* year

https://api.themoviedb.org/3/search/multi  
Use multi search when you want to search for movies, TV shows and people in a single request.  
* query ?q=
* include_adult
* language

https://api.themoviedb.org/3/find/{external_id}
with ?external_source query parameter (required)
imdb_id / tvdb_id / wikidata_id
and optionnal language query parameter

https://api.themoviedb.org/3/movie/{movie_id}
Get the top level details of a movie by ID.
query parameters:
* language=fr-FR

## RAWG (video games)

https://api.rawg.io/api/games
Get a list of games.

query parameters:
* search string
* search_precise bool
* search_exact bool
* platforms string
* stores string
* ordering string (Available fields: name, released, added, created, updated, rating, metacritic)
* exclude_additions bool (no dlc)

platforms: 4=PC / 21=Android
stores: 1=Steam / 8=Google play

https://api.rawg.io/api/games/{game_id}
Get details of the game. 

## BGG API (board games)

https://boardgamegeek.com/xmlapi2/search
Get a list of board games

query parameters:
* query string (the search)
* type string  (TYPE might be rpgitem, videogame, boardgame, boardgameaccessory, boardgameexpansion or boardgamedesigner)
* exact int (1 to limit results to exact match)

https://boardgamegeek.com/xmlapi2/thing
Get details for a thing (can be boardgame, boardgameexpansion, boardgameaccessory, videogame, rpgitem, rpgissue (for periodicals))

query parameters:
* id
* type