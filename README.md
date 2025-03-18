# Golang onedrive navigator
This is onedrive folder/files navigator used with 'cd' to change directory and 'ls' to list current directory information. `cd <directory>` will move you to that directory and so forth.


## Environment variables

In the projects root directory create a `.env` file, add the `MS_OPENGRAPH_APP_ID` and `MS_OPENGRAPH_APP_SECRET` and `MS_OPENGRAPH_APP_TENANT_ID` variables with their respective values, also do not forget to establish the correct Redirect URI if it is not `https://localhost:3000`.
