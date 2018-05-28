Linker OAuth Service
===========================

Requirement
-----------

- make
- golang
- [Govendor](https://github.com/kardianos/govendor)




Bootstrap Go and Install Go
---------------------------

Install Go1.4 for bootstraping:

    git clone https://go.googlesource.com/go ~/go1.4
    cd ~/go1.4
    git checkout -t origin/release-branch.go1.4
    cd src && ./all.bash
    export GOROOT_BOOTSTRAP=~/go1.4

Install Go1.8

    cp -a ~/go1.4 ~/go
    cd ~/go
    git checkout -t origin/release-branch.go1.8
    cd src && ./all.bash

Setting up Paths

    export GOROOT=~/go
    export PATH=~/go/bin:$PATH


Build
---------------------------

1. Sync golang dependencies

    ```
    govendor sync
    ```

2. Build the binary

    ```
    make
    ```

    To build with go build flags:

    ```
    make GO_BUILD_FLAGS="-x"
    ```


3. Change properties in `config/default.json`

    ```
    config config/default.json
    ```

4. Setup portal static files


    scripts/setup-static-files


5. Run

    ```
    ./lnk-auth 8.8.8.8:3031
    ./lnk-auth :3031
    ./lnk-auth -h 8.8.8.8 -p 3031
    ```

Run Tests
---------
1. first of all, you shoud provide a MongoDB env. I suggest using Docker.
2. config mongodb host, port and other configs in `oauth.properties`.
3. run `make test`

To run tests, please setup your config file and run `make test` with the environment variable `OAUTH_CONFIG_PATH`:

    cp -v $PWD/config/test.json $PWD/config/testing_local.json
    vim $PWD/config/testing_local.json
    make test OAUTH_CONFIG_PATH="$PWD/config/testing_local.json"

API
---------

`/me`: Returns the current user info by the given access token. This API can be used
for the external services.

`/v1/signup`:

`/v1/signin`:

`/v1/dev/auth`

Schema
---------

### User Schema
- `_id` (bson.ObjectId): mongo default object id
- `serial_number` (string): device serial number
- `email` (string)
- `password` (string)
- `first_name` (string)
- `last_name` (string)
- `country_code` (string)
- `cellphone` (string)
- `verified` (boolean)
- `verification_code` (string)
- `jwt` (string): authentication token (tbd ??)
- `access_token` (string)
- `access_token_expiry_time` (string)
- `refresh_token` (string)
- `created_at` (timestamp)
- `updated_at` (timestamp)

### Applications Schema
- `_id` (bson.ObjectId): mongo default object id
- `name` (string): The name of the application.
- `description` (string): The description of the application.
- `redirect_uri` (string): oauth app redirect uri
- `client_id` (string): The client ID, could be generated as an UUID
- `client_secret` (string): The secret of the client (algorithm TBD)
- `created_at` - (timestamp)
- `updated_at` - (timestamp)
- `user_data` (interface{})

## Reference

- [理解OAuth 2.0](http://www.ruanyifeng.com/blog/2014/05/oauth_2_0.html)
- [OAuth 认证流程详解](http://www.jianshu.com/p/0db71eb445c8)
- [使用 OAuth 2.0 访问豆瓣 API](https://developers.douban.com/wiki/?title=oauth2)
- [golang rsa generate key](https://gist.github.com/sdorra/1c95de8cb80da31610d2ad767cd6f251)
- [oauth: The Client ID and Secret](https://www.oauth.com/oauth2-servers/client-registration/client-id-secret/)
- [DigitalOcean OAuth API Overview](https://developers.digitalocean.com/documentation/oauth/)
