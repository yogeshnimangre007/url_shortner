The goal of this project is to create an http.Handler that will look at the path of any incoming web request and determine if it should redirect the user to a new page, much like URL shortener would.

For instance, if we have a redirect setup for /doc to https://www.somesite.com/a-technology-blog we would look for any incoming web requests with the path /doc and redirect them.

I have done this using both map and yaml.

After I  got the YAML parsing down, try to convert the data into a map and then use the MapHandler to finish the YAMLHandler implementation.



