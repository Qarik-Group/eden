* Switched env vars to match alternate CLI [sb-cli](https://github.com/cppforlife/sb-cli)

    ```
    export SB_BROKER_URL=https://mybroker.com
    export SB_BROKER_USERNAME=username
    export SB_BROKER_PASSWORD=password
    ```

* [CI](https://ci.starkandwayne.com/teams/main/pipelines/eden) now includes `test-pr` job to automatically build + test pull requests
