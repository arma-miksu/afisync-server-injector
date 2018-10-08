# AFISync Server Injector
This project is designed for **AWS Lambda function** that returns an modified JSON
file containing your configured server server included.

For instance, if the original [AFI Sync](https://github.com/haikion/afi-sync)repository has three repositories, you'll
see the fourth one included.

Configure your injected server with following environment variables:

- `AFISYNC_SRV_INJ_SOURCE_UPDATE_URL`:
    The URL address of original `repositories.json` file
- `AFISYNC_SRV_INJ_REPLACE_UPDATE_URL`:
    The URL address of modified `repositories.json` file
- `AFISYNC_SRV_INJ_SOURCE_REPOSITORY_NAME`:
    Name of the repository, which will be copied as a base for injected server.
- `AFISYNC_SRV_INJ_TARGET_REPOSITORY_NAME`:
    Name of the injected repository name.
- `AFISYNC_SRV_INJ_TARGET_REPOSITORY_SERVER_ADDRESS`:
    Server address for injected server.
- `AFISYNC_SRV_INJ_TARGET_REPOSITORY_SERVER_PORT`:
    Server port for injected server. Numerical value.
- `AFISYNC_SRV_INJ_TARGET_REPOSITORY_PASSWORD`:
    Server password for injected server.
- `AFISYNC_SRV_INJ_TARGET_REPOSITORY_BATTL_EYE_ENABLED`:
    Server battle eye setting for injected server. Boolean true or false.

Example:

```
AFISYNC_SRV_INJ_SOURCE_UPDATE_URL=http://armafinland.fi/afisync/repositories.json
AFISYNC_SRV_INJ_REPLACE_UPDATE_URL=https://api-id.execute-api.region.amazonaws.com/default/afisyncServerInjector
AFISYNC_SRV_INJ_SOURCE_REPOSITORY_NAME=armafinland.fi Primary
AFISYNC_SRV_INJ_TARGET_REPOSITORY_NAME=Random rähinät servu
AFISYNC_SRV_INJ_TARGET_REPOSITORY_SERVER_ADDRESS=127.0.0.1
AFISYNC_SRV_INJ_TARGET_REPOSITORY_SERVER_PORT=2302
AFISYNC_SRV_INJ_TARGET_REPOSITORY_PASSWORD=rahka123
AFISYNC_SRV_INJ_TARGET_REPOSITORY_BATTL_EYE_ENABLED=false
```