# Spotify Authorization Flow Server

The Spotify oauth2 authorization flow requires both the client id and client secret.
As we need to keep the client secret secure, we need a server to manage the fetching
of access_token, refresh_token on behalf of the [spotify-filter](https://github.com/atinm/spotify-filter)
application. This is the authorization flow server code for handling the code to token exchange on behalf
of the spotify-filter application.

# Building and Running

Start by registering your application at [Spotify Application Registration](https://developer.spotify.com/my-applications/).

Set the Redirect URI to be `http://localhost:5009/callback` (this can be changed if you intend
to run the server at a different address in the config.json file described later) which is the port
where you will run this authorization server on, and Save. You'll get a client ID
and secret key for your application.

Export the `SPOTIFY_ID` and `SPOTIFY_SECRET` environment variables set to the client id and
client secret you created above at application registration to make them available to the program
(or you can use the configration file `config.json` to save it as described under the Configuration section).

The server uses HTTPS, therefore you need to provide a cert.pem and key.pem for it to work. You can generate
these using the generate_cert program in crypto/tls.

    go get github.com/atinm/spotify-auth-server
    go build
    export SPOTIFY_ID=<the client id from the Spotify application registration>
    export SPOTIFY_SECRET=<the client secret from the Spotify application registration>
    ./spotify-auth-server

The program will run waiting for the Spotify authorization server to redirect authorization requests
to it. The redirected request will contain the code and state that were generated by the spotify-filer app for
your application and the server will exchange the code for an access_token and refresh_token that will be sent back
to the spotify-filter application.

# Configuration

You may have a config.json file in the same directory as the program:

    {
        "client_id": "<the client id from the Spotify application registration>,
        "client_secret": "<the client secret from the Spotify application registration>,
        "log_level": "DEBUG|INFO|WARN(default)|ERROR",
        "application_uri": "<the uri where the application is listening runs - default https://localhost:5007/callback>",
        "my_uri":  "<the uri where this server is listening - default https://localhost:5009/callback>",
        "cert": "cert.pem",
        "key": "key.pem"
    }
